Las estrategias comunes para soft deletion en PostgreSQL (válidas para tu stack Go + React) son:

- **Columna `deleted_at` (timestamp nullable)**  
  La más usada. Agregás un campo `TIMESTAMPTZ`. Cuando está `NULL`, el registro está activo; al "eliminar", ponés la fecha actual. Permite saber cuándo se eliminó, ordenar, y recuperar fácilmente. Índice parcial recomendado: `WHERE deleted_at IS NULL` para búsquedas activas.

- **Columna `is_deleted` (boolean)**  
  Simple, valor por defecto `false`. Al eliminar se marca `true`. No registra el momento, pero ocupa menos espacio y es muy rápida de filtrar. También se beneficia de un índice parcial.

- **Tabla de auditoría o archivado**  
  Movés el registro eliminado a una tabla espejo (ej. `clientes_archivo`) con un trigger o a nivel aplicación. La tabla principal solo tiene activos, evitás filtros en todas las queries pero complicás las relaciones y el restore.

- **Vista que oculta los eliminados**  
  Creás una vista con `WHERE deleted_at IS NULL` y usás esa vista en lugar de la tabla directamente en la mayoría de las consultas. Simplifica el código Go pero puede complicar escrituras (necesitás `INSTEAD OF` triggers o manejarlo a nivel aplicación).

- **Esquema de versionado o histórico**  
  Mantenés un registro de cambios (tipo `clientes_historial`) con fechas de vigencia. La eliminación es un cierre de vigencia. Complejo, pero ideal para auditoría completa.

Para un ERP que necesitás consistencia y posible restauración, la columna `deleted_at` con índice parcial suele ser el mejor equilibrio. En Go, tu lógica de eliminación actualiza ese campo y todas las consultas normales llevan `WHERE deleted_at IS NULL`. Después, podés crear middlewares o repositorios base que añadan esa condición automáticamente.
