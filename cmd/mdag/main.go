// Command mdag implements a content-addressed file system CLI.
//
// Usage:
//
//	mdag add <local-path>          Import a file/directory, print root CID
//	mdag resolve <root> <path>     Print target CID and type
//	mdag cat <root> <path>         Print file content to stdout
//	mdag ls <root> <path>          List directory children
package main

import (
	"fmt"
	"mdag/importer"
	"mdag/resolver"
	"mdag/store"
	"os"
	"text/tabwriter"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	// Determine store type.
	// Default: file store under ./data/.
	// Override with --mem flag for in-memory (useful in tests).
	var s store.Store
	var err error
	useMem := false
	for i, a := range args {
		if a == "--mem" {
			useMem = true
			args = append(args[:i], args[i+1:]...)
			break
		}
	}
	if useMem {
		s = store.NewMemoryStore()
	} else {
		s, err = store.NewFileStore("./data")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}

	switch cmd {
	case "add":
		runAdd(args, s)
	case "resolve":
		runResolve(args, s)
	case "cat":
		runCat(args, s)
	case "ls":
		runLs(args, s)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", cmd)
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Println(`Usage:
  mdag add <local-path>          Import a file/directory, print root CID
  mdag resolve <root> <path>     Resolve path, print target CID and type
  mdag cat <root> <path>         Print file content
  mdag ls <root> <path>          List directory children

Flags:
  --mem     Use in-memory store instead of file store`)
}

func runAdd(args []string, s store.Store) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Usage: mdag add <local-path>")
		os.Exit(1)
	}
	rootCID, err := importer.AddPath(args[0], s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Root CID: %s\n", rootCID)
}

func runResolve(args []string, s store.Store) {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: mdag resolve <root-cid> <path>")
		os.Exit(1)
	}
	rootCID := args[0]
	path := args[1]
	res, err := resolver.Resolve(rootCID, path, s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Target CID: %s\n", res.CID)
	fmt.Printf("Type: %s\n", res.Type)
}

func runCat(args []string, s store.Store) {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: mdag cat <root-cid> <path>")
		os.Exit(1)
	}
	data, err := resolver.ReadFile(args[0], args[1], s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	os.Stdout.Write(data)
}

func runLs(args []string, s store.Store) {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: mdag ls <root-cid> <path>")
		os.Exit(1)
	}
	links, err := resolver.ListDir(args[0], args[1], s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	if len(links) == 0 {
		fmt.Println("(empty directory)")
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	for _, link := range links {
		typ := ""
		if obj, err := s.GetObject(link.CID); err == nil {
			typ = string(obj.Type)
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%d\n", link.Name, typ, link.CID, link.Size)
	}
	w.Flush()
}
