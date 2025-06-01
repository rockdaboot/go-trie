// Package trie provides functionality to efficiently store and manipulate stack traces
// using a trie data structure.
// It supports operations such as adding stack traces, converting the trie to arrays,
// and reconstructing stacks from arrays.
package trie

import (
	"fmt"
	"strings"
)

// Trie is a data structure that stores unique stack traces and their relationships.
type Trie[T comparable] struct {
	uniqueStacks    map[string]stackData
	uniqueLocations map[T]int
	locationTable   []T
	stackIndex      []int // stackIndex keeps track of the added stacks.
	rootID          string
}

// New creates a new Trie instance. It initializes the trie with an artificial root frame.
func New[T comparable]() *Trie[T] {
	var root T
	trie := Trie[T]{
		uniqueStacks:    make(map[string]stackData),
		uniqueLocations: make(map[T]int),
		stackIndex:      make([]int, 0),
		rootID:          mkStackID([]T{root}),
	}

	// Add an artificial root frame.
	trie.locationTable = append(trie.locationTable, root)
	trie.uniqueLocations[root] = 0
	trie.uniqueStacks[trie.rootID] = stackData{parentStackID: "", parentArrayIdx: 0, locationIdx: 0}

	return &trie
}

// NewFromStacks creates a new Trie instance from a slice of stack traces.
func NewFromStacks[T comparable](stacks [][]T) *Trie[T] {
	trie := New[T]()
	for _, stack := range stacks {
		// Add the stack trace to the trie.
		// The leaf frame is at position 0, so we add the stack in natural order.
		trie.AddStack(stack)
	}
	return trie
}

// Len returns the number of added unique stacks in the trie plus one for the artificial root frame.
func (t *Trie[T]) Len() int {
	return len(t.uniqueStacks)
}

// Index returns the index of the stack in the trie or -1 if the stack is not found.
func (t *Trie[T]) Index(stack []T) int {
	if stackItem, ok := t.uniqueStacks[mkStackID(stack)]; ok {
		return stackItem.parentArrayIdx
	}
	return -1
}

// Exists checks if a stack trace exists in the trie.
func (t *Trie[T]) Exists(stack []T) bool {
	_, ok := t.uniqueStacks[mkStackID(stack)]
	return ok
}

// AddStack adds a stack trace to the trie. The stack trace is expected to be in
// natural order, meaning the leaf frame is at position 0.
func (t *Trie[T]) AddStack(stack []T) int {
	stackID := mkStackID(stack)
	if stackItem, exists := t.uniqueStacks[stackID]; exists {
		return stackItem.parentArrayIdx
	}

	parentStackID := t.rootID
	for i := len(stack) - 1; i >= 0; i-- {
		stackID = mkStackID(stack[i:])
		index, ok := t.uniqueStacks[stackID]
		if !ok {
			index.parentArrayIdx = len(t.uniqueStacks)
			index.locationIdx, ok = t.uniqueLocations[stack[i]]
			if !ok {
				index.locationIdx = len(t.uniqueLocations)
				t.uniqueLocations[stack[i]] = index.locationIdx
			}
			index.parentStackID = parentStackID
			t.uniqueStacks[stackID] = index
		}
		parentStackID = stackID
	}

	index := t.uniqueStacks[stackID].parentArrayIdx
	t.stackIndex = append(t.stackIndex, index)
	return index
}

// ToArrays converts the trie to arrays that can be used for further processing.
func (t *Trie[T]) ToArrays() (locationTable []T,
	stackParentArray, stackLocationIndex, stackIndex []int) {
	// Create the location table with a single allocation.
	locationTable = make([]T, len(t.uniqueLocations))
	for elem, idx := range t.uniqueLocations {
		locationTable[idx] = elem
	}

	// Create the stack arrays, each with a single allocation.
	stackLocationIndex = make([]int, len(t.uniqueStacks))
	stackParentArray = make([]int, len(t.uniqueStacks))
	for _, v := range t.uniqueStacks {
		stackLocationIndex[v.parentArrayIdx] = v.locationIdx
		if v.locationIdx == 0 {
			stackParentArray[v.parentArrayIdx] = 0
		} else {
			parentStackTrace := t.uniqueStacks[v.parentStackID]
			stackParentArray[v.parentArrayIdx] = parentStackTrace.parentArrayIdx
		}
	}

	return locationTable, stackParentArray, stackLocationIndex, t.stackIndex
}

// Indices returns the indices of the added stacks in the trie.
func (t *Trie[T]) Indices() []int {
	return t.stackIndex
}

type stackData struct {
	parentStackID  string
	parentArrayIdx int
	locationIdx    int
}

// mkStackID creates a unique ID for a stack trace.
func mkStackID[T any](stack []T) string {
	var builder strings.Builder
	for i := range stack {
		builder.WriteString(fmt.Sprintf("%v", stack[i]))
		builder.WriteString("|")
	}
	return builder.String()
}

// BuildStacks reconstructs the stacks from the arrays.
func BuildStacks[T comparable](locationTable []T, stackParentArray []int, stackLocationIndex []int,
	stackIndex []int) [][]T {
	stacks := make([][]T, len(stackIndex))
	for stackIdx, i := range stackIndex {
		stack := make([]T, 0)
		for i != 0 {
			stack = append(stack, locationTable[stackLocationIndex[i]])
			i = stackParentArray[i]
		}
		stacks[stackIdx] = stack
	}
	return stacks
}
