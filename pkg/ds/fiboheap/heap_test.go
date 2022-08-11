package fiboheap_test

import (
	"github.com/PieterD/brot/pkg/ds/fiboheap"
	"github.com/PieterD/brot/pkg/ds/fiboheap/poolalloc"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFiboHeap(t *testing.T) {
	var tests = []struct {
		desc     string
		adds     [][]int
		expected []int
	}{
		{
			desc: "add 3",
			adds: [][]int{
				{7, 3, 5},
			},
			expected: []int{3, 5, 7},
		},
		{
			desc: "add 3+1 end",
			adds: [][]int{
				{7, 3, 5},
				{8},
			},
			expected: []int{3, 5, 7, 8},
		},
		{
			desc: "add 3+1 newmin",
			adds: [][]int{
				{7, 3, 5},
				{2},
			},
			expected: []int{2, 3, 5, 7},
		},
		{
			desc: "add 3+3",
			adds: [][]int{
				{7, 3, 5},
				{8, 2, 6},
			},
			expected: []int{2, 3, 5, 6, 7, 8},
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			less := func(l, r int) bool {
				return l < r
			}
			allocator := poolalloc.New[int]()
			h := fiboheap.NewHeap[int](less, allocator, 3)
			for _, adds := range test.adds {
				h.Add(adds...)
			}
			var gots []int
			for {
				v, ok := h.Pop()
				if !ok {
					break
				}
				gots = append(gots, v)
			}
			require.Equal(t, test.expected, gots)
			_, ok := h.Pop()
			require.False(t, ok)
		})
	}
}
