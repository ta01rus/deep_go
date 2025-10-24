package main

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

func TestCOWBuffer(t *testing.T) {
	data := []byte{'a', 'b', 'c', 'd'}
	buffer := NewCOWBuffer(data)
	defer buffer.Close()

	copy1 := buffer.Clone()
	copy2 := buffer.Clone()

	assert.Equal(t, unsafe.SliceData(data), unsafe.SliceData(buffer.data), "check 1")
	assert.Equal(t, unsafe.SliceData(buffer.data), unsafe.SliceData(copy1.data), "check 2")
	assert.Equal(t, unsafe.SliceData(copy1.data), unsafe.SliceData(copy2.data), "check 3")

	assert.True(t, (*byte)(unsafe.SliceData(data)) == unsafe.StringData(buffer.String()), "check 4")
	assert.True(t, (*byte)(unsafe.StringData(buffer.String())) == unsafe.StringData(copy1.String()), "check 5")
	assert.True(t, (*byte)(unsafe.StringData(copy1.String())) == unsafe.StringData(copy2.String()), "check 6")

	assert.True(t, buffer.Update(0, 'g'), "check 7")
	assert.False(t, buffer.Update(-1, 'g'), "check 8")
	assert.False(t, buffer.Update(4, 'g'), "check 9")

	assert.True(t, reflect.DeepEqual([]byte{'g', 'b', 'c', 'd'}, buffer.data), "check 10")
	assert.True(t, reflect.DeepEqual([]byte{'a', 'b', 'c', 'd'}, copy1.data), "check 11")
	assert.True(t, reflect.DeepEqual([]byte{'a', 'b', 'c', 'd'}, copy2.data), "check 12")

	assert.NotEqual(t, unsafe.SliceData(buffer.data), unsafe.SliceData(copy1.data), "check 13")
	assert.Equal(t, unsafe.SliceData(copy1.data), unsafe.SliceData(copy2.data), "check 14")

	copy1.Close()

	previous := copy2.data
	copy2.Update(0, 'f')
	current := copy2.data

	// 1 reference - don't need to copy buffer during update
	assert.Equal(t, unsafe.SliceData(previous), unsafe.SliceData(current), "check 15")

	copy2.Close()
}
