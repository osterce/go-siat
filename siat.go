package siat

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/ron86i/go-siat/internal/adapter/service"
	"github.com/ron86i/go-siat/internal/core/port"
)

// Map es un alias para map[string]interface{} que proporciona métodos de utilidad
// para trabajar con datos JSON de forma más cómda.
// Es especialmente útil al trabajar con respuestas heterogéneas del SIAT.
type Map map[string]interface{}

// ToJSON convierte el Map a su representación en string JSON.
// Retorna un error si la codificación falla.
func (m Map) ToJSON() (string, error) {
	bytes, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// Sum retorna la suma de todos los valores numéricos en el Map.
// Soporta tipos float64, float32, int, int64 e int32.
// Los valores no numéricos se ignoran.
func (m Map) Sum() float64 {
	var total float64
	for _, v := range m {
		switch val := v.(type) {
		case float64:
			total += val
		case float32:
			total += float64(val)
		case int:
			total += float64(val)
		case int64:
			total += float64(val)
		case int32:
			total += float64(val)
		}
	}
	return total
}

// ToStruct convierte el Map en la estructura Go especificada.
// Utiliza encoding/json internamente, por lo que se requiere que v sea un puntero
// a una estructura con etiquetas json apropiadas.
func (m Map) ToStruct(v interface{}) error {
	bytes, err := m.ToJSON()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(bytes), v)
}

// SiatServices es el punto de entrada principal del SDK.
// Agrupa todas las implementaciones de los servicios del SIAT
// (Códigos, Sincronización, Operaciones, Compra-Venta, Computarizada, Electrónica)
// y proporciona acceso a ellos a través de métodos orientados a objetivos.
// Los usuarios deben crear una instancia usando New().
type SiatServices struct {
	operaciones    port.SiatOperacionesPort
	sincronizacion port.SiatSincronizacionService
	codigos        port.SiatCodigosService
	compraVenta    port.SiatCompraVentaService
	computarizada  port.SiatComputarizadaService
	electronica    port.SiatElectronicaService
}

// Operaciones retorna el servicio para la gestión de puntos de venta (PV),
// cierre de períodos de facturación y eventos significativos (cambios de modalidad, etc.).
func (s *SiatServices) Operaciones() port.SiatOperacionesPort {
	return s.operaciones
}

// Sincronizacion retorna el servicio que proporciona acceso a catálogos maestros:
// actividades económicas, documentos fiscales, monedas, tipos de cambio, etc.
// Estos catálogos son esenciales para validar datos antes de emitir facturas.
func (s *SiatServices) Sincronizacion() port.SiatSincronizacionService {
	return s.sincronizacion
}

// Codigos retorna el servicio para:
// - Solicitud de códigos CUIS (Código Único de Identificación de Sistemas)
// - Solicitud de códigos CUFD (Código Único de Facturación por Dirección)
// - Validación de números NIT (Rol Tributario)
// Los códigos CUIS y CUFD son obligatorios para emitir facturas.
func (s *SiatServices) Codigos() port.SiatCodigosService {
	return s.codigos
}

// CompraVenta retorna el servicio para el sector de compra-venta (Sector 1).
// Permite enviar, recibir y anular facturas comerciales estándar.
// Este es el sector más común para comercios generales.
func (s *SiatServices) CompraVenta() port.SiatCompraVentaService {
	return s.compraVenta
}

// Computarizada retorna el servicio para facturación computarizada
// (sin firma digital, basada en máquinas registradoras fiscales).
// Permite enviar, recibir y anular facturas de este tipo.
func (s *SiatServices) Computarizada() port.SiatComputarizadaService {
	return s.computarizada
}

// Electronica retorna el servicio para facturación electrónica (con firma digital).
// Permite enviar, recibir y anular facturas electrónicas de todos los sectores.
// Este es el tipo de facturación más moderno y flexible del SIAT.
func (s *SiatServices) Electronica() port.SiatElectronicaService {
	return s.electronica
}

// New crea e inicializa una nueva instancia de SiatServices.
//
// Parámetros:
//   - baseUrl: URL base de los servicios SIAT (ej: https://pilotosiatservicios.impuestos.gob.bo/v2)
//   - httpClient: Cliente HTTP personalizado (opcional). Si es nil, se crea uno con configuración segura.
//
// La función configura automáticamente:
//   - Timeouts apropiados (15s handshake, 45s total)
//   - TLS 1.2+ para seguridad
//   - Pools de conexión para alto rendimiento
//   - Proxy desde variables de entorno si están configuradas
//
// Retorna un error si baseUrl está vacía o si alguno de los servicios falla al inicializarse.
//
// Ejemplo:
//
//	s, err := siat.New("https://pilotosiatservicios.impuestos.gob.bo/v2", nil)
//	if err != nil {
//		log.Fatal(err)
//	}
func New(baseUrl string, httpClient *http.Client) (*SiatServices, error) {
	if httpClient != nil {
		clonedClient := *httpClient
		httpClient = &clonedClient
	} else {
		httpClient = &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxConnsPerHost:     100,
				MaxIdleConnsPerHost: 20,
				IdleConnTimeout:     90 * time.Second,
				Proxy:               http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   15 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				TLSHandshakeTimeout: 15 * time.Second,
				TLSClientConfig: &tls.Config{
					MinVersion: tls.VersionTLS12,
				},
			},
			Timeout: 45 * time.Second,
		}
	}

	baseUrl = strings.TrimSpace(baseUrl)
	if baseUrl == "" {
		return nil, fmt.Errorf("baseUrl is empty")
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

	computarizada, err := service.NewSiatComputarizadaService(baseUrl, httpClient)
	if err != nil {
		return nil, err
	}
	electronica, err := service.NewSiatElectronicaService(baseUrl, httpClient)
	if err != nil {
		return nil, err
	}
	return &SiatServices{
		operaciones:    operaciones,
		sincronizacion: sincronizacion,
		codigos:        codigos,
		compraVenta:    compraVenta,
		computarizada:  computarizada,
		electronica:    electronica,
	}, nil
}
