# ADR 0001 · Arquitectura y tecnologías

- **Estado:** aceptada (21/10/2026, junto con el cronograma)
- **Responsable:** Branko Almeira

## Contexto

La consigna exige que el núcleo y las reglas de negocio estén en Go, con TDD sobre métricas, estimación, reglas y
validaciones, y escenarios BDD automatizados cuando sea posible. Tenemos 12 días hasta la presentación final, en
mini-sprints de 2–3 días, y el Sprint 2 pide una interfaz usable y el Sprint 3 un dashboard con gráficos.

## Decisión

| Aspecto | Decisión | Por qué |
|---|---|---|
| Tipo de aplicación | **Web**, servida por Go | Un dashboard con gráficos y una interfaz usable son mucho más simples en web que en consola |
| Servidor | `net/http` de la biblioteca estándar (Go 1.22+, con patrones de ruta `GET /x/{id}`) | Sin dependencias; se testea con `httptest` |
| Interfaz | `html/template` + **HTMX** + **Chart.js** por CDN | Todo en Go, sin un segundo stack ni paso de build |
| Arquitectura | Capas: `domain` → `metrics` / `service` → `store` / `server` | El dominio no depende de la web ni del almacenamiento: sus tests son simples y rápidos |
| Persistencia | Repositorios como **interfaces**; implementación en memoria que guarda un archivo JSON en `data/` | Sin instalar bases de datos. Si hace falta, se cambia por SQLite sin tocar el dominio |
| BDD | **godog** con escenarios en español (`# language: es`) | Es el estándar de Cucumber para Go |
| Tests | `testing` de la biblioteca estándar, tablas de casos, `-race` y cobertura | Lo que pide la consigna, sin magia |
| PDF | `github.com/go-pdf/fpdf` (Sprint 4) | Puro Go, sin dependencias del sistema |
| CI | GitHub Actions: `gofmt`, `go vet`, tests unitarios, BDD y cobertura | Hace cumplir la Definition of Done en cada PR |

## Fuera de alcance

Login de usuarios, base de datos en la nube y deploy público: no son requerimientos de la consigna y priorizamos
los requerimientos mínimos. La demo se hace con la aplicación corriendo localmente.

## Consecuencias

- Cada épica tiene su paquete en `internal/domain/<épica>`, así cada integrante trabaja en su carpeta y hay menos conflictos.
- Los escenarios BDD se ejecutan contra la capa de servicio (y contra HTTP cuando aporta), no contra el navegador.
- Si se agrega una dependencia nueva, se documenta en un ADR.
