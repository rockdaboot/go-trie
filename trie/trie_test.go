package trie

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_buildArrays(t *testing.T) {
	tests := []struct {
		name   string
		stacks [][]string
		want   [][]string
	}{
		{
			name:   "empty",
			stacks: [][]string{},
		},
		{
			name: "single",
			stacks: [][]string{
				{"bar", "foo", "main"},
			},
		},
		{
			name: "multiple",
			stacks: [][]string{
				{"bar", "foo", "main"},
				{"baz1", "foo", "main"},
				{"baz1", "bar", "foo", "main"},
				{"baz2", "bar", "foo", "main"},
				{"baz2", "foo", "main"},
				{"why", "why", "what"},
			},
		},
		{
			name: "de-duplicated",
			stacks: [][]string{
				{"bar", "foo", "main"},
				{"bar", "foo", "main"},
				{"baz1", "foo", "main"},
			},
			want: [][]string{
				{"bar", "foo", "main"},
				{"baz1", "foo", "main"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := NewFromStacks(tt.stacks)
			locationTable, stackParentArray, stackLocationIndex, stackIndex := data.ToArrays()
			stacks := BuildStacks(locationTable, stackParentArray, stackLocationIndex, stackIndex)
			if tt.want == nil {
				tt.want = tt.stacks
			}
			if !reflect.DeepEqual(stacks, tt.want) {
				t.Errorf("buildArrays() %v, want %v", stacks, tt.want)
			}
			require.Equal(t, tt.want, stacks)
		})
	}
}
