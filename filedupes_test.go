package filedupes

import (
	"math/rand"
	"os"
	"sort"
	"testing"
)

const (
	file1 = "filedupes.go"
	file2 = "README.md"
	file3 = "filedupes_test.go"
)

func TestDupes(t *testing.T) {
	files := []*os.File{}

	for _, filename := range []string{file1, file2, file3} {
		for range make([]struct{}, 5) {
			file, err := os.Open(filename)

			if err != nil {
				t.Fatal(err)
			}

			files = append(files, file)
		}
	}

	for i := range files {
		j := rand.Intn(i + 1)
		files[i], files[j] = files[j], files[i]
	}

	duplicates, err := Dupes(files)

	if err != nil {
		t.Fatal(err)
	}

	if len(duplicates) != 3 {
		t.Fatal("Wrong number of duplicate groups:", len(duplicates))
	}

	all_names, err := names(duplicates)

	if err != nil {
		t.Fatal(err)
	}

	actual := []string{}

	for _, names := range all_names {
		name := names[0]
		actual = append(actual, name)
		if !(name == file1 || name == file2 || name == file3) {
			t.Fatal("File with name that wasn't supplied found:", name)
		}
		for _, single := range names {
			if name != single {
				t.Fatal("files should be sorted but arent:", name, single)
			}
		}
	}

	expected := []string{file1, file2, file3}
	sort.Strings(expected)
	sort.Strings(actual)

	if expected[0] != actual[0] ||
		expected[1] != actual[1] ||
		expected[2] != actual[2] {
		t.Fatal("Filenames were not equal to expected. Expected:",
			expected, "actual:", actual)
	}

	for _, file := range files {
		file.Close()
	}
}

func names(files [][]*os.File) ([][]string, error) {
	strings := [][]string{}

	for _, duplicates := range files {
		dups := []string{}
		for _, file := range duplicates {
			stats, err := file.Stat()
			if err != nil {
				return nil, err
			}

			dups = append(dups, stats.Name())
		}
		strings = append(strings, dups)
	}

	return strings, nil
}
