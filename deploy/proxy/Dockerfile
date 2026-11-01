# Imagen del proxy (TEC-24, ADR 0005): Caddy con el Caddyfile de deploy/ adentro, como recomienda Microsoft para
# Caddy en contenedores. No se monta nada en tiempo de ejecución. El workflow la construye con contexto deploy/:
#   docker build -f deploy/proxy/Dockerfile -t grupo02-proxy deploy
#
# Variables que necesita al arrancar (las pone deploy/containerapp.yaml): WRITE_PASSWORD_HASH y CSP_SCRIPT_HASH.

# Caddy 2.11.7, fijado por digest (A7). Para actualizar: docker pull caddy:2 y copiar el digest nuevo.
FROM caddy:2@sha256:36c24ad91a7ff7bc4c87fe22bfa8df3c4b18da0773d73fc12ee2c52b98c0bb7e

COPY Caddyfile /etc/caddy/Caddyfile

# Sin root, como la app (A7): escucha en 8081, un puerto que no necesita privilegios. Caddy solo necesita escribir
# en sus carpetas de configuración y de datos.
RUN chown -R 65532:65532 /config /data
USER 65532:65532

EXPOSE 8081
