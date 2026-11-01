# Imagen de la aplicación (TEC-23, ADR 0005). Una sola imagen con dos programas:
#   /server  la aplicación web (contenedor "app")
#   /demo    carga los datos de nuestro proyecto en DATA_DIR (init container "cargar-demo")
# Las plantillas y los estáticos van embebidos en el binario (web.FS), así que la imagen final solo lleva los dos
# binarios. Sin shell ni gestor de paquetes (distroless) y sin usuario root (TEC-24).

# ---------- Compilación ----------
FROM golang:1.26 AS compilar
WORKDIR /src

# Primero las dependencias, así una compilación que solo cambia código reutiliza esta capa.
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /salida/server ./cmd/server \
 && CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /salida/demo ./cmd/demo \
 && mkdir -p /salida/data

# ---------- Imagen final ----------
FROM gcr.io/distroless/static:nonroot

COPY --from=compilar /salida/server /server
COPY --from=compilar /salida/demo /demo
# La carpeta de datos pertenece al usuario sin privilegios: ahí escriben la demo y la app.
COPY --from=compilar --chown=nonroot:nonroot /salida/data /data

ENV PORT=8080 DATA_DIR=/data
EXPOSE 8080
USER nonroot:nonroot

# Por defecto arranca la app; el init container usa el comando /demo (deploy/containerapp.yaml).
ENTRYPOINT ["/server"]
