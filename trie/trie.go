package trie

import (
	"strings"
)

type Trie struct {
	uniqueStacks    map[string]stackData
	uniqueLocations map[string]int
	locationTable   []Location
	stackIndex      []int // stackIndex keeps track of the added stacks.
	rootID          string
}

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

func (t *Trie) ToArrays() (locationTable []Location, stackParentArray []int, stackLocationIndex []int, stackIndex []int) {
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

	stackIndex = t.stackIndex

	return
}

func (t *Trie) LeafIndices() []int {
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

func buildArrays(stacks [][]string) (locationTable []Location, stackParentArray []int, stackLocationIndex []int, stackIndex []int) {
	uniqueStacks := make(map[string]stackData)
	uniqueLocations := make(map[string]int)

	// Add an artificial root frame.
	uniqueLocations[""] = 0
	rootID := mkStackID([]string{""})
	uniqueStacks[rootID] = stackData{parentStackID: "", parentArrayIdx: 0, locationIdx: 0}

	for _, stack := range stacks {
		parentStackID := rootID
		for i := len(stack) - 1; i >= 0; i-- {
			stackID := mkStackID(stack[i:])
			index, ok := uniqueStacks[stackID]
			if !ok {
				index.parentArrayIdx = len(uniqueStacks)
				index.locationIdx, ok = uniqueLocations[stack[i]]
				if !ok {
					index.locationIdx = len(uniqueLocations)
					uniqueLocations[stack[i]] = index.locationIdx
				}
				index.parentStackID = parentStackID
				uniqueStacks[stackID] = index
			}
			parentStackID = stackID
		}
	}

	// Create the location table with a single allocation.
	locationTable = make([]Location, len(uniqueLocations))
	for name, idx := range uniqueLocations {
		locationTable[idx] = Location{name: name}
	}

	// Create the stack arrays, each with a single allocation.
	stackLocationIndex = make([]int, len(uniqueStacks))
	stackParentArray = make([]int, len(uniqueStacks))
	for _, v := range uniqueStacks {
		stackLocationIndex[v.parentArrayIdx] = v.locationIdx
		if v.locationIdx == 0 {
			stackParentArray[v.parentArrayIdx] = 0
		} else {
			parentStackTrace := uniqueStacks[v.parentStackID]
			stackParentArray[v.parentArrayIdx] = parentStackTrace.parentArrayIdx
		}
	}

	// Keep track of the leaf frames, allows reconstructing the input stacks.
	stackIndex = make([]int, len(stacks))
	for i, stack := range stacks {
		stackIndex[i] = uniqueStacks[mkStackID(stack)].parentArrayIdx
	}

	return
}

// mkStackID creates a unique ID for a stack trace.
func mkStackID(stack []string) string {
	var builder strings.Builder
	for i := 0; i < len(stack); i++ {
		builder.WriteString(stack[i])
		builder.WriteString("|")
	}
	return builder.String()
}

// buildStacks reconstructs the stacks from the arrays.
func buildStacks(locationTable []Location, stackParentArray []int, stackLocationIndex []int, stackIndex []int) [][]string {
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
