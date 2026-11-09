package store

import (
	"fmt"
	"mdag/object"
)

// MemoryStore holds all objects in an in-memory map.
// CID → object; the map makes Put idempotent (same CID always maps to the same object).
type MemoryStore struct {
	objects map[string]object.Object
}

// NewMemoryStore creates an empty in-memory store.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{objects: make(map[string]object.Object)}
}

// PutObject stores obj and returns its CID.
// Storing the same object twice is idempotent.
func (s *MemoryStore) PutObject(obj object.Object) (string, error) {
	cid, err := object.ComputeCID(obj)
	if err != nil {
		return "", err
	}
	s.objects[cid] = obj
	return cid, nil
}

// GetObject retrieves the object identified by cid.
// Returns an error if the CID is not found.
func (s *MemoryStore) GetObject(cid string) (object.Object, error) {
	obj, ok := s.objects[cid]
	if !ok {
		return object.Object{}, fmt.Errorf("object %s: not found", cid)
	}
	return obj, nil
}
