package main

import (
	"unsafe"
)

func Convert[T int32 | uint32 | int64 | uint64 | int16 | uint16](in T) T {

	var result T
	// размер типа данных
	size := int(unsafe.Sizeof(in))

	resPointer := unsafe.Pointer(&result)
	pointerIn := unsafe.Pointer(&in)

	for i := 0; i < size; i++ {
		*(*int8)(unsafe.Add(resPointer, i)) = *(*int8)(unsafe.Add(pointerIn, i))
	}
	return result

}
