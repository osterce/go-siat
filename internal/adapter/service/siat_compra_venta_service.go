package service

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/ron86i/go-siat/internal/core/domain/datatype/soap"
	"github.com/ron86i/go-siat/internal/core/domain/facturacion/compra_venta"
	"github.com/ron86i/go-siat/internal/core/port"
	"github.com/ron86i/go-siat/pkg/config"
	"github.com/ron86i/go-siat/pkg/models"
)

type SiatCompraVentaService struct {
	url        string
	HttpClient *http.Client
}

// AnulacionFactura permite anular una factura previamente aceptada por el SIAT.
// Recibe una solicitud opaca de tipo AnulacionFacturaRequest construida vía Builder.
func (s *SiatCompraVentaService) AnulacionFactura(ctx context.Context, config config.Config, opaqueReq any) (*soap.EnvelopeResponse[compra_venta.AnulacionFacturaResponse], error) {
	req := models.GetInternalRequest[compra_venta.AnulacionFactura](opaqueReq)
	xmlBody, err := buildRequest(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", s.url, bytes.NewReader(xmlBody))
	if err != nil {
		return nil, fmt.Errorf("error al crear petición HTTP: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/xml")
	httpReq.Header.Set("apiKey", fmt.Sprintf("TokenApi %s", config.Token))

	resp, err := s.HttpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("error al hacer request HTTP anulacion factura: %w", err)
	}
	return parseSoapResponse[compra_venta.AnulacionFacturaResponse](resp)
}

// RecepcionFactura envía una factura firmada, comprimida y codificada al SIAT para su procesamiento.
// Recibe una solicitud opaca de tipo RecepcionFacturaRequest construida vía Builder.
func (s *SiatCompraVentaService) RecepcionFactura(ctx context.Context, config config.Config, opaqueReq any) (*soap.EnvelopeResponse[compra_venta.RecepcionFacturaResponse], error) {
	req := models.GetInternalRequest[compra_venta.RecepcionFactura](opaqueReq)
	xmlBody, err := buildRequest(req)
	if err != nil {
		return nil, err
	}
	log.Printf("Request XML: %s", string(xmlBody))

	httpReq, err := http.NewRequestWithContext(ctx, "POST", s.url, bytes.NewReader(xmlBody))
	if err != nil {
		return nil, fmt.Errorf("error al crear petición HTTP: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/xml")
	httpReq.Header.Set("apiKey", fmt.Sprintf("TokenApi %s", config.Token))

	resp, err := s.HttpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("error al hacer request HTTP recepcion factura: %w", err)
	}
	return parseSoapResponse[compra_venta.RecepcionFacturaResponse](resp)
}

func NewSiatCompraVentaService(baseUrl string, httpClient *http.Client) (*SiatCompraVentaService, error) {
	baseUrl = strings.TrimSpace(baseUrl)
	if baseUrl == "" {
		return nil, fmt.Errorf("la URL base del SIAT no puede estar vacía")
	}

	// Si no se inyecta un cliente, creamos uno con configuraciones seguras por defecto
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 15 * time.Second,
		}
	}
	return &SiatCompraVentaService{
		url:        fullURL(baseUrl, SiatCompraVenta),
		HttpClient: httpClient,
	}, nil
}

var _ port.SiatCompraVentaService = (*SiatCompraVentaService)(nil)
