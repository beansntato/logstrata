package skiplist_test

import (
	"testing"

	"beanstato.dev/logstrata/skiplist"
)

func TestInsertAndSearch(t *testing.T) {
	sl := skiplist.New(4)
	values := []int{3, 1, 7, 5, 2}
	for _, v := range values {
		sl.Insert(v)
	}
	for _, v := range values {
		node := sl.Search(v)
		if node == nil {
			t.Errorf("Search(%d) = nil, want node", v)
			continue
		}
		if node.Value != v {
			t.Errorf("Search(%d).value = %d, want %d", v, node.Value, v)
		}
	}
}

func TestSearchMissing(t *testing.T) {
	sl := skiplist.New(4)
	sl.Insert(1)
	sl.Insert(3)
	sl.Insert(5)

	for _, v := range []int{0, 2, 4, 6, 100} {
		if node := sl.Search(v); node != nil {
			t.Errorf("Search(%d) = %v, want nil", v, node)
		}
	}
}

func TestSearchEmpty(t *testing.T) {
	sl := skiplist.New(4)
	if node := sl.Search(42); node != nil {
		t.Errorf("Search on empty list = %v, want nil", node)
	}
}

func TestInsertDuplicates(t *testing.T) {
	sl := skiplist.New(4)
	sl.Insert(5)
	sl.Insert(5)
	// At least one node with value 5 must be findable
	node := sl.Search(5)
	if node == nil {
		t.Error("Search(5) = nil after inserting duplicates")
	}
}
