package port

import (
	"context"

	"github.com/ron86i/go-siat/internal/core/domain/datatype/soap"
	"github.com/ron86i/go-siat/internal/core/domain/facturacion/compra_venta"
	"github.com/ron86i/go-siat/pkg/config"
)

type SiatCompraVentaService interface {
	// AnulacionFactura anula una factura previamente enviada al SIAT.
	AnulacionFactura(ctx context.Context, config config.Config, req any) (*soap.EnvelopeResponse[compra_venta.AnulacionFacturaResponse], error)

	// RecepcionFactura envía una factura al SIAT para su procesamiento y validación.
	RecepcionFactura(ctx context.Context, config config.Config, req any) (*soap.EnvelopeResponse[compra_venta.RecepcionFacturaResponse], error)
}
