# TEC-23 · Despliegue en Azure

- **Issue:** #138 · **Responsable:** Federico Sosa · **Story Points:** 5
- **Redactada por:** Branko Almeira (Agile Enabler) · **Validada por:** el equipo al aprobar el PR de las specs
- **Depende de:** ADR 0005. La seguridad del despliegue está en TEC-24, que entra en el mismo PR
- **Estado:** propuesta (01/11/2026). **No cambia el código congelado:** solo agrega archivos de despliegue

## Objetivo

Que la aplicación de `main` quede publicada en Azure, con HTTPS, y que se vuelva a desplegar sola en cada merge, con
los datos de la demo cargados.

## Datos que entrega la cátedra

| Dato | Valor | Dónde va |
|---|---|---|
| `CONTAINER_APP_NAME` | `grupo02-seguimiento-medicion` | Variable del workflow |
| `ACR_NAME` / `ACR_NAME_URL` | `ingcalidad2026` / `ingcalidad2026.azurecr.io` | Variables del workflow |
| `AZURE_RESOURCE_GROUP` | `frsrutn` | Variable del workflow |
| `AZURE_CONTAINERAPPS_ENVIRONMENT` | `grupo02-env` (la cátedra indicó `utnapp-env`; ver ADR 0005) | Variable del workflow |
| `AZURE_CREDENTIALS` | JSON del service principal | **Secreto** del Environment `azure` (TEC-24) |
| `ACR_PASSWORD` | Clave del registro | **Secreto** del Environment `azure` (TEC-24) |

Región: `switzerlandnorth` (política de regiones permitidas y misma región que el registro).

## Entregables

| Archivo | Qué hace |
|---|---|
| `Dockerfile` | Compila `cmd/server` y `cmd/demo` con **Go 1.26.x** (`CGO_ENABLED=0`) y los copia a una imagen `distroless/static:nonroot`. Las plantillas ya van embebidas en el binario (`web.FS`) |
| `deploy/proxy/Dockerfile` y `deploy/proxy/Caddyfile` | Imagen del sidecar Caddy (TEC-24) |
| `deploy/containerapp.yaml` | Plantilla de la app: init container, app, sidecar, volumen, ingreso y escala |
| `.github/workflows/cd-azure.yml` | Construye, publica y despliega (abajo) |
| `docs/despliegue.md` | URL, cómo se despliega, cómo se cargan los secretos, cómo se vuelve atrás |

## Cómo queda armada la app

| Pieza | Imagen | Recursos | Detalle |
|---|---|---|---|
| Init container `cargar-demo` | la de la app, comando `/demo` | 0,25 vCPU · 0,5 GiB | `DATA_DIR=/data`. Carga el proyecto de la demo y termina |
| Contenedor `app` | la de la app, comando `/server` | 0,25 vCPU · 0,5 GiB | `DATA_DIR=/data`, `PORT=8080`. No se expone |
| Sidecar `proxy` | `grupo02-proxy` | 0,25 vCPU · 0,5 GiB | Escucha en 8081 y reenvía a `localhost:8080` |
| Volumen `datos` | `EmptyDir` | — | Montado en `/data` en el init container y en la app |
| Ingreso | — | — | Externo, `targetPort: 8081`, **solo HTTPS** (`allowInsecure: false`) |
| Escala | — | — | `minReplicas: 0`, **`maxReplicas: 1`** |

## Workflow `cd-azure.yml`

- **Cuándo corre:** push a `main` que cambie código, y a mano (`workflow_dispatch`). Nunca en PRs.
- **Pasos:**
  1. Verifica que estén los secretos. Si falta alguno, construye las imágenes para validarlas y termina en verde sin
     publicar.
  2. Construye las dos imágenes con la etiqueta `sha-<commit>` y las publica en el registro como
     `grupo02-seguimiento-medicion` y `grupo02-proxy`.
  3. Toma el **digest** de cada imagen recién publicada.
  4. Si el entorno propio `grupo02-env` no existe, lo crea en `switzerlandnorth` en **modo estándar** (perfil
     `Consumption`) con `az rest --method PUT`, espera a que termine y **falla si quedó en modo express**. Si existe,
     no lo toca. (Actualización del 01/11: `az containerapp env create` creó `utnapp-env` en modo express, que no
     admite init containers, sidecars ni volúmenes. Ver ADR 0005.)
  5. Crea o actualiza `grupo02-seguimiento-medicion` con `deploy/containerapp.yaml` y los digests.
  6. **Prueba de humo:** espera hasta 3 minutos a que `https://<fqdn>/health` responda `{"estado":"ok"}` y a que
     `/backlog` muestre el proyecto de la demo. Si falla, el job queda en rojo y muestra los logs de la réplica.
  7. Deja la URL en el resumen del job.
  8. Borra del registro las imágenes `grupo02-*` salvo las 5 últimas.
- **Volver atrás:** volver a correr el workflow sobre un commit anterior, o activar la revisión anterior con
  `az containerapp revision activate`.

## Reglas

- **Solo se crea, cambia o borra lo que empieza con `grupo02-`**, incluido el entorno `grupo02-env`. La única
  excepción es leer el registro. Nunca se toca `utnapp-env` ni se borra el registro o el grupo de recursos.
- **Sin Terraform** (ADR 0005).
- **La demo local no cambia:** `make demo` y `make run` siguen igual. La presentación se hace en local y la URL
  pública es el respaldo, o al revés si lo pide la cátedra.

## Casos límite

- **La app estuvo un rato sin tráfico:** escaló a cero. El primer pedido tarda unos segundos (arranque en frío) y los
  datos vuelven a ser los de la demo.
- **Otro grupo borra o pisa nuestra app:** se vuelve a correr el workflow.
- **Una etiqueta del registro fue pisada por otro grupo:** no afecta, porque se despliega por digest.

## Criterio de terminado

- [ ] `docker build` local de las dos imágenes; la app arranca con `docker run` y responde `/health`
- [ ] El workflow en verde sobre `main`, con la URL en el resumen
- [ ] La URL abre por HTTPS; `http://` redirige a `https://`
- [ ] Las pantallas del guion de la demo (`docs/demo.md`) responden con los datos de la demo
- [ ] Con dos pestañas abiertas sigue habiendo una sola réplica (`az containerapp replica list`)
- [ ] `docs/despliegue.md` escrito y enlazado desde el README y el anexo A.5 del manual
- [ ] Criterios de TEC-24 cumplidos
