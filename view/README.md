# View Package

Package view menyediakan template rendering yang powerful, fleksibel, dan production-ready untuk Go applications.

## Features

- ✅ **Template Caching** - Automatic caching untuk production performance
- ✅ **Flexible Structure** - Support layouts, components, dan nested views
- ✅ **Separate Paths** - Template bisa tersebar di berbagai lokasi
- ✅ **Custom Path Resolver** - Full control atas path resolution
- ✅ **Custom Functions** - Built-in functions + custom function support
- ✅ **Global Data** - Shared data across all templates
- ✅ **Thread-Safe** - Concurrent rendering dengan sync.RWMutex
- ✅ **Hot Reload** - Disable cache untuk development
- ✅ **Custom Delimiters** - Support custom template delimiters
- ✅ **Module-Based** - Perfect untuk modular application structure

## Installation

```go
import "github.com/fatkulnurk/foundation/view"
```

## Quick Start

### Basic Usage

```go
v := view.New(view.Config{
    LayoutsPath:    "./view/layouts",
    ComponentsPath: "./view/component",
    ViewsPath:      "./module",
    EnableCache:    false,
})

// Render: ./module/gold/view/price.html
html, err := v.Render(context.Background(), "gold/view/price", map[string]any{
    "Title": "Gold Price",
    "Price": 1000000,
})
```

### With Layout

```go
html, err := v.RenderWithLayout(
    context.Background(),
    "app",              // layout: ./view/layouts/app.html
    "gold/view/price",  // view: ./module/gold/view/price.html
    data,
)
```

## Directory Structure

```
project/
├── view/
│   ├── layouts/
│   │   ├── app.html
│   │   └── admin.html
│   └── component/
│       ├── header.html
│       ├── footer.html
│       └── table.html
└── module/
    ├── gold/
    │   └── view/
    │       ├── price.html
    │       └── chart.html
    └── user/
        └── view/
            ├── profile.html
            └── settings.html
```

## Configuration

```go
type Config struct {
    // Full path ke directory layouts (e.g., "./view/layouts")
    LayoutsPath string

    // Full path ke directory components (e.g., "./view/component")
    ComponentsPath string

    // Full path ke directory views (e.g., "./module")
    ViewsPath string

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

    // Custom path resolver function
    PathResolver func(templateType, name string) string
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

### 1. Separate Paths Configuration

```go
v := view.New(view.Config{
    LayoutsPath:    "./view/layouts",
    ComponentsPath: "./view/component",
    ViewsPath:      "./module",
    EnableCache:    true,
})

// Layout: ./view/layouts/app.html
// Components: ./view/component/*.html
// View: ./module/gold/view/price.html
html, err := v.RenderWithLayout(ctx, "app", "gold/view/price", data)
```

### 2. Custom Path Resolver

```go
v := view.New(view.Config{
    EnableCache: false,
    PathResolver: func(templateType, name string) string {
        switch templateType {
        case "layout":
            return filepath.Join("./view/layouts", name+".html")
        case "component":
            return filepath.Join("./view/component", name+".html")
        case "view":
            // Custom logic: modulename/viewname -> ./module/modulename/view/viewname.html
            parts := strings.Split(name, "/")
            if len(parts) >= 2 {
                moduleName := parts[0]
                viewName := strings.Join(parts[1:], "/")
                return filepath.Join("./module", moduleName, "view", viewName+".html")
            }
            return filepath.Join("./views", name+".html")
        default:
            return name + ".html"
        }
    },
})

// Dengan custom resolver, "gold/price" akan resolve ke:
// ./module/gold/view/price.html
html, err := v.Render(ctx, "gold/price", data)
```

### 3. Global Data

```go
v := view.New(view.Config{
    LayoutsPath:    "./view/layouts",
    ComponentsPath: "./view/component",
    ViewsPath:      "./module",
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
