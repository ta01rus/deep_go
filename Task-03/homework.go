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
	return COWBuffer{
		data: data,
		refs: new(int),
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
		b.data = []byte{}
	}
	*b.refs = *b.refs - 1
}

// изменить определенный байт в буффере
func (b *COWBuffer) Update(index int, value byte) bool {

	if 0 <= index && index < len(b.data)-1 {

		if *b.refs > 0 {
			*b.refs = *b.refs - 1

			dest := make([]byte, len(b.data))
			copy(dest, b.data)

			b.data = dest
			b.refs = new(int)
		}

		b.data[index] = value
		return true
	}
	return false
}

// сконвертировать буффер в строку
func (b *COWBuffer) String() string {
	return unsafe.String(unsafe.SliceData(b.data), len(b.data))
}
