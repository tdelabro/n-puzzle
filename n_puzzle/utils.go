package main

func FindPos(v int, state []int) int {
	for i := 0; i < len(state); i++ {
		if state[i] == v {
			return i
		}
	}
	return -1
}
