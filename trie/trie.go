// Package trie provides functionality to efficiently store and manipulate stack traces
// using a trie data structure.
// It supports operations such as adding stack traces, converting the trie to arrays,
// and reconstructing stacks from arrays.
package trie

import (
	"strings"
)

// Trie is a data structure that stores unique stack traces and their relationships.
type Trie struct {
	uniqueStacks    map[string]stackData
	uniqueLocations map[string]int
	locationTable   []Location
	stackIndex      []int // stackIndex keeps track of the added stacks.
	rootID          string
}

// New creates a new Trie instance. It initializes the trie with an artificial root frame.
func New() *Trie {
	trie := Trie{
		uniqueStacks:    make(map[string]stackData),
		uniqueLocations: make(map[string]int),
		stackIndex:      make([]int, 0),
		rootID:          mkStackID([]string{""}),
	}

	// Add an artificial root frame.
	trie.locationTable = append(trie.locationTable, Location{name: ""})
	trie.uniqueLocations[""] = 0
	trie.uniqueStacks[trie.rootID] = stackData{parentStackID: "", parentArrayIdx: 0, locationIdx: 0}

	return &trie
}

// NewFromStacks creates a new Trie instance from a slice of stack traces.
func NewFromStacks(stacks [][]string) *Trie {
	trie := New()
	for _, stack := range stacks {
		// Add the stack trace to the trie.
		// The leaf frame is at position 0, so we add the stack in natural order.
		trie.AddStack(stack)
	}
	return trie
}

// Len returns the number of added unique stacks in the trie plus one for the artificial root frame.
func (t *Trie) Len() int {
	return len(t.uniqueStacks)
}

// Index returns the index of the stack in the trie or -1 if the stack is not found.
func (t *Trie) Index(stack []string) int {
	if stackItem, ok := t.uniqueStacks[mkStackID(stack)]; ok {
		return stackItem.parentArrayIdx
	}
	return -1
}

// Exists checks if a stack trace exists in the trie.
func (t *Trie) Exists(stack []string) bool {
	_, ok := t.uniqueStacks[mkStackID(stack)]
	return ok
}

// AddStack adds a stack trace to the trie. The stack trace is expected to be in
// natural order, meaning the leaf frame is at position 0.
func (t *Trie) AddStack(stack []string) int {
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
func (t *Trie) ToArrays() (locationTable []Location,
	stackParentArray, stackLocationIndex, stackIndex []int) {
	// Create the location table with a single allocation.
	locationTable = make([]Location, len(t.uniqueLocations))
	for name, idx := range t.uniqueLocations {
		locationTable[idx] = Location{name: name}
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
func (t *Trie) Indices() []int {
	return t.stackIndex
}

// Location is a fake type to represent a location.
type Location struct {
	name string
}

type stackData struct {
	parentStackID  string
	parentArrayIdx int
	locationIdx    int
}

// mkStackID creates a unique ID for a stack trace.
func mkStackID(stack []string) string {
	var builder strings.Builder
	for i := range stack {
		builder.WriteString(stack[i])
		builder.WriteString("|")
	}
	return builder.String()
}

// BuildStacks reconstructs the stacks from the arrays.
func BuildStacks(locationTable []Location, stackParentArray []int, stackLocationIndex []int,
	stackIndex []int) [][]string {
	stacks := make([][]string, len(stackIndex))
	for stackIdx, i := range stackIndex {
		stack := make([]string, 0)
		for i != 0 {
			stack = append(stack, locationTable[stackLocationIndex[i]].name)
			i = stackParentArray[i]
		}
		stacks[stackIdx] = stack
	}
	return stacks
}
