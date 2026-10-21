# TEC-01 · Estado del servicio (`/health`)

- **Issue:** TEC-01 · **Responsable:** Branko Almeira
- **Escenarios BDD:** [`features/tec-01-health.feature`](../features/tec-01-health.feature)
- **Estado:** implementada

> Ítem técnico del esqueleto. Sirve como **ejemplo completo de la cadena de trazabilidad**
> (spec → escenarios → tests → código) para copiar en las historias de usuario.

## Objetivo

Permitir verificar que la aplicación está levantada, para el CI y para comprobarlo antes de la demo.

## Entradas

| Dato | Tipo | Obligatorio | Validación |
|---|---|---|---|
| Método HTTP | `GET` | Sí | Solo se acepta `GET` |
| Ruta | `/health` | Sí | — |

## Salidas esperadas

Respuesta `200 OK` con `Content-Type: application/json` y el cuerpo `{"estado":"ok"}`.

## Reglas de negocio

- **RN-TEC01.1:** el endpoint no requiere datos cargados: responde aunque no exista ningún proyecto.

## Restricciones

- No expone información interna (versiones, rutas, configuración).

## Casos límite

- Ruta inexistente (`/salud`): responde `404 Not Found`.

## Condiciones de error

| Situación | Error esperado |
|---|---|
| Método distinto de `GET` (por ejemplo `POST /health`) | `405 Method Not Allowed` |

## Criterios de aceptación

| ID | Criterio | Escenario BDD |
|---|---|---|
| CA-TEC01.1 | `GET /health` responde 200 con `{"estado":"ok"}` | `Escenario: el servicio está levantado` |
| CA-TEC01.2 | Un método distinto de `GET` responde 405 | `Escenario: método no permitido` |
