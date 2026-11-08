// Package importer walks local files and directories and encodes them as
// Merkle DAG objects (Blob for files, Tree for directories).
package importer

import (
	"fmt"
	"mdag/object"
	"mdag/store"
	"os"
	"path/filepath"
	"sort"
)

// AddPath imports the file or directory at localPath into the given Store
// and returns the root CID of the resulting Merkle DAG.
//
//   - A regular file is encoded as a Blob.
//   - A directory is encoded as a Tree whose Links reference the imported
//     children.  Directories are imported recursively.
func AddPath(localPath string, s store.Store) (string, error) {
	info, err := os.Stat(localPath)
	if err != nil {
		return "", fmt.Errorf("stat %s: %w", localPath, err)
	}

	if info.IsDir() {
		return importDir(localPath, s)
	}
	return importFile(localPath, s)
}

// importFile reads a regular file and stores it as a Blob.
func importFile(path string, s store.Store) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read %s: %w", path, err)
	}
	obj := object.Object{
		Type: object.BlobType,
		Data: data,
	}
	return s.PutObject(obj)
}

// importDir walks a directory and stores it as a Tree.
func importDir(path string, s store.Store) (string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return "", fmt.Errorf("read dir %s: %w", path, err)
	}

	// Sort entries by name to ensure deterministic ordering.
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	var links []object.Link
	for _, entry := range entries {
		childPath := filepath.Join(path, entry.Name())
		childCID, err := AddPath(childPath, s)
		if err != nil {
			return "", fmt.Errorf("import %s: %w", childPath, err)
		}
		// For Blob children, record the file size.
		childObj, err := s.GetObject(childCID)
		if err != nil {
			return "", fmt.Errorf("get %s: %w", childCID, err)
		}
		sz := object.Size(childObj)
		links = append(links, object.Link{
			Name: entry.Name(),
			CID:  childCID,
			Size: sz,
		})
	}

	obj := object.Object{
		Type:  object.TreeType,
		Links: links,
	}
	return s.PutObject(obj)
}

// AddPathWithChunkSize imports a file, splitting it into chunks of at most
// chunkSize bytes, and ordering the chunks under a List root.
//
// Small files (≤ chunkSize) are stored as a single Blob.  Larger files are
// split into chunkSize Blobs and the root is a List referencing those chunks.
//
// This is an optional ("List 与文件分块") enhancement.
func AddPathWithChunkSize(localPath string, s store.Store, chunkSize int) (string, error) {
	info, err := os.Stat(localPath)
	if err != nil {
		return "", fmt.Errorf("stat %s: %w", localPath, err)
	}
	if info.IsDir() {
		// Directories still use Tree encoding.
		return importDir(localPath, s)
	}

	data, err := os.ReadFile(localPath)
	if err != nil {
		return "", fmt.Errorf("read %s: %w", localPath, err)
	}

	if len(data) <= chunkSize {
		return s.PutObject(object.Object{Type: object.BlobType, Data: data})
	}

	// Split into chunks.
	var links []object.Link
	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}
		chunkCID, err := s.PutObject(object.Object{
			Type: object.BlobType,
			Data: data[i:end],
		})
		if err != nil {
			return "", err
		}
		links = append(links, object.Link{
			CID:  chunkCID,
			Size: int64(end - i),
		})
	}
	return s.PutObject(object.Object{Type: object.ListType, Links: links})
}
