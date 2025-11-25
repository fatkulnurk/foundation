package httprouter

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// =============== HELPER FUNCTIONS ===============

func makeRequest(t *testing.T, router HttpRouter, method, path string, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, nil)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("status code = %d, want %d", got, want)
	}
}

func assertBody(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("body = %q, want %q", got, want)
	}
}

func assertContains(t *testing.T, got, want string) {
	t.Helper()
	if !strings.Contains(got, want) {
		t.Errorf("body = %q, want to contain %q", got, want)
	}
}

// =============== BASIC ROUTING TESTS ===============

func TestRouter_GET(t *testing.T) {
	r := New()
	r.GET("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})

	w := makeRequest(t, r, "GET", "/hello", nil)
	assertStatus(t, w.Code, http.StatusOK)
	assertBody(t, w.Body.String(), "hello world")
}

func TestRouter_POST(t *testing.T) {
	r := New()
	r.POST("/users", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("user created"))
	})

	w := makeRequest(t, r, "POST", "/users", nil)
	assertStatus(t, w.Code, http.StatusOK)
	assertBody(t, w.Body.String(), "user created")
}

func TestRouter_PUT(t *testing.T) {
	r := New()
	r.PUT("/users/1", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("user updated"))
	})

	w := makeRequest(t, r, "PUT", "/users/1", nil)
	assertStatus(t, w.Code, http.StatusOK)
	assertBody(t, w.Body.String(), "user updated")
}

func TestRouter_PATCH(t *testing.T) {
	r := New()
	r.PATCH("/users/1", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("user patched"))
	})

	w := makeRequest(t, r, "PATCH", "/users/1", nil)
	assertStatus(t, w.Code, http.StatusOK)
	assertBody(t, w.Body.String(), "user patched")
}

func TestRouter_DELETE(t *testing.T) {
	r := New()
	r.DELETE("/users/1", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("user deleted"))
	})

	w := makeRequest(t, r, "DELETE", "/users/1", nil)
	assertStatus(t, w.Code, http.StatusOK)
	assertBody(t, w.Body.String(), "user deleted")
}

func TestRouter_PathParams(t *testing.T) {
	r := New()
	r.GET("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		w.Write([]byte("user id: " + id))
	})

	w := makeRequest(t, r, "GET", "/users/123", nil)
	assertStatus(t, w.Code, http.StatusOK)
	assertBody(t, w.Body.String(), "user id: 123")
}

func TestRouter_MultiplePathParams(t *testing.T) {
	r := New()
	r.GET("/users/{userId}/posts/{postId}", func(w http.ResponseWriter, r *http.Request) {
		userId := r.PathValue("userId")
		postId := r.PathValue("postId")
		w.Write([]byte("user: " + userId + ", post: " + postId))
	})

	w := makeRequest(t, r, "GET", "/users/42/posts/99", nil)
	assertStatus(t, w.Code, http.StatusOK)
	assertBody(t, w.Body.String(), "user: 42, post: 99")
}

// =============== MIDDLEWARE TESTS ===============

func TestRouter_GlobalMiddleware(t *testing.T) {
	r := New()

	// Middleware yang menambahkan header
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Custom", "global")
			next.ServeHTTP(w, r)
		})
	})

	r.GET("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	w := makeRequest(t, r, "GET", "/test", nil)
	assertStatus(t, w.Code, http.StatusOK)
	if w.Header().Get("X-Custom") != "global" {
		t.Errorf("header X-Custom = %q, want %q", w.Header().Get("X-Custom"), "global")
	}
}

func TestRouter_RouteSpecificMiddleware(t *testing.T) {
	r := New()

	authMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Authorization") == "" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}

	r.GET("/public", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("public"))
	})

	r.GET("/private", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("private"))
	}, authMiddleware)

	// Public route should work
	w := makeRequest(t, r, "GET", "/public", nil)
	assertStatus(t, w.Code, http.StatusOK)

	// Private route without auth should fail
	w = makeRequest(t, r, "GET", "/private", nil)
	assertStatus(t, w.Code, http.StatusUnauthorized)

	// Private route with auth should work
	w = makeRequest(t, r, "GET", "/private", map[string]string{"Authorization": "Bearer token"})
	assertStatus(t, w.Code, http.StatusOK)
}

func TestRouter_MiddlewareChaining(t *testing.T) {
	r := New()
	var order []string

	mw1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "mw1-before")
			next.ServeHTTP(w, r)
			order = append(order, "mw1-after")
		})
	}

	mw2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "mw2-before")
			next.ServeHTTP(w, r)
			order = append(order, "mw2-after")
		})
	}

	r.Use(mw1)
	r.Use(mw2)

	r.GET("/test", func(w http.ResponseWriter, r *http.Request) {
		order = append(order, "handler")
		w.Write([]byte("ok"))
	})

	makeRequest(t, r, "GET", "/test", nil)

	expected := []string{"mw1-before", "mw2-before", "handler", "mw2-after", "mw1-after"}
	if len(order) != len(expected) {
		t.Fatalf("middleware order length = %d, want %d", len(order), len(expected))
	}
	for i, v := range expected {
		if order[i] != v {
			t.Errorf("order[%d] = %q, want %q", i, order[i], v)
		}
	}
}

// =============== GROUP TESTS ===============

func TestRouter_Group(t *testing.T) {
	r := New()

	r.Group("/api", func(api HttpRouter) {
		api.GET("/users", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("users"))
		})
	})

	w := makeRequest(t, r, "GET", "/api/users", nil)
	assertStatus(t, w.Code, http.StatusOK)
	assertBody(t, w.Body.String(), "users")
}

func TestRouter_NestedGroup(t *testing.T) {
	r := New()

	r.Group("/api", func(api HttpRouter) {
		api.Group("/v1", func(v1 HttpRouter) {
			v1.GET("/users", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("v1 users"))
			})
		})
	})

	w := makeRequest(t, r, "GET", "/api/v1/users", nil)
	assertStatus(t, w.Code, http.StatusOK)
	assertBody(t, w.Body.String(), "v1 users")
}

func TestRouter_GroupMiddleware(t *testing.T) {
	r := New()

	r.Group("/api", func(api HttpRouter) {
		api.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-API", "true")
				next.ServeHTTP(w, r)
			})
		})

		api.GET("/users", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("users"))
		})
	})

	r.GET("/public", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("public"))
	})

	// Group route should have middleware header
	w := makeRequest(t, r, "GET", "/api/users", nil)
	if w.Header().Get("X-API") != "true" {
		t.Errorf("X-API header not set on group route")
	}

	// Non-group route should not have middleware header
	w = makeRequest(t, r, "GET", "/public", nil)
	if w.Header().Get("X-API") != "" {
		t.Errorf("X-API header should not be set on non-group route")
	}
}

func TestRouter_GroupInheritsParentMiddleware(t *testing.T) {
	r := New()
	var headers []string

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			headers = append(headers, "global")
			next.ServeHTTP(w, r)
		})
	})

	r.Group("/api", func(api HttpRouter) {
		api.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				headers = append(headers, "api")
				next.ServeHTTP(w, r)
			})
		})

		api.GET("/users", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("users"))
		})
	})

	makeRequest(t, r, "GET", "/api/users", nil)

	expected := []string{"global", "api"}
	if len(headers) != len(expected) {
		t.Fatalf("headers length = %d, want %d", len(headers), len(expected))
	}
	for i, v := range expected {
		if headers[i] != v {
			t.Errorf("headers[%d] = %q, want %q", i, headers[i], v)
		}
	}
}

func TestRouter_NestedGroupInheritsMiddleware(t *testing.T) {
	r := New()
	var order []string

	r.Group("/api", func(api HttpRouter) {
		api.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				order = append(order, "api")
				next.ServeHTTP(w, r)
			})
		})

		api.Group("/v1", func(v1 HttpRouter) {
			v1.Use(func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					order = append(order, "v1")
					next.ServeHTTP(w, r)
				})
			})

			v1.GET("/users", func(w http.ResponseWriter, r *http.Request) {
				order = append(order, "handler")
				w.Write([]byte("users"))
			})
		})
	})

	makeRequest(t, r, "GET", "/api/v1/users", nil)

	expected := []string{"api", "v1", "handler"}
	if len(order) != len(expected) {
		t.Fatalf("order length = %d, want %d", len(order), len(expected))
	}
	for i, v := range expected {
		if order[i] != v {
			t.Errorf("order[%d] = %q, want %q", i, order[i], v)
		}
	}
}

// =============== STATIC FILE TESTS ===============

func TestRouter_Static(t *testing.T) {
	r := New()

	// Create temp directory with test file
	tmpDir := t.TempDir()
	testFile := tmpDir + "/test.txt"
	content := "test content"
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	r.Static("/static", tmpDir)

	w := makeRequest(t, r, "GET", "/static/test.txt", nil)
	assertStatus(t, w.Code, http.StatusOK)
	assertBody(t, w.Body.String(), content)
}

func TestRouter_Static_NotFound(t *testing.T) {
	r := New()
	tmpDir := t.TempDir()
	r.Static("/static", tmpDir)

	w := makeRequest(t, r, "GET", "/static/notfound.txt", nil)
	assertStatus(t, w.Code, http.StatusNotFound)
}

func TestRouter_Static_MethodNotAllowed(t *testing.T) {
	r := New()
	tmpDir := t.TempDir()
	r.Static("/static", tmpDir)

	w := makeRequest(t, r, "POST", "/static/test.txt", nil)
	assertStatus(t, w.Code, http.StatusMethodNotAllowed)
}

func TestRouter_Static_InGroup(t *testing.T) {
	r := New()
	tmpDir := t.TempDir()
	testFile := tmpDir + "/app.css"
	content := "body { color: red; }"
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	r.Group("/api", func(api HttpRouter) {
		api.Static("/assets", tmpDir)
	})

	w := makeRequest(t, r, "GET", "/api/assets/app.css", nil)
	assertStatus(t, w.Code, http.StatusOK)
	assertBody(t, w.Body.String(), content)
}

func TestRouter_Static_WithMiddleware(t *testing.T) {
	r := New()
	tmpDir := t.TempDir()
	testFile := tmpDir + "/test.txt"
	if err := os.WriteFile(testFile, []byte("content"), 0644); err != nil {
		t.Fatal(err)
	}

	authMw := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("X-Auth") == "" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}

	r.Static("/static", tmpDir, authMw)

	// Without auth
	w := makeRequest(t, r, "GET", "/static/test.txt", nil)
	assertStatus(t, w.Code, http.StatusUnauthorized)

	// With auth
	w = makeRequest(t, r, "GET", "/static/test.txt", map[string]string{"X-Auth": "token"})
	assertStatus(t, w.Code, http.StatusOK)
}

// =============== HANDLE/HANDLEFUNC TESTS ===============

func TestRouter_Handle(t *testing.T) {
	r := New()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("custom handler"))
	})

	r.Handle("GET /custom", handler)

	w := makeRequest(t, r, "GET", "/custom", nil)
	assertStatus(t, w.Code, http.StatusOK)
	assertBody(t, w.Body.String(), "custom handler")
}

func TestRouter_HandleFunc(t *testing.T) {
	r := New()
	r.HandleFunc("POST /data", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("data posted"))
	})

	w := makeRequest(t, r, "POST", "/data", nil)
	assertStatus(t, w.Code, http.StatusOK)
	assertBody(t, w.Body.String(), "data posted")
}

// =============== PATH UTILITY TESTS ===============

func TestClean(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", "/"},
		{"/", "/"},
		{"  /  ", "/"},
		{"/hello", "/hello"},
		{"hello", "/hello"},
		{"/hello/", "/hello"},
		{"/hello/world", "/hello/world"},
		{"/hello/world/", "/hello/world"},
		{"  /hello  ", "/hello"},
	}

	for _, tt := range tests {
		got := clean(tt.input)
		if got != tt.want {
			t.Errorf("clean(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestJoin(t *testing.T) {
	tests := []struct {
		prefix string
		path   string
		want   string
	}{
		{"/", "/", "/"},
		{"/api", "/users", "/api/users"},
		{"/api/", "/users", "/api/users"},
		{"/api", "users", "/api/users"},
		{"api", "users", "/api/users"},
		{"/api/v1", "/users", "/api/v1/users"},
		{"/", "/users", "/users"},
		{"/api", "/", "/api"},
	}

	for _, tt := range tests {
		got := join(tt.prefix, tt.path)
		if got != tt.want {
			t.Errorf("join(%q, %q) = %q, want %q", tt.prefix, tt.path, got, tt.want)
		}
	}
}

// =============== EDGE CASES ===============

func TestRouter_404(t *testing.T) {
	r := New()
	r.GET("/exists", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	w := makeRequest(t, r, "GET", "/notfound", nil)
	assertStatus(t, w.Code, http.StatusNotFound)
}

func TestRouter_MethodNotAllowed(t *testing.T) {
	r := New()
	r.GET("/resource", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	w := makeRequest(t, r, "POST", "/resource", nil)
	assertStatus(t, w.Code, http.StatusMethodNotAllowed) // Go 1.22+ ServeMux returns 405 for method mismatch
}

func TestRouter_EmptyPath(t *testing.T) {
	r := New()
	r.GET("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("root"))
	})

	w := makeRequest(t, r, "GET", "/", nil)
	assertStatus(t, w.Code, http.StatusOK)
	assertBody(t, w.Body.String(), "root")
}
