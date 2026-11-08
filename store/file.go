package store

import (
	"encoding/json"
	"fmt"
	"mdag/object"
	"os"
	"path/filepath"
)

// FileStore persists objects as individual JSON files under a data directory.
// File name = <cid>.json, e.g. ./data/abcd1234...json
type FileStore struct {
	dataDir string
}

// NewFileStore creates a FileStore rooted at dataDir.
// The directory is created if it does not exist.
func NewFileStore(dataDir string) (*FileStore, error) {
	abs, err := filepath.Abs(dataDir)
	if err != nil {
		return nil, fmt.Errorf("resolve data dir: %w", err)
	}
	if err := os.MkdirAll(abs, 0755); err != nil {
		return nil, fmt.Errorf("create data dir: %w", err)
	}
	return &FileStore{dataDir: abs}, nil
}

// PutObject serialises obj to JSON, writes to <cid>.json, and returns the CID.
func (s *FileStore) PutObject(obj object.Object) (string, error) {
	cid, err := object.ComputeCID(obj)
	if err != nil {
		return "", err
	}
	path := filepath.Join(s.dataDir, cid+".json")

	// If the file already exists, the content must match — no need to rewrite.
	if _, err := os.Stat(path); err == nil {
		return cid, nil
	}

	data, err := json.Marshal(obj)
	if err != nil {
		return "", fmt.Errorf("marshal object: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return "", fmt.Errorf("write %s: %w", path, err)
	}
	return cid, nil
}

// GetObject reads and deserialises the object stored at <cid>.json.
func (s *FileStore) GetObject(cid string) (object.Object, error) {
	path := filepath.Join(s.dataDir, cid+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return object.Object{}, fmt.Errorf("object %s: not found", cid)
		}
		return object.Object{}, fmt.Errorf("read %s: %w", path, err)
	}
	var obj object.Object
	if err := json.Unmarshal(data, &obj); err != nil {
		return object.Object{}, fmt.Errorf("unmarshal %s: %w", path, err)
	}
	return obj, nil
}
