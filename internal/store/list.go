package store

type List struct {
	items []any
}

func NewList() *List {
	return &List{items: []any{}}
}

func (l *List) Append(newItems []any) {
	l.items = append(l.items, newItems...)
}

func (l *List) Size() int {
	return len(l.items)
}
