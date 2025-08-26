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

func (l *List) Size() int {
	return len(l.items)
}

func (l *List) Range(start, stop int) []resp.Value {
	return l.items[start:stop]
}
