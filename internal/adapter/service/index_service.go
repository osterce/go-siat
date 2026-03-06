package service

import (
	"encoding/xml"
	"fmt"

	"io"
	"net/http"

	"github.com/ron86i/go-siat/internal/core/domain/datatype/soap"
)

// SiatService define los diferentes servicios disponibles en el SIAT.
type SiatService string

const (
	SiatCodigos        SiatService = "FacturacionCodigos"
	SiatOperaciones    SiatService = "FacturacionOperaciones"
	SiatSincronizacion SiatService = "FacturacionSincronizacion"
	SiatCompraVenta    SiatService = "ServicioFacturacionCompraVenta"
)

// fullURL construye la URL completa para acceder a un servicio específico del SIAT,
// concatenando la URL base del ambiente con el endpoint del servicio solicitado.
func fullURL(baseURL string, service SiatService) string {
	return baseURL + "/" + string(service)
}

// buildRequest encapsula un objeto de solicitud genérico dentro de un sobre SOAP estándar (Envelope),
// añadiendo los namespaces requeridos por el SIAT y serializando el resultado a formato XML.
func buildRequest(req any) ([]byte, error) {
	requestBody := soap.Envelope[any]{
		XmlnsSoapenv: "http://schemas.xmlsoap.org/soap/envelope/",
		XmlnsNs:      "https://siat.impuestos.gob.bo/",
		Body: soap.EnvelopeBody[any]{
			Content: req,
		},
	}

	xmlBody, err := xml.MarshalIndent(requestBody, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error al serializar body SOAP: %w", err)
	}
	return []byte(xml.Header + string(xmlBody)), nil
}

// parseSoapResponse procesa y valida una respuesta HTTP proveniente del servicio para extraer el contenido SOAP esperado.
func parseSoapResponse[T any](resp *http.Response) (*soap.EnvelopeResponse[T], error) {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error al leer el cuerpo de la respuesta: %w", err)
	}

	var result soap.EnvelopeResponse[T]

	// Intentar parsear la respuesta XML en la estructura de respuesta SOAP
	errUnmarshal := xml.Unmarshal(body, &result)

	// Si el servicio devolvió un SOAP Fault, priorizar este error descriptivo de negocio
	if errUnmarshal == nil && result.Body.Fault != nil {
		return nil, fmt.Errorf("SOAP Fault [%s]: %s", result.Body.Fault.FaultCode, result.Body.Fault.FaultString)
	}

	// Si el código de estado HTTP no es 200, informar el error de estado.
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP status inesperado: %d", resp.StatusCode)
	}

	// Si el status es 200 pero hubo un error de parseo, informar el error de XML
	if errUnmarshal != nil {
		return nil, fmt.Errorf("error al parsear respuesta SOAP: %w", errUnmarshal)
	}

	return &result, nil
}
