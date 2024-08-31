package caching

import (
	"bytes"
	"html/template"
	"os"
	"testing"
	"zehd/pkg"
)

func TestCachePages(t *testing.T) {
	// Setup
	err := os.Setenv("LOGLOCATION", "/tmp/")
	if err != nil {
		t.Fatalf("Failed to set LOGLOCATION environment variable: %v", err)
	}

	testDir := "/tmp/test_templates_dir"
	err = os.MkdirAll(testDir+"/html", os.ModePerm)
	if err != nil {
		t.Fatalf("Failed to create test templates directory: %v", err)
	}
	defer os.RemoveAll(testDir)

	pkg.TemplatesDir = testDir + "/"
	pkg.TemplateType = "html"

	layoutHTML := `{{ define "layout" }}<!DOCTYPE html><html><body>{{template "content"}}</body></html>{{ end }}`
	err = os.WriteFile(testDir+"/layout.html", []byte(layoutHTML), 0644)
	if err != nil {
		t.Fatalf("Failed to create layout.html file: %v", err)
	}

	templateHTML := `{{define "content"}}<h1>Test Template</h1>{{end}}`
	err = os.WriteFile(testDir+"/html/test.html", []byte(templateHTML), 0644)
	if err != nil {
		t.Fatalf("Failed to create test.html file: %v", err)
	}

	pages := &Pages{
		RouteMap: make(map[string]*template.Template),
	}

	// Test execution
	err = pages.CachePages()
	// Assertions
	if err != nil {
		t.Errorf("CachePages returned an error: %v", err)
	}

	// Print out all cached templates
	t.Logf("Cached templates:")
	for name := range pages.RouteMap {
		t.Logf("- %s", name)
	}

	got := len(pages.RouteMap)
	want := 1 // We expect one template to be cached
	if got != want {
		t.Errorf("CachePages cached %d templates, want %d", got, want)
	}

	if _, exists := pages.RouteMap["test"]; !exists {
		t.Errorf("Expected 'test' template to be cached, but it wasn't")
	}

	// Check if there are any unexpected templates
	for name := range pages.RouteMap {
		if name != "test" {
			t.Errorf("Unexpected template cached: %s", name)
		}
	}
}

func TestTemplateBuilder(t *testing.T) {
	// Setup
	testDir := "/tmp/test_templates_dir"
	err := os.MkdirAll(testDir+"/html", os.ModePerm)
	if err != nil {
		t.Fatalf("Failed to create test templates directory: %v", err)
	}
	defer os.RemoveAll(testDir)

	pkg.TemplatesDir = testDir + "/"
	pkg.TemplateType = "html"

	layoutHTML := `{{ define "layout" }}<!DOCTYPE html><html><body>{{template "content"}}</body></html>{{ end }}`
	err = os.WriteFile(testDir+"/layout.html", []byte(layoutHTML), 0644)
	if err != nil {
		t.Fatalf("Failed to create layout.html file: %v", err)
	}

	templateHTML := `{{define "content"}}<h1>Test Template</h1>{{end}}`
	err = os.WriteFile(testDir+"/html/test.html", []byte(templateHTML), 0644)
	if err != nil {
		t.Fatalf("Failed to create test.html file: %v", err)
	}

	// Test execution
	got, err := templateBuilder("test.html", ".html")
	// Assertions
	if err != nil {
		t.Errorf("templateBuilder returned an error: %v", err)
	}

	if got == nil {
		t.Errorf("templateBuilder returned nil template")
	}

	// Test template execution
	want := "<!DOCTYPE html><html><body><h1>Test Template</h1></body></html>"
	var buf bytes.Buffer
	err = got.ExecuteTemplate(&buf, "layout", nil)
	if err != nil {
		t.Errorf("Failed to execute template: %v", err)
	}

	result := buf.String()
	if result != want {
		t.Errorf("Template execution result = %q, want %q", result, want)
	}
}

// =================================================================================================================
// =========================================  Benchmarks  ==========================================================
// =================================================================================================================

func BenchmarkCachePages(b *testing.B) {
	// Setup
	err := os.Setenv("LOGLOCATION", "/tmp/")
	if err != nil {
		b.Fatalf("Failed to set LOGLOCATION environment variable: %v", err)
	}

	testDir := "/tmp/test_templates_dir"
	err = os.MkdirAll(testDir+"/html", os.ModePerm)
	if err != nil {
		b.Fatalf("Failed to create test templates directory: %v", err)
	}
	defer os.RemoveAll(testDir)

	pkg.TemplatesDir = testDir + "/"
	pkg.TemplateType = "html"

	layoutHTML := `{{ define "layout" }}<!DOCTYPE html><html><body>{{template "content"}}</body></html>{{ end }}`
	err = os.WriteFile(testDir+"/layout.html", []byte(layoutHTML), 0644)
	if err != nil {
		b.Fatalf("Failed to create layout.html file: %v", err)
	}

	templateHTML := `{{define "content"}}<h1>Test Template</h1>{{end}}`
	err = os.WriteFile(testDir+"/html/test.html", []byte(templateHTML), 0644)
	if err != nil {
		b.Fatalf("Failed to create test.html file: %v", err)
	}

	pages := &Pages{
		RouteMap: make(map[string]*template.Template),
	}

	// Benchmark execution
	for i := 0; i < b.N; i++ {
		err := pages.CachePages()
		if err != nil {
			b.Errorf("CachePages returned an error: %v", err)
		}
	}
}

func BenchmarkTemplateBuilder(b *testing.B) {
	// Setup
	testDir := "/tmp/test_templates_dir"
	err := os.MkdirAll(testDir+"/html", os.ModePerm)
	if err != nil {
		b.Fatalf("Failed to create test templates directory: %v", err)
	}
	defer os.RemoveAll(testDir)

	pkg.TemplatesDir = testDir + "/"
	pkg.TemplateType = "html"

	layoutHTML := `{{ define "layout" }}<!DOCTYPE html><html><body>{{template "content"}}</body></html>{{ end }}`
	err = os.WriteFile(testDir+"/layout.html", []byte(layoutHTML), 0644)
	if err != nil {
		b.Fatalf("Failed to create layout.html file: %v", err)
	}

	templateHTML := `{{define "content"}}<h1>Test Template</h1>{{end}}`
	err = os.WriteFile(testDir+"/html/test.html", []byte(templateHTML), 0644)
	if err != nil {
		b.Fatalf("Failed to create test.html file: %v", err)
	}

	// Benchmark execution
	for i := 0; i < b.N; i++ {
		got, err := templateBuilder("test.html", ".html")
		if err != nil {
			b.Errorf("templateBuilder returned an error: %v", err)
		}

		if got == nil {
			b.Errorf("templateBuilder returned nil template")
		}

		var buf bytes.Buffer
		err = got.ExecuteTemplate(&buf, "layout", nil)
		if err != nil {
			b.Errorf("Failed to execute template: %v", err)
		}
	}
}
