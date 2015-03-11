package filedupes

import (
	"io"
	"os"
)

const dupeTreeBufferSize = 256

func Dupes(files []*os.File) ([][]*os.File, error) {
	trees := make(map[int64]*dupeTree)

	for _, file := range files {
		stat, err := file.Stat()
		if err != nil {
			return nil, err
		}

		filelen := stat.Size()

		var tree *dupeTree

		if t, ok := trees[filelen]; ok {
			tree = t
		} else {
			trees[filelen] = newDupeTree()
			tree = trees[filelen]
		}

		tree.insert(file)
	}

	dupes := [][]*os.File{}

	for _, tree := range trees {
		for _, dupeSet := range tree.dupes() {
			if len(dupeSet) != 0 {
				dupes = append(dupes, dupeSet)
			}
		}
	}

	return dupes, nil
}

type dupeTree struct {
	children   map[[dupeTreeBufferSize]byte]*dupeTree
	duplicates []*os.File
	leaf       bool
	deferred   *os.File
}

func newDupeTree() *dupeTree {
	return &dupeTree{make(map[[dupeTreeBufferSize]byte]*dupeTree), []*os.File{}, true, nil}
}

func (d *dupeTree) dupes() [][]*os.File {
	dupes := [][]*os.File{}

	for _, node := range d.nodes() {
		if len(node.duplicates) != 0 {
			dupes = append(dupes, node.duplicates)
		}
	}

	return dupes
}

func (d *dupeTree) nodes() []*dupeTree {
	nodes := []*dupeTree{d}

	for _, child := range d.children {
		if len(child.nodes()) != 0 {
			nodes = append(nodes, child.nodes()...)
		}
	}

	return nodes
}

func (d *dupeTree) insert(file *os.File) {
	if d.leaf && d.deferred != nil {
		// Push the deferred node, then push the file
		deferred := d.deferred

		d.deferred = nil
		d.leaf = false

		d.push(deferred)
		d.push(file)
	} else if d.leaf && d.deferred == nil {
		// Defer the file
		d.deferred = file
	} else {
		// Push the file
		d.push(file)
	}
}

// Preconditions: d.leaf, d.deferred != nil
func (d *dupeTree) push(file *os.File) error {
	var chunk [dupeTreeBufferSize]byte

	_, err := file.Read(chunk[:])

	if err == io.EOF {
		d.duplicates = append(d.duplicates, file)
		return nil
	} else if err != nil {
		return err
	}

	// More to follow
	if child, ok := d.children[chunk]; ok {
		child.insert(file)
	} else {
		d.children[chunk] = newDupeTree()
		d.children[chunk].insert(file)
	}

	return nil
}
