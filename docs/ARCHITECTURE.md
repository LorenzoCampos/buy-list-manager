# Arquitectura del Sistema - BuyList Manager

## 📐 Visión General

BuyList Manager sigue una arquitectura **cliente-servidor tradicional** con separación clara de responsabilidades:

```
┌─────────────────────────────────────────────────────────────┐
│                         FRONTEND                             │
│                    React + Tailwind CSS                      │
│                                                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │   Pages/     │  │  Components/ │  │   Services/  │      │
│  │   Views      │◄─┤  UI Elements │◄─┤  API Client  │      │
│  └──────────────┘  └──────────────┘  └───────┬──────┘      │
│                                               │              │
└───────────────────────────────────────────────┼──────────────┘
                                                │
                                            HTTP/JSON
                                                │
┌───────────────────────────────────────────────┼──────────────┐
│                         BACKEND                              │
│                      Go + Fiber                              │
│                                                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │   Handlers/  │  │  Services/   │  │ Repository/  │      │
│  │ Controllers  │─►│   Business   │─►│  Data Access │      │
│  └──────────────┘  │    Logic     │  └───────┬──────┘      │
│                    └──────────────┘          │              │
│                                              │              │
└──────────────────────────────────────────────┼──────────────┘
                                               │
                                          SQL Queries
                                               │
                                    ┌──────────▼──────────┐
                                    │   PostgreSQL DB     │
                                    │  (or SQLite local)  │
                                    └─────────────────────┘
```

---

## 🎯 Decisiones de Arquitectura

### 1. **¿Por qué separar Frontend y Backend?**

**Decisión:** Monorepo con frontend y backend separados, pero en el mismo repositorio.

**Razones:**
- ✅ **Flexibilidad de deployment**: Podés deployar frontend (Vercel/Netlify) y backend (Railway/Fly.io) independientemente
- ✅ **Escalabilidad**: Si el día de mañana querés hacer una app mobile, reutilizás el backend
- ✅ **Claridad**: Cada parte tiene sus dependencias y build process aislados
- ✅ **Testing independiente**: Podés testear API y UI por separado

**Tradeoff aceptado:**
- ❌ Más complejo que un monolito (pero marginalmente)
- ❌ CORS headers necesarios en desarrollo

---

### 2. **¿Por qué Go (Fiber) en el Backend?**

**Alternativas consideradas:**
- Node.js (Express/Fastify)
- Laravel (PHP)
- Python (FastAPI)

**Por qué Go + Fiber ganó:**
- ✅ **Performance**: Go es compilado, extremadamente rápido
- ✅ **Single binary deployment**: Compilás y tenés un ejecutable, nada de `node_modules` o dependencias de runtime
- ✅ **Concurrencia nativa**: Goroutines hacen que manejar múltiples requests sea trivial
- ✅ **Tipado estático**: Menos bugs en producción
- ✅ **Fiber**: Sintaxis similar a Express, fácil de aprender, benchmarks excelentes
- ✅ **Tooling para scraping**: Si en el futuro agregamos scraping, Go tiene Colly y Chromedp que son lo mejor que hay

**Tradeoff aceptado:**
- ❌ ORM menos maduro que Eloquent (Laravel) o TypeORM
- ❌ Más verboso que Python/Node para ciertas cosas

---

### 3. **¿Por qué React (no Vue/Svelte/Angular)?**

**Decisión:** React 18+ con Vite

**Razones:**
- ✅ **Ecosistema masivo**: shadcn/ui, Radix UI, TanStack Query, etc.
- ✅ **Developer experience**: Vite es increíblemente rápido
- ✅ **Conocimiento previo del usuario**: Aunque seas más de backend, React es el más común
- ✅ **Hiring/Contributors**: Más fácil que alguien contribuya si conoce React

**Por qué NO los otros:**
- Vue: Excelente, pero ecosistema de componentes menos rico
- Svelte: Increíble performance, pero comunidad más chica
- Angular: Overkill para este proyecto, demasiado opinado

---

### 4. **¿Por qué PostgreSQL?**

**Alternativas consideradas:**
- SQLite
- MySQL
- MongoDB

**Por qué PostgreSQL:**
- ✅ **Relacional**: Este proyecto tiene relaciones claras (productos → categorías → subcategorías)
- ✅ **JSONB**: Si querés guardar specs de productos de forma flexible, PostgreSQL tiene soporte excelente
- ✅ **Open source y maduro**: Cero vendor lock-in
- ✅ **Extensiones**: TimescaleDB para series de tiempo (útil para tracking de precios histórico)
- ✅ **Deploy gratuito**: Supabase, Railway, Neon ofrecen tiers free generosos

**SQLite como alternativa:**
- ✅ Perfecto para desarrollo local (zero config)
- ✅ GORM soporta ambos, podés switchear con config
- ❌ No recomendado para producción con múltiples usuarios

---

## 🏗️ Estructura de Capas

### Backend (Go) - Arquitectura en Capas

```
┌─────────────────────────────────────────────┐
│           HTTP LAYER (Handlers)             │
│  • Recibe requests                          │
│  • Valida input                             │
│  • Devuelve responses HTTP                  │
└──────────────────┬──────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────┐
│         SERVICE LAYER (Business Logic)      │
│  • Lógica de negocio                        │
│  • Cálculo de precios totales               │
│  • Validaciones complejas                   │
│  • Orquestación de múltiples repos          │
└──────────────────┬──────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────┐
│       REPOSITORY LAYER (Data Access)        │
│  • CRUD operations                          │
│  • Queries a la DB                          │
│  • Abstracción de GORM                      │
└──────────────────┬──────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────┐
│              DATABASE (PostgreSQL)          │
└─────────────────────────────────────────────┘
```

**Ejemplo de flujo:**

```go
// 1. Handler recibe request
func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
    var req CreateProductRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(...)
    }
    
    // 2. Llama al Service
    product, err := h.productService.Create(req)
    if err != nil {
        return c.Status(500).JSON(...)
    }
    
    return c.Status(201).JSON(product)
}

// 3. Service ejecuta lógica de negocio
func (s *ProductService) Create(req CreateProductRequest) (*Product, error) {
    // Calcula precio total
    total := req.BasePrice + req.Shipping + req.Taxes
    
    product := &Product{
        Name: req.Name,
        // ... otros campos
        TotalPrice: total,
    }
    
    // 4. Repository guarda en DB
    return s.productRepo.Create(product)
}

// 5. Repository interactúa con GORM
func (r *ProductRepository) Create(product *Product) (*Product, error) {
    result := r.db.Create(product)
    return product, result.Error
}
```

**¿Por qué esta separación?**
- ✅ **Testeable**: Podés mockear cada capa independientemente
- ✅ **Mantenible**: Cambios en DB no afectan lógica de negocio
- ✅ **Reusable**: Services pueden ser llamados desde handlers HTTP, CLI, jobs, etc.

---

### Frontend (React) - Arquitectura por Features

```
src/
├── pages/              # Páginas/Rutas principales
│   ├── Dashboard.tsx   # Vista general
│   ├── Products.tsx    # Lista de productos
│   └── Categories.tsx  # Gestión de categorías
│
├── components/         # Componentes reutilizables
│   ├── ui/            # shadcn/ui components
│   ├── layout/        # Header, Sidebar, Footer
│   └── products/      # ProductCard, ProductForm, etc.
│
├── services/          # API client
│   └── api.ts         # Axios/Fetch wrapper
│
├── hooks/             # Custom hooks
│   ├── useProducts.ts # TanStack Query hooks
│   └── useCategories.ts
│
├── lib/               # Utilidades
│   ├── utils.ts       # Helpers generales
│   └── formatters.ts  # Formateo de moneda, fechas
│
└── types/             # TypeScript types
    └── index.ts       # Interfaces compartidas
```

**Patrón de componentes:**

```tsx
// Container Component (maneja estado y lógica)
function ProductsPage() {
  const { data: products, isLoading } = useProducts()
  
  if (isLoading) return <Spinner />
  
  return <ProductList products={products} />
}

// Presentational Component (solo UI)
function ProductList({ products }) {
  return (
    <div className="grid gap-4">
      {products.map(p => <ProductCard key={p.id} product={p} />)}
    </div>
  )
}
```

---

## 🔄 Flujo de Datos

### Ejemplo: Crear un Producto

```
┌─────────┐
│  USER   │
└────┬────┘
     │ 1. Completa formulario
     ▼
┌─────────────────┐
│ ProductForm.tsx │
└────┬────────────┘
     │ 2. onSubmit → llama a mutation
     ▼
┌──────────────────┐
│ useCreateProduct │  (TanStack Query mutation)
│   (hook)         │
└────┬─────────────┘
     │ 3. POST /api/products
     ▼
┌──────────────────┐
│  API Client      │
│  (axios/fetch)   │
└────┬─────────────┘
     │ 4. HTTP Request
     ▼
┌──────────────────┐
│ Backend Handler  │
└────┬─────────────┘
     │ 5. Valida + llama Service
     ▼
┌──────────────────┐
│ Product Service  │
└────┬─────────────┘
     │ 6. Calcula total + llama Repo
     ▼
┌──────────────────┐
│ Product Repo     │
└────┬─────────────┘
     │ 7. INSERT INTO products
     ▼
┌──────────────────┐
│   PostgreSQL     │
└──────────────────┘
```

---

## 🔐 Seguridad (Fase Futura)

Para el MVP no hay autenticación (uso personal local), pero cuando se agregue:

- **JWT tokens**: Stateless authentication
- **Refresh tokens**: Para renovar sessions
- **CORS**: Configurado correctamente entre frontend y backend
- **Rate limiting**: Fiber middleware para prevenir abuse
- **Sanitización de inputs**: Validación en handler + prepared statements (GORM ya lo hace)

---

## 🚀 Deployment Strategy

### Opción 1: Separado (Recomendado)
- **Frontend**: Vercel/Netlify (build estático de React)
- **Backend**: Railway/Fly.io/Render (binario de Go)
- **DB**: Supabase/Neon/Railway (PostgreSQL managed)

### Opción 2: Todo junto
- **Docker Compose** en un VPS (DigitalOcean, Linode, etc.)
- Nginx como reverse proxy
- Certbot para SSL

### Opción 3: Self-hosted local
- Binario de Go compilado
- `npm run build` del frontend servido por Fiber como static files
- PostgreSQL instalado localmente

---

## 📊 Consideraciones de Performance

### Backend
- **Connection pooling**: GORM maneja el pool de conexiones a PostgreSQL
- **Indexes en DB**: En columnas frecuentemente consultadas (ver DATABASE.md)
- **Eager loading**: Cargar relaciones (categorías) con GORM preloading

### Frontend
- **Code splitting**: Vite lo hace automáticamente
- **Lazy loading**: React.lazy() para rutas no críticas
- **Memoization**: React.memo() para componentes pesados
- **TanStack Query caching**: Reduce requests innecesarias

---

## 🧪 Testing Strategy (Fase Futura)

### Backend
```bash
go test ./...
```
- Unit tests para Services (lógica de negocio)
- Integration tests para Repositories (con DB en memoria)
- E2E tests con httptest para Handlers

### Frontend
- **Vitest**: Unit tests de componentes y hooks
- **Testing Library**: Tests de integración de UI
- **Playwright/Cypress**: E2E tests (opcional)

---

## 🔮 Extensibilidad Futura

Decisiones que facilitan agregar features después:

1. **API versionada**: `/api/v1/products` permite cambios sin romper clientes viejos
2. **Repository pattern**: Fácil cambiar de GORM a otro ORM o queries raw
3. **Service layer**: Lógica compleja encapsulada, fácil de testear
4. **shadcn/ui**: Componentes customizables, fácil agregar theming/dark mode
5. **Monorepo**: Fácil agregar `/mobile` o `/cli` después

---

## 📝 Notas de Implementación

### Variables de Entorno

**Backend (.env)**
```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=buylist_db
DB_SSLMODE=disable

# Server
PORT=8080
ENV=development

# CORS (para desarrollo)
FRONTEND_URL=http://localhost:5173
```

**Frontend (.env)**
```env
VITE_API_URL=http://localhost:8080
```

---

## 🤔 Preguntas Frecuentes de Arquitectura

**Q: ¿Por qué no usar tRPC/GraphQL?**  
**A:** Para este proyecto, REST es suficiente. tRPC requiere TypeScript en backend (no Go), y GraphQL es overkill para un CRUD simple.

**Q: ¿Por qué no un framework full-stack como Next.js?**  
**A:** El usuario quiere backend en Go por performance y experiencia. Next.js fuerza Node.js en el backend.

**Q: ¿Usarían microservicios?**  
**A:** NO. Este es un proyecto pequeño, un monolito bien estructurado es infinitamente más simple y suficiente.

---

**Última actualización:** 2026-01-01
