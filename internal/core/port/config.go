package port

// Config agrupa la configuración para realizar solicitudes autenticadas al SIAT.
// Debe ser instanciada para cada operación que requiera autenticación.
//
// Campos:
//   - Token: Token de autenticación obtenido del SIAT (obligatorio)
//   - UserAgent: Identificador opcional del cliente HTTP para registro y debugging
type Config struct {
	// Token es el código de autenticación proporcionado por el SIAT
	Token string
	// UserAgent es el identificador opcional del cliente para propósitos de logging
	UserAgent string
}
