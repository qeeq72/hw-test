package hw04lrucache

const (
	UnknownPosition = iota
	FrontPosition
	MiddlePosition
	BackPosition
)

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	len   int
	front *ListItem
	back  *ListItem
}

func NewList() List {
	return new(list)
}

func (l list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	i := &ListItem{}
	i.Value = v
	i.Prev = nil
	if l.len == 0 {
		i.Next = nil
		l.back = i
	} else {
		i.Next = l.front
		l.front.Prev = i
	}
	l.front = i
	l.len++
	return i
}

func (l *list) PushBack(v interface{}) *ListItem {
	i := &ListItem{}
	i.Value = v
	i.Next = nil
	if l.len == 0 {
		i.Prev = nil
		l.front = i
	} else {
		i.Prev = l.back
		l.back.Next = i
	}
	l.back = i
	l.len++
	return i
}

func (l *list) Remove(i *ListItem) {
	if l.len != 0 && i != nil {
		switch getItemPosition(i) {
		case MiddlePosition:
			i.Next.Prev = i.Prev
			i.Prev.Next = i.Next
		case FrontPosition:
			i.Next.Prev = nil
			l.front = i.Next
		case BackPosition:
			i.Prev.Next = nil
			l.back = i.Prev
		default:
			return
		}
		i.Value = nil
		l.len--
	}
}

func (l *list) MoveToFront(i *ListItem) {
	if l.len > 1 && i != nil {
		if i.Prev != nil {
			if i.Next != nil {
				i.Next.Prev = i.Prev
				i.Prev.Next = i.Next
			} else {
				i.Prev.Next = nil
				l.back = i.Prev
			}
			i.Prev = nil
			i.Next = l.front
			l.front.Prev = i
			l.front = i
		}
	}
}

func getItemPosition(i *ListItem) int {
	if i.Next != nil && i.Prev != nil {
		return MiddlePosition
	}
	if i.Next != nil {
		return FrontPosition
	}
	if i.Prev != nil {
		return BackPosition
	}
	return UnknownPosition
}
