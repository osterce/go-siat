package siat

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ron86i/go-siat/internal/adapter/service"
)

// siatServices agrupa todas las implementaciones de los servicios del SIAT
// accesibles a través de un único punto de entrada.
type siatServices struct {
	Operaciones    *service.SiatOperacionesService
	Sincronizacion *service.SiatSincronizacionService
	Codigos        *service.SiatCodigosService
	CompraVenta    *service.SiatCompraVentaService
}

// New crea e inicializa una nueva instancia de los servicios del SIAT.
// Requiere la URL base del servicio (Pruebas o Producción) y un cliente HTTP opcional.
// Si httpClient es nil, se utilizará uno por defecto con un timeout de 15 segundos.
func New(baseUrl string, httpClient *http.Client) (*siatServices, error) {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 15 * time.Second,
		}
	}

	baseUrl = strings.TrimSpace(baseUrl)
	if baseUrl == "" {
		return nil, fmt.Errorf("la URL base del SIAT no puede estar vacía")
	}

	operaciones, err := service.NewSiatOperacionesService(baseUrl, httpClient)
	if err != nil {
		return nil, err
	}
	sincronizacion, err := service.NewSiatSincronizacionService(baseUrl, httpClient)
	if err != nil {
		return nil, err
	}
	codigos, err := service.NewSiatCodigosService(baseUrl, httpClient)
	if err != nil {
		return nil, err
	}
	compraVenta, err := service.NewSiatCompraVentaService(baseUrl, httpClient)
	if err != nil {
		return nil, err
	}
	return &siatServices{
		Operaciones:    operaciones,
		Sincronizacion: sincronizacion,
		Codigos:        codigos,
		CompraVenta:    compraVenta,
	}, nil
}
