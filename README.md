# go-siat

[![Status](https://img.shields.io/badge/status-active-success?style=flat-square)](https://github.com/ron86i/go-siat)
[![Go Version](https://img.shields.io/badge/go-1.23+-00ADD8?style=flat-square)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-green?style=flat-square)](LICENSE)

**go-siat** es un SDK profesional desarrollado en Go, diseñado para simplificar la integración con los servicios web SOAP del **SIAT (Sistema de Facturación de Impuestos Nacionales de Bolivia)**.

---

## Características Principales

*   **Type-Safe**: Estructuras de datos rigurosas para todas las solicitudes y respuestas del SIAT.
*   **Builder Pattern**: Construcción intuitiva de solicitudes complejas mediante interfaces fluidas.
*   **Abstracción SOAP**: Gestión transparente de la capa SOAP y seguridad (XMLDSig).
*   **Modular**: Separación clara entre modelos, servicios y adaptadores.

---

## Tabla de Contenidos

1. [Capacidades Implementadas](#capacidades-implementadas)
2. [Guía de Inicio Rápido](#guía-de-inicio-rápido)
3. [Referencia de Uso (Tests)](#referencia-de-uso-tests)
4. [Licencia](#licencia)

---

## Capacidades Implementadas

El SDK cubre los servicios críticos del ecosistema SIAT:

| Servicios | Funcionalidades Clave |
| :--- | :--- |
| **Códigos** | Solicitud de CUIS/CUFD (Individual y Masivo), Validación de NIT, Comunicación. |
| **Sincronización** | Catálogos de actividades, paramétricas, productos, servicios y documentos sector. |
| **Operaciones** | Registro/Cierre de Puntos de Venta, Gestión de Eventos Significativos. |
| **Compra y Venta** | Generación de CUF, Firma Digital XML, Recepción y Anulación de Facturas. |

---

## Guía de Inicio Rápido

### Instalación

```bash
go get github.com/ron86i/go-siat@v0.3.1
```

### Uso Básico

El siguiente ejemplo demuestra cómo inicializar el cliente y realizar una solicitud de código CUIS:

```go
package main

import (
    "context"
    "log"
    "github.com/ron86i/go-siat"
    "github.com/ron86i/go-siat/pkg/config"
    "github.com/ron86i/go-siat/pkg/models"
)

func main() {
    // 1. Configurar cliente unificado
    s, err := siat.New("https://pilotosiatservicios.impuestos.gob.bo/v2", nil)
    if err != nil {
        log.Fatal("Error al inicializar SDK:", err)
    }

    // 2. Construir solicitud usando el Builder
    req := models.Codigos().NewCuisBuilder().
		WithCodigoAmbiente(1).
		WithCodigoModalidad(1).
		WithCodigoPuntoVenta(0).
		WithCodigoSucursal(0).
		WithCodigoSistema("ABC123DEF").
		WithNit(123456789).
		Build()

    // 3. Ejecutar operación
    ctx := context.Background()
    cfg := config.Config{Token: "TU_TOKEN_API"}
    
    resp, err := s.Codigos().SolicitudCuis(ctx, cfg, req)
    if err != nil {
        log.Fatal("Error en la solicitud:", err)
    }
    
    log.Println("Código CUIS obtenido:", resp.Body.Content.RespuestaCuis.Codigo)
}
```

---

## Referencia de Uso (Tests)

Para una comprensión profunda de cada servicio, los **Tests de Integración** actúan como la documentación técnica principal.

| Categoría | Archivo de Test |
| :--- | :--- |
| **Códigos** | [`siat_codigos_service_test.go`](./internal/adapter/service/siat_codigos_service_test.go) |
| **Sincronización** | [`siat_sincronizacion_service_test.go`](./internal/adapter/service/siat_sincronizacion_service_test.go) |
| **Operaciones** | [`siat_operaciones_service_test.go`](./internal/adapter/service/siat_operaciones_service_test.go) |
| **Compra y Venta** | [`siat_compra_venta_service_test.go`](./internal/adapter/service/siat_compra_venta_service_test.go) |

> **Configuración de Ambiente**
> Antes de ejecutar los tests, asegúrese de crear un archivo `.env` configurado con sus credenciales del ambiente de pruebas del SIAT.

---

## Licencia

Distribuido bajo la **Licencia MIT**. Consulte el archivo `LICENSE` para más detalles.
