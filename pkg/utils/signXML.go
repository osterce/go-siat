package utils

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/beevik/etree"
	dsig "github.com/russellhaering/goxmldsig"
)

// pemKeyStore implementa dsig.X509KeyStore para cargar llaves desde archivos PEM.
type pemKeyStore struct {
	PrivateKey *rsa.PrivateKey
	Cert       []byte
}

func (ks *pemKeyStore) GetKeyPair() (*rsa.PrivateKey, []byte, error) {
	return ks.PrivateKey, ks.Cert, nil
}

// SignXMLBytes firma un documento XML recibiendo los certificados y la llave directamente en bytes.
func SignXMLBytes(xmlBytes, keyBytes, certBytes []byte) ([]byte, error) {
	// 1. Parsear clave privada desde los bytes proporcionados
	privKey, err := parseRSAPrivateKey(keyBytes)
	if err != nil {
		return nil, err
	}

	// 2. Decodificar certificado PEM
	blockCert, _ := pem.Decode(certBytes)
	if blockCert == nil {
		return nil, fmt.Errorf("error decoding PEM certificate")
	}

	// 3. Configure KeyStore
	ks := &pemKeyStore{
		PrivateKey: privKey,
		Cert:       blockCert.Bytes,
	}

	// 4. Configurar contexto de firma
	ctx := dsig.NewDefaultSigningContext(ks)
	ctx.Canonicalizer = dsig.MakeC14N10WithCommentsCanonicalizer()
	ctx.SetSignatureMethod(dsig.RSASHA256SignatureMethod)

	// 5. Parsear XML
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(xmlBytes); err != nil {
		return nil, err
	}

	// 6. Firmar XML (Enveloped Signature)
	signedElement, err := ctx.SignEnveloped(doc.Root())
	if err != nil {
		return nil, err
	}

	signedDoc := etree.NewDocument()
	signedDoc.SetRoot(signedElement)

	// 7. Renderizar a bytes
	var buf bytes.Buffer
	if _, err := signedDoc.WriteTo(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// SignXML firma un documento XML recibiendo los certificados y la llave desde archivos.
func SignXML(xmlBytes []byte, keyPath, certPath string) ([]byte, error) {
	keyData, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	certData, err := os.ReadFile(certPath)
	if err != nil {
		return nil, err
	}

	return SignXMLBytes(xmlBytes, keyData, certData)
}

// parseRSAPrivateKey procesa los bytes de una clave PEM (PKCS#1 o PKCS#8).
func parseRSAPrivateKey(keyData []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, fmt.Errorf("invalid PEM format in private key")
	}

	// Intentar parsear como PKCS#1
	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return key, nil
	}

	// Intentar parsear como PKCS#8
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("loaded private key is not of type RSA")
	}

	return rsaKey, nil
}
