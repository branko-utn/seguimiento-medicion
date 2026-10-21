// Package web contiene las plantillas HTML y los archivos estáticos, embebidos en el binario
// para que la aplicación funcione sin internet.
package web

import "embed"

// FS tiene las carpetas templates/ y static/.
//
//go:embed static
var FS embed.FS
