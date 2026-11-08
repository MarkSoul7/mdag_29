// Package store defines the object storage interface and implementations.
package store

import "mdag/object"

// Store is the persistent key-value interface for Merkle DAG objects.
//
//   - PutObject serialises obj, computes its CID, stores it, and returns the CID.
//   - GetObject retrieves a previously stored object by CID.
//
// PutObject is idempotent: storing the same object twice succeeds and
// returns the same CID both times.
type Store interface {
	PutObject(obj object.Object) (string, error)
	GetObject(cid string) (object.Object, error)
}
