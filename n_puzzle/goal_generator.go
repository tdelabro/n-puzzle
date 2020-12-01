package main

import (
	"math/rand"
	"time"
)

func Generator(size int) []int {
	goal := make([]int, size*size, size*size)
	cur := 1
	x := 0
	ix := 1
	y := 0
	iy := 0
	for true {
		goal[x+y*size] = cur
		cur++
		if cur == size*size {
			break
		}
		if x+ix == size || x+ix < 0 || (ix != 0 && goal[x+ix+y*size] != 0) {
			iy = ix
			ix = 0
		} else if y+iy == size || y+iy < 0 || (iy != 0 && goal[x+(y+iy)*size] != 0) {
			ix = -iy
			iy = 0
		}
		x += ix
		y += iy
	}
	return goal
}

func Shuffle(size int, iteration int, state []int) []int {
	rand.Seed(time.Now().UTC().UnixNano())
	for iteration > 0 {
		idx := FindPos(0, state)
		poss := make([]int, 0, 4)
		if idx%size > 0 {
			poss = append(poss, idx-1)
		}
		if idx%size < size-1 {
			poss = append(poss, idx+1)
		}
		if idx/size > 0 {
			poss = append(poss, idx-size)
		}
		if idx/size < size-1 {
			poss = append(poss, idx+size)
		}
		swi := poss[rand.Intn(len(poss))]
		state[idx] = state[swi]
		state[swi] = 0
		iteration--
	}
	return state
}
