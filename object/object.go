// Package object defines the core Merkle DAG data model.
//
// It implements the course-required Blob, Tree, and optional List
// object types with JSON serialization and SHA-256 content addressing.
package object

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// ObjectType identifies the kind of a Merkle DAG node.
type ObjectType string

const (
	BlobType ObjectType = "blob"
	TreeType ObjectType = "tree"
	ListType ObjectType = "list"
)

// Link is a named reference from a Tree (or List) to a child node.
type Link struct {
	Name string `json:"name,omitempty"`
	CID  string `json:"cid"`
	Size int64  `json:"size,omitempty"`
}

// Object is the unified Merkle DAG node.
//
//   - Blob: carries file data in the Data field.
//   - Tree: carries directory entries in the Links field.
//   - List: carries ordered chunk references in the Links field.
type Object struct {
	Type  ObjectType `json:"type"`
	Data  []byte     `json:"data,omitempty"`
	Links []Link     `json:"links,omitempty"`
}

// ComputeCID serialises obj to JSON and returns its SHA-256 hex digest.
func ComputeCID(obj Object) (string, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return "", fmt.Errorf("serialise object: %w", err)
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}

// MustComputeCID is like ComputeCID but panics on error.
func MustComputeCID(obj Object) string {
	cid, err := ComputeCID(obj)
	if err != nil {
		panic(err)
	}
	return cid
}

// Size returns the logical size of an object:
//   - Blob: len(Data)
//   - Tree: 0 (directories have no intrinsic size)
//   - List: sum of child Link sizes (when known)
func Size(obj Object) int64 {
	switch obj.Type {
	case BlobType:
		return int64(len(obj.Data))
	case ListType:
		var total int64
		for _, l := range obj.Links {
			if l.Size > 0 {
				total += l.Size
			}
		}
		return total
	default:
		return 0
	}
}
