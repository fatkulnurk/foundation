package main

import (
	"context"
	"fmt"
	"html/template"
	"log"

	"github.com/fatkulnurk/foundation/view"
)

func main() {
	// Example 1: Basic usage
	basicExample()

	// Example 2: With layout
	layoutExample()

	// Example 3: With global data
	globalDataExample()

	// Example 4: With custom functions
	customFuncExample()

	// Example 5: Production mode with caching
	productionExample()
}

func basicExample() {
	fmt.Println("=== Basic Example ===")

	v := view.New(view.Config{
		TemplateDir: "./templates",
		EnableCache: false, // Disable cache untuk development
	})

	// Render simple view
	html, err := v.Render(context.Background(), "home", map[string]any{
		"Title":   "Welcome",
		"Message": "Hello, World!",
	})
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	fmt.Println(html)
	fmt.Println()
}

func layoutExample() {
	fmt.Println("=== Layout Example ===")

	v := view.New(view.Config{
		TemplateDir: "./templates",
		EnableCache: false,
	})

	// Render dengan layout
	html, err := v.RenderWithLayout(
		context.Background(),
		"main",                    // layout name
		"gold/gold-bullion-price", // view name (support nested path)
		map[string]any{
			"Title": "Gold Price",
			"Prices": []map[string]any{
				{"Weight": "1g", "Price": 1000000},
				{"Weight": "5g", "Price": 5000000},
				{"Weight": "10g", "Price": 10000000},
			},
		},
	)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	fmt.Println(html)
	fmt.Println()
}

func globalDataExample() {
	fmt.Println("=== Global Data Example ===")

	v := view.New(view.Config{
		TemplateDir: "./templates",
		EnableCache: false,
		GlobalData: map[string]any{
			"SiteName": "Harga Emas",
			"Version":  "1.0.0",
		},
	})

	// Tambah global data secara dinamis
	v.SetGlobal("Year", 2025)
	v.SetGlobal("Author", "Your Name")

	html, err := v.Render(context.Background(), "about", nil)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	fmt.Println(html)
	fmt.Println()
}

func customFuncExample() {
	fmt.Println("=== Custom Function Example ===")

	v := view.New(view.Config{
		TemplateDir: "./templates",
		EnableCache: false,
		FuncMap: template.FuncMap{
			"formatCurrency": func(amount int) string {
				return fmt.Sprintf("Rp %d", amount)
			},
		},
	})

	// Tambah custom function secara dinamis
	v.AddFunc("greet", func(name string) string {
		return "Hello, " + name + "!"
	})

	html, err := v.Render(context.Background(), "product", map[string]any{
		"ProductName": "Gold Bar",
		"Price":       1000000,
		"UserName":    "John",
	})
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	fmt.Println(html)
	fmt.Println()
}

func productionExample() {
	fmt.Println("=== Production Example (with caching) ===")

	v := view.New(view.Config{
		TemplateDir:   "./templates",
		LayoutsDir:    "layouts",
		ComponentsDir: "components",
		ViewsDir:      "views",
		Extension:     ".html",
		EnableCache:   true, // Enable cache untuk production
		GlobalData: map[string]any{
			"SiteName": "Harga Emas",
			"Year":     2025,
		},
	})

	// First render - akan parse template
	html1, _ := v.Render(context.Background(), "home", map[string]any{
		"Title": "First Render",
	})
	fmt.Println("First render:", len(html1), "bytes")

	// Second render - akan menggunakan cache
	html2, _ := v.Render(context.Background(), "home", map[string]any{
		"Title": "Second Render (from cache)",
	})
	fmt.Println("Second render:", len(html2), "bytes")

	// Clear cache jika perlu
	v.ClearCache()
	fmt.Println("Cache cleared")

	fmt.Println()
}
