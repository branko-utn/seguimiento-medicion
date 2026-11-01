# Imagen de la aplicación (TEC-23, ADR 0005). Una sola imagen con dos programas:
#   /server  la aplicación web (contenedor "app")
#   /demo    carga los datos de nuestro proyecto en DATA_DIR (init container "cargar-demo")
# Las plantillas y los estáticos van embebidos en el binario (web.FS), así que la imagen final solo lleva los dos
# binarios. Sin shell ni gestor de paquetes (distroless) y sin usuario root (TEC-24).

# Las imágenes base se fijan por digest (TEC-24, A7): una etiqueta se puede mover, un digest no. La etiqueta queda
# solo para leer qué versión es. Para actualizar: docker pull <imagen:etiqueta> y copiar el digest nuevo.

# ---------- Compilación (Go 1.26.10) ----------
FROM golang:1.26@sha256:670bb5aa5349ae03b8516764718486f4f16399134bdc9acd0c46a23f043ab05c AS compilar
WORKDIR /src

# Primero las dependencias, así una compilación que solo cambia código reutiliza esta capa.
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /salida/server ./cmd/server \
 && CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /salida/demo ./cmd/demo \
 && mkdir -p /salida/data

# ---------- Imagen final ----------
FROM gcr.io/distroless/static-debian12:nonroot@sha256:afa5c872c891853ca7fcf1f12c3edb23f7eeef36189728842dd51042ff57f7ab

COPY --from=compilar /salida/server /server
COPY --from=compilar /salida/demo /demo
# La carpeta de datos pertenece al usuario sin privilegios (65532, nonroot): ahí escriben la demo y la app. La imagen
# no tiene shell para hacer mkdir, así que se copia vacía desde la compilación; si no, un volumen nuevo queda de root.
COPY --from=compilar --chown=65532:65532 /salida/data /data

ENV PORT=8080 DATA_DIR=/data
EXPOSE 8080
USER 65532:65532

# Por defecto arranca la app; el init container usa el comando /demo (deploy/containerapp.yaml).
ENTRYPOINT ["/server"]
