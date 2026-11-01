# Despliegue en Azure

La aplicación de `main` se publica en **Azure Container Apps**, con HTTPS y los datos de la demo cargados, y se
vuelve a desplegar sola en cada merge que cambie código. Decisión: [ADR 0005](adr/0005-despliegue-en-azure.md).
Specs: [TEC-23](../specs/TEC-23-despliegue-en-azure.md) (despliegue) y
[TEC-24](../specs/TEC-24-seguridad-del-despliegue.md) (seguridad).

## URL

**<https://grupo02-seguimiento-medicion.wonderfulwave-ff1bf196.switzerlandnorth.azurecontainerapps.io>**

Primer despliegue: 01/11/2026, en el entorno `grupo02-env`. Si la app se recrea, la URL de cada despliegue también
queda en el resumen del job `desplegar` (pestaña *Actions* → *CD Azure* → la corrida).

- **Leer es público:** cualquiera puede recorrer la app sin credenciales.
- **Escribir pide usuario y clave** (usuario `equipo`). La clave la reparte el Agile Enabler por privado: nunca va en
  el repo, en issues, en PRs ni en logs. Para cargarla una sola vez antes de la demo, abrí `/ingresar`: el navegador
  la pide, vuelve a la portada y la recuerda para las escrituras siguientes.
- **Arranque en frío:** sin tráfico la app escala a cero. El primer pedido después de un rato tarda unos segundos y
  los datos vuelven a ser los de la demo.

## Cómo está armada

Cada réplica (hay como máximo una) tiene tres contenedores que comparten la red y un volumen temporal:

| Contenedor | Imagen | Qué hace |
|---|---|---|
| `cargar-demo` (init) | `grupo02-seguimiento-medicion`, comando `/demo` | Carga el proyecto de la demo en `/data` y termina |
| `app` | `grupo02-seguimiento-medicion`, comando `/server` | La aplicación, en `localhost:8080`. **No tiene ingreso propio** |
| `proxy` | `grupo02-proxy` (Caddy 2.11.7) | El único expuesto (puerto 8081). Pide la clave para escribir, limita los pedidos a 1 MB y agrega las cabeceras de seguridad |

- **HTTPS** lo da la plataforma; `http://` redirige a `https://`.
- **Los datos viven en un volumen temporal** (`EmptyDir`): se reinician cuando la app escala a cero o se redespliega.
- **Las imágenes no tienen shell ni corren como root** (distroless para la app; Caddy con el usuario 65532).

Archivos:

| Archivo | Qué es |
|---|---|
| [`Dockerfile`](../Dockerfile) | Imagen de la app: compila `cmd/server` y `cmd/demo` con Go 1.26.10 sobre distroless, todo fijado por digest |
| [`deploy/proxy/Dockerfile`](../deploy/proxy/Dockerfile) y [`deploy/Caddyfile`](../deploy/Caddyfile) | Imagen y configuración del proxy |
| [`deploy/containerapp.yaml`](../deploy/containerapp.yaml) | Plantilla de la Container App, con marcadores `${…}` que completa el workflow |
| [`.github/workflows/cd-azure.yml`](../.github/workflows/cd-azure.yml) | Construye, publica, despliega y prueba |

## Cómo se despliega

**Solo:** en cada push a `main` que cambie código (`cmd/`, `internal/`, `web/`, `go.mod`, el `Dockerfile`, `deploy/` o el
propio workflow). **A mano:** *Actions* → *CD Azure* → *Run workflow*, eligiendo la rama. **Nunca corre en un PR.**

El workflow hace, en orden:

1. **`configuracion`:** verifica que estén los secretos. Si falta alguno, construye las imágenes para validarlas y
   termina en verde **sin publicar**.
2. **`imagen`:** calcula el hash de la CSP del dashboard (`deploy/csp-hash.sh`), construye las dos imágenes, las
   publica como `sha-<commit>` y toma sus **digests**.
3. **`desplegar`:**
   1. crea nuestro entorno `grupo02-env` en `switzerlandnorth` **solo si no existe** (nunca lo borra), en **modo
      estándar** con el perfil `Consumption`, y **falla si quedó en modo express** (ver abajo);
   2. arma la plantilla con los digests y los secretos en un archivo temporal, que **se borra siempre**;
   3. crea o actualiza `grupo02-seguimiento-medicion` **por digest** (otro grupo no puede pisarnos la imagen);
   4. **prueba de humo:** hasta 3 minutos esperando que `/health` responda `{"estado":"ok"}` y que `/backlog` muestre
      el proyecto de la demo; si falla, muestra los logs de la app y del proxy;
   5. deja la URL en el resumen;
   6. borra del registro las imágenes `grupo02-*` salvo las 5 últimas.

**Regla de oro:** el workflow solo crea, cambia o borra lo que empieza con `grupo02-`. El grupo de recursos, el
registro y la cuenta son de la cátedra y los comparten los 16 grupos.

**Por qué un entorno propio:** el que indicó la cátedra, `utnapp-env`, quedó en **modo express** (lo creó nuestro
workflow con `az containerapp env create` el 01/11), y express no admite init containers, sidecars ni volúmenes. Por
eso usamos `grupo02-env`, creado con la API de Azure en modo estándar. `utnapp-env` no se toca. Detalle: [ADR
0005](adr/0005-despliegue-en-azure.md).

## Secretos

Van en el **Environment `azure`** del repo (*Settings* → *Environments* → `azure`), que solo puede usar `main` (y la
rama del PR del despliegue mientras está abierto). **Nunca en archivos, issues, PRs, logs ni el chat.**

| Secreto | Qué es | De dónde sale |
|---|---|---|
| `AZURE_CREDENTIALS` | JSON del service principal de la cátedra | Lo entrega la cátedra |
| `ACR_PASSWORD` | Clave del registro `ingcalidad2026.azurecr.io` | Lo entrega la cátedra |
| `WRITE_PASSWORD_HASH` | Hash bcrypt de la clave de escritura (nunca la clave en claro) | `docker run --rm caddy:2 caddy hash-password --plaintext '<clave>'` |

Los datos que no son secretos (nombre de la app, registro, grupo de recursos, entorno y región) están como variables
al principio del workflow.

## Cómo volver atrás

- **Al commit anterior:** *Run workflow* eligiendo `main` en un commit anterior, o revertir el PR (el merge del revert
  despliega solo).
- **A la revisión anterior sin reconstruir:**

  ```bash
  az containerapp revision list -g frsrutn -n grupo02-seguimiento-medicion -o table
  az containerapp revision activate -g frsrutn -n grupo02-seguimiento-medicion --revision <nombre-de-la-revisión>
  ```

- **Si otro grupo borró o pisó la app:** volver a correr el workflow. La recrea en un par de minutos.

## Ver qué pasa

```bash
az containerapp show -g frsrutn -n grupo02-seguimiento-medicion --query properties.configuration.ingress.fqdn -o tsv
az containerapp logs show -g frsrutn -n grupo02-seguimiento-medicion --container app --tail 100
az containerapp logs show -g frsrutn -n grupo02-seguimiento-medicion --container proxy --tail 100
az containerapp replica list -g frsrutn -n grupo02-seguimiento-medicion -o table   # debe haber una sola
```

## Probar en local antes de subir

Con Docker. Simula la réplica: la demo en un volumen, la app y el proxy compartiendo la red.

> **En Windows con Git Bash**, corré antes `export MSYS_NO_PATHCONV=1`. Si no, Git Bash convierte `/demo` en una ruta
> de Windows (`C:/Program Files/Git/demo`) y el init container falla. En Linux, macOS o PowerShell no hace falta.

```bash
docker build -t grupo02-seguimiento-medicion:local .
docker build -f deploy/proxy/Dockerfile -t grupo02-proxy:local deploy

docker volume create sme-datos
docker run --rm -v sme-datos:/data --entrypoint /demo grupo02-seguimiento-medicion:local         # init container
docker run -d --name sme-app -p 8099:8081 -v sme-datos:/data grupo02-seguimiento-medicion:local   # app

# Proxy con una clave de prueba (no la real) y el hash de la CSP
HASH=$(docker run --rm grupo02-proxy:local caddy hash-password --plaintext prueba-local)
docker run -d --name sme-proxy --network container:sme-app \
  -e WRITE_PASSWORD_HASH="$HASH" -e CSP_SCRIPT_HASH="$(sh deploy/csp-hash.sh)" grupo02-proxy:local
```

Abrí <http://localhost:8099>: leer funciona sin clave; escribir pide el usuario `equipo` y la clave `prueba-local`.
Para limpiar: `docker rm -f sme-proxy sme-app && docker volume rm sme-datos`.

La demo para la presentación sigue siendo la local (`make demo` y `make run`, ver [demo.md](demo.md)); la URL pública
es el respaldo.
