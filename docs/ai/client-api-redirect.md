Durante el desarrollo, la forma más sencilla de tener ambos en el mismo dominio (ej. `localhost:3000`) es usar el proxy inverso integrado de Next.js mediante **rewrites**.

**Configuración en Next.js (`next.config.ts`):**
```ts
import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  async rewrites() {
    return [
      {
        source: "/api/:path*",
        destination: "http://localhost:8080/api/:path*", // Tu API Go
      },
    ];
  },
};

export default nextConfig;
```

**Ventajas:**
- El frontend y la API comparten `localhost:3000`, por lo que las cookies (incluidas `httpOnly`) se envían y establecen sin problemas de CORS o `SameSite`.
- No necesitás CORS en el backend Go durante el desarrollo.
- Las rutas `/api/*` en el frontend van directamente a tu servidor Go como si estuvieran en el mismo origen.

**En producción:**
Usá un reverse proxy (Nginx, Traefik) para enrutar todas las peticiones bajo el dominio real: Next.js para rutas de página y Go para `/api/*`. La misma lógica se traslada sin cambios en el código cliente.

Así evitás configurar HTTPS por separado o manejar dominios distintos en desarrollo.
