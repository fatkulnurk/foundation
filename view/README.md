# View Package

A powerful, flexible template rendering engine for Go applications with caching, layouts, components, and modular structure support.

## Table of Contents

- [What is View?](#what-is-view)
- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Configuration](#configuration)
- [Directory Structure](#directory-structure)
- [Built-in Functions](#built-in-functions)
- [Complete Examples](#complete-examples)
- [API Reference](#api-reference)
- [Best Practices](#best-practices)

---

## What is View?

The view package helps you **render HTML templates** in your Go applications. It's like a smart template engine that can combine layouts, components, and views to create complete web pages.

**Simple Analogy:**
- **View** = The main content of a page (like a blog post)
- **Layout** = The frame around content (header, footer, navigation)
- **Component** = Reusable pieces (buttons, cards, tables)
- **Template** = The blueprint with placeholders for data
- **Render** = Fill in the placeholders with actual data

---

## Features

-  **Template Caching** - Automatic caching for production performance
-  **Flexible Structure** - Support for layouts, components, and nested views
-  **Separate Paths** - Templates can be in different locations
-  **Custom Path Resolver** - Full control over path resolution
-  **20+ Built-in Functions** - Date, math, string, and utility functions
-  **Custom Functions** - Add your own template functions
-  **Global Data** - Shared data across all templates
-  **Thread-Safe** - Concurrent rendering with sync.RWMutex
-  **Hot Reload** - Disable cache for development
-  **Custom Delimiters** - Support custom template delimiters
-  **Module-Based** - Perfect for modular application structure

---

## Installation

```bash
go get github.com/fatkulnurk/foundation/view
```

**Dependencies:**
- Go 1.25 or higher
- Standard library only (no external dependencies)

---

## Quick Start

### 1. Basic Usage (No Layout)

```go
package main

import (
    "context"
    "github.com/fatkulnurk/foundation/view"
)

func main() {
    // Create view instance
    v := view.New(view.Config{
        ViewsPath:   "./templates",
        EnableCache: false, // Disable for development
    })
    
    // Render template
    html, err := v.Render(context.Background(), "home", map[string]any{
        "Title":   "Welcome",
        "Message": "Hello, World!",
    })
    
    if err != nil {
        panic(err)
    }
    
    println(html)
}
```

**templates/home.html:**
```html
<!DOCTYPE html>
<html>
<head>
    <title>{{ .Title }}</title>
</head>
<body>
    <h1>{{ .Title }}</h1>
    <p>{{ .Message }}</p>
</body>
</html>
```

### 2. With Layout

```go
v := view.New(view.Config{
    LayoutsPath:    "./view/layouts",
    ComponentsPath: "./view/components",
    ViewsPath:      "./view/pages",
    EnableCache:    false,
})

// Render with layout
html, err := v.RenderWithLayout(
    context.Background(),
    "main",  // Layout: ./view/layouts/main.html
    "home",  // View: ./view/pages/home.html
    data,
)
```

### 3. Modular Structure

```go
v := view.New(view.Config{
    LayoutsPath:    "./view/layouts",
    ComponentsPath: "./view/components",
    ViewsPath:      "./modules",
    EnableCache:    true, // Enable for production
})

// Render: ./modules/gold/view/price.html
html, err := v.RenderWithLayout(
    ctx,
    "app",              // Layout
    "gold/view/price",  // Module view
    data,
)
```

---

## Configuration

```go
type Config struct {
    // Full path to layouts directory (e.g., "./view/layouts")
    LayoutsPath string

    // Full path to components directory (e.g., "./view/components")
    ComponentsPath string

    // Full path to views directory (e.g., "./views" or "./modules")
    ViewsPath string

    // File extension (default: ".html")
    Extension string

    // Enable/disable caching (default: false)
    // Set to true for production
    EnableCache bool

    // Custom template delimiters (optional)
    // Default: "{{" and "}}"
    LeftDelim  string
    RightDelim string

    // Custom template functions
    FuncMap template.FuncMap

    // Global data available in all templates
    GlobalData map[string]any

    // Custom path resolver function
    // Signature: func(templateType, name string) string
    // templateType: "layout", "component", "view"
    PathResolver func(templateType, name string) string
}
```

### Configuration Examples

#### Development Setup
```go
v := view.New(view.Config{
    LayoutsPath:    "./view/layouts",
    ComponentsPath: "./view/components",
    ViewsPath:      "./view/pages",
    EnableCache:    false, // Hot reload
})
```

#### Production Setup
```go
v := view.New(view.Config{
    LayoutsPath:    "./view/layouts",
    ComponentsPath: "./view/components",
    ViewsPath:      "./view/pages",
    EnableCache:    true, // Performance
    GlobalData: map[string]any{
        "SiteName": "My Website",
        "Version":  "1.0.0",
    },
})
```

#### Modular Application
```go
v := view.New(view.Config{
    LayoutsPath:    "./view/layouts",
    ComponentsPath: "./view/components",
    ViewsPath:      "./modules", // Each module has its own views
    EnableCache:    true,
})

// Render: ./modules/user/view/profile.html
v.RenderWithLayout(ctx, "app", "user/view/profile", data)
```

---

## Directory Structure

### Standard Structure

```
project/
├── view/
│   ├── layouts/
│   │   ├── main.html
│   │   ├── admin.html
│   │   └── auth.html
│   ├── components/
│   │   ├── header.html
│   │   ├── footer.html
│   │   ├── navbar.html
│   │   └── card.html
│   └── pages/
│       ├── home.html
│       ├── about.html
│       └── contact.html
```

### Modular Structure

```
project/
├── view/
│   ├── layouts/
│   │   └── app.html
│   └── components/
│       ├── header.html
│       └── footer.html
└── modules/
    ├── user/
    │   └── view/
    │       ├── profile.html
    │       ├── settings.html
    │       └── dashboard.html
    ├── product/
    │   └── view/
    │       ├── list.html
    │       ├── detail.html
    │       └── create.html
    └── admin/
        └── view/
            ├── dashboard.html
            └── users.html
```

---

## Built-in Functions

### HTML/JS/URL Safety

#### `raw`
Render HTML without escaping.

```html
{{ raw .HTMLContent }}
```

#### `safeHTML`
Escape HTML for safe display.

```html
{{ safeHTML .UserInput }}
```

#### `safeJS`
Escape JavaScript strings.

```html
<script>
var message = {{ safeJS .Message }};
</script>
```

#### `safeURL`
Safe URL rendering.

```html
<a href="{{ safeURL .Link }}">Click here</a>
```

### Date/Time Functions

#### `now`
Get current time.

```html
<p>Current time: {{ now }}</p>
```

#### `formatDate`
Format time with layout.

```html
<p>Date: {{ formatDate .CreatedAt "2006-01-02" }}</p>
<p>Time: {{ formatDate .CreatedAt "15:04:05" }}</p>
<p>Full: {{ formatDate .CreatedAt "2006-01-02 15:04:05" }}</p>
```

#### `year`
Get current year.

```html
<footer>&copy; {{ year }} My Company</footer>
```

### Math Functions

#### `add`
Addition.

```html
<p>Total: {{ add .Price .Tax }}</p>
```

#### `sub`
Subtraction.

```html
<p>Discount: {{ sub .OriginalPrice .CurrentPrice }}</p>
```

#### `mul`
Multiplication.

```html
<p>Total: {{ mul .Price .Quantity }}</p>
```

#### `div`
Division.

```html
<p>Average: {{ div .Total .Count }}</p>
```

### String Functions

#### `join`
Join strings with separator.

```html
<p>Tags: {{ join ", " .Tags }}</p>
```

#### `upper`
Convert to uppercase.

```html
<h1>{{ upper .Title }}</h1>
```

#### `lower`
Convert to lowercase.

```html
<p>{{ lower .Email }}</p>
```

#### `title`
Convert to title case.

```html
<h2>{{ title .Name }}</h2>
```

### Utility Functions

#### `default`
Provide default value if nil/empty.

```html
<p>Name: {{ default "Guest" .UserName }}</p>
```

#### `global`
Access global data.

```html
<title>{{ .Title }} - {{ global "SiteName" }}</title>
<footer>&copy; {{ global "Year" }} {{ global "Company" }}</footer>
```

---

## Complete Examples

### Example 1: Simple Page

**main.go:**
```go
package main

import (
    "context"
    "github.com/fatkulnurk/foundation/view"
)

func main() {
    v := view.New(view.Config{
        ViewsPath:   "./templates",
        EnableCache: false,
    })
    
    html, _ := v.Render(context.Background(), "home", map[string]any{
        "Title":   "Welcome",
        "Message": "Hello, World!",
    })
    
    println(html)
}
```

**templates/home.html:**
```html
<!DOCTYPE html>
<html>
<head>
    <title>{{ .Title }}</title>
</head>
<body>
    <h1>{{ .Title }}</h1>
    <p>{{ .Message }}</p>
</body>
</html>
```

### Example 2: With Layout and Components

**main.go:**
```go
v := view.New(view.Config{
    LayoutsPath:    "./view/layouts",
    ComponentsPath: "./view/components",
    ViewsPath:      "./view/pages",
    EnableCache:    false,
    GlobalData: map[string]any{
        "SiteName": "My Website",
    },
})

html, _ := v.RenderWithLayout(
    context.Background(),
    "main",
    "home",
    map[string]any{
        "Title":   "Home Page",
        "Message": "Welcome to our website!",
    },
)
```

**view/layouts/main.html:**
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

**view/components/header.html:**
```html
{{ define "header" }}
<header>
    <h1>{{ global "SiteName" }}</h1>
    <nav>
        <a href="/">Home</a>
        <a href="/about">About</a>
        <a href="/contact">Contact</a>
    </nav>
</header>
{{ end }}
```

**view/components/footer.html:**
```html
{{ define "footer" }}
<footer>
    <p>&copy; {{ year }} {{ global "SiteName" }}</p>
</footer>
{{ end }}
```

**view/pages/home.html:**
```html
{{ define "content" }}
<h1>{{ .Title }}</h1>
<p>{{ .Message }}</p>
{{ end }}
```

### Example 3: Custom Functions

```go
v := view.New(view.Config{
    ViewsPath: "./templates",
    FuncMap: template.FuncMap{
        "formatCurrency": func(amount int) string {
            return fmt.Sprintf("$%d", amount)
        },
        "isEven": func(n int) bool {
            return n%2 == 0
        },
    },
})

// Add more functions dynamically
v.AddFunc("greet", func(name string) string {
    return "Hello, " + name + "!"
})

html, _ := v.Render(ctx, "product", map[string]any{
    "Name":  "Laptop",
    "Price": 1000,
})
```

**templates/product.html:**
```html
<div class="product">
    <h2>{{ .Name }}</h2>
    <p class="price">{{ formatCurrency .Price }}</p>
    <p>{{ greet "Customer" }}</p>
</div>
```

### Example 4: Global Data

```go
v := view.New(view.Config{
    ViewsPath: "./templates",
    GlobalData: map[string]any{
        "SiteName": "My Shop",
        "Version":  "1.0.0",
        "Year":     2024,
    },
})

// Update global data dynamically
v.SetGlobal("UserCount", 1000)
v.SetGlobal("LastUpdate", time.Now())
```

**templates/page.html:**
```html
<footer>
    <p>{{ global "SiteName" }} v{{ global "Version" }}</p>
    <p>&copy; {{ global "Year" }}</p>
    <p>Users: {{ global "UserCount" }}</p>
    <p>Last update: {{ formatDate (global "LastUpdate") "2006-01-02 15:04" }}</p>
</footer>
```

### Example 5: Custom Path Resolver

```go
v := view.New(view.Config{
    EnableCache: false,
    PathResolver: func(templateType, name string) string {
        switch templateType {
        case "layout":
            return filepath.Join("./view/layouts", name+".html")
        case "component":
            return filepath.Join("./view/components", name+".html")
        case "view":
            // Custom logic for module-based structure
            // "user/profile" -> "./modules/user/view/profile.html"
            parts := strings.Split(name, "/")
            if len(parts) >= 2 {
                moduleName := parts[0]
                viewName := strings.Join(parts[1:], "/")
                return filepath.Join("./modules", moduleName, "view", viewName+".html")
            }
            return filepath.Join("./views", name+".html")
        default:
            return name + ".html"
        }
    },
})

// "user/profile" resolves to "./modules/user/view/profile.html"
html, _ := v.Render(ctx, "user/profile", data)
```

### Example 6: Custom Delimiters

```go
v := view.New(view.Config{
    ViewsPath:  "./templates",
    LeftDelim:  "[[",
    RightDelim: "]]",
})
```

**templates/page.html:**
```html
<h1>[[ .Title ]]</h1>
<p>[[ .Message ]]</p>
<p>Price: [[ formatCurrency .Price ]]</p>
```

### Example 7: HTTP Handler Integration

```go
package main

import (
    "context"
    "net/http"
    "github.com/fatkulnurk/foundation/view"
)

var v view.View

func init() {
    v = view.New(view.Config{
        LayoutsPath:    "./view/layouts",
        ComponentsPath: "./view/components",
        ViewsPath:      "./view/pages",
        EnableCache:    true,
        GlobalData: map[string]any{
            "SiteName": "My Website",
        },
    })
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    html, err := v.RenderWithLayout(
        r.Context(),
        "main",
        "home",
        map[string]any{
            "Title":   "Home",
            "Message": "Welcome!",
        },
    )
    
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    w.Write([]byte(html))
}

func main() {
    http.HandleFunc("/", homeHandler)
    http.ListenAndServe(":8080", nil)
}
```

### Example 8: Conditional Rendering

**templates/user.html:**
```html
{{ define "content" }}
<div class="user-profile">
    <h1>{{ .User.Name }}</h1>
    
    {{ if .User.IsAdmin }}
        <span class="badge">Admin</span>
    {{ end }}
    
    {{ if .User.Email }}
        <p>Email: {{ .User.Email }}</p>
    {{ else }}
        <p>No email provided</p>
    {{ end }}
    
    <h2>Posts</h2>
    {{ if .Posts }}
        <ul>
        {{ range .Posts }}
            <li>{{ .Title }} - {{ formatDate .CreatedAt "2006-01-02" }}</li>
        {{ end }}
        </ul>
    {{ else }}
        <p>No posts yet</p>
    {{ end }}
</div>
{{ end }}
```

---

## API Reference

### View Interface

```go
type View interface {
    // Render renders a template without layout
    Render(ctx context.Context, name string, data any) (string, error)
    
    // RenderWithLayout renders a template with a specific layout
    RenderWithLayout(ctx context.Context, layout, name string, data any) (string, error)
    
    // AddFunc adds a custom function to templates
    // Returns self for method chaining
    AddFunc(name string, fn any) View
    
    // SetGlobal sets global data available in all templates
    // Returns self for method chaining
    SetGlobal(key string, value any) View
    
    // ClearCache clears the template cache
    // Useful for hot reload or template updates
    ClearCache()
}
```

### Methods

#### `Render(ctx context.Context, name string, data any) (string, error)`

Renders a template without a layout.

**Parameters:**
- `ctx`: Context for cancellation
- `name`: Template name (e.g., "home", "user/profile")
- `data`: Data to pass to the template

**Returns:**
- `string`: Rendered HTML
- `error`: Error if rendering fails

**Example:**
```go
html, err := v.Render(ctx, "home", map[string]any{
    "Title": "Home Page",
})
```

#### `RenderWithLayout(ctx context.Context, layout, name string, data any) (string, error)`

Renders a template with a specific layout.

**Parameters:**
- `ctx`: Context for cancellation
- `layout`: Layout name (e.g., "main", "admin")
- `name`: Template name (e.g., "home", "user/profile")
- `data`: Data to pass to the template

**Returns:**
- `string`: Rendered HTML
- `error`: Error if rendering fails

**Example:**
```go
html, err := v.RenderWithLayout(ctx, "main", "home", data)
```

#### `AddFunc(name string, fn any) View`

Adds a custom function to templates.

**Parameters:**
- `name`: Function name to use in templates
- `fn`: Function implementation

**Returns:**
- `View`: Self for method chaining

**Example:**
```go
v.AddFunc("formatPrice", func(price int) string {
    return fmt.Sprintf("$%d.00", price)
}).AddFunc("isPositive", func(n int) bool {
    return n > 0
})
```

#### `SetGlobal(key string, value any) View`

Sets global data available in all templates.

**Parameters:**
- `key`: Data key
- `value`: Data value

**Returns:**
- `View`: Self for method chaining

**Example:**
```go
v.SetGlobal("SiteName", "My Website").
  SetGlobal("Year", 2024).
  SetGlobal("Version", "1.0.0")
```

#### `ClearCache()`

Clears the template cache. Useful when templates are updated or for hot reload.

**Example:**
```go
v.ClearCache()
```

---

## Best Practices

### 1. Enable Caching in Production

```go
// Development
v := view.New(view.Config{
    ViewsPath:   "./templates",
    EnableCache: false, // Hot reload
})

// Production
v := view.New(view.Config{
    ViewsPath:   "./templates",
    EnableCache: true, // Performance
})
```

### 2. Use Layouts for Consistent Design

```go
// All pages use the same layout
html, _ := v.RenderWithLayout(ctx, "main", "home", data)
html, _ := v.RenderWithLayout(ctx, "main", "about", data)
html, _ := v.RenderWithLayout(ctx, "main", "contact", data)

// Admin pages use different layout
html, _ := v.RenderWithLayout(ctx, "admin", "dashboard", data)
```

### 3. Organize Templates by Feature

```
view/
├── layouts/
│   ├── main.html       # Public layout
│   ├── admin.html      # Admin layout
│   └── auth.html       # Authentication layout
├── components/
│   ├── header.html
│   ├── footer.html
│   ├── navbar.html
│   └── sidebar.html
└── pages/
    ├── home.html
    ├── about.html
    ├── auth/
    │   ├── login.html
    │   └── register.html
    └── admin/
        ├── dashboard.html
        └── users.html
```

### 4. Use Global Data for Site-Wide Information

```go
v := view.New(view.Config{
    ViewsPath: "./templates",
    GlobalData: map[string]any{
        "SiteName":    "My Website",
        "Version":     "1.0.0",
        "SupportEmail": "support@example.com",
    },
})

// Update dynamically
v.SetGlobal("UserCount", getUserCount())
v.SetGlobal("LastUpdate", time.Now())
```

### 5. Handle Errors Properly

```go
html, err := v.Render(ctx, "home", data)
if err != nil {
    log.Printf("Template error: %v", err)
    // Show error page or fallback
    return
}
```

### 6. Use Custom Functions for Reusable Logic

```go
v.AddFunc("formatCurrency", func(amount int) string {
    return fmt.Sprintf("$%d.00", amount)
}).AddFunc("truncate", func(s string, length int) string {
    if len(s) <= length {
        return s
    }
    return s[:length] + "..."
})
```

### 7. Separate Concerns with Components

```html
<!-- Layout uses components -->
{{ template "header" . }}
{{ template "navbar" . }}
{{ template "content" . }}
{{ template "footer" . }}
```

### 8. Use Context for Cancellation

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

html, err := v.Render(ctx, "home", data)
```

### 9. Clear Cache When Templates Change

```go
// After updating templates
v.ClearCache()

// Or reload templates
v = view.New(config)
```

### 10. Test Your Templates

```go
func TestHomeTemplate(t *testing.T) {
    v := view.New(view.Config{
        ViewsPath:   "./templates",
        EnableCache: false,
    })
    
    html, err := v.Render(context.Background(), "home", map[string]any{
        "Title": "Test",
    })
    
    if err != nil {
        t.Fatalf("Render failed: %v", err)
    }
    
    if !strings.Contains(html, "Test") {
        t.Error("Title not found in rendered HTML")
    }
}
```

---

## Performance

### Caching

- **Development**: Set `EnableCache: false` for hot reload
- **Production**: Set `EnableCache: true` for performance
- Templates are parsed once and cached in memory
- Subsequent renders use cached templates

### Concurrency

- Thread-safe with `sync.RWMutex`
- Multiple goroutines can render templates simultaneously
- Cache reads are concurrent
- Cache writes are exclusive

### Memory

- Templates are stored in memory for fast access
- Clear cache with `ClearCache()` if memory is a concern
- Use `EnableCache: false` to avoid caching

---

## Troubleshooting

### Error: "template not found"

**Problem:** Template file doesn't exist

**Solution:**
- Check file path and name
- Verify `ViewsPath`, `LayoutsPath`, `ComponentsPath` configuration
- Ensure file has correct extension (default: `.html`)

### Error: "template: ... is undefined"

**Problem:** Referenced template (layout/component) not found

**Solution:**
- Ensure layout file exists in `LayoutsPath`
- Ensure component files exist in `ComponentsPath`
- Check template `{{ define "name" }}` matches `{{ template "name" }}`

### Templates Not Updating

**Problem:** Cache is enabled

**Solution:**
- Set `EnableCache: false` for development
- Call `v.ClearCache()` after template changes
- Restart application

### Function Not Found

**Problem:** Custom function not registered

**Solution:**
- Register function with `AddFunc()` before rendering
- Check function name matches template usage
- Ensure function is registered before first render

---

## Extending

You can create custom template functions and path resolvers to extend the view package functionality.

### Custom Template Functions

```go
// Add custom functions when creating view
v := view.New(view.Config{
    ViewsPath: "./templates",
    FuncMap: template.FuncMap{
        "myFunction": func(s string) string {
            // Your logic
            return strings.ToUpper(s)
        },
    },
})

// Or add dynamically
v.AddFunc("formatPrice", func(price float64) string {
    return fmt.Sprintf("$%.2f", price)
})
```

### Custom Path Resolver

```go
v := view.New(view.Config{
    PathResolver: func(templateType, name string) string {
        switch templateType {
        case "layout":
            return filepath.Join("./custom/layouts", name+".html")
        case "component":
            return filepath.Join("./custom/components", name+".html")
        case "view":
            // Custom logic for views
            return filepath.Join("./custom/views", name+".html")
        default:
            return name + ".html"
        }
    },
})
```

### Example: Custom Markdown Function

```go
import "github.com/russross/blackfriday/v2"

v := view.New(view.Config{
    ViewsPath: "./templates",
})

v.AddFunc("markdown", func(content string) template.HTML {
    output := blackfriday.Run([]byte(content))
    return template.HTML(output)
})
```

**Template usage:**
```html
<div class="content">
    {{ markdown .ArticleContent }}
</div>
```

### Example: Custom Date Formatting

```go
v.AddFunc("formatDateTime", func(t time.Time, format string) string {
    layouts := map[string]string{
        "short": "2006-01-02",
        "long":  "January 2, 2006",
        "time":  "15:04:05",
        "full":  "2006-01-02 15:04:05",
    }
    
    layout, ok := layouts[format]
    if !ok {
        layout = format
    }
    
    return t.Format(layout)
})
```

**Template usage:**
```html
<p>Date: {{ formatDateTime .CreatedAt "short" }}</p>
<p>Full: {{ formatDateTime .CreatedAt "full" }}</p>
```

---

## Summary

The view package provides a **powerful template rendering engine** for Go:
- Flexible structure with layouts, components, and views
- Template caching for production performance
- 20+ built-in functions for common operations
- Custom functions and global data support
- Thread-safe concurrent rendering
- Perfect for modular applications

**Key Features:**
- Render templates with or without layouts
- Reusable components
- Custom path resolution
- Hot reload for development
- Production-ready caching

Now you can easily render beautiful HTML templates in your Go applications! 
