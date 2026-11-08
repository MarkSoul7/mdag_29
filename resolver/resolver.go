// Package resolver provides path resolution and file reading over a Merkle
// DAG Store.  Resolve walks from a root Tree through child Links by name;
// ReadFile returns the full byte content of a Blob at a given path.
package resolver

import (
	"fmt"
	"mdag/object"
	"mdag/store"
	"strings"
)

// Result holds the outcome of a successful Resolve call.
type Result struct {
	CID  string            // target object CID
	Type object.ObjectType // target object type
}

// Resolve walks the Merkle DAG starting at rootCID, following path segments
// through Tree Links, and returns the CID and type of the terminal object.
//
// An empty or "/"-only path returns rootCID itself (useful for ls on root).
func Resolve(rootCID, path string, s store.Store) (Result, error) {
	// Normalise: strip leading slash, split into segments.
	clean := strings.Trim(path, "/")
	if clean == "" {
		root, err := s.GetObject(rootCID)
		if err != nil {
			return Result{}, err
		}
		return Result{CID: rootCID, Type: root.Type}, nil
	}
	parts := strings.Split(clean, "/")

	currentCID := rootCID
	for i, name := range parts {
		obj, err := s.GetObject(currentCID)
		if err != nil {
			return Result{}, fmt.Errorf("resolve %s: %w", currentCID, err)
		}
		if obj.Type != object.TreeType {
			return Result{}, fmt.Errorf("%q is a %s, not a tree: cannot descend into %q",
				currentCID, obj.Type, name)
		}
		found := false
		for _, link := range obj.Links {
			if link.Name == name {
				currentCID = link.CID
				found = true
				break
			}
		}
		if !found {
			return Result{}, fmt.Errorf("path %q: name %q not found in tree %s",
				path, name, currentCID)
		}
		// If this was the last segment, return the target's type.
		if i == len(parts)-1 {
			target, err := s.GetObject(currentCID)
			if err != nil {
				return Result{}, err
			}
			return Result{CID: currentCID, Type: target.Type}, nil
		}
	}

	// Should not reach here for non-empty path.
	return Result{}, fmt.Errorf("resolve %q: incomplete resolution", path)
}

// ReadFile resolves path from rootCID and returns the full content bytes of
// the target Blob.  Returns an error if the target is not a Blob.
func ReadFile(rootCID, path string, s store.Store) ([]byte, error) {
	res, err := Resolve(rootCID, path, s)
	if err != nil {
		return nil, err
	}
	obj, err := s.GetObject(res.CID)
	if err != nil {
		return nil, err
	}
	switch obj.Type {
	case object.BlobType:
		return obj.Data, nil
	case object.ListType:
		// Optional: recursively concatenate list children.
		var buf []byte
		for _, link := range obj.Links {
			child, err := s.GetObject(link.CID)
			if err != nil {
				return nil, fmt.Errorf("read list chunk %s: %w", link.CID, err)
			}
			if child.Type != object.BlobType {
				return nil, fmt.Errorf("list child %s is %s, expected blob", link.CID, child.Type)
			}
			buf = append(buf, child.Data...)
		}
		return buf, nil
	default:
		return nil, fmt.Errorf("%q (%s) is a %s, not a file", path, res.CID, obj.Type)
	}
}

// ListDir resolves path from rootCID and returns the direct child links of
// the target Tree.  Returns an error if the target is not a Tree.
func ListDir(rootCID, path string, s store.Store) ([]object.Link, error) {
	res, err := Resolve(rootCID, path, s)
	if err != nil {
		return nil, err
	}
	obj, err := s.GetObject(res.CID)
	if err != nil {
		return nil, err
	}
	if obj.Type != object.TreeType {
		return nil, fmt.Errorf("%q (%s) is a %s, not a directory", path, res.CID, obj.Type)
	}
	return obj.Links, nil
}
