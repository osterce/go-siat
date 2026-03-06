package models

// requestWrapper es una envoltura genérica utilizada para ocultar la implementación concreta
// de una solicitud y satisfacer las interfaces opacas del SDK.
type requestWrapper[T any] struct {
	request *T
}

// commonRequest es el método marcador que permite a requestWrapper satisfacer
// las interfaces opacas definidas en los paquetes de modelos.
func (r requestWrapper[T]) commonRequest() {
}

// GetInternalRequest desempaqueta la estructura de solicitud concreta desde una interfaz opaca.
// Este método es utilizado internamente por los servicios para acceder a los campos de la solicitud.
// Soporta tanto envolturas (wrappers) como punteros directos para mayor flexibilidad.
func GetInternalRequest[T any](req any) *T {
	if wrapper, ok := req.(requestWrapper[T]); ok {
		return wrapper.request
	}
	if res, ok := req.(*T); ok {
		return res
	}
	return nil
}
