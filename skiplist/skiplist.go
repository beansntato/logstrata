package skiplist

import (
	"math"
	"math/rand"
)

type Node struct {
	next  []*Node
	Value int
}

type SkipList struct {
	head     *Node
	maxLevel int
}

func New(maxLevel int) *SkipList {
	return &SkipList{
		maxLevel: maxLevel,
		head: &Node{
			next:  make([]*Node, maxLevel),
			Value: math.MinInt,
		},
	}
}

func (sl *SkipList) Insert(value int) *Node {
	current := sl.head
	level := sl.maxLevel - 1
	update := make([]*Node, sl.maxLevel)

	// check each level for the last traversed element
	// the last traversed element can possibly be the predecessor depending on nodeLevel luck
	for i := level; i >= 0; i-- {
		for current.next[i] != nil && current.next[i].Value < value {
			current = current.next[i]
		}
		// current is now the rightmost node on level i that is < value
		update[i] = current

	}

	// test luck, p here is hardcoded as 0.5, make it configurable if needed
	nodeLevel := 1
	for rand.Float64() > 0.5 && nodeLevel < sl.maxLevel {
		nodeLevel++
	}

	newNode := &Node{
		next:  make([]*Node, nodeLevel),
		Value: value,
	}

	for i := 0; i < nodeLevel; i++ {
		newNode.next[i] = update[i].next[i]
		update[i].next[i] = newNode
	}

	return newNode
}

func (sl *SkipList) Search(target int) *Node {
	current := sl.head
	level := sl.maxLevel - 1

	for i := level; i >= 0; i-- {
		for current != nil {
			if target > current.Value {
				current = current.next[i]
			} else if target == current.Value {
				return current
			} else {
				current = sl.head
				break
			}
		}
	}

	return nil
}
