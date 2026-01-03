# Database Schema - BuyList Manager

## 🗄️ Tecnología

- **Motor**: PostgreSQL 15+
- **ORM**: GORM (Go)
- **Migraciones**: GORM AutoMigrate (MVP) → Migrate tool (producción)

---

## 📊 Diagrama de Relaciones

```
┌─────────────────────────────────────────────────────────────────┐
│                         CATEGORIES                              │
│─────────────────────────────────────────────────────────────────│
│ id            SERIAL PRIMARY KEY                                │
│ name          VARCHAR(100) NOT NULL                             │
│ type          VARCHAR(20) NOT NULL  -- 'one_time' | 'recurring' │
│ created_at    TIMESTAMP                                         │
│ updated_at    TIMESTAMP                                         │
└────────────┬────────────────────────────────────────────────────┘
             │
             │ 1:N
             │
┌────────────▼────────────────────────────────────────────────────┐
│                       SUBCATEGORIES                             │
│─────────────────────────────────────────────────────────────────│
│ id            SERIAL PRIMARY KEY                                │
│ category_id   INTEGER REFERENCES categories(id) ON DELETE CASCADE│
│ name          VARCHAR(100) NOT NULL                             │
│ created_at    TIMESTAMP                                         │
│ updated_at    TIMESTAMP                                         │
└────────────┬────────────────────────────────────────────────────┘
             │
             │ 1:N
             │
┌────────────▼────────────────────────────────────────────────────┐
│                         PRODUCTS                                │
│─────────────────────────────────────────────────────────────────│
│ id                SERIAL PRIMARY KEY                            │
│ name              VARCHAR(255) NOT NULL                         │
│ description       TEXT                                          │
│ base_price        DECIMAL(10,2) NOT NULL                        │
│ shipping_cost     DECIMAL(10,2) DEFAULT 0                       │
│ taxes             DECIMAL(10,2) DEFAULT 0                       │
│ total_price       DECIMAL(10,2) GENERATED ALWAYS AS             │
│                   (base_price + shipping_cost + taxes) STORED   │
│ source_url        VARCHAR(500)   -- Link de dónde sacaste el precio │
│ price_date        TIMESTAMP      -- Cuándo registraste el precio│
│ category_id       INTEGER REFERENCES categories(id)             │
│ subcategory_id    INTEGER REFERENCES subcategories(id)          │
│ recurrence_type   VARCHAR(20)    -- NULL | 'monthly' | 'yearly' │
│ is_purchased      BOOLEAN DEFAULT FALSE                         │
│ purchase_date     TIMESTAMP                                     │
│ notes             TEXT                                          │
│ created_at        TIMESTAMP DEFAULT NOW()                       │
│ updated_at        TIMESTAMP DEFAULT NOW()                       │
└────────────┬────────────────────────────────────────────────────┘
             │
             │ 1:N (Fase 2 - Histórico de precios)
             │
┌────────────▼────────────────────────────────────────────────────┐
│                      PRICE_HISTORY                              │
│                    (FASE 2 - FUTURO)                            │
│─────────────────────────────────────────────────────────────────│
│ id            SERIAL PRIMARY KEY                                │
│ product_id    INTEGER REFERENCES products(id) ON DELETE CASCADE │
│ base_price    DECIMAL(10,2) NOT NULL                            │
│ shipping_cost DECIMAL(10,2)                                     │
│ taxes         DECIMAL(10,2)                                     │
│ total_price   DECIMAL(10,2)                                     │
│ source_url    VARCHAR(500)                                      │
│ recorded_at   TIMESTAMP DEFAULT NOW()                           │
└─────────────────────────────────────────────────────────────────┘
```

---

## 📋 Descripción de Tablas

### 1. `categories`

Categorías principales que dividen productos en **Compra Única** vs **Recurrente**.

| Columna      | Tipo         | Descripción                                  | Ejemplo                |
|--------------|--------------|----------------------------------------------|------------------------|
| `id`         | SERIAL       | Primary key auto-incremental                 | 1, 2, 3...             |
| `name`       | VARCHAR(100) | Nombre de la categoría                       | "Compra Única"         |
| `type`       | VARCHAR(20)  | Tipo: `one_time` o `recurring`               | "one_time"             |
| `created_at` | TIMESTAMP    | Fecha de creación                            | 2026-01-01 10:00:00    |
| `updated_at` | TIMESTAMP    | Última actualización                         | 2026-01-01 10:00:00    |

**Constraints:**
- `type` debe ser uno de: `'one_time'`, `'recurring'`

**Datos iniciales (seed):**
```sql
INSERT INTO categories (name, type) VALUES
  ('Compra Única', 'one_time'),
  ('Suscripción Mensual', 'recurring'),
  ('Suscripción Anual', 'recurring');
```

---

### 2. `subcategories`

Subcategorías que organizan productos dentro de cada categoría principal.

| Columna       | Tipo         | Descripción                                  | Ejemplo                       |
|---------------|--------------|----------------------------------------------|-------------------------------|
| `id`          | SERIAL       | Primary key                                  | 1, 2, 3...                    |
| `category_id` | INTEGER      | FK a `categories.id`                         | 1                             |
| `name`        | VARCHAR(100) | Nombre de la subcategoría                    | "Reparación de Electrónicos"  |
| `created_at`  | TIMESTAMP    | Fecha de creación                            | 2026-01-01 10:00:00           |
| `updated_at`  | TIMESTAMP    | Última actualización                         | 2026-01-01 10:00:00           |

**Constraints:**
- `category_id` → `ON DELETE CASCADE` (si borrás una categoría, se borran sus subcategorías)

**Datos iniciales (seed):**
```sql
INSERT INTO subcategories (category_id, name) VALUES
  -- Subcategorías de "Compra Única" (category_id = 1)
  (1, 'Reparación de Electrónicos'),
  (1, 'Trabajo/Productividad'),
  (1, 'Gaming'),
  (1, 'Hogar'),
  
  -- Subcategorías de suscripciones (category_id = 2 o 3)
  (2, 'IA y Herramientas'),
  (2, 'Entretenimiento'),
  (3, 'Software Profesional');
```

---

### 3. `products`

Tabla principal que almacena los productos a comprar o suscripciones.

| Columna            | Tipo          | Descripción                                             | Ejemplo                                |
|--------------------|---------------|---------------------------------------------------------|----------------------------------------|
| `id`               | SERIAL        | Primary key                                             | 1, 2, 3...                             |
| `name`             | VARCHAR(255)  | Nombre del producto                                     | "Mouse Vertical Logitech MX"           |
| `description`      | TEXT          | Descripción detallada                                   | "Mouse ergonómico inalámbrico..."      |
| `base_price`       | DECIMAL(10,2) | Precio base del producto                                | 75.99                                  |
| `shipping_cost`    | DECIMAL(10,2) | Costo de envío                                          | 12.50                                  |
| `taxes`            | DECIMAL(10,2) | Impuestos/tasas                                         | 5.00                                   |
| `total_price`      | DECIMAL(10,2) | **COMPUTED**: `base_price + shipping_cost + taxes`      | 93.49 (calculado automáticamente)      |
| `source_url`       | VARCHAR(500)  | Link del producto                                       | "https://mercadolibre.com.ar/..."      |
| `price_date`       | TIMESTAMP     | Cuándo registraste este precio                          | 2026-01-01 15:30:00                    |
| `category_id`      | INTEGER       | FK a `categories.id`                                    | 1                                      |
| `subcategory_id`   | INTEGER       | FK a `subcategories.id`                                 | 2                                      |
| `recurrence_type`  | VARCHAR(20)   | NULL para one-time, 'monthly' o 'yearly' para recurring | "monthly"                              |
| `is_purchased`     | BOOLEAN       | Si ya lo compraste                                      | false                                  |
| `purchase_date`    | TIMESTAMP     | Cuándo lo compraste (NULL si no lo compraste aún)       | NULL                                   |
| `notes`            | TEXT          | Notas adicionales                                       | "Esperar Black Friday"                 |
| `created_at`       | TIMESTAMP     | Fecha de creación del registro                          | 2026-01-01 10:00:00                    |
| `updated_at`       | TIMESTAMP     | Última modificación                                     | 2026-01-01 10:00:00                    |

**Constraints:**
- `category_id` → FK con validación
- `subcategory_id` → FK con validación
- `total_price` → **Generated column** (PostgreSQL calcula automáticamente)
- `recurrence_type` → CHECK: debe ser NULL, 'monthly', o 'yearly'

**Notas importantes:**
- Si `category.type = 'one_time'` → `recurrence_type` debe ser NULL
- Si `category.type = 'recurring'` → `recurrence_type` debe ser 'monthly' o 'yearly'

---

### 4. `price_history` (FASE 2 - No implementar en MVP)

Tabla para trackear cambios de precio en el tiempo.

| Columna        | Tipo          | Descripción                          | Ejemplo                      |
|----------------|---------------|--------------------------------------|------------------------------|
| `id`           | SERIAL        | Primary key                          | 1, 2, 3...                   |
| `product_id`   | INTEGER       | FK a `products.id`                   | 5                            |
| `base_price`   | DECIMAL(10,2) | Precio base en ese momento           | 80.00                        |
| `shipping_cost`| DECIMAL(10,2) | Costo de envío en ese momento        | 10.00                        |
| `taxes`        | DECIMAL(10,2) | Impuestos en ese momento             | 5.00                         |
| `total_price`  | DECIMAL(10,2) | Total calculado                      | 95.00                        |
| `source_url`   | VARCHAR(500)  | Link donde se registró el precio     | "https://..."                |
| `recorded_at`  | TIMESTAMP     | Cuándo se registró este precio       | 2026-01-15 10:00:00          |

**Uso futuro:**
- Cada vez que actualizás el precio de un producto, guardás el anterior acá
- Permite graficar evolución de precios
- Útil para saber si conviene comprar ahora o esperar

---

## 🔍 Indexes Recomendados

Para optimizar queries frecuentes:

```sql
-- Búsqueda de productos por categoría
CREATE INDEX idx_products_category ON products(category_id);

-- Búsqueda de productos por subcategoría
CREATE INDEX idx_products_subcategory ON products(subcategory_id);

-- Filtrar productos no comprados
CREATE INDEX idx_products_purchased ON products(is_purchased);

-- Buscar por fecha de precio (para reportes)
CREATE INDEX idx_products_price_date ON products(price_date);

-- FASE 2: Histórico de precios por producto
CREATE INDEX idx_price_history_product ON price_history(product_id, recorded_at DESC);
```

---

## 📝 Queries Comunes

### 1. Listar todos los productos con su categoría y subcategoría

```sql
SELECT 
  p.id,
  p.name,
  p.base_price,
  p.shipping_cost,
  p.taxes,
  p.total_price,
  p.source_url,
  p.price_date,
  c.name AS category_name,
  c.type AS category_type,
  s.name AS subcategory_name,
  p.recurrence_type,
  p.is_purchased
FROM products p
LEFT JOIN categories c ON p.category_id = c.id
LEFT JOIN subcategories s ON p.subcategory_id = s.id
ORDER BY p.created_at DESC;
```

---

### 2. Total a gastar en compras únicas pendientes

```sql
SELECT 
  SUM(total_price) AS total_one_time_pending
FROM products p
JOIN categories c ON p.category_id = c.id
WHERE c.type = 'one_time' 
  AND p.is_purchased = FALSE;
```

---

### 3. Gasto mensual en suscripciones

```sql
SELECT 
  SUM(
    CASE 
      WHEN p.recurrence_type = 'monthly' THEN p.total_price
      WHEN p.recurrence_type = 'yearly' THEN p.total_price / 12
      ELSE 0
    END
  ) AS monthly_recurring_cost
FROM products p
JOIN categories c ON p.category_id = c.id
WHERE c.type = 'recurring';
```

---

### 4. Productos más caros por subcategoría

```sql
SELECT 
  s.name AS subcategory,
  p.name AS product_name,
  p.total_price
FROM products p
JOIN subcategories s ON p.subcategory_id = s.id
WHERE p.total_price = (
  SELECT MAX(p2.total_price)
  FROM products p2
  WHERE p2.subcategory_id = p.subcategory_id
)
ORDER BY p.total_price DESC;
```

---

### 5. Histórico de precios de un producto (FASE 2)

```sql
SELECT 
  ph.recorded_at,
  ph.base_price,
  ph.shipping_cost,
  ph.taxes,
  ph.total_price,
  ph.source_url
FROM price_history ph
WHERE ph.product_id = $1  -- El ID del producto
ORDER BY ph.recorded_at DESC
LIMIT 10;
```

---

## 🛠️ Migraciones con GORM

### Setup inicial (AutoMigrate para MVP)

En tu código Go:

```go
package main

import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

type Category struct {
    ID        uint      `gorm:"primaryKey"`
    Name      string    `gorm:"size:100;not null"`
    Type      string    `gorm:"size:20;not null"` // one_time, recurring
    CreatedAt time.Time
    UpdatedAt time.Time
}

type Subcategory struct {
    ID         uint      `gorm:"primaryKey"`
    CategoryID uint      `gorm:"not null"`
    Category   Category  `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE"`
    Name       string    `gorm:"size:100;not null"`
    CreatedAt  time.Time
    UpdatedAt  time.Time
}

type Product struct {
    ID             uint         `gorm:"primaryKey"`
    Name           string       `gorm:"size:255;not null"`
    Description    string       `gorm:"type:text"`
    BasePrice      float64      `gorm:"type:decimal(10,2);not null"`
    ShippingCost   float64      `gorm:"type:decimal(10,2);default:0"`
    Taxes          float64      `gorm:"type:decimal(10,2);default:0"`
    TotalPrice     float64      `gorm:"type:decimal(10,2);->"` // Read-only, generated
    SourceURL      string       `gorm:"size:500"`
    PriceDate      *time.Time
    CategoryID     uint
    Category       Category     `gorm:"foreignKey:CategoryID"`
    SubcategoryID  uint
    Subcategory    Subcategory  `gorm:"foreignKey:SubcategoryID"`
    RecurrenceType *string      `gorm:"size:20"` // monthly, yearly
    IsPurchased    bool         `gorm:"default:false"`
    PurchaseDate   *time.Time
    Notes          string       `gorm:"type:text"`
    CreatedAt      time.Time
    UpdatedAt      time.Time
}

func main() {
    dsn := "host=localhost user=postgres password=yourpass dbname=buylist_db port=5432 sslmode=disable"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }

    // Auto migrate (solo para desarrollo)
    db.AutoMigrate(&Category{}, &Subcategory{}, &Product{})
}
```

---

### Computed Column para `total_price`

GORM no soporta generated columns nativamente, así que tenés dos opciones:

**Opción 1: Hook de GORM (BeforeSave)**
```go
func (p *Product) BeforeSave(tx *gorm.DB) error {
    p.TotalPrice = p.BasePrice + p.ShippingCost + p.Taxes
    return nil
}
```

**Opción 2: Trigger SQL directo (mejor)**
```sql
CREATE OR REPLACE FUNCTION update_total_price()
RETURNS TRIGGER AS $$
BEGIN
    NEW.total_price := NEW.base_price + NEW.shipping_cost + NEW.taxes;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_total_price
BEFORE INSERT OR UPDATE ON products
FOR EACH ROW
EXECUTE FUNCTION update_total_price();
```

---

## 🔐 Consideraciones de Seguridad

1. **Nunca confiar en el total_price del frontend**: Siempre recalcular en backend
2. **Validar FK**: GORM valida automáticamente que `category_id` y `subcategory_id` existan
3. **Evitar SQL Injection**: GORM usa prepared statements automáticamente
4. **Sanitizar URLs**: Validar que `source_url` sea una URL válida

---

## 📦 Seed Data (Datos Iniciales)

Archivo: `backend/seeds/initial_data.go`

```go
func SeedDatabase(db *gorm.DB) {
    // Categorías
    categories := []Category{
        {Name: "Compra Única", Type: "one_time"},
        {Name: "Suscripción Mensual", Type: "recurring"},
        {Name: "Suscripción Anual", Type: "recurring"},
    }
    db.Create(&categories)

    // Subcategorías
    subcategories := []Subcategory{
        {CategoryID: 1, Name: "Reparación de Electrónicos"},
        {CategoryID: 1, Name: "Trabajo/Productividad"},
        {CategoryID: 1, Name: "Gaming"},
        {CategoryID: 2, Name: "IA y Herramientas"},
        {CategoryID: 2, Name: "Entretenimiento"},
        {CategoryID: 3, Name: "Software Profesional"},
    }
    db.Create(&subcategories)

    // Productos de ejemplo
    now := time.Now()
    products := []Product{
        {
            Name:          "Claude Pro",
            Description:   "Suscripción mensual a Claude AI",
            BasePrice:     20.00,
            ShippingCost:  0,
            Taxes:         0,
            SourceURL:     "https://claude.ai/pro",
            PriceDate:     &now,
            CategoryID:    2,
            SubcategoryID: 4,
            RecurrenceType: stringPtr("monthly"),
        },
        {
            Name:          "Kit Destornilladores de Precisión",
            Description:   "Set de 24 piezas para reparación de electrónicos",
            BasePrice:     15.99,
            ShippingCost:  5.00,
            Taxes:         2.50,
            SourceURL:     "https://mercadolibre.com/...",
            PriceDate:     &now,
            CategoryID:    1,
            SubcategoryID: 1,
        },
    }
    db.Create(&products)
}

func stringPtr(s string) *string {
    return &s
}
```

---

## 🚀 Próximos Pasos (Fase 2)

1. Implementar `price_history` para tracking histórico
2. Agregar índices adicionales según métricas de performance
3. Considerar partitioning de `price_history` por fecha si crece mucho
4. Agregar full-text search en `products.name` y `description`

---

**Última actualización:** 2026-01-01
