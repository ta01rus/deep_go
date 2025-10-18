package main

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

func TestCircularQueue(t *testing.T) {
	const queueSize = 3
	queue := NewCircularQueue[int16](queueSize)

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())

	assert.Equal(t, -1, queue.Front())
	assert.Equal(t, -1, queue.Back())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Push(1))
	assert.True(t, queue.Push(2))
	assert.True(t, queue.Push(3))
	assert.False(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int16{1, 2, 3}, queue.values))

	assert.False(t, queue.Empty())
	assert.True(t, queue.Full())

	assert.Equal(t, 1, queue.Front())
	assert.Equal(t, 3, queue.Back())

	assert.True(t, queue.Pop())
	assert.False(t, queue.Empty())
	assert.False(t, queue.Full())
	assert.True(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int16{4, 2, 3}, queue.values))

	assert.Equal(t, 2, queue.Front())
	assert.Equal(t, 4, queue.Back())

	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())
}

/*
тест который проверяет действия при определенных
состояниях очереди
*/
func Test_CircularQueueState(t *testing.T) {
	tests := []struct {
		name       string
		queue      *CircularQueue[int16]
		isFull     bool
		isEmpty    bool
		Front      int16
		Back       int16
		operations func(*CircularQueue[int16]) []error // предварительные операции
		debug      bool                                // для debug
	}{
		{
			"1. очередь пуста",
			WithValues(0, 0, make([]int16, 0, 4)...),
			false,
			true,
			-1,
			-1,
			func(*CircularQueue[int16]) []error {
				return nil
			},
			true,
		},

		{
			"2. очередь заполнена",
			WithValues(0, 3, []int16{1, 2, 3, 4}...),
			true,
			false,
			1,
			4,
			func(*CircularQueue[int16]) []error {
				return nil
			},
			false,
		},

		{
			"3. добавим 1 элемент в пустую очередь ",
			WithValues(0, 0, make([]int16, 0, 4)...),
			false,
			false,
			5,
			5,
			func(q *CircularQueue[int16]) []error {
				errs := make([]error, 0, 5)
				p := q.Push(5)
				if p == false {
					errs = append(errs, fmt.Errorf("Push(5) = false wont = %t  ", true))
				}
				return errs
			},
			false,
		},

		{
			"4. добавим 2 элемента в пустую очередь ",
			WithValues(0, 0, make([]int16, 0, 4)...),
			false,
			false,
			5,
			6,
			func(q *CircularQueue[int16]) []error {
				errs := make([]error, 0, 5)
				p := q.Push(5)
				if p == false {
					errs = append(errs, fmt.Errorf("Push(5) = false wont = %t  ", true))
				}

				p = q.Push(6)
				if p == false {
					errs = append(errs, fmt.Errorf("Push(6) = false wont = %t  ", true))
				}
				return errs
			},
			false,
		},
		{
			"5. добавим 1 элемента в полную очередь ",
			WithValues(0, 3, []int16{1, 2, 3, 4}...),
			true,
			false,
			1,
			4,
			func(q *CircularQueue[int16]) []error {
				errs := make([]error, 0, 5)
				p := q.Push(5)
				if p == true {
					errs = append(errs, fmt.Errorf("Push(5) = true wont = %t  ", false))
				}
				return errs
			},
			false,
		},
		{
			"6. уберем и добавим 1 элемент в не полной очереди при этом end находится перед стартом",
			WithValues(2, 1, []int16{1, 2, 3, 4}...),
			false,
			false,
			1,
			5,
			func(q *CircularQueue[int16]) []error {
				errs := make([]error, 0, 5)
				p := q.Pop()
				if !p {
					errs = append(errs, fmt.Errorf("Pop() = false wont = %t  ", true))
				}

				p = q.Pop()
				if !p {
					errs = append(errs, fmt.Errorf("Pop() = false wont = %t  ", true))
				}

				p = q.Push(5)
				if p == false {
					errs = append(errs, fmt.Errorf("Push(5) = false wont = %t  ", true))
				}
				return errs
			},
			false,
		},
		{
			"7. уберем и добавим 1 элемент в полной очереди",
			WithValues(0, 3, []int16{1, 2, 3, 4}...),
			true,
			false,
			2,
			5,
			func(q *CircularQueue[int16]) []error {
				errs := make([]error, 0, 5)
				p := q.Pop()
				if !p {
					errs = append(errs, fmt.Errorf("Pop() = false wont = %t  ", true))
				}

				p = q.Push(5)
				if p == false {
					errs = append(errs, fmt.Errorf("Push(5) = false wont = %t  ", true))
				}
				return errs
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := make([]error, 0, 4)
			if tt.debug {
				// сюда ставим точку останова
				log.Println("debug")
			}

			t.Log(tt.name)
			{
				opErrors := tt.operations(tt.queue)
				if len(opErrors) > 0 {
					errs = append(errs, opErrors...)
				}

				isEmpty := tt.queue.Empty()
				isFull := tt.queue.Full()
				front := tt.queue.Front()
				back := tt.queue.Back()

				if isEmpty != tt.isEmpty {
					errs = append(errs, fmt.Errorf("queue.Empty() = %t wont: %t", isEmpty, tt.isEmpty))
				}

				if isFull != tt.isFull {
					errs = append(errs, fmt.Errorf("queue.Full() = %t wont: %t", isFull, tt.isFull))
				}

				if front != int(tt.Front) {
					errs = append(errs, fmt.Errorf("queue.Front() = %d wont: %d", front, tt.Front))
				}

				if back != int(tt.Back) {
					errs = append(errs, fmt.Errorf("queue.Back() = %d wont: %d", front, tt.Back))
				}

				if len(errs) > 0 {
					t.Fatal(errors.Join(errs...))
				}
			}
		})
	}
}
