package util

type Item[T any] struct {
	Value *T
	Next  *Item[T]
}

type Queue[T any] struct {
	Head *Item[T]
}

func (q *Queue[T]) Add(value *T) {
	item := &Item[T]{Value: value}

	if q.Head == nil {
		q.Head = item
		return
	}

	item.Next = q.Head
	q.Head = item
}

func (q *Queue[T]) Remove(matcher func(*T) bool) {
	var prev *Item[T]
	item := q.Head

	for item != nil {
		if matcher(item.Value) {
			if item == q.Head {
				q.Head = item.Next
			} else {
				prev.Next = item.Next
			}

			break
		}

		prev = item
		item = item.Next
	}
}
