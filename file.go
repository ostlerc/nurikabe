package main

import (
	"io/ioutil"
	"sort"
	"strconv"
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

type numericint struct {
	strs []string
}

// Len is part of sort.Interface.
func (n *numericint) Len() int {
	return len(n.strs)
}

// Swap is part of sort.Interface.
func (n *numericint) Swap(i, j int) {
	n.strs[i], n.strs[j] = n.strs[j], n.strs[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (n *numericint) Less(i, j int) bool {
	iint, err := strconv.Atoi(n.strs[i][:len(n.strs[i])-5])
	if err != nil {
		panic(err)
	}
	jint, err := strconv.Atoi(n.strs[j][:len(n.strs[j])-5])
	if err != nil {
		panic(err)
	}
	return iint < jint
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

	n := &numericint{strs: names}
	sort.Sort(n)
	return n.strs
}
