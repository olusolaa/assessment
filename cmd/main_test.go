package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"assessment/domain/model"
)

func TestPrintProducts(t *testing.T) {

	date, _ := time.Parse("2006-01-02", "2020-01-01")
	products := model.ProductList{
		{
			ID:         1,
			Name:       "Test Product",
			Price:      10.99,
			Created:    date,
			SalesCount: 100,
			ViewsCount: 1000,
		},
	}

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	printProducts(products)

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, err := io.Copy(&buf, r)
	if err != nil {
		t.Fatalf("Failed to copy from pipe: %v", err)
	}
	output := buf.String()

	expected := "ID: 1, Name: Test Product, Price: $10.99, Created: 2020-01-01, Sales/View: 0.100000"
	if !strings.Contains(output, expected) {
		t.Errorf("printProducts output mismatch: got %s, want %s", output, expected)
	}
}

func TestMainFunction(t *testing.T) {

	if testing.Short() {
		t.Skip("Skipping test in short mode")
	}

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	main()

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, err := io.Copy(&buf, r)
	if err != nil {
		t.Fatalf("Failed to copy from pipe: %v", err)
	}
	output := buf.String()

	expectedStrings := []string{

		"Using default configuration",
		"Available sorters:",
		"Products sorted by price (ascending):",
		"Products sorted by sales per view (descending):",
		"Products sorted by creation date (ascending):",
		"Products sorted by name (ascending):",
		"Demonstrating pagination",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("main output missing expected string: %s", expected)
		}
	}
}
