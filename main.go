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

	trieData := trie.NewFromStacks(stacks)

	// Convert the trie to arrays.
	locationTable, stackParentArray, stackLocationIndex, stackIndex := trieData.ToArrays()

	// Print the arrays
	fmt.Println("locationTable:", locationTable)
	fmt.Println("stackParentArray:", stackParentArray)
	fmt.Println("stackLocationIndex:", stackLocationIndex)
	fmt.Println("stackIndex:", stackIndex)

	// Print the stacks, demonstrates how to construct the stacks from the arrays.
	for _, i := range stackIndex {
		for i != 0 {
			fmt.Printf(" %s", locationTable[stackLocationIndex[i]])
			i = stackParentArray[i]
		}
		fmt.Println()
	}
}
