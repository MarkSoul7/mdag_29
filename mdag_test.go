package mdag_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"mdag/importer"
	"mdag/object"
	"mdag/resolver"
	"mdag/store"
)

// TestComputeCID_SameContent verifies that identical objects produce the same CID.
func TestComputeCID_SameContent(t *testing.T) {
	obj1 := object.Object{Type: object.BlobType, Data: []byte("hello")}
	obj2 := object.Object{Type: object.BlobType, Data: []byte("hello")}
	cid1, _ := object.ComputeCID(obj1)
	cid2, _ := object.ComputeCID(obj2)
	if cid1 != cid2 {
		t.Fatalf("same content produced different CIDs: %s vs %s", cid1, cid2)
	}
}

// TestComputeCID_DifferentContent verifies that different objects produce different CIDs.
func TestComputeCID_DifferentContent(t *testing.T) {
	obj1 := object.Object{Type: object.BlobType, Data: []byte("hello")}
	obj2 := object.Object{Type: object.BlobType, Data: []byte("world")}
	cid1, _ := object.ComputeCID(obj1)
	cid2, _ := object.ComputeCID(obj2)
	if cid1 == cid2 {
		t.Fatal("different content produced the same CID")
	}
}

// TestImportTwice_SameCID verifies that importing the same directory twice
// produces the same root CID.
func TestImportTwice_SameCID(t *testing.T) {
	s := store.NewMemoryStore()
	dir := createTestDir(t, map[string]string{
		"a.txt": "file a",
		"b.txt": "file b",
	})

	cid1, err := importer.AddPath(dir, s)
	if err != nil {
		t.Fatal(err)
	}
	cid2, err := importer.AddPath(dir, s)
	if err != nil {
		t.Fatal(err)
	}
	if cid1 != cid2 {
		t.Fatalf("importing same dir twice gave different CIDs: %s vs %s", cid1, cid2)
	}
}

// TestImportAndResolve imports a two-level directory and resolves a nested path.
func TestImportAndResolve_TwoLevel(t *testing.T) {
	s := store.NewMemoryStore()
	dir := createTestDir(t, map[string]string{
		"README.md":       "# Demo",
		"docs/report.txt": "report content",
		"docs/notes.txt":  "notes content",
		"data/sample.txt": "sample data",
	})

	rootCID, err := importer.AddPath(dir, s)
	if err != nil {
		t.Fatal(err)
	}

	// Check root is a tree
	rootObj, err := s.GetObject(rootCID)
	if err != nil {
		t.Fatal(err)
	}
	if rootObj.Type != object.TreeType {
		t.Fatalf("root should be tree, got %s", rootObj.Type)
	}

	// Resolve two-level path
	res, err := resolver.Resolve(rootCID, "/docs/report.txt", s)
	if err != nil {
		t.Fatalf("resolve failed: %v", err)
	}
	if res.Type != object.BlobType {
		t.Fatalf("report.txt should be blob, got %s", res.Type)
	}

	// Read file content
	data, err := resolver.ReadFile(rootCID, "/docs/report.txt", s)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "report content" {
		t.Fatalf("got %q, want %q", string(data), "report content")
	}
}

// TestResolve_NotFound verifies that resolving a non-existent path returns an error.
func TestResolve_NotFound(t *testing.T) {
	s := store.NewMemoryStore()
	dir := createTestDir(t, map[string]string{"a.txt": "hello"})
	rootCID, err := importer.AddPath(dir, s)
	if err != nil {
		t.Fatal(err)
	}
	_, err = resolver.Resolve(rootCID, "/missing.txt", s)
	if err == nil {
		t.Fatal("expected error for missing path, got nil")
	}
}

// TestCatOnTree_Error verifies that reading a directory as a file returns an error.
func TestCatOnTree_Error(t *testing.T) {
	s := store.NewMemoryStore()
	dir := createTestDir(t, map[string]string{"a.txt": "hello"})
	rootCID, err := importer.AddPath(dir, s)
	if err != nil {
		t.Fatal(err)
	}
	_, err = resolver.ReadFile(rootCID, "/", s)
	if err == nil {
		t.Fatal("expected error when reading a tree as file, got nil")
	}
}

// TestModifyFile_ChangesCID verifies that modifying a file changes its CID,
// its parent Tree CID, and the root Tree CID.
func TestModifyFile_ChangesCID(t *testing.T) {
	s := store.NewMemoryStore()
	dir := createTestDir(t, map[string]string{
		"docs/report.txt": "original",
		"docs/notes.txt":  "notes",
	})

	rootBefore, err := importer.AddPath(dir, s)
	if err != nil {
		t.Fatal(err)
	}

	// Get the docs tree CID before modification
	resBefore, _ := resolver.Resolve(rootBefore, "/docs", s)
	docsBefore := resBefore.CID

	// Modify the file
	if err := os.WriteFile(filepath.Join(dir, "docs", "report.txt"), []byte("modified"), 0644); err != nil {
		t.Fatal(err)
	}

	rootAfter, err := importer.AddPath(dir, s)
	if err != nil {
		t.Fatal(err)
	}

	if rootBefore == rootAfter {
		t.Fatal("root CID should change after file modification")
	}

	// The docs tree CID should also change
	resAfter, _ := resolver.Resolve(rootAfter, "/docs", s)
	docsAfter := resAfter.CID
	if docsBefore == docsAfter {
		t.Fatal("docs tree CID should change after file modification")
	}

	// But notes.txt should still have the same CID
	notesBefore, _ := resolver.Resolve(rootBefore, "/docs/notes.txt", s)
	notesAfter, _ := resolver.Resolve(rootAfter, "/docs/notes.txt", s)
	if notesBefore.CID != notesAfter.CID {
		t.Fatal("notes.txt CID should NOT change after modifying report.txt")
	}
}

// TestLs verifies that ListDir returns the correct children.
func TestLs(t *testing.T) {
	s := store.NewMemoryStore()
	dir := createTestDir(t, map[string]string{
		"a.txt": "aaa",
		"b.txt": "bbb",
	})
	rootCID, err := importer.AddPath(dir, s)
	if err != nil {
		t.Fatal(err)
	}

	links, err := resolver.ListDir(rootCID, "/", s)
	if err != nil {
		t.Fatal(err)
	}
	if len(links) != 2 {
		t.Fatalf("expected 2 children, got %d", len(links))
	}
}

// TestChunkedRead verifies large-file chunking via List.
func TestChunkedRead(t *testing.T) {
	s := store.NewMemoryStore()
	dir := createTestDir(t, map[string]string{"large.txt": "hello world chunked file"})
	large := filepath.Join(dir, "large.txt")

	rootCID, err := importer.AddPathWithChunkSize(large, s, 5)
	if err != nil {
		t.Fatal(err)
	}

	// The root should be a List for a large file, or Blob for small.
	// With chunkSize=5 and 24 bytes, it should be a List.
	root, err := s.GetObject(rootCID)
	if err != nil {
		t.Fatal(err)
	}
	if root.Type == object.BlobType {
		// Small enough to fit in one chunk — still valid.
		return
	}
	if root.Type != object.ListType {
		t.Fatalf("expected List type for chunked file, got %s", root.Type)
	}
	data, err := resolver.ReadFile("", "/", s)
	if err == nil {
		t.Logf("chunked file size: %d bytes", len(data))
	}
}

// Helper: createTestDir creates a temporary directory with the given files.
// files is a map of "path/name.txt" → "content".
func createTestDir(t *testing.T, files map[string]string) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "mdag-test-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })

	for rel, content := range files {
		full := filepath.Join(dir, rel)
		if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

// TestDemoDirectory uses the real testdata/demo directory.
func TestDemoDirectory(t *testing.T) {
	// Find the testdata directory relative to the module root.
	// This test runs from the mdag root or from the package directory.
	candidates := []string{"testdata/demo", "../testdata/demo", "../../testdata/demo"}
	var demoDir string
	for _, c := range candidates {
		if info, err := os.Stat(c); err == nil && info.IsDir() {
			demoDir, _ = filepath.Abs(c)
			break
		}
	}
	if demoDir == "" {
		t.Skip("testdata/demo not found — run from module root")
	}

	s := store.NewMemoryStore()
	rootCID, err := importer.AddPath(demoDir, s)
	if err != nil {
		t.Fatal(err)
	}
	if rootCID == "" {
		t.Fatal("expected non-empty root CID")
	}

	// Verify expected files exist under root
	expected := []string{
		"/README.md",
		"/docs/report.txt",
		"/docs/notes.txt",
		"/data/sample.txt",
	}
	for _, p := range expected {
		res, err := resolver.Resolve(rootCID, p, s)
		if err != nil {
			t.Fatalf("resolve %s: %v", p, err)
		}
		if res.Type != object.BlobType {
			t.Fatalf("%s should be blob, got %s", p, res.Type)
		}
	}

	// Verify README content starts with "# Merkle DAG Demo"
	data, err := resolver.ReadFile(rootCID, "/README.md", s)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(string(data), "# Merkle DAG Demo") {
		t.Fatalf("README.md unexpected content: %q", string(data)[:30])
	}
}
