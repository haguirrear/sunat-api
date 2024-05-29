package sunat

import (
	"archive/zip"
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestCreateSingleFileZip(t *testing.T) {
	filePath := "/home/folder/example.txt"
	fileContents := `<xml/>`

	fileReader := strings.NewReader(fileContents)
	t.Log("Creating zip")
	s := Sunat{}
	zFile, err := s.createSingleFileZip(filePath, fileReader)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Zip created")

	zReader, err := zip.NewReader(bytes.NewReader(zFile.Bytes()), int64(zFile.Len()))

	if len(zReader.File) != 1 {
		t.Fatalf("expected 1 file, got %d", len(zReader.File))
	}

	zf := zReader.File[0]

	if zf.Name != "example.txt" {
		t.Fatalf("expected 'example.txt', got %s", zf.Name)
	}

	rc, err := zf.Open()
	if err != nil {
		t.Fatal(err)
	}

	defer rc.Close()
	c, err := io.ReadAll(rc)
	if err != nil {
		t.Fatal(err)
	}

	if string(c) != fileContents {
		t.Fatalf("Expecting filecontents '%s', got '%s'", fileContents, string(c))
	}
}
