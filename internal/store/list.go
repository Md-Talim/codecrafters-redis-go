package store

import "github.com/md-talim/codecrafters-redis-go/internal/resp"

type List struct {
	items []resp.Value
}

func NewList() *List {
	return &List{items: []resp.Value{}}
}

func (l *List) Append(newItems []resp.Value) {
	l.items = append(l.items, newItems...)
}

func (l *List) Prepend(newItems []resp.Value) {
	for _, item := range newItems {
		l.items = append([]resp.Value{item}, l.items...)
	}
}

func (l *List) Size() int {
	return len(l.items)
}

func (l *List) Range(start, stop int) []resp.Value {
	return l.items[start:stop]
}

// Pop removes the first element of the list and returns it
func (l *List) Pop() resp.Value {
	firstElement := l.items[0]
	l.items = l.items[1:]
	return firstElement
}

func (l *List) IsEmpty() bool {
	return len(l.items) == 0
}
