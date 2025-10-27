package main

import (
	"unsafe"
)

type COWBuffer struct {
	data []byte
	refs *int
}

// создать буфер с определенными данными
func NewCOWBuffer(data []byte) COWBuffer {
	var v int = 1
	return COWBuffer{
		data: data,
		// fix "Я бы единицей инициализировал бы, ведь 1 объект создался уже"
		refs: &v,
	}
}

// создать новую копию буфера
func (b *COWBuffer) Clone() COWBuffer {
	*b.refs = *b.refs + 1
	return COWBuffer{
		data: b.data,
		refs: b.refs,
	}
}

// перестать использовать копию буффера
func (b *COWBuffer) Close() {
	if *b.refs > 0 {
		b.data = nil
	}
	*b.refs = *b.refs - 1
}

// изменить определенный байт в буффере
func (b *COWBuffer) Update(index int, value byte) bool {

	// fix: "Сложно, зачем такая вложенность условий?"
	if 0 > index || index > len(b.data)-1 {
		return false
	}

	if *b.refs > 0 {
		*b.refs = *b.refs - 1

		dest := make([]byte, len(b.data))
		copy(dest, b.data)

		// fix: "Внутри можно конструктор переиспользовать"
		(*b) = NewCOWBuffer(dest)
	}

	b.data[index] = value
	return true

}

// сконвертировать буффер в строку
func (b *COWBuffer) String() string {
	return unsafe.String(unsafe.SliceData(b.data), len(b.data))
}
