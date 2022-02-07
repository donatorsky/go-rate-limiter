package ratelimiter

import "errors"

var (
	ErrListIsEmpty = errors.New("the list is empty")
)

type item struct {
	next  *item
	prev  *item
	value interface{}
}

type DoublyLinkedList struct {
	head *item
	tail *item
}

func (l *DoublyLinkedList) PushBack(v interface{}) {
	if l.tail == nil {
		l.tail = &item{
			next:  nil,
			prev:  nil,
			value: v,
		}

		l.head = l.tail
	} else {
		i := &item{
			next:  nil,
			prev:  l.tail,
			value: v,
		}

		l.tail.next = i
		l.tail = i
	}
}

func (l *DoublyLinkedList) PopFront() (interface{}, error) {
	currentHead := l.head

	if currentHead == nil {
		return nil, ErrListIsEmpty
	}

	l.head = currentHead.next

	if l.head == nil {
		l.tail = nil
	} else {
		l.head.prev = nil
	}

	return currentHead.value, nil
}

func (l *DoublyLinkedList) IsEmpty() bool {
	return l.head == nil
}
