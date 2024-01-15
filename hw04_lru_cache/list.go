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
	length int

	firstItem *ListItem
	lastItem  *ListItem
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *ListItem {
	return l.firstItem
}

func (l *list) Back() *ListItem {
	return l.lastItem
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := NewListItem(v)
	firstItem := l.Front()

	if firstItem == nil {
		l.append(&item)
		return &item
	}

	l.appendBefore(firstItem, &item)
	return &item
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := NewListItem(v)
	lastItem := l.Back()

	if lastItem == nil {
		l.append(&item)
		return &item
	}

	l.appendAfter(lastItem, &item)
	return &item
}

func (l *list) Remove(src *ListItem) {
	l.length--

	if src.Prev == nil {
		l.firstItem = src.Next
	} else {
		src.Prev.Next = src.Next
	}

	if src.Next == nil {
		l.lastItem = src.Prev
	} else {
		src.Next.Prev = src.Prev
	}
}

func (l *list) MoveToFront(src *ListItem) {
	if src.Prev == nil {
		return
	}

	src.Prev.Next = src.Next

	if src.Next != nil {
		src.Next.Prev = src.Prev
	} else {
		l.lastItem = src.Prev
	}

	l.firstItem.Prev = src

	src.Prev = nil
	src.Next = l.firstItem

	l.firstItem = src
}

func (l *list) append(i *ListItem) {
	l.length++

	l.firstItem = i
	l.lastItem = i
}

func (l *list) appendBefore(src *ListItem, i *ListItem) {
	l.length++

	i.Next = src

	if src.Prev == nil {
		i.Prev = nil
		l.firstItem = i
	} else {
		i.Prev = src.Prev
		src.Prev.Next = i
	}

	src.Prev = i
}

func (l *list) appendAfter(src *ListItem, i *ListItem) {
	l.length++

	i.Prev = src

	if src.Next == nil {
		i.Next = nil
		l.lastItem = i
	} else {
		i.Next = src.Next
		src.Next.Prev = i
	}

	src.Next = i
}

func NewList() List {
	return new(list)
}

func NewListItem(v interface{}) ListItem {
	return ListItem{
		Value: v,
	}
}
