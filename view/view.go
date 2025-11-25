package view

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io/fs"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// View interface untuk rendering template
type View interface {
	Render(ctx context.Context, name string, data any) (string, error)
	RenderWithLayout(ctx context.Context, layout, name string, data any) (string, error)
	AddFunc(name string, fn any) View
	SetGlobal(key string, value any) View
	ClearCache()
}

// Config untuk konfigurasi view
type Config struct {
	// LayoutsPath adalah full path ke directory layouts (e.g., "./view/layouts")
	LayoutsPath string

	// ComponentsPath adalah full path ke directory components (e.g., "./view/component")
	ComponentsPath string

	// ViewsPath adalah full path ke directory views (e.g., "./module")
	ViewsPath string

	// Extension adalah file extension untuk template (default: ".html")
	Extension string

	// EnableCache untuk enable/disable template caching (default: true untuk production)
	EnableCache bool

	// Delimiters untuk custom template delimiters (optional)
	LeftDelim  string
	RightDelim string

	// FuncMap untuk custom template functions
	FuncMap template.FuncMap

	// GlobalData untuk data yang tersedia di semua template
	GlobalData map[string]any

	// PathResolver adalah custom function untuk resolve path template
	// Jika nil, akan menggunakan default resolver
	// Signature: func(templateType, name string) string
	// templateType: "layout", "component", "view"
	PathResolver func(templateType, name string) string
}

type view struct {
	config   Config
	cache    map[string]*template.Template
	cacheMu  sync.RWMutex
	funcMap  template.FuncMap
	globalMu sync.RWMutex
}

// New membuat instance View baru
func New(config Config) View {
	// Set defaults
	if config.Extension == "" {
		config.Extension = ".html"
	}
	if config.GlobalData == nil {
		config.GlobalData = make(map[string]any)
	}

	v := &view{
		config:  config,
		cache:   make(map[string]*template.Template),
		funcMap: make(template.FuncMap),
	}

	// Register default functions
	v.registerDefaultFuncs()

	// Register custom functions dari config
	if config.FuncMap != nil {
		for name, fn := range config.FuncMap {
			v.funcMap[name] = fn
		}
	}

	return v
}

// registerDefaultFuncs mendaftarkan fungsi-fungsi default
func (v *view) registerDefaultFuncs() {
	v.funcMap["raw"] = func(s string) template.HTML {
		return template.HTML(s)
	}

	v.funcMap["safeHTML"] = func(s string) template.HTML {
		return template.HTML(template.HTMLEscapeString(s))
	}

	v.funcMap["safeJS"] = func(s string) template.JS {
		return template.JS(template.JSEscapeString(s))
	}

	v.funcMap["safeURL"] = func(s string) template.URL {
		return template.URL(s)
	}

	v.funcMap["now"] = func() time.Time {
		return time.Now()
	}

	v.funcMap["formatDate"] = func(t time.Time, layout string) string {
		return t.Format(layout)
	}

	v.funcMap["year"] = func() int {
		return time.Now().Year()
	}

	v.funcMap["add"] = func(a, b int) int {
		return a + b
	}

	v.funcMap["sub"] = func(a, b int) int {
		return a - b
	}

	v.funcMap["mul"] = func(a, b int) int {
		return a * b
	}

	v.funcMap["div"] = func(a, b int) int {
		if b == 0 {
			return 0
		}
		return a / b
	}

	v.funcMap["join"] = func(sep string, items []string) string {
		return strings.Join(items, sep)
	}

	v.funcMap["upper"] = func(s string) string {
		return strings.ToUpper(s)
	}

	v.funcMap["lower"] = func(s string) string {
		return strings.ToLower(s)
	}

	v.funcMap["title"] = func(s string) string {
		// Simple title case: capitalize first letter of each word
		words := strings.Fields(s)
		for i, word := range words {
			if len(word) > 0 {
				words[i] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
			}
		}
		return strings.Join(words, " ")
	}

	v.funcMap["default"] = func(defaultVal, val any) any {
		if val == nil || val == "" {
			return defaultVal
		}
		return val
	}

	// Global data accessor
	v.funcMap["global"] = func(key ...string) any {
		v.globalMu.RLock()
		defer v.globalMu.RUnlock()

		if len(key) == 0 {
			return v.config.GlobalData
		}

		if val, ok := v.config.GlobalData[key[0]]; ok {
			return val
		}
		return nil
	}
}

// AddFunc menambahkan custom function ke template
func (v *view) AddFunc(name string, fn any) View {
	v.funcMap[name] = fn
	v.ClearCache() // Clear cache karena funcMap berubah
	return v
}

// SetGlobal menambahkan/update global data
func (v *view) SetGlobal(key string, value any) View {
	v.globalMu.Lock()
	defer v.globalMu.Unlock()
	v.config.GlobalData[key] = value
	return v
}

// ClearCache membersihkan template cache
func (v *view) ClearCache() {
	v.cacheMu.Lock()
	defer v.cacheMu.Unlock()
	v.cache = make(map[string]*template.Template)
}

// Render me-render template dengan layout default
func (v *view) Render(ctx context.Context, name string, data any) (string, error) {
	return v.RenderWithLayout(ctx, "", name, data)
}

// RenderWithLayout me-render template dengan layout spesifik
func (v *view) RenderWithLayout(ctx context.Context, layout, name string, data any) (string, error) {
	cacheKey := v.getCacheKey(layout, name)

	// Check cache
	if v.config.EnableCache {
		v.cacheMu.RLock()
		if tmpl, ok := v.cache[cacheKey]; ok {
			v.cacheMu.RUnlock()
			return v.executeTemplate(tmpl, name, data)
		}
		v.cacheMu.RUnlock()
	}

	// Parse template
	tmpl, err := v.parseTemplate(layout, name)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	// Cache template
	if v.config.EnableCache {
		v.cacheMu.Lock()
		v.cache[cacheKey] = tmpl
		v.cacheMu.Unlock()
	}

	return v.executeTemplate(tmpl, name, data)
}

// parseTemplate mem-parse template files
func (v *view) parseTemplate(layout, name string) (*template.Template, error) {
	tmpl := template.New("")

	// Set custom delimiters jika ada
	if v.config.LeftDelim != "" && v.config.RightDelim != "" {
		tmpl = tmpl.Delims(v.config.LeftDelim, v.config.RightDelim)
	}

	// Add functions
	tmpl = tmpl.Funcs(v.funcMap)

	// Collect all template files
	files := []string{}

	// 1. Add layout if specified
	if layout != "" {
		layoutPath := v.resolvePath("layout", layout)
		files = append(files, layoutPath)
	}

	// 2. Add all components (optional, bisa di-skip jika tidak ada)
	if v.config.ComponentsPath != "" {
		componentsPattern := filepath.Join(v.config.ComponentsPath, "*"+v.config.Extension)
		componentFiles, _ := filepath.Glob(componentsPattern)
		files = append(files, componentFiles...)
	}

	// 3. Add the main view
	viewPath := v.resolvePath("view", name)
	files = append(files, viewPath)

	// Parse all files
	if len(files) == 0 {
		return nil, fmt.Errorf("no template files found for: %s", name)
	}

	tmpl, err := tmpl.ParseFiles(files...)
	if err != nil {
		return nil, fmt.Errorf("failed to parse files: %w", err)
	}

	return tmpl, nil
}

// resolvePath me-resolve path template
// templateType: "layout", "component", "view"
func (v *view) resolvePath(templateType, name string) string {
	// Gunakan custom PathResolver jika ada
	if v.config.PathResolver != nil {
		return v.config.PathResolver(templateType, name)
	}

	// Default resolver
	var basePath string
	switch templateType {
	case "layout":
		basePath = v.config.LayoutsPath
	case "component":
		basePath = v.config.ComponentsPath
	case "view":
		basePath = v.config.ViewsPath
	default:
		basePath = v.config.ViewsPath
	}

	// Jika name sudah include extension, gunakan langsung
	if strings.HasSuffix(name, v.config.Extension) {
		return filepath.Join(basePath, name)
	}

	// Jika name mengandung subdirectory (e.g., "gold/price" atau "modulename/view/price")
	if strings.Contains(name, "/") {
		return filepath.Join(basePath, name+v.config.Extension)
	}

	// Default: name + extension
	return filepath.Join(basePath, name+v.config.Extension)
}

// executeTemplate mengeksekusi template
func (v *view) executeTemplate(tmpl *template.Template, name string, data any) (string, error) {
	// Extract template name dari path
	templateName := v.extractTemplateName(name)

	var buf bytes.Buffer
	err := tmpl.ExecuteTemplate(&buf, templateName, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", templateName, err)
	}

	return buf.String(), nil
}

// extractTemplateName mengekstrak nama template dari path
func (v *view) extractTemplateName(name string) string {
	// Remove extension
	name = strings.TrimSuffix(name, v.config.Extension)

	// Get base name (last part of path)
	parts := strings.Split(name, "/")
	return parts[len(parts)-1]
}

// getCacheKey membuat cache key
func (v *view) getCacheKey(layout, name string) string {
	if layout == "" {
		return name
	}
	return layout + ":" + name
}

// WalkTemplates walks through all template files (useful for preloading)
func (v *view) WalkTemplates(callback func(path string) error) error {
	// Walk through all configured paths
	paths := []string{}

	if v.config.LayoutsPath != "" {
		paths = append(paths, v.config.LayoutsPath)
	}
	if v.config.ComponentsPath != "" {
		paths = append(paths, v.config.ComponentsPath)
	}
	if v.config.ViewsPath != "" {
		paths = append(paths, v.config.ViewsPath)
	}

	for _, basePath := range paths {
		err := filepath.WalkDir(basePath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !d.IsDir() && strings.HasSuffix(path, v.config.Extension) {
				return callback(path)
			}

			return nil
		})

		if err != nil {
			return err
		}
	}

	return nil
}
