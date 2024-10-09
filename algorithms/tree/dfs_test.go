package main

import (
	"reflect"
	"testing"
)

func TestDFS(t *testing.T) {
	tests := []struct {
		input    *Node
		expected []rune
	}{
		{
			//
			input:    nil,
			expected: []rune{},
		},
		{
			//   S
			input:    &Node{val: 'S'},
			expected: []rune{'S'},
		},
		{
			//   S
			//   |
			//   A
			input: &Node{val: 'S', children: []*Node{
				{val: 'A'},
			}},
			expected: []rune{'S', 'A'},
		},
		{
			//       S
			//      / \
			//     A   B
			//    /
			//   C
			input: &Node{val: 'S', children: []*Node{
				{val: 'A', children: []*Node{
					{val: 'C'},
				}},
				{val: 'B'},
			}},
			expected: []rune{'S', 'A', 'C', 'B'},
		},
		{
			//       S
			//      /|\
			//     A B C
			//    / \   \
			//   D   E   F
			//       |
			//       G
			input: &Node{val: 'S', children: []*Node{
				{val: 'A', children: []*Node{
					{val: 'D'},
					{val: 'E', children: []*Node{
						{val: 'G'},
					}},
				}},
				{val: 'B'},
				{val: 'C', children: []*Node{
					{val: 'F'},
				}},
			}},
			expected: []rune{'S', 'A', 'D', 'E', 'G', 'B', 'C', 'F'},
		},
		{
			//       S
			//      /|\
			//     A B C
			//    / \
			//   D   E
			//       |
			//       F
			input: &Node{val: 'S', children: []*Node{
				{val: 'A', children: []*Node{
					{val: 'D'},
					{val: 'E', children: []*Node{
						{val: 'F'},
					}},
				}},
				{val: 'B'},
				{val: 'C'},
			}},
			expected: []rune{'S', 'A', 'D', 'E', 'F', 'B', 'C'},
		},
		{
			//        S
			//      / | \
			//     A  B  C
			//    / \   / \
			//   D   E G   F
			//       |
			//       G
			input: &Node{val: 'S', children: []*Node{
				{val: 'A', children: []*Node{
					{val: 'D', children: []*Node{
						{val: 'G', children: nil},
					}},
					{val: 'E'},
				}},
				{val: 'B'},
				{val: 'C', children: []*Node{
					{val: 'F'},
					{val: 'G'},
				}},
			}},
			expected: []rune{'S', 'A', 'D', 'G', 'E', 'B', 'C', 'F', 'G'},
		},
	}

	name := "DFS"
	algo := DFS[rune]
	for _, test := range tests {
		result := algo(test.input)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("[%s] expected %v, got %v", name, test.expected, result)
		}
	}
}
