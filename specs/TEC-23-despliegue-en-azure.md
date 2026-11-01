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
| `AZURE_CONTAINERAPPS_ENVIRONMENT` | `utnapp-env` | Variable del workflow |
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
  4. Si `utnapp-env` no existe, lo crea en `switzerlandnorth`. Si existe, no lo toca.
  5. Crea o actualiza `grupo02-seguimiento-medicion` con `deploy/containerapp.yaml` y los digests.
  6. **Prueba de humo:** espera hasta 3 minutos a que `https://<fqdn>/health` responda `{"estado":"ok"}` y a que
     `/backlog` muestre el proyecto de la demo. Si falla, el job queda en rojo y muestra los logs de la réplica.
  7. Deja la URL en el resumen del job.
  8. Borra del registro las imágenes `grupo02-*` salvo las 5 últimas.
- **Volver atrás:** volver a correr el workflow sobre un commit anterior, o activar la revisión anterior con
  `az containerapp revision activate`.

## Cómo se trabaja (Sosa y Lara)

- **Una rama y un PR para TEC-23 y TEC-24:** `feature/TEC-23-24-despliegue`. El PR dice `Closes #138` y
  `Closes #139`. Cada uno hace commits solo sobre lo suyo; revisa el otro y audita el Agile Enabler.
- **Orden, para no pisarse en `cd-azure.yml`:**
  1. **Sosa** crea la rama y sube el `Dockerfile`, `deploy/containerapp.yaml` y el workflow con los jobs
     `configuracion`, `imagen` y `desplegar`. Hasta acá el proxy puede ser un Caddy que solo reenvía.
  2. **Lara**, sobre la misma rama, sube `deploy/proxy/` (la imagen y el Caddyfile completo), agrega al workflow el job
     `vulnerabilidades` (que `imagen` necesita para publicar) y las pruebas de seguridad de la prueba de humo.
  3. Después de cada `git pull`, el que sigue continúa. Si los dos tienen que tocar el workflow a la vez, se avisan por
     el grupo.
- **Probar el despliegue antes del merge:** el Environment `azure` permite también la rama
  `feature/TEC-23-24-despliegue` mientras dure el PR. El workflow se corre a mano (`workflow_dispatch`) eligiendo esa
  rama. Al mergear, el Agile Enabler quita la rama del Environment.
- **Probar en local antes de subir** (Docker):

  ```bash
  docker build -t sme:local .
  docker volume create sme-datos
  docker run --rm -v sme-datos:/data -e DATA_DIR=/data --entrypoint /demo sme:local      # init container
  docker run --rm -p 8080:8080 -v sme-datos:/data -e DATA_DIR=/data sme:local            # app
  ```

  Con el proxy, agregar una red de Docker y el contenedor de Caddy apuntando a `app:8080`. En Azure es
  `localhost:8080` porque los contenedores de una réplica comparten la red.
- **La clave de escritura** se la pasa el Agile Enabler a cada uno por privado. No se escribe en el repo, issues,
  PR ni logs.

## Guía técnica

- **Imagen de la app:** etapa `golang:1.26` fijada por digest, `CGO_ENABLED=0 go build -trimpath -ldflags="-s -w"` de
  `./cmd/server` y `./cmd/demo`. Etapa final `gcr.io/distroless/static-debian12:nonroot` fijada por digest, con
  `/server` y `/demo`, `ENTRYPOINT ["/server"]` y `EXPOSE 8080`. El init container cambia el comando a `/demo`.
  La imagen final no tiene shell para hacer `mkdir`. Por eso se copia desde la etapa de compilación una carpeta
  `/data` vacía con `COPY --chown=65532:65532`. Si no, el usuario sin privilegios no puede escribir en un volumen
  nuevo de Docker.
- **`deploy/containerapp.yaml`:** formato de `az containerapp create --yaml`. Tiene `initContainers` (`cargar-demo`),
  `containers` (`app` y `proxy`), un volumen `datos` con `storageType: EmptyDir` montado en `/data`, `ingress` con
  `external: true`, `targetPort: 8081`, `transport: auto`, `allowInsecure: false`, y `scale` con `minReplicas: 0` y
  `maxReplicas: 1`. Las imágenes y los secretos van como marcadores (`${IMAGEN_APP}`, `${IMAGEN_PROXY}`) que el workflow
  reemplaza con `envsubst` en un archivo de `$RUNNER_TEMP`, que se borra al terminar. Ese archivo nunca se sube al repo
  ni se imprime.
- **Secretos de la Container App:** `registry-password` (de `ACR_PASSWORD`, para bajar las imágenes del registro con
  el usuario `ingcalidad2026`) y `write-password-hash` (de `WRITE_PASSWORD_HASH`). El proxy recibe el hash como
  variable `WRITE_PASSWORD_HASH` con `secretRef: write-password-hash`.
- **Crear el entorno si falta:**

  ```bash
  az containerapp env show -g frsrutn -n utnapp-env >/dev/null 2>&1 || \
    az containerapp env create -g frsrutn -n utnapp-env -l switzerlandnorth --logs-destination none
  ```

  `--logs-destination none` evita el costo de Log Analytics. Los logs se ven con `az containerapp logs show`.
- **Crear o actualizar la app:** si existe, `az containerapp update -g frsrutn -n grupo02-seguimiento-medicion --yaml
  …`; si no, `az containerapp create … --environment utnapp-env --yaml …`.
- **Digest:** sale de la salida `digest` de `docker/build-push-action`. La imagen se referencia como
  `ingcalidad2026.azurecr.io/grupo02-seguimiento-medicion@sha256:…`.
- **Limpieza del registro:** con el usuario y la clave del registro (`az acr repository show-tags` y
  `az acr repository delete`), solo sobre `grupo02-seguimiento-medicion` y `grupo02-proxy`.
- **Acciones del workflow:** `actions/checkout`, `actions/setup-go`, `azure/login`, `docker/setup-buildx-action`,
  `docker/login-action` y `docker/build-push-action`, todas fijadas por hash de commit (TEC-24, A4).

## Reglas

- **Solo se crea, cambia o borra lo que empieza con `grupo02-`.** Las excepciones son crear `utnapp-env` si no existe
  y leer el registro. Nunca se borra el entorno, el registro ni el grupo de recursos.
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
