package main

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/fatkulnurk/foundation/view"
)

func main() {
	fmt.Println("=== View Package Examples ===\n")

	// Example 1: Gold Price Page (dengan layout)
	goldPriceExample()

	// Example 2: User Profile Page
	userProfileExample()

	// Example 3: Custom path resolver
	customPathResolverExample()

	// Example 4: Production mode with caching
	productionExample()

	fmt.Println("\n✅ All examples completed!")
	fmt.Println("\nFile structure:")
	fmt.Println("view/")
	fmt.Println("├── layouts/app.html")
	fmt.Println("├── component/")
	fmt.Println("│   ├── header.html")
	fmt.Println("│   ├── footer.html")
	fmt.Println("│   ├── table.html")
	fmt.Println("│   └── alert.html")
	fmt.Println("module/")
	fmt.Println("├── gold/view/price.html")
	fmt.Println("└── user/view/profile.html")
}

func goldPriceExample() {
	fmt.Println("=== 1. Gold Price Page Example ===")

	// Setup view dengan global data
	v := view.New(view.Config{
		LayoutsPath:    "./view/layouts",
		ComponentsPath: "./view/component",
		ViewsPath:      "./module",
		EnableCache:    false,
		GlobalData: map[string]any{
			"SiteName": "Harga Emas Indonesia",
			"Version":  "1.0.0",
		},
	})

	// Data untuk halaman harga emas
	data := map[string]any{
		"Title": "Harga Emas Hari Ini",
		"Prices": []map[string]any{
			{"Weight": "1 gram", "BuyPrice": "Rp 1.050.000", "SellPrice": "Rp 1.000.000"},
			{"Weight": "5 gram", "BuyPrice": "Rp 5.250.000", "SellPrice": "Rp 5.000.000"},
			{"Weight": "10 gram", "BuyPrice": "Rp 10.500.000", "SellPrice": "Rp 10.000.000"},
			{"Weight": "25 gram", "BuyPrice": "Rp 26.250.000", "SellPrice": "Rp 25.000.000"},
			{"Weight": "50 gram", "BuyPrice": "Rp 52.500.000", "SellPrice": "Rp 50.000.000"},
			{"Weight": "100 gram", "BuyPrice": "Rp 105.000.000", "SellPrice": "Rp 100.000.000"},
		},
		"HighestPrice": "Rp 1.100.000",
		"LowestPrice":  "Rp 950.000",
		"AveragePrice": "Rp 1.025.000",
	}

	// Render dengan layout
	html, err := v.RenderWithLayout(
		context.Background(),
		"app",             // layout: ./view/layouts/app.html
		"gold/view/price", // view: ./module/gold/view/price.html
		data,
	)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	fmt.Printf("✓ Rendered gold price page (%d bytes)\n", len(html))
	fmt.Println("  Layout: ./view/layouts/app.html")
	fmt.Println("  Components: header.html, footer.html, table.html")
	fmt.Println("  View: ./module/gold/view/price.html")
	fmt.Println()
}

func userProfileExample() {
	fmt.Println("=== 2. User Profile Page Example ===")

	v := view.New(view.Config{
		LayoutsPath:    "./view/layouts",
		ComponentsPath: "./view/component",
		ViewsPath:      "./module",
		EnableCache:    false,
		GlobalData: map[string]any{
			"SiteName": "Harga Emas Indonesia",
			"Version":  "1.0.0",
		},
	})

	// Data untuk halaman profile
	data := map[string]any{
		"Title":    "Profil Pengguna",
		"Name":     "Fatkul Nur Koirudin",
		"Email":    "fatkul@example.com",
		"Phone":    "+62 812-3456-7890",
		"Status":   "Active",
		"JoinDate": time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
	}

	html, err := v.RenderWithLayout(
		context.Background(),
		"app",               // layout: ./view/layouts/app.html
		"user/view/profile", // view: ./module/user/view/profile.html
		data,
	)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	fmt.Printf("✓ Rendered user profile page (%d bytes)\n", len(html))
	fmt.Println("  Layout: ./view/layouts/app.html")
	fmt.Println("  Components: header.html, footer.html")
	fmt.Println("  View: ./module/user/view/profile.html")
	fmt.Println()
}

func customPathResolverExample() {
	fmt.Println("=== 3. Custom Path Resolver Example ===")

	// Example custom resolver function
	customResolver := func(templateType, name string) string {
		switch templateType {
		case "layout":
			return filepath.Join("./view/layouts", name+".html")
		case "component":
			return filepath.Join("./view/component", name+".html")
		case "view":
			// Custom logic: bisa disesuaikan dengan struktur project
			// Contoh: "gold/price" bisa di-resolve ke path apapun
			return filepath.Join("./module", name+".html")
		default:
			return name + ".html"
		}
	}

	v := view.New(view.Config{
		LayoutsPath:    "./view/layouts",
		ComponentsPath: "./view/component",
		ViewsPath:      "./module",
		EnableCache:    false,
		GlobalData: map[string]any{
			"SiteName": "Harga Emas Indonesia",
			"Version":  "1.0.0",
		},
		PathResolver: customResolver,
	})

	fmt.Println("✓ Custom path resolver configured")
	fmt.Println("  You can implement any path resolution logic")
	fmt.Println("  Example: Different paths for different modules")
	fmt.Println("  Example: Dynamic path based on environment")
	fmt.Println("  Example: Multi-tenant path resolution")

	// Demonstrate that it still works
	data := map[string]any{
		"Title": "Harga Emas",
		"Prices": []map[string]any{
			{"Weight": "1 gram", "BuyPrice": "Rp 1.050.000", "SellPrice": "Rp 1.000.000"},
		},
		"HighestPrice": "Rp 1.100.000",
		"LowestPrice":  "Rp 950.000",
		"AveragePrice": "Rp 1.025.000",
	}

	html, err := v.RenderWithLayout(context.Background(), "app", "gold/view/price", data)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	fmt.Printf("✓ Rendered successfully (%d bytes)\n", len(html))
	fmt.Println()
}

func productionExample() {
	fmt.Println("=== 4. Production Mode (with caching) ===")

	v := view.New(view.Config{
		LayoutsPath:    "./view/layouts",
		ComponentsPath: "./view/component",
		ViewsPath:      "./module",
		Extension:      ".html",
		EnableCache:    true, // Enable cache untuk production
		GlobalData: map[string]any{
			"SiteName": "Harga Emas Indonesia",
			"Version":  "1.0.0",
		},
	})

	data := map[string]any{
		"Title": "Harga Emas",
		"Prices": []map[string]any{
			{"Weight": "1 gram", "BuyPrice": "Rp 1.050.000", "SellPrice": "Rp 1.000.000"},
		},
		"HighestPrice": "Rp 1.100.000",
		"LowestPrice":  "Rp 950.000",
		"AveragePrice": "Rp 1.025.000",
	}

	// First render - akan parse template
	html1, _ := v.RenderWithLayout(context.Background(), "app", "gold/view/price", data)
	fmt.Printf("✓ First render: %d bytes (parsed from files)\n", len(html1))

	// Second render - akan menggunakan cache
	html2, _ := v.RenderWithLayout(context.Background(), "app", "gold/view/price", data)
	fmt.Printf("✓ Second render: %d bytes (from cache - faster!)\n", len(html2))

	// Clear cache jika perlu (misalnya saat hot reload)
	v.ClearCache()
	fmt.Println("✓ Cache cleared")
	fmt.Println()
}
