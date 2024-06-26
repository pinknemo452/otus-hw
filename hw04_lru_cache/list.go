package hw04lrucache

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
	head *ListItem
	tail *ListItem
	len  int
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.head
}

func (l *list) Back() *ListItem {
	return l.tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	li := ListItem{
		Value: v,
		Next:  l.head,
		Prev:  nil,
	}
	if l.head != nil {
		l.head.Prev = &li
	} else {
		l.head = &li
		l.tail = &li
		li.Next = nil
		li.Prev = nil
		l.len++
		return &li
	}
	l.head = &li
	l.len++
	return &li
}

func (l *list) PushBack(v interface{}) *ListItem {
	li := ListItem{
		Value: v,
		Next:  nil,
		Prev:  l.tail,
	}
	if l.tail != nil {
		l.tail.Next = &li
	} else {
		return l.PushFront(v)
	}

	l.tail = &li
	l.len++
	return &li
}

func (l *list) Remove(i *ListItem) {
	if i.Prev == nil {
		l.head = i.Next
	} else {
		i.Prev.Next = i.Next
	}
	if i.Next == nil {
		l.tail = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}
	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	l.Remove(i)
	l.PushFront(i.Value)
}

func NewList() List {
	return new(list)
}
