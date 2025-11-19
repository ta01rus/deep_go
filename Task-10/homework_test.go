package main

import (
	"reflect"
	"slices"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type mem struct {
	point unsafe.Pointer
	index int
}

// go test -v homework_test.go

func Defragment(memory []byte, objects ...[]unsafe.Pointer) {

	pointers := make(map[unsafe.Pointer][]unsafe.Pointer)

	for i := 0; i < len(objects); i++ {
		objSlice := objects[i]
		for i := 0; i < len(objSlice); i++ {
			ptr := objSlice[i]
			list, ok := pointers[ptr]
			if !ok {
				list = []unsafe.Pointer{}
			}
			list = append(list, unsafe.Pointer(&objSlice[i]))
			pointers[ptr] = list
		}
	}

	freeMem := make([]mem, 0, len(memory))
	for i := 0; i < len(memory); i++ {
		p := unsafe.Pointer(&memory[i])
		v := memory[i]

		if v == 0 {
			freeMem = append(freeMem, mem{
				point: p,
				index: i,
			})
		} else {
			// наверно нужна оценка(размер) самого объекта
			// и перетаскивать его в свободное место цельно (не стал Вы уж простите)
			if len(freeMem) > 0 {
				m := freeMem[0]
				memory[m.index], memory[i] = memory[i], 0

				links, ok := pointers[p]
				if ok && len(links) > 0 {
					for i := 0; i < len(links); i++ {
						pl := (*unsafe.Pointer)(links[i])
						*pl = unsafe.Pointer(&memory[m.index])
					}
				}

				pointers[p] = nil
				freeMem = slices.Delete(freeMem, 0, 1)

				// освободили добавим
				freeMem = append(freeMem, mem{
					point: p,
					index: i,
				})
			}
		}
	}
}

func TestDefragmentation(t *testing.T) {
	var fragmentedMemory = []byte{
		0xFF, 0x00, 0x00, 0x00,
		0x00, 0xFF, 0x00, 0x00,
		0x00, 0x00, 0xFF, 0x00,
		0x00, 0x00, 0x00, 0xFF,
	}

	var fragmentedPointers = []unsafe.Pointer{
		unsafe.Pointer(&fragmentedMemory[0]),
		unsafe.Pointer(&fragmentedMemory[5]),
		unsafe.Pointer(&fragmentedMemory[10]),
		unsafe.Pointer(&fragmentedMemory[15]),
	}

	var defragmentedPointers = []unsafe.Pointer{
		unsafe.Pointer(&fragmentedMemory[0]),
		unsafe.Pointer(&fragmentedMemory[1]),
		unsafe.Pointer(&fragmentedMemory[2]),
		unsafe.Pointer(&fragmentedMemory[3]),
	}

	var defragmentedMemory = []byte{
		0xFF, 0xFF, 0xFF, 0xFF,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}

	Defragment(fragmentedMemory, fragmentedPointers)
	assert.True(t, reflect.DeepEqual(defragmentedMemory, fragmentedMemory))
	assert.True(t, reflect.DeepEqual(defragmentedPointers, fragmentedPointers))
}
