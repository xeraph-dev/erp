Para una arquitectura React (SPA) + API Go, las estrategias más comunes son:

**1. JWT (Access Token + Refresh Token)**
- **Flujo**: el cliente envía credenciales, recibe un `access_token` (corta vida, minutos) y un `refresh_token` (larga vida, días/semanas) que guarda en memoria o en una cookie `httpOnly`.
- **Ventajas**: stateless, escalable, sin estado en el servidor, ideal para APIs REST. Los refresh token permiten renovar el acceso sin re-login.
- **Consideraciones**: revocación compleja (necesitás una blacklist o rotación de refresh tokens). Para un ERP, combiná refresh token con rotación y almacenamiento seguro.

**2. Sesiones con cookies (`httpOnly`, `secure`, `SameSite`)**
- **Flujo**: autenticación tradicional. El servidor crea una sesión (guardada en DB, Redis o en memoria) y envía una cookie de sesión al cliente.
- **Ventajas**: muy seguro contra XSS (la cookie no es accesible desde JS), revocación inmediata al eliminar la sesión en servidor. Familiar.
- **Desventajas**: requiere estado en el backend (no stateless), puede complicar el escalado horizontal (necesitás compartir sesiones). En Go, con `gorilla/sessions` o similar.

**3. OAuth2 / OpenID Connect (terceros o propio)**
- **Flujo**: delegás la autenticación a un proveedor (Google, GitHub, Azure AD) o montás tu propio servidor de autorización (ej. con Hydra, Keycloak).
- **Ventajas**: Single Sign-On, estándar robusto, separación de responsabilidades, ideal para ERP con múltiples servicios o clientes.
- **Desventajas**: más complejidad de infraestructura para montarlo propio; si usas un proveedor externo, dependencia de terceros.

**4. Autenticación basada en API Keys (para integraciones)**
- Útil para exponer endpoints a otros sistemas (facturación electrónica, automatizaciones). Simple: una clave estática en header.
- No es adecuada para usuarios humanos interactivos.

**Recomendación para tu ERP (Go + Next.js + JWT):**
Usá **JWT con access token en memoria + refresh token en cookie `httpOnly` y `Secure`**. Esto te da seguridad contra XSS, revocación vía rotación de refresh token y una API stateless. Con `golang-jwt/jwt` ya tenés lo necesario; almacená los refresh tokens (o su hash) en PostgreSQL para poder invalidarlos. En Next.js, con SWR podés interceptar errores 401 para refrescar automáticamente el token.

---

Al tener ambos en el mismo dominio (Next.js proxy → Go) la estrategia se simplifica enormemente y queda robusta así:

**Autenticación**
- **Access Token** (JWT de vida corta, ~15 min) viaja en memoria del frontend y se envía manualmente en el header `Authorization: Bearer ...` en cada petición.
- **Refresh Token** (vida larga, ~7 días) se guarda en una cookie `httpOnly`, `Secure` (en producción), `SameSite=Strict` y `Path=/api/auth`. El navegador la adjunta automáticamente solo en peticiones a `/api/auth` y únicamente desde tu propio dominio.
- Flujo de login: `POST /api/auth/login` → el servidor valida credenciales, genera ambos tokens, devuelve el access token en el cuerpo JSON y establece la cookie con el refresh token. El cliente guarda el access token en una variable (Zustand / React state).
- Refresco silencioso: cuando el access token expira, el cliente llama a `POST /api/auth/refresh` (cookie se envía sola) y recibe un nuevo access token en la respuesta.

**Protección de datos en tránsito**
- En producción usás **HTTPS**, que cifra todo el tráfico (login, cambio de contraseña, formularios). No se necesita cifrado adicional en el cuerpo; TLS es suficiente.
- Durante el desarrollo en `localhost`, HTTP es seguro porque el tráfico no sale de tu máquina.

**¿Necesitás CSRF?**  
No con este diseño:
- El access token no está en una cookie, así que no se envía automáticamente; un atacante no puede forzar una petición autenticada a un endpoint protegido.
- El refresh token está en una cookie `httpOnly` y `SameSite=Strict`, por lo que el navegador no la incluye en peticiones de otros sitios. Aunque alguien intente un CSRF hacia `/api/auth/refresh`, no funcionará. Incluso si ocurriera, el atacante no puede leer la respuesta (bloqueada por same-origin policy) ni obtener el nuevo access token.

**Formularios (cambio de contraseña, datos sensibles)**
- Se envían con `POST` y el access token en `Authorization`. El servidor valida el token y procesa la operación.
- Misma protección HTTPS + JWT. No se requiere nada especial extra.

Este enfoque aprovecha al máximo el mismo dominio, elimina la complejidad de CORS y CSRF, y mantiene la seguridad usando los estándares modernos.
