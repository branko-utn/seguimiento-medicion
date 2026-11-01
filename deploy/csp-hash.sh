#!/bin/sh
# Imprime el hash del único script en línea de la app (el del dashboard) para la CSP del proxy: sha256-<base64>.
# Lo usa el workflow en cada build y lo pasa al proxy en CSP_SCRIPT_HASH (TEC-24, A11).
#
# Uso: deploy/csp-hash.sh [plantilla]   (por defecto, web/templates/paginas/dashboard.html)
set -eu

plantilla="${1:-web/templates/paginas/dashboard.html}"

# Si aparece otro script en línea en cualquier plantilla, la CSP lo bloquearía: mejor fallar acá que en la demo.
cantidad=$(cat web/templates/*.html web/templates/*/*.html | grep -c '<script>' || true)
if [ "$cantidad" -ne 1 ] || ! grep -q '<script>' "$plantilla"; then
	echo "se esperaba un solo script en línea, en $plantilla; hay $cantidad en web/templates" >&2
	exit 1
fi

# El hash se calcula sobre el texto exacto entre <script> y </script>. Se quitan los \r porque el navegador
# convierte CRLF en LF antes de calcularlo (pasa si la plantilla se sacó con CRLF, como en Windows).
hash=$(perl -0777 -ne 'print $1 if /<script>(.*?)<\/script>/s' "$plantilla" | tr -d '\r' |
	openssl dgst -sha256 -binary | openssl base64 -A)
echo "sha256-$hash"
