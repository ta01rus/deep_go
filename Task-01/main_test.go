package main

import "testing"

func Test_ConvertInt32(t *testing.T) {
	tests := []struct {
		name string
		in   int32
		want uint32
	}{
		{
			"test 1",
			0x12345600,
			0x56341200,
		},

		{
			"test 2",
			0xFFF00,
			0x00FFF,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Convert(tt.in)
			if got != tt.in {
				t.Errorf("convert() = %v, want %v", got, tt.want)
			}
		})
	}
}
