package compra_venta

import "encoding/xml"

type AnulacionFactura struct {
	XMLName                           xml.Name                          `xml:"ns:anulacionFactura" json:"-"`
	SolicitudServicioAnulacionFactura SolicitudServicioAnulacionFactura `xml:"SolicitudServicioAnulacionFactura" json:"solicitudservicioanulacionfactura"`
}

type SolicitudServicioAnulacionFactura struct {
	CodigoAmbiente        int    `xml:"codigoAmbiente" json:"codigoambiente"`
	CodigoDocumentoSector int    `xml:"codigoDocumentoSector" json:"codigoDocumentosector"`
	CodigoEmision         int    `xml:"codigoEmision" json:"codigoemision"`
	CodigoModalidad       int    `xml:"codigoModalidad" json:"codigomodalidad"`
	CodigoPuntoVenta      int    `xml:"codigoPuntoVenta" json:"codigopuntoventa"`
	CodigoSistema         string `xml:"codigoSistema" json:"codigosistema"`
	CodigoSucursal        int    `xml:"codigoSucursal" json:"codigosucursal"`
	Cufd                  string `xml:"cufd" json:"cufd"`
	Cuf                   string `xml:"cuf" json:"cuf"`
	Cuis                  string `xml:"cuis" json:"cuis"`
	Nit                   int64  `xml:"nit" json:"nit"`
	TipoFacturaDocumento  int    `xml:"tipoFacturaDocumento" json:"tipofacturadocumento"`
	CodigoMotivo          int    `xml:"codigoMotivo" json:"codigomotivo"`
}

type AnulacionFacturaResponse struct {
	XMLName                      xml.Name                     `xml:"anulacionFacturaResponse" json:"-"`
	XmlnsNs2                     string                       `xml:"xmlns:ns2,attr" json:"-"`
	RespuestaServicioFacturacion RespuestaServicioFacturacion `xml:"RespuestaServicioFacturacion" json:"respuestadelserviciofacturacion"`
}
