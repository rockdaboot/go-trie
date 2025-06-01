// Package main demonstrates the usage of the trie package to efficiently store,
// manipulate, and reconstruct stack traces using a trie data structure.
package main

import (
	"fmt"

	"github.com/rockdaboot/go-trie/trie"
)

func main() {
	// stacks contains the stack traces with leaf frame at position 0.
	// That's how we store them in ebpf-profiler (natural order).
	//
	// After encoding, we have
	//
	//	locationTable: [{} {main} {foo} {bar} {baz1} {baz2} {what} {why}]
	//	stackParentArray: [0 0 1 2 2 3 3 2 0 8 9]
	//	stackLocationIndex: [0 1 2 3 4 4 5 5 6 7 7]
	//	stackIndex: [3 3 4 5 6 7 10]
	var stacks = [][]string{
		{"bar", "foo", "main"},
		{"bar", "foo", "main"},
		{"baz1", "foo", "main"},
		{"baz1", "bar", "foo", "main"},
		{"baz2", "bar", "foo", "main"},
		{"baz2", "foo", "main"},
		{"why", "why", "what"},
	}

	printStacks(trie.NewFromStacks(stacks).ToArrays())
	fmt.Println()

	type frame struct {
		fileID  uint64
		address uint64
	}
	//nolint:mnd
	var stacksOfFrames = [][]frame{
		{{fileID: 1, address: 0x100}, {fileID: 2, address: 0x200}, {fileID: 3, address: 0x300}},
		{{fileID: 4, address: 0x400}, {fileID: 2, address: 0x200}, {fileID: 3, address: 0x300}},
	}
	printStacks(trie.NewFromStacks(stacksOfFrames).ToArrays())
}

func printStacks[T comparable](locationTable []T,
	stackParentArray, stackLocationIndex, stackIndex []int) {
	fmt.Println("locationTable:", locationTable)
	fmt.Println("stackParentArray:", stackParentArray)
	fmt.Println("stackLocationIndex:", stackLocationIndex)
	fmt.Println("stackIndex:", stackIndex)

	for _, i := range stackIndex {
		for i != 0 {
			fmt.Printf(" %v", locationTable[stackLocationIndex[i]])
			i = stackParentArray[i]
		}
		fmt.Println()
	}
}
