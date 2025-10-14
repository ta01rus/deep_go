package main

import (
	"unsafe"
)

func convert[T int32 | int64 | int16 | int8](in T) T {

	// размер типа данных
	size := int(unsafe.Sizeof(in))

	p := unsafe.Pointer(&in)

	var result T

	for i := 0; i < size; i++ {
		b := *(*int8)(unsafe.Add(p, i))
		p := (size - i - 1) * 8
		result |= T(int(b) << p)
	}

	return result
}
