package facturas

import (
	"encoding/xml"
	"strconv"
	"time"

	"github.com/ron86i/go-siat/internal/core/domain/datatype"
	"github.com/ron86i/go-siat/internal/core/domain/models_facturas"
)

// FacturaCompraVenta representa la estructura completa de una factura lista para ser procesada.
// Esta interfaz oculta los detalles de implementación del dominio interno.
type FacturaCompraVenta interface{}

// FacturaCompraVentaCabecera representa la sección de cabecera de una factura de compra y venta.
type FacturaCompraVentaCabecera interface{}

// FacturaCompraVentaDetalle representa un ítem individual dentro del detalle de una factura.
type FacturaCompraVentaDetalle interface{}

// NewFacturaCompraVentaBuilder inicia el proceso de construcción de una FacturaCompraVenta,
// configurando los nombres de nodo XML necesarios para el SIAT.
//
// Por defecto, la factura se configura para ser emitida electrónica.
func NewFacturaCompraVentaBuilder() *FacturaCompraVentaBuilder {
	return &FacturaCompraVentaBuilder{
		factura: &models_facturas.FacturaCompraVenta{
			XMLName:           xml.Name{Local: "facturaElectronicaCompraVenta"},
			XmlnsXsi:          "http://www.w3.org/2001/XMLSchema-instance",
			XsiSchemaLocation: "facturaElectronicaCompraVenta.xsd",
		},
	}
}

// NewFacturaCompraVentaCabeceraBuilder crea una nueva instancia del constructor para la cabecera
// de facturas de compra y venta.
func NewFacturaCompraVentaCabeceraBuilder() *FacturaCompraVentaCabeceraBuilder {
	return &FacturaCompraVentaCabeceraBuilder{
		cabecera: &models_facturas.Cabecera{},
	}
}

// NewFacturaCompraVentaDetalleBuilder crea una nueva instancia del constructor para los ítems
// de detalle de la factura.
func NewFacturaCompraVentaDetalleBuilder() *DetalleBuilder {
	return &DetalleBuilder{
		detalle: &models_facturas.Detalle{},
	}
}

type FacturaCompraVentaBuilder struct {
	factura *models_facturas.FacturaCompraVenta
}

// WithCabecera asocia la cabecera construida previamente a la factura.
func (b *FacturaCompraVentaBuilder) WithCabecera(req FacturaCompraVentaCabecera) *FacturaCompraVentaBuilder {
	if c := getInternalRequest[models_facturas.Cabecera](req); c != nil {
		b.factura.Cabecera = *c
	}
	return b
}

// AddDetalle añade un ítem de detalle a la lista de detalles de la factura.
func (b *FacturaCompraVentaBuilder) AddDetalle(req FacturaCompraVentaDetalle) *FacturaCompraVentaBuilder {
	if d := getInternalRequest[models_facturas.Detalle](req); d != nil {
		b.factura.Detalle = append(b.factura.Detalle, *d)
	}
	return b
}

// WithModalidad configura los metadatos XML de la factura según la modalidad (Electrónica o Computarizada).
func (b *FacturaCompraVentaBuilder) WithModalidad(tipo int) *FacturaCompraVentaBuilder {
	switch tipo {
	case ModalidadElectronica:
		b.factura.XMLName = xml.Name{Local: "facturaElectronicaCompraVenta"}
		b.factura.XsiSchemaLocation = "facturaElectronicaCompraVenta.xsd"
	case ModalidadComputarizada:
		b.factura.XMLName = xml.Name{Local: "facturaComputarizadaCompraVenta"}
		b.factura.XsiSchemaLocation = "facturaComputarizadaCompraVenta.xsd"
	}
	return b
}

// Build finaliza la construcción y retorna la interfaz opaca lista para ser firmada y enviada.
func (b *FacturaCompraVentaBuilder) Build() FacturaCompraVenta {
	return requestWrapper[models_facturas.FacturaCompraVenta]{request: b.factura}
}

// FacturaCompraVentaCabeceraBuilder ayuda a configurar la cabecera de la factura de compra y venta.
type FacturaCompraVentaCabeceraBuilder struct {
	cabecera *models_facturas.Cabecera
}

// WithNitEmisor configura el NIT del emisor de la factura.
func (b *FacturaCompraVentaCabeceraBuilder) WithNitEmisor(nitEmisor int64) *FacturaCompraVentaCabeceraBuilder {
	b.cabecera.NitEmisor = nitEmisor
	return b
}

// WithRazonSocialEmisor configura la razón social o nombre del emisor.
func (b *FacturaCompraVentaCabeceraBuilder) WithRazonSocialEmisor(razonSocialEmisor string) *FacturaCompraVentaCabeceraBuilder {
	b.cabecera.RazonSocialEmisor = razonSocialEmisor
	return b
}

// WithMunicipio configura el municipio donde se emite la factura.
func (b *FacturaCompraVentaCabeceraBuilder) WithMunicipio(municipio string) *FacturaCompraVentaCabeceraBuilder {
	b.cabecera.Municipio = municipio
	return b
}

// WithTelefono configura el número de teléfono (opcional).
func (b *FacturaCompraVentaCabeceraBuilder) WithTelefono(telefono string) *FacturaCompraVentaCabeceraBuilder {
	b.cabecera.Telefono = datatype.Nilable[string]{Value: &telefono}
	return b
}

// WithNumeroFactura configura el número correlativo de la factura.
func (b *FacturaCompraVentaCabeceraBuilder) WithNumeroFactura(numeroFactura int64) *FacturaCompraVentaCabeceraBuilder {
	b.cabecera.NumeroFactura = numeroFactura
	return b
}

// WithCuf configura el Código Único de Factura (CUF).
func (b *FacturaCompraVentaCabeceraBuilder) WithCuf(cuf string) *FacturaCompraVentaCabeceraBuilder {
	b.cabecera.Cuf = cuf
	return b
}

// WithCufd configura el Código Único de Facturación Diaria (CUFD).
func (b *FacturaCompraVentaCabeceraBuilder) WithCufd(cufd string) *FacturaCompraVentaCabeceraBuilder {
	b.cabecera.Cufd = cufd
	return b
}

// WithCodigoSucursal configura el código de la sucursal emisora (0 para casa matriz).
func (b *FacturaCompraVentaCabeceraBuilder) WithCodigoSucursal(codigoSucursal int) *FacturaCompraVentaCabeceraBuilder {
	b.cabecera.CodigoSucursal = codigoSucursal
	return b
}

// WithDireccion configura la dirección física del establecimiento emisor.
func (b *FacturaCompraVentaCabeceraBuilder) WithDireccion(direccion string) *FacturaCompraVentaCabeceraBuilder {
	b.cabecera.Direccion = direccion
	return b
}

// WithCodigoPuntoVenta configura el código del punto de venta registrado ante el SIAT.
func (b *FacturaCompraVentaCabeceraBuilder) WithCodigoPuntoVenta(codigoPuntoVenta int) *FacturaCompraVentaCabeceraBuilder {
	b.cabecera.CodigoPuntoVenta = codigoPuntoVenta
	return b
}

// WithFechaEmision configura la fecha y hora de emisión de la factura en formato string (ISO8601).
func (b *FacturaCompraVentaCabeceraBuilder) WithFechaEmision(fechaEmision time.Time) *FacturaCompraVentaCabeceraBuilder {
	b.cabecera.FechaEmision = datatype.NewTimeSiat(fechaEmision)
	return b
}

// WithNombreRazonSocial configura el nombre o razón social del cliente (opcional).
func (b *FacturaCompraVentaCabeceraBuilder) WithNombreRazonSocial(nombreRazonSocial string) *FacturaCompraVentaCabeceraBuilder {
	b.cabecera.NombreRazonSocial = datatype.Nilable[string]{Value: &nombreRazonSocial}
	return b
}

// WithCodigoTipoDocumentoIdentidad configura el código del tipo de documento de identidad del cliente.
func (b *FacturaCompraVentaCabeceraBuilder) WithCodigoTipoDocumentoIdentidad(codigoTipoDocumentoIdentidad int) *FacturaCompraVentaCabeceraBuilder {
	b.cabecera.CodigoTipoDocumentoIdentidad = codigoTipoDocumentoIdentidad
	return b
}

// WithNumeroDocumento configura el número de documento de identidad del cliente.
func (b *FacturaCompraVentaCabeceraBuilder) WithNumeroDocumento(numeroDocumento string) *FacturaCompraVentaCabeceraBuilder {
	b.cabecera.NumeroDocumento = numeroDocumento
	return b
}

// WithComplemento configura el complemento del documento de identidad (opcional).
func (b *FacturaCompraVentaCabeceraBuilder) WithComplemento(complemento string) *FacturaCompraVentaCabeceraBuilder {
	b.cabecera.Complemento = datatype.Nilable[string]{Value: &complemento}
	return b
}

// WithCodigoCliente configura un código interno único para identificar al cliente.
func (b *FacturaCompraVentaCabeceraBuilder) WithCodigoCliente(codigoCliente string) *FacturaCompraVentaCabeceraBuilder {
	b.cabecera.CodigoCliente = codigoCliente
	return b
}

// WithCodigoMetodoPago configura el código de la paramétrica del método de pago utilizado.
func (b *FacturaCompraVentaCabeceraBuilder) WithCodigoMetodoPago(codigoMetodoPago int) *FacturaCompraVentaCabeceraBuilder {
	b.cabecera.CodigoMetodoPago = codigoMetodoPago
	return b
}

// WithNumeroTarjeta configura los primeros y últimos dígitos de la tarjeta (opcional, enmascarado).
func (b *FacturaCompraVentaCabeceraBuilder) WithNumeroTarjeta(numeroTarjeta int64) *FacturaCompraVentaCabeceraBuilder {
	b.cabecera.NumeroTarjeta = datatype.Nilable[int64]{Value: &numeroTarjeta}
	return b
}

// WithMontoTotal configura el monto total de la factura, redondeado automáticamente a 2 decimales.
func (b *FacturaCompraVentaCabeceraBuilder) WithMontoTotal(montoTotal float64) *FacturaCompraVentaCabeceraBuilder {
	// Asegurar que el valor sea redondeado a 2 decimales
	montoTotal, _ = strconv.ParseFloat(strconv.FormatFloat(montoTotal, 'f', 2, 64), 64)
	b.cabecera.MontoTotal = montoTotal
	return b
}

// WithMontoTotalSujetoIva configura el monto base para el crédito fiscal IVA, redondeado automáticamente.
func (b *FacturaCompraVentaCabeceraBuilder) WithMontoTotalSujetoIva(montoTotalSujetoIva float64) *FacturaCompraVentaCabeceraBuilder {
	// Asegurar que el valor sea redondeado a 2 decimales
	montoTotalSujetoIva, _ = strconv.ParseFloat(strconv.FormatFloat(montoTotalSujetoIva, 'f', 2, 64), 64)
	b.cabecera.MontoTotalSujetoIva = montoTotalSujetoIva
	return b
}

// WithCodigoMoneda configura el código de la moneda utilizada (ej. 1 para Bolivianos).
func (b *FacturaCompraVentaCabeceraBuilder) WithCodigoMoneda(codigoMoneda int) *FacturaCompraVentaCabeceraBuilder {
	b.cabecera.CodigoMoneda = codigoMoneda
	return b
}

// WithTipoCambio configura el tipo de cambio respecto al Boliviano, redondeado automáticamente.
func (b *FacturaCompraVentaCabeceraBuilder) WithTipoCambio(tipoCambio float64) *FacturaCompraVentaCabeceraBuilder {
	// Asegurar que el valor sea redondeado a 2 decimales
	tipoCambio, _ = strconv.ParseFloat(strconv.FormatFloat(tipoCambio, 'f', 2, 64), 64)
	b.cabecera.TipoCambio = tipoCambio
	return b
}

// WithMontoTotalMoneda configura el monto total expresado en la moneda original, redondeado automáticamente.
func (b *FacturaCompraVentaCabeceraBuilder) WithMontoTotalMoneda(montoTotalMoneda float64) *FacturaCompraVentaCabeceraBuilder {
	// Asegurar que el valor sea redondeado a 2 decimales
	montoTotalMoneda, _ = strconv.ParseFloat(strconv.FormatFloat(montoTotalMoneda, 'f', 2, 64), 64)
	b.cabecera.MontoTotalMoneda = montoTotalMoneda
	return b
}

// WithMontoGiftCard configura el monto pagado con tarjeta de regalo o prepago (opcional, redondeado).
func (b *FacturaCompraVentaCabeceraBuilder) WithMontoGiftCard(montoGiftCard float64) *FacturaCompraVentaCabeceraBuilder {
	// Asegurar que el valor sea redondeado a 2 decimales
	montoGiftCard, _ = strconv.ParseFloat(strconv.FormatFloat(montoGiftCard, 'f', 2, 64), 64)
	b.cabecera.MontoGiftCard = datatype.Nilable[float64]{Value: &montoGiftCard}
	return b
}

// WithDescuentoAdicional configura un descuento global aplicado a toda la factura (opcional, redondeado).
func (b *FacturaCompraVentaCabeceraBuilder) WithDescuentoAdicional(descuentoAdicional float64) *FacturaCompraVentaCabeceraBuilder {
	// Asegurar que el valor sea redondeado a 2 decimales
	descuentoAdicional, _ = strconv.ParseFloat(strconv.FormatFloat(descuentoAdicional, 'f', 2, 64), 64)
	b.cabecera.DescuentoAdicional = datatype.Nilable[float64]{Value: &descuentoAdicional}
	return b
}

// WithCodigoExcepcion configura un código de excepción si los datos del cliente son inválidos pero autorizados (opcional).
func (b *FacturaCompraVentaCabeceraBuilder) WithCodigoExcepcion(codigoExcepcion int64) *FacturaCompraVentaCabeceraBuilder {
	b.cabecera.CodigoExcepcion = datatype.Nilable[int64]{Value: &codigoExcepcion}
	return b
}

// WithCafc configura el Código de Autorización de Facturas por Contingencia (opcional).
func (b *FacturaCompraVentaCabeceraBuilder) WithCafc(cafc string) *FacturaCompraVentaCabeceraBuilder {
	b.cabecera.Cafc = datatype.Nilable[string]{Value: &cafc}
	return b
}

// WithLeyenda configura la leyenda obligatoria según la normativa de Impuestos Nacionales.
func (b *FacturaCompraVentaCabeceraBuilder) WithLeyenda(leyenda string) *FacturaCompraVentaCabeceraBuilder {
	b.cabecera.Leyenda = leyenda
	return b
}

// WithUsuario configura el identificador del usuario que emite la factura.
func (b *FacturaCompraVentaCabeceraBuilder) WithUsuario(usuario string) *FacturaCompraVentaCabeceraBuilder {
	b.cabecera.Usuario = usuario
	return b
}

// WithCodigoDocumentoSector configura el código que identifica el diseño o sector de la factura.
func (b *FacturaCompraVentaCabeceraBuilder) WithCodigoDocumentoSector(codigoDocumentoSector int) *FacturaCompraVentaCabeceraBuilder {
	b.cabecera.CodigoDocumentoSector = codigoDocumentoSector
	return b
}

// Build finaliza la construcción de la cabecera retornando la interfaz opaca.
func (b *FacturaCompraVentaCabeceraBuilder) Build() FacturaCompraVentaCabecera {
	return requestWrapper[models_facturas.Cabecera]{request: b.cabecera}
}

// DetalleBuilder ayuda a configurar un ítem individual de detalle de la factura.
type DetalleBuilder struct {
	detalle *models_facturas.Detalle
}

// WithActividadEconomica configura el código de actividad económica asociado al producto/servicio.
func (b *DetalleBuilder) WithActividadEconomica(actividadEconomica string) *DetalleBuilder {
	b.detalle.ActividadEconomica = actividadEconomica
	return b
}

// WithCodigoProductoSin configura el código de producto según el catálogo del SIAT.
func (b *DetalleBuilder) WithCodigoProductoSin(codigoProductoSin string) *DetalleBuilder {
	b.detalle.CodigoProductoSin = codigoProductoSin
	return b
}

// WithCodigoProducto configura un código interno propio de la empresa para el producto.
func (b *DetalleBuilder) WithCodigoProducto(codigoProducto string) *DetalleBuilder {
	b.detalle.CodigoProducto = codigoProducto
	return b
}

// WithDescripcion configura la descripción detallada del artículo o servicio.
func (b *DetalleBuilder) WithDescripcion(descripcion string) *DetalleBuilder {
	b.detalle.Descripcion = descripcion
	return b
}

// WithCantidad configura la cantidad vendida, redondeada automáticamente a 2 decimales.
func (b *DetalleBuilder) WithCantidad(cantidad float64) *DetalleBuilder {
	// Asegurar que el valor sea redondeado a 2 decimales
	cantidad, _ = strconv.ParseFloat(strconv.FormatFloat(cantidad, 'f', 2, 64), 64)
	b.detalle.Cantidad = cantidad
	return b
}

// WithUnidadMedida configura el código de la paramétrica de unidad de medida (ej. 1 para unidades).
func (b *DetalleBuilder) WithUnidadMedida(unidadMedida int) *DetalleBuilder {
	b.detalle.UnidadMedida = unidadMedida
	return b
}

// WithPrecioUnitario configura el precio unitario del ítem, redondeado automáticamente.
func (b *DetalleBuilder) WithPrecioUnitario(precioUnitario float64) *DetalleBuilder {
	// Asegurar que el valor sea redondeado a 2 decimales
	precioUnitario, _ = strconv.ParseFloat(strconv.FormatFloat(precioUnitario, 'f', 2, 64), 64)
	b.detalle.PrecioUnitario = precioUnitario
	return b
}

// WithMontoDescuento configura un descuento aplicado específicamente a este ítem (opcional, redondeado).
func (b *DetalleBuilder) WithMontoDescuento(montoDescuento float64) *DetalleBuilder {
	// Asegurar que el valor sea redondeado a 2 decimales
	montoDescuento, _ = strconv.ParseFloat(strconv.FormatFloat(montoDescuento, 'f', 2, 64), 64)
	b.detalle.MontoDescuento = datatype.Nilable[float64]{Value: &montoDescuento}
	return b
}

// WithSubTotal configura el subtotal (Cantidad * PrecioUnitario - MontoDescuento), redondeado automáticamente.
func (b *DetalleBuilder) WithSubTotal(subTotal float64) *DetalleBuilder {
	// Asegurar que el valor sea redondeado a 2 decimales
	subTotal, _ = strconv.ParseFloat(strconv.FormatFloat(subTotal, 'f', 2, 64), 64)
	b.detalle.SubTotal = subTotal
	return b
}

// WithNumeroSerie configura el número de serie del producto (opcional).
func (b *DetalleBuilder) WithNumeroSerie(numeroSerie string) *DetalleBuilder {
	b.detalle.NumeroSerie = datatype.Nilable[string]{Value: &numeroSerie}
	return b
}

// WithNumeroImei configura el número IMEI si se trata de equipos telefónicos (opcional).
func (b *DetalleBuilder) WithNumeroImei(numeroImei string) *DetalleBuilder {
	b.detalle.NumeroImei = datatype.Nilable[string]{Value: &numeroImei}
	return b
}

// Build finaliza la construcción del detalle retornando la interfaz opaca.
func (b *DetalleBuilder) Build() FacturaCompraVentaDetalle {
	return requestWrapper[models_facturas.Detalle]{request: b.detalle}
}
