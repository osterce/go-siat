package facturas

import (
	"encoding/xml"
)

// requestWrapper es una envoltura genérica utilizada para ocultar la implementación concreta
// de una solicitud y satisfacer las interfaces opacas del SDK.
type requestWrapper[T any] struct {
	request *T
}

// MarshalXML implementa la interfaz xml.Marshaler para delegar la serialización
// al objeto interno, evitando que la etiqueta raíz sea "requestWrapper".
func (r requestWrapper[T]) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.Encode(r.request)
}

// getInternalRequest desempaqueta la estructura de solicitud concreta desde una interfaz opaca.
// Este método es utilizado internamente por los servicios para acceder a los campos de la solicitud.
// Soporta tanto envolturas (wrappers) como punteros directos para mayor flexibilidad.
func getInternalRequest[T any](req any) *T {
	if wrapper, ok := req.(requestWrapper[T]); ok {
		return wrapper.request
	}
	if res, ok := req.(*T); ok {
		return res
	}
	return nil
}

const (
	// ModalidadElectronica requiere firma digital de los documentos XML.
	ModalidadElectronica = 1
	// ModalidadComputarizada no requiere firma digital, usa un código de control.
	ModalidadComputarizada = 2
	// AmbienteProduccion para operaciones reales con validez tributaria.
	AmbienteProduccion = 1
	// AmbientePruebas para entornos de desarrollo y certificación.
	AmbientePruebas = 2
	// EmisionOnline emisión se realizó en línea
	EmisionOnline = 1
	// EmisionOffline emisión se realizó fuera de línea
	EmisionOffline = 2
	// EmisionMasiva para emisión masiva de factura
	EmisionMasiva = 3
)
