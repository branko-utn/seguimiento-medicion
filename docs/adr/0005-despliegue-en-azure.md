# ADR 0005 · Despliegue en Azure Container Apps con un proxy delante

- **Estado:** propuesta (01/11/2026)
- **Decide:** el equipo, al aprobar el PR de las specs TEC-23 y TEC-24

## Contexto

La cátedra pide desplegar la aplicación y entrega, para todos los grupos, la misma cuenta de Azure:

- un service principal con rol **Contributor sobre el grupo de recursos `frsrutn`**, compartido por los 16 grupos;
- un registro de imágenes compartido (`ingcalidad2026.azurecr.io`, usuario administrador con una clave común);
- un nombre para nuestra app (`grupo02-seguimiento-medicion`) y un entorno de Container Apps (`utnapp-env`);
- créditos limitados y compartidos.

Lo que comprobamos el 01/11 con consultas de solo lectura:

| Hecho | Consecuencia |
|---|---|
| El entorno `utnapp-env` no existe en la suscripción | Hay que crearlo (o usar ACI, como el Grupo03) |
| Una política solo permite `italynorth`, `germanywestcentral`, `chilecentral`, `newzealandnorth` y `switzerlandnorth`; el grupo de recursos está en `southcentralus` | La región se indica siempre a mano: `switzerlandnorth`, la del registro |
| El registro es compartido y otros grupos pueden escribir en él | Se despliega por **digest** (`@sha256:…`), nunca por una etiqueta que otro pueda pisar |
| La aplicación no tiene login (ADR 0004) y guarda los datos en archivos JSON | Hace falta un control de acceso para escribir y **una sola réplica** |
| El código está congelado desde el 31/10 | La seguridad se agrega afuera del código: proxy, imagen y plataforma |

## Opciones

| | Azure Container Apps | Azure Container Instances (ACI) |
|---|---|---|
| HTTPS | Lo da la plataforma, con certificado administrado | Hay que montarlo (Caddy + Let's Encrypt, con riesgo de límites de emisión al redesplegar) |
| Costo | **Escala a cero:** sin tráfico no consume. Entra en la asignación gratuita mensual | Cobra cada segundo que está encendido: hay que apagarlo a mano o con un cron |
| Contenedores auxiliares | Sidecar e init container | Sidecar |
| Requisito | El entorno `utnapp-env` (no existe) | Ninguno |

## Decisión

**Azure Container Apps en `switzerlandnorth`**, con:

1. **Init container** que carga los datos de la demo (`cmd/demo`) en un volumen temporal antes de arrancar.
2. **Contenedor de la app**, que solo escucha dentro de la réplica.
3. **Sidecar Caddy**, el único expuesto: agrega las cabeceras de seguridad, limita el tamaño de los pedidos y
   **pide usuario y clave para cualquier escritura** (todas las escrituras son `POST`).
4. Mínimo 0 y **máximo 1 réplica**.

Si `utnapp-env` no existe, el workflow lo crea con ese nombre (el que fijó la cátedra) y nunca lo borra; lo decidió el equipo el 01/11. Se le avisa
al profesor.

**Plan B:** ACI en `switzerlandnorth` con el mismo sidecar Caddy sirviendo HTTPS, y un workflow para encenderlo y
apagarlo. Se usa solo si Container Apps no se puede crear.

**Terraform no se usa.** El grupo de recursos y el registro son de la cátedra y los comparten 16 grupos. Terraform
necesita guardar su estado, y un `plan`/`destroy` sobre recursos que no son nuestros es un riesgo. El workflow usa
`az` y toca solo lo que lleva el prefijo `grupo02-`.

## Consecuencias

- **Los datos se reinician** cada vez que la app escala a cero o se redespliega. Para una demo es una ventaja: lo
  que alguien rompa se arregla solo. Lo que se cargue durante la presentación dura mientras haya tráfico.
- **Leer es público; escribir pide clave.** El profesor puede recorrer todo sin credenciales.
- **La credencial compartida sigue siendo un riesgo que no controlamos:** otro grupo podría borrar nuestra app. La
  mitigación es poder redesplegar en minutos desde el workflow.
