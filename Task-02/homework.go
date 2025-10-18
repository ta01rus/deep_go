package main

func abc(v int) int {
	if v < 0 {
		return v * (-1)
	}
	return v
}

type CircularQueue[T int8 | int16 | int32 | int64] struct {
	values []T

	// позиция начала очереди
	start int
	// позиция конца
	end int
	// количество свободных ячеек
	countFreeCells int
}

func NewCircularQueue[T int8 | int16 | int32 | int64](size int) CircularQueue[T] {
	return CircularQueue[T]{
		values:         make([]T, 0, size),
		countFreeCells: size,
	}

}

// дает изначальное состояние очереди для последующей возможности строить операции и тестировать
func WithValues[T int8 | int16 | int32 | int64](start, end int, values ...T) *CircularQueue[T] {

	countFreeCells := cap(values) - len(values)
	return &CircularQueue[T]{
		start:          start,
		end:            end,
		values:         values,
		countFreeCells: countFreeCells,
	}
}

// добавить значение в конец очереди (false, если очередь заполнена)
func (q *CircularQueue[T]) Push(value T) bool {
	_ = q.values[q.start:len(q.values)]
	if q.Full() {
		return false
	}
	// резервируем
	q.countFreeCells--

	if len(q.values) < cap(q.values) {
		q.values = append(q.values, value)
		q.end = len(q.values) - 1
		return true
	}
	nexPoz := q.end + 1
	if nexPoz >= cap(q.values) {
		nexPoz = 0
	}
	q.end = nexPoz
	q.values[nexPoz] = value

	return true

}

// удалить значение из начала очереди (false, если очередь пустая)
func (q *CircularQueue[T]) Pop() bool {
	_ = q.values[q.start:len(q.values)]

	if q.Empty() {
		return false
	}
	// добавим пустые места
	q.countFreeCells++

	if q.start+1 >= cap(q.values) {
		q.start = 0
		return true
	}
	// это двигаем, так как конец очереди не может быть "удален"
	if q.start == q.end {
		q.end++
	}
	q.start++

	return true
}

// получить значение из начала очереди (-1, если очередь пустая)
func (q *CircularQueue[T]) Front() int {
	if q.Empty() {
		return -1
	}
	return int(q.values[q.start])
}

// получить значение из конца очереди (-1, если очередь пустая)
func (q *CircularQueue[T]) Back() int {
	if q.Empty() {
		return -1
	}
	return int(q.values[q.end])
}

// проверить пустая ли очередь
func (q *CircularQueue[T]) Empty() bool {
	// если количество свободных мест равно длине списка, то пусто
	return q.countFreeCells == cap(q.values)
}

// проверить заполнена ли очередь
func (q *CircularQueue[T]) Full() bool {
	// мест нет
	return q.countFreeCells == 0
}
