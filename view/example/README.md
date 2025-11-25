# View Package - Complete Example

Contoh lengkap penggunaan view package dengan struktur template yang terpisah.

## 📁 Struktur Direktori

```
example/
├── main.go                          # Example code
├── view/
│   ├── layouts/
│   │   └── app.html                # Layout utama
│   └── component/
│       ├── header.html             # Component header
│       ├── footer.html             # Component footer
│       ├── table.html              # Component table
│       └── alert.html              # Component alert
└── module/
    ├── gold/
    │   └── view/
    │       └── price.html          # Halaman harga emas
    └── user/
        └── view/
            └── profile.html        # Halaman profile user
```

## 🚀 Cara Menjalankan

```bash
cd pkg/view/example
go run main.go
```

## 📝 Contoh yang Tersedia

### 1. Gold Price Page
Menampilkan halaman harga emas dengan:
- Layout lengkap (header, footer)
- Component table untuk menampilkan harga
- Global data (SiteName, Version)
- Template functions (formatDate, year, global)

### 2. User Profile Page
Menampilkan halaman profil user dengan:
- Layout yang sama
- Data user (nama, email, phone, status)
- Custom styling dengan inline CSS
- Template functions (formatDate, default, upper)

### 3. Custom Path Resolver
Demonstrasi penggunaan custom path resolver untuk:
- Flexible path resolution
- Dynamic path berdasarkan kondisi
- Multi-tenant support
- Environment-based paths

### 4. Production Mode
Demonstrasi caching untuk production:
- First render: parse dari file
- Second render: dari cache (lebih cepat)
- Clear cache untuk hot reload

## 🎨 Template Features

### Layout (app.html)
- Responsive design
- Modern gradient header
- Clean footer
- Card-based content area

### Components
- **header.html**: Navigation bar dengan menu
- **footer.html**: Footer dengan copyright dan version
- **table.html**: Reusable table untuk harga emas
- **alert.html**: Alert component dengan custom styling

### Views
- **gold/view/price.html**: Halaman harga emas dengan statistik
- **user/view/profile.html**: Halaman profil user dengan form

## 💡 Key Concepts

### 1. Separate Paths
```go
v := view.New(view.Config{
    LayoutsPath:    "./view/layouts",
    ComponentsPath: "./view/component",
    ViewsPath:      "./module",
})
```

### 2. Global Data
```go
GlobalData: map[string]any{
    "SiteName": "Harga Emas Indonesia",
    "Version":  "1.0.0",
}
```

### 3. Template Rendering
```go
html, err := v.RenderWithLayout(
    ctx,
    "app",              // layout
    "gold/view/price",  // view
    data,               // data
)
```

### 4. Built-in Functions
- `{{ global "SiteName" }}` - Access global data
- `{{ formatDate now "02 Jan 2006" }}` - Format date
- `{{ year }}` - Current year
- `{{ default "value" .Field }}` - Default value
- `{{ upper .Text }}` - Uppercase
- `{{ raw .HTML }}` - Render HTML

## 📊 Output

Running the example will produce:
- Gold price page: ~4KB HTML
- User profile page: ~3KB HTML
- Custom resolver demo: ~2.5KB HTML
- Production caching demo

## 🔧 Customization

Anda bisa memodifikasi:
1. **Layouts**: Edit `view/layouts/app.html`
2. **Components**: Tambah/edit di `view/component/`
3. **Views**: Buat module baru di `module/yourmodule/view/`
4. **Styling**: Update inline CSS atau tambahkan external CSS
5. **Data**: Ubah data di `main.go`

## 📚 Learn More

Lihat dokumentasi lengkap di `../README.md`
