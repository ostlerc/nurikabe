package main

import (
	"io/ioutil"
	"sort"
)

func dirs(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	l := len(files)
	names := make([]string, 0, l)
	for _, f := range files {
		if f.IsDir() {
			names = append(names, f.Name())
		}
	}
	sort.Strings(names)
	return names
}

func files(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	l := len(files)
	names := make([]string, l, l)
	for i, f := range files {
		names[i] = f.Name()
	}

	sort.Strings(names)
	return names
}
