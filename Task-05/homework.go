package main

func Map(data []int, action func(int) int) []int {
	var (
		l = len(data)
	)
	if l == 0 {
		return data
	}

	for i := 0; i < len(data); i++ {
		data[i] = action(data[i])
	}

	return data
}

func Filter(data []int, action func(int) bool) []int {
	l := len(data)

	if l == 0 {
		return data
	}

	result := make([]int, 0, l)

	for i := 0; i < l; i++ {
		if action(data[i]) {
			result = append(result, data[i])
		}
	}

	return result
}

func Reduce(data []int, initial int, action func(int, int) int) int {

	l := len(data)

	if l == 0 {
		return 0
	}

	initial = action(initial, data[0])
	for i := 1; i < l; i++ {
		initial = action(initial, data[i])
	}
	return initial
}
