# TEC-24 · Seguridad del despliegue

- **Issue:** #139 · **Responsable:** Agustín Lara · **Story Points:** 3
- **Redactada por:** Branko Almeira (Agile Enabler) · **Validada por:** el equipo al aprobar el PR de las specs
- **Depende de:** ADR 0005 y TEC-23, que entra en el mismo PR
- **Estado:** implementada (Sprint 4, PR #142, mergeado el 01/11/2026). No cambia el código congelado: solo agrega
  archivos de despliegue y configuración

## Objetivo

Publicar la app sin que nadie de afuera pueda modificar los datos, robar las credenciales de la cátedra o usar nuestra
app para gastar los créditos compartidos. Solo entra lo **indispensable para desplegar**; lo deseable queda anotado
al final.

## Amenazas y controles

| # | Amenaza | Control | Dónde |
|---|---|---|---|
| A1 | Cualquiera que tenga la URL crea, edita o borra datos: no hay login (ADR 0004) | **Usuario y clave para toda escritura:** el proxy pide autenticación básica en todo `POST` (las 29 rutas que escriben son `POST`; las 30 `GET` solo leen). La ruta `/ingresar` también la pide y redirige a la portada, así el navegador guarda la clave antes de la demo | `Caddyfile` |
| A2 | La clave o la credencial viajan en texto plano | **Solo HTTPS:** certificado de la plataforma, `allowInsecure: false` y `Strict-Transport-Security` | `containerapp.yaml`, `Caddyfile` |
| A3 | Se filtran `AZURE_CREDENTIALS`, `ACR_PASSWORD` o la clave de escritura | Secretos en el **Environment `azure` de GitHub, que solo puede usar `main`**. El repo es público: nunca van en archivos, logs ni el chat. La clave del registro se pasa por `--password-stdin`. La clave de escritura se guarda solo como hash bcrypt, en un secreto de la app | GitHub, workflow |
| A4 | Un PR (de un fork o propio) ejecuta código con los secretos | El workflow **no corre en `pull_request` ni `pull_request_target`**. `permissions: contents: read`. Las acciones de terceros se fijan por **hash de commit** | `cd-azure.yml` |
| A5 | Otro grupo sube una imagen con nuestra etiqueta (el registro es compartido) | Se despliega por **digest** `@sha256:…`, nunca por `latest` ni por etiqueta | `cd-azure.yml` |
| A6 | Vulnerabilidades conocidas: Go 1.22.2 tiene 45 en la biblioteca estándar que nuestro código alcanza (`govulncheck`, 01/11) | La imagen se compila con **Go 1.26.x**. Un job del workflow corre `govulncheck` con ese Go y **no publica si encuentra alguna** | `Dockerfile`, `cd-azure.yml` |
| A7 | Ejecutar código o escalar privilegios dentro del contenedor | Imagen **`distroless/static:nonroot`**: sin shell, sin gestor de paquetes, usuario sin privilegios. Las imágenes base se fijan por digest | `Dockerfile` |
| A8 | La app queda expuesta sin pasar por el proxy | El ingreso apunta solo al puerto del proxy (8081). La app escucha en 8080 dentro de la réplica, sin ingreso propio | `containerapp.yaml` |
| A9 | Pedidos gigantes o lentos para tumbar la app (el servidor no tiene límite de cuerpo ni timeouts de lectura) | El proxy limita el cuerpo a **1 MB** y corta las lecturas lentas | `Caddyfile` |
| A10 | Un bot consume los créditos compartidos o llena los archivos de datos | **Máximo 1 réplica** y 0,75 vCPU en total. Escala a cero. Los datos viven en un volumen temporal que se reinicia | `containerapp.yaml` |
| A11 | Robo de clics, scripts inyectados o fuga de datos por el `Referer` | Cabeceras de seguridad (abajo) | `Caddyfile` |
| A12 | El workflow toca recursos de otros grupos | Solo crea, cambia o borra nombres `grupo02-*` (TEC-23). Nunca hay un `delete` sobre el entorno, el registro ni el grupo de recursos | `cd-azure.yml` |

### Cabeceras que agrega el proxy

| Cabecera | Valor |
|---|---|
| `Strict-Transport-Security` | `max-age=31536000` |
| `Content-Security-Policy` | `default-src 'self'; script-src 'self' 'sha256-<hash del script del dashboard>'; style-src 'self' 'sha256-<hash del estilo que inyecta htmx>'; img-src 'self' data:; frame-ancestors 'none'; form-action 'self'; base-uri 'self'` |
| `X-Content-Type-Options` | `nosniff` |
| `X-Frame-Options` | `DENY` |
| `Referrer-Policy` | `same-origin` |
| `Server` | se quita |

El único script en línea es el del dashboard (`web/templates/paginas/dashboard.html`). Es texto fijo, sin acciones de
plantilla adentro, así que su hash no cambia entre pedidos. Se autoriza por ese hash, que el build del proxy calcula
en cada corrida. Así los scripts no necesitan `'unsafe-inline'`.

Los estilos tampoco usan `'unsafe-inline'`. Las plantillas no tienen atributos `style=`, pero HTMX 2.0.4 agrega un
`<style>` propio para los indicadores de carga, siempre con el mismo texto. Se autoriza por su hash, que queda fijo en
el Caddyfile. Si se actualiza HTMX, el hash cambia y el navegador muestra el bloqueo en la consola.

## Guía técnica

- **Imagen del proxy:** `deploy/proxy/Dockerfile` parte de `caddy:2` fijada por digest y copia el Caddyfile a
  `/etc/caddy/Caddyfile`, como recomienda Microsoft para Caddy en contenedores. No se monta nada en tiempo de
  ejecución.
- **Caddyfile:**
  - Bloque global con `auto_https off` (el HTTPS lo pone Container Apps) y `servers { timeouts { read_header 5s
    read_body 10s idle 2m } }` (A9).
  - Sitio `:8081` con `request_body { max_size 1MB }`, `header` con la tabla de cabeceras (y `-Server`), y
    `reverse_proxy localhost:8080`.
  - `handle /ingresar { basic_auth { equipo {$WRITE_PASSWORD_HASH} } redir / 303 }`.
  - Un `handle` para todo lo demás, con el matcher `@escritura method POST`,
    `basic_auth @escritura { equipo {$WRITE_PASSWORD_HASH} }` y `reverse_proxy localhost:8080`.
  - Un Caddyfile con esta estructura pasa `caddy validate` con Caddy 2.11.7 (verificado el 01/11).
- **Las 25 escrituras de las plantillas son formularios comunes (`method="post"`, ninguna usa `hx-post`).** El
  navegador muestra su ventana de usuario y clave y reenvía el formulario. Después de `/ingresar`, el navegador manda
  la clave sola en todas las rutas, porque la guarda para la raíz.
- **Hash de la CSP:** `sha256` en base64 del texto exacto entre `<script>` y `</script>` del dashboard. Se calcula
  en el build y se escribe en el Caddyfile. Se verifica con la consola del navegador en `/dashboard?proyecto=P1`.
- **Job `vulnerabilidades`:** `actions/setup-go` con `go-version: '1.26.x'`, después
  `go run golang.org/x/vuln/cmd/govulncheck@v1.9.0 ./...`. El job `imagen` declara `needs: vulnerabilidades`.
- **Pruebas de seguridad en la prueba de humo** (con `curl`, contra la URL desplegada):
  - `POST` sin clave a una ruta que escribe (por ejemplo `/backlog`) → `401`.
  - `GET /ingresar` sin clave → `401`.
  - `curl -I` de la portada → las cabeceras de la tabla, sin `Server` de Caddy.
  - `http://` → redirección a `https://`.
  - Con la clave no se prueba en el workflow, para no tener la clave en claro en GitHub. Esa prueba se hace a mano
    (criterio de terminado).
- **Usuario sin privilegios:** se verifica con `docker run --rm --entrypoint sh sme:local`, que tiene que fallar.

### La clave de escritura

- Un usuario (`equipo`) y una clave larga y aleatoria, generada para el despliegue.
- La clave en claro se reparte por un canal privado (no en el repo ni en issues). En GitHub se guarda solo el hash
  bcrypt (`WRITE_PASSWORD_HASH`).
- **Para el profesor, la app es de solo lectura sin clave.** Si quiere probar las escrituras, se le pasa la clave.

## Fuera del alcance (deseable, no indispensable)

- Límite de pedidos por minuto: Caddy no lo trae sin un plugin. `maxReplicas: 1` ya acota el costo.
- Protección CSRF: con la clave de escritura y `form-action 'self'` el riesgo baja. Agregarla requiere tocar el código.
- Rotar la credencial compartida de la cátedra después de la presentación: se le recomienda al profesor.
- Escaneo de la imagen con Trivy.

## Criterio de terminado

- [ ] Sin la clave, cualquier `POST` responde 401 y los datos no cambian. Hay un test con `curl` en la prueba de humo
      del workflow
- [ ] Con la clave, se completa el recorrido del guion de la demo (crear, votar, registrar horas, registrar defecto)
- [ ] Con la clave, `/ingresar` redirige a la portada y después las escrituras no vuelven a pedirla
- [ ] `http://` redirige a `https://`; las cabeceras de la tabla aparecen en la respuesta (`curl -I`)
- [ ] El dashboard muestra los gráficos con la CSP puesta (sin errores en la consola)
- [ ] `govulncheck` en verde con Go 1.26 dentro del workflow
- [ ] La imagen corre como usuario sin privilegios y no tiene shell (`docker run --entrypoint sh` falla)
- [ ] `git grep` no encuentra en el repo ninguna de las credenciales de la cátedra ni la clave de escritura
- [ ] El Environment `azure` solo permite la rama `main`
