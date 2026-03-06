package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"log"
	"time"

	"github.com/ron86i/go-siat"
	"github.com/ron86i/go-siat/internal/core/domain/datatype"
	"github.com/ron86i/go-siat/pkg/config"
	"github.com/ron86i/go-siat/pkg/models"
)

func main() {
	// 1. Inicializar el servicio SIAT
	s, err := siat.New("https://pilotosiatservicios.impuestos.gob.bo/v2", nil)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	cvService := s.CompraVenta

	ctx := context.Background()
	cfg := config.Config{Token: "TU_TOKEN_API"}

	nit := int64(123456789)
	fechaEmision := time.Now()
	codigoControl := "XYZ789"

	// 2. Generar CUF usando el helper en el namespace de modelos
	cuf, err := models.CompraVenta.GenerarCUF(nit, fechaEmision, 0, 1, 1, 1, 1, 1, 0, codigoControl)
	if err != nil {
		log.Fatalf("Error generando CUF: %v", err)
	}

	// 3. Construir la Factura usando Builders (Siempre)
	cabeceraReq := models.CompraVenta.NewCabecera().
		WithNitEmisor(nit).
		WithRazonSocialEmisor("Mi Empresa S.A.").
		WithMunicipio("La Paz").
		WithNumeroFactura(1).
		WithCuf(cuf).
		WithCufd("CODIGO_CUFD_EJEMPLO").
		WithCodigoSucursal(0).
		WithDireccion("Av. Principal 123").
		WithCodigoPuntoVenta(0).
		WithFechaEmision(fechaEmision.Format("2006-01-02T15:04:05.000")).
		WithNombreRazonSocial("JUAN PEREZ").
		WithCodigoTipoDocumentoIdentidad(1).
		WithNumeroDocumento("5544332").
		WithCodigoCliente("CLI-001").
		WithCodigoMetodoPago(1).
		WithMontoTotal(100.0).
		WithMontoTotalSujetoIva(100.0).
		WithCodigoMoneda(1).
		WithTipoCambio(1.0).
		WithMontoTotalMoneda(100.0).
		WithLeyenda("Ley N° 453: El proveedor deberá suministrar el servicio...").
		WithUsuario("admin").
		WithCodigoDocumentoSector(1).
		Build()

	detalleReq := models.CompraVenta.NewDetalle().
		WithActividadEconomica("461000").
		WithCodigoProductoSin("12345").
		WithCodigoProducto("PROD-001").
		WithDescripcion("Producto de prueba").
		WithCantidad(1.0).
		WithUnidadMedida(57).
		WithPrecioUnitario(100.0).
		WithSubTotal(100.0).
		Build()

	facturaReq := models.CompraVenta.NewFactura().
		WithCabecera(cabeceraReq).
		AddDetalle(detalleReq).
		Build()

	// 4. Serialización, Firma, Compresión y Hash
	xmlData, err := xml.Marshal(facturaReq)
	if err != nil {
		log.Fatalf("Error al serializar XML: %v", err)
	}

	// A. Firmar XML usando el helper en el namespace de modelos
	signedXML, err := models.CompraVenta.SignXML(xmlData, "key.pem", "cert.crt")
	if err != nil {
		fmt.Println("Omitiendo firma real por falta de certificados físicos...")
		signedXML = xmlData
	}

	// B. Comprimir con Gzip
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	if _, err := zw.Write(signedXML); err != nil {
		log.Fatalf("Error comprimiendo XML: %v", err)
	}
	zw.Close()
	compressedBytes := buf.Bytes()

	// C. Calcular Hash SHA256 sobre bytes COMPRIMIDOS
	hash := sha256.Sum256(compressedBytes)
	hashString := hex.EncodeToString(hash[:])

	// D. Codificar a Base64 para el envío
	encodedArchivo := base64.StdEncoding.EncodeToString(compressedBytes)

	// 5. Construir Solicitud de Recepción usando el Wrapper
	recepcionReq := models.CompraVenta.NewRecepcionFacturaRequest().
		WithCodigoAmbiente(2).
		WithCodigoDocumentoSector(1).
		WithCodigoEmision(1).
		WithCodigoModalidad(1).
		WithCodigoPuntoVenta(0).
		WithCodigoSistema("ABC123DEF").
		WithCodigoSucursal(0).
		WithCufd("CODIGO_CUFD_EJEMPLO").
		WithCuis("C2FC682B").
		WithNit(nit).
		WithTipoFacturaDocumento(1).
		WithArchivo([]byte(encodedArchivo)).
		WithFechaEnvio(datatype.TimeSiat(fechaEmision)).
		WithHashArchivo(hashString).
		Build()

	// 6. Enviar al SIAT
	resp, err := cvService.RecepcionFactura(ctx, cfg, recepcionReq)
	if err != nil {
		log.Fatalf("Error en recepción factura: %v", err)
	}

	if resp != nil && resp.Body.Content.RespuestaServicioFacturacion.Transaccion {
		fmt.Printf("Éxito! Código Recepción: %s\n", resp.Body.Content.RespuestaServicioFacturacion.CodigoRecepcion)
	} else {
		fmt.Printf("Error en recepción factura: %v\n", resp.Body.Content.RespuestaServicioFacturacion)
	}

}
