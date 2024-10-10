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
			input:    &Node{Value: 'S'},
			expected: []rune{'S'},
		},
		{
			//   S
			//   |
			//   A
			input: &Node{Value: 'S', Children: []*Node{
				{Value: 'A'},
			}},
			expected: []rune{'S', 'A'},
		},
		{
			//       S
			//      / \
			//     A   B
			//    /
			//   C
			input: &Node{Value: 'S', Children: []*Node{
				{Value: 'A', Children: []*Node{
					{Value: 'C'},
				}},
				{Value: 'B'},
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
			input: &Node{Value: 'S', Children: []*Node{
				{Value: 'A', Children: []*Node{
					{Value: 'D'},
					{Value: 'E', Children: []*Node{
						{Value: 'G'},
					}},
				}},
				{Value: 'B'},
				{Value: 'C', Children: []*Node{
					{Value: 'F'},
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
			input: &Node{Value: 'S', Children: []*Node{
				{Value: 'A', Children: []*Node{
					{Value: 'D'},
					{Value: 'E', Children: []*Node{
						{Value: 'F'},
					}},
				}},
				{Value: 'B'},
				{Value: 'C'},
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
			input: &Node{Value: 'S', Children: []*Node{
				{Value: 'A', Children: []*Node{
					{Value: 'D', Children: []*Node{
						{Value: 'G', Children: nil},
					}},
					{Value: 'E'},
				}},
				{Value: 'B'},
				{Value: 'C', Children: []*Node{
					{Value: 'F'},
					{Value: 'G'},
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
