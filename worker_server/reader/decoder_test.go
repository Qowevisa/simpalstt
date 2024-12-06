package reader_test

import (
	"log"
	"os"
	"testing"
	"worker_server/reader"
)

func TestStartJsonDecoder_ValidFile(t *testing.T) {
	// Create a temp file
	tmpfile, err := os.CreateTemp("", "test*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	jsonContent := `[{"_id":"1234","categories":{"subcategory":"0099"},"title":{"ro":"LoremIpsum","ru":"LoremIpsum2"},"type":"NaTaNaTa","posted":3.14}]`
	if _, err := tmpfile.WriteString(jsonContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	tmpfile.Close()

	dataChannel, err := reader.StartJsonDecoder(tmpfile.Name(), reader.WithSuppressedLogging())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	r := <-dataChannel
	if r.XId != "1234" ||
		r.Categories.Subcategory != "0099" ||
		r.Title.Ro != "LoremIpsum" ||
		r.Title.Ru != "LoremIpsum2" ||
		r.Type != "NaTaNaTa" ||
		r.Posted != 3.14 {
		t.Errorf("Unexpected data received: %+v", r)
	}
}

func TestStartJsonDecoder_ValidFileWithThreeData(t *testing.T) {
	// Create a temp file
	tmpfile, err := os.CreateTemp("", "test*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	jsonContent := `[` +
		`{"_id":"1234","categories":{"subcategory":"0099"},"title":{"ro":"LoremIpsum","ru":"LoremIpsum2"},"type":"NaTaNaTa","posted":3.14},` +
		`{"_id":"1234","categories":{"subcategory":"0099"},"title":{"ro":"LoremIpsum","ru":"LoremIpsum2"},"type":"NaTaNaTa","posted":3.14},` +
		`{"_id":"1234","categories":{"subcategory":"0099"},"title":{"ro":"LoremIpsum","ru":"LoremIpsum2"},"type":"NaTaNaTa","posted":3.14}` +
		`]`
	if _, err := tmpfile.WriteString(jsonContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	tmpfile.Close()

	dataChannel, err := reader.StartJsonDecoder(tmpfile.Name(), reader.WithSuppressedLogging())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	for i := 0; i < 3; i++ {
		r := <-dataChannel
		if r == nil {
			log.Printf("R is nil\n")
			break
		}
		if r.XId != "1234" ||
			r.Categories.Subcategory != "0099" ||
			r.Title.Ro != "LoremIpsum" ||
			r.Title.Ru != "LoremIpsum2" ||
			r.Type != "NaTaNaTa" ||
			r.Posted != 3.14 {
			t.Errorf("Unexpected data received: %+v", r)
		}
	}
}

func TestStartJsonDecoder_InvalidFile(t *testing.T) {
	_, err := reader.StartJsonDecoder("nonexistent.json")
	if err == nil {
		t.Fatal("Expected error for nonexistent file, got nil")
	}
}
