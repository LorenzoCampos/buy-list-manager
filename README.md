# BuyList Manager

Sistema de gestión y tracking de productos para comprar, con seguimiento de precios históricos, categorización inteligente y cálculo automático de costos totales.

## 🎯 Descripción

**BuyList Manager** es una aplicación web diseñada para ayudarte a gestionar productos que querés comprar, ya sean compras únicas (gadgets, herramientas, etc.) o suscripciones recurrentes (servicios, membresías). 

### Características principales:

- ✅ **Gestión de productos** con precio base, envío e impuestos
- 📊 **Categorización flexible**: Compras únicas vs. Recurrentes (mensual/anual)
- 🏷️ **Subcategorías customizables** (AI, Entretenimiento, Trabajo, Reparación, etc.)
- 📈 **Tracking de precios en el tiempo** - Guardás de dónde y cuándo sacaste cada precio
- 💰 **Cálculo automático de costo total**
- 🎨 **Interfaz moderna y responsiva**

---

## 🛠️ Stack Tecnológico

### Frontend
- **React 18+** (con Vite para desarrollo rápido)
- **Tailwind CSS** (styling utilitario)
- **shadcn/ui** (componentes UI pre-diseñados y accesibles)
- **React Router** (navegación)
- **TanStack Query** (manejo de estado del servidor)

### Backend
- **Go 1.21+**
- **Fiber** (web framework minimalista y ultra-rápido)
- **GORM** (ORM para Go)
- **PostgreSQL 15+** driver

### Base de Datos
- **PostgreSQL 15+** (producción y desarrollo)

### DevOps & Tools
- **Docker & Docker Compose** (containerización)
- **Air** (hot-reload para Go en desarrollo)
- **Git** (control de versiones)

---

## 📋 Prerequisitos

Antes de instalar, asegurate de tener:

- **Node.js 18+** y **npm/pnpm**
- **Go 1.21+**
- **PostgreSQL 15+** (o Docker para correrlo containerizado)
- **Git**

---

## 🚀 Instalación y Setup

### 1. Clonar el repositorio

```bash
git clone https://github.com/tu-usuario/buylist-manager.git
cd buylist-manager
```

### 2. Setup del Backend (Go)

```bash
cd backend

# Instalar dependencias
go mod download

# Copiar el archivo de configuración
cp .env.example .env

# Editar .env con tus credenciales de PostgreSQL
# Ejemplo:
# DB_HOST=localhost
# DB_PORT=5432
# DB_USER=postgres
# DB_PASSWORD=tu_password
# DB_NAME=buylist_db

# Correr migraciones
go run cmd/migrate/main.go

# Iniciar el servidor (modo desarrollo con hot-reload)
air
# O sin hot-reload:
go run cmd/api/main.go
```

El backend estará corriendo en `http://localhost:8080`

### 3. Setup del Frontend (React)

```bash
cd frontend

# Instalar dependencias
npm install
# o con pnpm:
pnpm install

# Copiar archivo de configuración
cp .env.example .env

# Editar .env con la URL del backend
# VITE_API_URL=http://localhost:8080

# Iniciar el servidor de desarrollo
npm run dev
```

El frontend estará corriendo en `http://localhost:5173`

### 4. Setup con Docker (Alternativa)

```bash
# Desde la raíz del proyecto
docker-compose up -d

# Esto levanta:
# - PostgreSQL en el puerto 5432
# - Backend en el puerto 8080
# - Frontend en el puerto 5173
```

---

## 📁 Estructura del Proyecto

```
buylist-manager/
├── backend/              # API en Go
│   ├── cmd/
│   │   ├── api/         # Entry point de la API
│   │   └── migrate/     # Script de migraciones
│   ├── internal/
│   │   ├── models/      # Modelos de DB (GORM)
│   │   ├── handlers/    # Controllers/Handlers
│   │   ├── repository/  # Data access layer
│   │   ├── services/    # Lógica de negocio
│   │   └── config/      # Configuración
│   ├── migrations/      # SQL migrations
│   └── go.mod
│
├── frontend/            # App React
│   ├── src/
│   │   ├── components/  # Componentes reutilizables
│   │   ├── pages/       # Páginas/Rutas
│   │   ├── hooks/       # Custom hooks
│   │   ├── services/    # API calls
│   │   ├── lib/         # Utilidades
│   │   └── App.tsx
│   ├── public/
│   └── package.json
│
├── docs/                # Documentación adicional
│   ├── ARCHITECTURE.md
│   └── DATABASE.md
│
├── docker-compose.yml
└── README.md
```

---

## 🎮 Comandos Principales

### Backend (Go)
```bash
# Desarrollo con hot-reload
air

# Correr tests
go test ./...

# Build para producción
go build -o bin/api cmd/api/main.go

# Ejecutar el binario
./bin/api
```

### Frontend (React)
```bash
# Desarrollo
npm run dev

# Build para producción
npm run build

# Preview del build
npm run preview

# Linting
npm run lint
```

---

## 🔌 API Endpoints

### Categories
```
GET    /api/v1/categories          - Listar todas las categorías
GET    /api/v1/categories/:id      - Obtener una categoría por ID
POST   /api/v1/categories          - Crear nueva categoría
PUT    /api/v1/categories/:id      - Actualizar categoría
DELETE /api/v1/categories/:id      - Eliminar categoría
```

### Subcategories
```
GET    /api/v1/subcategories              - Listar subcategorías
GET    /api/v1/subcategories?category_id=1 - Filtrar por categoría
GET    /api/v1/subcategories/:id          - Obtener una subcategoría
POST   /api/v1/subcategories              - Crear subcategoría
PUT    /api/v1/subcategories/:id          - Actualizar subcategoría
DELETE /api/v1/subcategories/:id          - Eliminar subcategoría
```

### Products
```
GET    /api/v1/products                   - Listar todos los productos
GET    /api/v1/products?pending=true      - Productos no comprados
GET    /api/v1/products?category_id=1     - Filtrar por categoría
GET    /api/v1/products/stats             - Estadísticas (totales, gastos)
GET    /api/v1/products/:id               - Obtener un producto
POST   /api/v1/products                   - Crear producto
PUT    /api/v1/products/:id               - Actualizar producto
DELETE /api/v1/products/:id               - Eliminar producto
```

### Health Check
```
GET    /api/v1/health                     - Estado del servidor
```

### Ejemplos de uso

**Crear una categoría:**
```bash
curl -X POST http://localhost:8080/api/v1/categories \
  -H "Content-Type: application/json" \
  -d '{"name": "Compra Única", "type": "one_time"}'
```

**Crear un producto:**
```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Mouse Vertical Logitech MX",
    "description": "Mouse ergonómico inalámbrico",
    "base_price": 75.99,
    "shipping_cost": 12.50,
    "taxes": 5.00,
    "source_url": "https://mercadolibre.com/...",
    "category_id": 1,
    "subcategory_id": 2,
    "notes": "Esperar Black Friday"
  }'
```

**Obtener estadísticas:**
```bash
curl http://localhost:8080/api/v1/products/stats
```

Respuesta:
```json
{
  "total_pending_one_time": 150.49,
  "monthly_recurring_cost": 35.00,
  "yearly_recurring_cost": 420.00
}
```

---

## 🗺️ Roadmap

### ✅ Fase 1 - MVP (Completado)
- ✅ CRUD completo de productos, categorías y subcategorías
- ✅ Sistema de categorización (compras únicas vs recurrentes)
- ✅ Cálculo automático de precio total (base + envío + impuestos)
- ✅ API REST completa con 19 endpoints
- ✅ Validaciones de negocio
- ✅ Estadísticas de gastos (pendientes, mensuales, anuales)
- ✅ Filtros por categoría, subcategoría y estado
- 🚧 Interfaz básica con React + Tailwind (En progreso)

### 🚧 Fase 2 - Features Avanzadas
- Tracking histórico de precios con gráficos
- Filtros y búsqueda avanzada
- Export de data (CSV/JSON)
- Dark mode

### 🔮 Fase 3 - Automatización
- Bookmarklet para importar productos desde páginas web
- Integración con APIs oficiales (MercadoLibre, eBay)
- Sistema de alertas cuando bajan precios
- Multi-usuario con autenticación

---

## 🤝 Contribuciones

Este proyecto está abierto a contribuciones. Si encontrás un bug o querés agregar una feature:

1. Fork el proyecto
2. Creá un branch para tu feature (`git checkout -b feature/nueva-feature`)
3. Commit tus cambios (`git commit -m 'Add: nueva feature'`)
4. Push al branch (`git push origin feature/nueva-feature`)
5. Abrí un Pull Request

---

## 📄 Licencia

MIT License - Sentite libre de usar este proyecto como quieras.

---

## 👨‍💻 Autor

Creado con 🧉 por un dev que estaba cansado de perder track de precios en 15 pestañas abiertas de MercadoLibre.

---

## 📚 Documentación Adicional

- [Arquitectura del Sistema](./docs/ARCHITECTURE.md)
- [Schema de Base de Datos](./docs/DATABASE.md)

---

## ⚠️ Disclaimer

Este proyecto NO realiza web scraping automatizado ni viola los Terms of Service de ninguna plataforma. Toda la información de productos es ingresada manualmente por el usuario o a través de APIs oficiales cuando estén disponibles.
