package internal

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDoublyLinkedList_IsEmpty(t *testing.T) {
	t.Run("List is empty", func(t *testing.T) {
		list := DoublyLinkedList{
			head: nil,
			tail: nil,
		}

		require.True(t, list.IsEmpty(), "List should be considered to be empty.")
	})

	t.Run("List is not empty", func(t *testing.T) {
		listItem := &item{
			next:  nil,
			prev:  nil,
			value: nil,
		}

		list := DoublyLinkedList{
			head: listItem,
			tail: listItem,
		}

		require.False(t, list.IsEmpty(), "List should be considered to be not empty.")
	})
}

func TestDoublyLinkedList_PopFront(t *testing.T) {
	t.Run("Popping when list is empty returns an error", func(t *testing.T) {
		list := DoublyLinkedList{
			head: nil,
			tail: nil,
		}

		value, err := list.PopFront()

		require.Nil(t, value, "The returned value is expected to be <nil>.")
		require.Same(t, ErrListIsEmpty, err, "The returned error is expected to be ErrListIsEmpty.")
	})

	t.Run("Popping from the list with one item correctly sets head and tail", func(t *testing.T) {
		listItem := item{
			next:  nil,
			prev:  nil,
			value: 123,
		}

		list := DoublyLinkedList{
			head: &listItem,
			tail: &listItem,
		}

		value, err := list.PopFront()

		require.Equal(t, 123, value, "The returned value is invalid.")
		require.Nil(t, err, "No error is expected to be returned.")
		require.Nil(t, list.head, "The list's head is expected to be <nil>.")
		require.Nil(t, list.tail, "The list's tail is expected to be <nil>.")
	})

	t.Run("Popping from the list with two items correctly sets head and tail", func(t *testing.T) {
		listItem1 := &item{
			next:  nil,
			prev:  nil,
			value: 123,
		}

		listItem2 := &item{
			next:  nil,
			prev:  listItem1,
			value: 456,
		}

		listItem1.next = listItem2

		list := DoublyLinkedList{
			head: listItem1,
			tail: listItem2,
		}

		value, err := list.PopFront()

		require.Equal(t, 123, value, "The returned value is unexpected.")
		require.Nil(t, err, "No error is expected to be returned.")
		require.Same(t, list.head, list.tail, "The list's head is expected to be the same as list's tail.")
		require.Same(t, listItem2, list.head, "The list's head is expected to be listItem2.")
		require.Nil(t, list.head.next, "The list's head.next property is expected to be <nil>.")
		require.Nil(t, list.head.prev, "The list's head.prev property is expected to be <nil>.")
	})

	t.Run("Popping from the list with three items correctly sets head and tail", func(t *testing.T) {
		listItem1 := &item{
			next:  nil,
			prev:  nil,
			value: 123,
		}

		listItem2 := &item{
			next:  nil,
			prev:  listItem1,
			value: 456,
		}

		listItem3 := &item{
			next:  nil,
			prev:  listItem2,
			value: 789,
		}

		listItem1.next = listItem2
		listItem2.next = listItem3

		list := DoublyLinkedList{
			head: listItem1,
			tail: listItem3,
		}

		value, err := list.PopFront()

		require.Equal(t, 123, value, "The returned value is unexpected.")
		require.Nil(t, err, "No error is expected to be returned.")
		require.NotSame(t, list.head, list.tail, "The list's head is expected to be different than list's tail.")

		require.Same(t, listItem2, list.head, "The list's head is expected to be listItem2.")
		require.Same(t, listItem3, list.head.next, "The list's head.next property is expected to be listItem3.")
		require.Nil(t, list.head.prev, "The list's head.prev property is expected to be <nil>.")

		require.Same(t, listItem3, list.tail, "The list's head is expected to be listItem3.")
		require.Nil(t, list.tail.next, "The list's tail.next property is expected to be <nil>.")
		require.Same(t, listItem2, list.tail.prev, "The list's tail.prev property is expected to be listItem2.")
	})
}

func TestDoublyLinkedList_PushBack(t *testing.T) {
	t.Run("Push new element to empty list", func(t *testing.T) {
		list := DoublyLinkedList{
			head: nil,
			tail: nil,
		}

		list.PushBack(123)

		require.NotNil(t, list.head, "The list's head is expected not to be <nil>.")
		require.Same(t, list.head, list.tail, "The list's head is expected to be the same as list's tail.")
		require.Nil(t, list.head.next, "The list's head.next property is expected to be <nil>.")
		require.Nil(t, list.head.prev, "The list's head.prev property is expected to be <nil>.")
		require.Equal(t, 123, list.head.value, "The list's head.value property is expected to be set properly.")
	})

	t.Run("Push new element to list with one item", func(t *testing.T) {
		listItem1 := &item{
			next:  nil,
			prev:  nil,
			value: 123,
		}

		list := DoublyLinkedList{
			head: listItem1,
			tail: listItem1,
		}

		list.PushBack(456)

		require.NotSame(t, list.head, list.tail, "The list's head is expected to be different than list's tail.")

		require.NotNil(t, list.head, "The list's head is expected not to be <nil>.")
		require.Same(t, listItem1, list.head, "The list's head is expected to be listItem1.")
		require.NotNil(t, list.head.next, "The list's head.next property is expected not to be <nil>.")
		require.Nil(t, list.head.prev, "The list's head.prev property is expected to be <nil>.")
		require.Equal(t, 123, list.head.value, "The list's head.value property is expected to be set properly.")

		require.NotNil(t, list.tail, "The list's tail is expected not to be <nil>.")
		require.NotSame(t, listItem1, list.tail, "The list's tail is expected not to be listItem1.")
		require.Nil(t, list.tail.next, "The list's tail.next property is expected to be <nil>.")
		require.Same(t, listItem1, list.tail.prev, "The list's tail.prev property is expected to be listItem1.")
		require.Equal(t, 456, list.tail.value, "The list's tail.value property is expected to be set properly.")
	})

	t.Run("Push new element to list with two items", func(t *testing.T) {
		listItem1 := &item{
			next:  nil,
			prev:  nil,
			value: 123,
		}

		listItem2 := &item{
			next:  nil,
			prev:  listItem1,
			value: 456,
		}

		listItem1.next = listItem2

		list := DoublyLinkedList{
			head: listItem1,
			tail: listItem2,
		}

		list.PushBack(789)

		require.NotSame(t, list.head, list.tail, "The list's head is expected to be different than list's tail.")

		require.NotNil(t, list.head, "The list's head is expected not to be <nil>.")
		require.Same(t, listItem1, list.head, "The list's head is expected to be listItem1.")
		require.NotNil(t, list.head.next, "The list's head.next property is expected not to be <nil>.")
		require.Nil(t, list.head.prev, "The list's head.prev property is expected to be <nil>.")
		require.Equal(t, 123, list.head.value, "The list's head.value property is expected to be set properly.")

		require.NotNil(t, list.tail, "The list's tail is expected not to be <nil>.")
		require.NotSame(t, listItem2, list.tail, "The list's tail is expected not to be listItem2.")
		require.Nil(t, list.tail.next, "The list's tail.next property is expected to be <nil>.")
		require.Same(t, listItem2, list.tail.prev, "The list's tail.prev property is expected to be listItem2.")
		require.Equal(t, 789, list.tail.value, "The list's tail.value property is expected to be set properly.")
	})
}
