# View Package

Package view menyediakan template rendering yang powerful, fleksibel, dan production-ready untuk Go applications.

## Features

- ✅ **Template Caching** - Automatic caching untuk production performance
- ✅ **Flexible Structure** - Support layouts, components, dan nested views
- ✅ **Custom Functions** - Built-in functions + custom function support
- ✅ **Global Data** - Shared data across all templates
- ✅ **Thread-Safe** - Concurrent rendering dengan sync.RWMutex
- ✅ **Hot Reload** - Disable cache untuk development
- ✅ **Custom Delimiters** - Support custom template delimiters
- ✅ **Path Resolution** - Smart path resolution untuk nested templates

## Installation

```go
import "github.com/fatkulnurk/foundation/view"
```

## Quick Start

### Basic Usage

```go
v := view.New(view.Config{
    TemplateDir: "./templates",
    EnableCache: false, // false untuk development
})

html, err := v.Render(context.Background(), "home", map[string]any{
    "Title": "Welcome",
    "Message": "Hello, World!",
})
```

### With Layout

```go
html, err := v.RenderWithLayout(
    context.Background(),
    "main",           // layout name
    "gold/price",     // view name (support nested path)
    data,
)
```

## Directory Structure

```
templates/
├── layouts/
│   ├── main.html
│   └── admin.html
├── components/
│   ├── header.html
│   ├── footer.html
│   └── table.html
└── views/
    ├── home.html
    ├── about.html
    └── gold/
        ├── price.html
        └── chart.html
```

## Configuration

```go
type Config struct {
    // Root directory untuk semua template
    TemplateDir string

    // Subdirectory untuk layouts (default: "layouts")
    LayoutsDir string

    // Subdirectory untuk components (default: "components")
    ComponentsDir string

    // Subdirectory untuk views (default: "views")
    ViewsDir string

    // File extension (default: ".html")
    Extension string

    // Enable/disable caching (default: false)
    EnableCache bool

    // Custom template delimiters (optional)
    LeftDelim  string
    RightDelim string

    // Custom template functions
    FuncMap template.FuncMap

    // Global data tersedia di semua template
    GlobalData map[string]any
}
```

## Built-in Functions

### HTML/JS/URL Safety
- `raw` - Render HTML tanpa escaping
- `safeHTML` - Escape HTML
- `safeJS` - Escape JavaScript
- `safeURL` - Safe URL

### Date/Time
- `now` - Current time
- `formatDate` - Format time dengan layout
- `year` - Current year

### Math
- `add` - Addition
- `sub` - Subtraction
- `mul` - Multiplication
- `div` - Division

### String
- `join` - Join strings
- `upper` - Uppercase
- `lower` - Lowercase
- `title` - Title case

### Utility
- `default` - Default value if nil/empty
- `global` - Access global data

## Examples

### 1. Global Data

```go
v := view.New(view.Config{
    TemplateDir: "./templates",
    GlobalData: map[string]any{
        "SiteName": "My Site",
        "Version": "1.0.0",
    },
})

// Tambah global data dinamis
v.SetGlobal("Year", 2025)
```

Template:
```html
<footer>
    <p>{{ global "SiteName" }} v{{ global "Version" }}</p>
    <p>&copy; {{ global "Year" }}</p>
</footer>
```

### 2. Custom Functions

```go
v := view.New(view.Config{
    TemplateDir: "./templates",
    FuncMap: template.FuncMap{
        "formatCurrency": func(amount int) string {
            return fmt.Sprintf("Rp %d", amount)
        },
    },
})

// Tambah function dinamis
v.AddFunc("greet", func(name string) string {
    return "Hello, " + name + "!"
})
```

Template:
```html
<p>Price: {{ formatCurrency .Price }}</p>
<p>{{ greet .UserName }}</p>
```

### 3. Layout Template

**layouts/main.html:**
```html
<!DOCTYPE html>
<html>
<head>
    <title>{{ .Title }} - {{ global "SiteName" }}</title>
</head>
<body>
    {{ template "header" . }}
    
    <main>
        {{ template "content" . }}
    </main>
    
    {{ template "footer" . }}
</body>
</html>
```

**views/home.html:**
```html
{{ define "content" }}
<h1>{{ .Title }}</h1>
<p>{{ .Message }}</p>
{{ end }}
```

**components/header.html:**
```html
{{ define "header" }}
<header>
    <h1>{{ global "SiteName" }}</h1>
</header>
{{ end }}
```

### 4. Production Setup

```go
v := view.New(view.Config{
    TemplateDir:   "./templates",
    LayoutsDir:    "layouts",
    ComponentsDir: "components",
    ViewsDir:      "views",
    Extension:     ".html",
    EnableCache:   true, // Enable untuk production
    GlobalData: map[string]any{
        "SiteName": "Harga Emas",
        "Year":     2025,
    },
})

// Render akan menggunakan cache
html, err := v.RenderWithLayout(ctx, "main", "gold/price", data)
```

### 5. Development Setup

```go
v := view.New(view.Config{
    TemplateDir: "./templates",
    EnableCache: false, // Disable untuk hot reload
})

// Template akan di-parse ulang setiap render
html, err := v.Render(ctx, "home", data)
```

### 6. Custom Delimiters

```go
v := view.New(view.Config{
    TemplateDir: "./templates",
    LeftDelim:   "[[",
    RightDelim:  "]]",
})
```

Template:
```html
<h1>[[ .Title ]]</h1>
<p>[[ .Message ]]</p>
```

## API Reference

### Methods

#### `Render(ctx context.Context, name string, data any) (string, error)`
Render template tanpa layout.

#### `RenderWithLayout(ctx context.Context, layout, name string, data any) (string, error)`
Render template dengan layout spesifik.

#### `AddFunc(name string, fn any) View`
Tambah custom function. Returns self untuk method chaining.

#### `SetGlobal(key string, value any) View`
Set global data. Returns self untuk method chaining.

#### `ClearCache()`
Clear template cache. Berguna saat hot reload atau update template.

## Best Practices

### 1. Production vs Development

```go
// Development
v := view.New(view.Config{
    TemplateDir: "./templates",
    EnableCache: false, // Hot reload
})

// Production
v := view.New(view.Config{
    TemplateDir: "./templates",
    EnableCache: true, // Performance
})
```

### 2. Error Handling

```go
html, err := v.Render(ctx, "home", data)
if err != nil {
    log.Printf("Template error: %v", err)
    // Fallback atau error page
    return
}
```

### 3. Global Data Management

```go
// Set saat initialization
v := view.New(view.Config{
    GlobalData: map[string]any{
        "SiteName": "My Site",
    },
})

// Update dinamis
v.SetGlobal("UserCount", getUserCount())
v.SetGlobal("LastUpdate", time.Now())
```

### 4. Template Organization

```
templates/
├── layouts/          # Reusable layouts
├── components/       # Reusable components
└── views/           # Page-specific views
    ├── auth/        # Grouped by feature
    ├── admin/
    └── public/
```

## Performance

- **Caching**: Enable `EnableCache: true` untuk production
- **Concurrency**: Thread-safe dengan RWMutex
- **Memory**: Template di-cache di memory untuk fast access

## License

MIT
