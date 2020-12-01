package main

import (
	"sync"
)

const intSize = 32 << (^uint(0) >> 63)

func abs(n int) int {
	y := n >> (intSize - 1)
	return (n ^ y) - y
}

func Manhattan(size int, state []int, goal []int) int {
	sum := 0
	l := size * size

	c1 := make(chan int, l)
	var wg sync.WaitGroup
	wg.Add(l)

	for i := 0; i < l; i++ {
		go func(i int) {
			defer wg.Done()
			for j := 0; j < l; j++ {
				if state[i] == goal[j] {
					c1 <- abs(j/size-i/size) + abs(j%size-i%size)
					break
				}
			}
		}(i)
	}
	wg.Wait()
	for i := 0; i < l; i++ {
		sum += <-c1
	}
	return sum
}

func Hamming(size int, state []int, goal []int) int {
	sum := 0

	for i := 0; i < len(state); i++ {
		if state[i] != goal[i] {
			sum++
		}
	}
	return sum
}

type tile struct {
	index       int
	target      int
	n_conflicts int
	conflicts   []int
}

func getConflicts(tiles []tile) {
	for i := 0; i < len(tiles)-1; i++ {
		for j := i + 1; j < len(tiles); j++ {
			if (tiles[i].index < tiles[j].index && tiles[i].target > tiles[j].target) || (tiles[i].index > tiles[j].index && tiles[i].target < tiles[j].target) {
				tiles[i].n_conflicts++
				tiles[j].conflicts = append(tiles[j].conflicts, i)
			}
			if (tiles[j].index < tiles[i].index && tiles[j].target > tiles[i].target) || (tiles[j].index > tiles[i].index && tiles[j].target < tiles[i].target) {
				tiles[j].n_conflicts++
				tiles[i].conflicts = append(tiles[i].conflicts, j)
			}
		}
	}
}

func getMostConflictual(tiles []tile) int {
	max := 0
	index := 0
	for i := 0; i < len(tiles); i++ {
		if tiles[i].n_conflicts > max {
			max = tiles[i].n_conflicts
			index = i
		}
	}
	return index
}

func stillConflicts(tiles []tile) bool {
	for i := 0; i < len(tiles); i++ {
		if tiles[i].n_conflicts > 0 {
			return true
		}
	}
	return false
}

func getConflicsCount(tiles []tile) int {
	c := 0
	getConflicts(tiles)
	for stillConflicts(tiles) {
		m := getMostConflictual(tiles)
		for j := 0; j < len(tiles[m].conflicts); j++ {
			tiles[tiles[m].conflicts[j]].n_conflicts--
		}
		tiles[m].n_conflicts = 0
		c++
	}
	return c
}

func LinearConflict(size int, state []int, goal []int) int {
	rc := 0

	c1 := make(chan int)
	go func() {
		c1 <- Manhattan(size, state, goal)
	}()

	for i := 0; i < size; i++ {
		rowTiles := make([]tile, 0)
		colTiles := make([]tile, 0)
		for j := 0; j < size; j++ {
			b := false
			for k := 0; k < size; k++ {
				if state[i*size+j] == goal[i*size+k] {
					rowTiles = append(rowTiles, tile{j, k, 0, make([]int, 0)})
					if b {
						break
					} else {
						b = true
					}
				}
				if state[i+size*j] == goal[i+size*k] {
					colTiles = append(colTiles, tile{j, k, 0, make([]int, 0)})
					if b {
						break
					} else {
						b = true
					}
				}
			}
		}
		if len(rowTiles) > 1 {
			rc += getConflicsCount(rowTiles)
		}
		if len(colTiles) > 1 {
			rc += getConflicsCount(colTiles)
		}
	}
	return rc*2 + <-c1
}

func isAlreadyIn(v int, slice []int) bool {
	for i := 0; i < len(slice); i++ {
		if slice[i] == v {
			return true
		}
	}
	return false
}

func ctSize3(corners []int, size int, state []int, goal []int) int {
	tiles := make([]int, 0)

	c1 := make(chan int)
	go func() {
		c1 <- LinearConflict(size, state, goal)
	}()

	for i := 0; i < len(corners); i++ {
		if state[corners[i]] == goal[corners[i]] {
			continue
		}
		//	p := FindPos(goal[corners[i]], state)
		//	if corners[i]%size == p%size || corners[i]/size == p/size {
		//		continue
		//	}
		if corners[i]%size == 0 {
			if state[corners[i]+1] == goal[corners[i]+1] && !isAlreadyIn(corners[i]+1, tiles) {
				tiles = append(tiles, corners[i]+1)
			}
		}
		if corners[i]%size == size-1 {
			if state[corners[i]-1] == goal[corners[i]-1] && !isAlreadyIn(corners[i]-1, tiles) {
				tiles = append(tiles, corners[i]-1)
			}
		}
		if corners[i]/size == 0 {
			if state[corners[i]+size] == goal[corners[i]+size] && !isAlreadyIn(corners[i]+size, tiles) {
				tiles = append(tiles, corners[i]+size)
			}
		}
		if corners[i]/size == size-1 {
			if state[corners[i]-size] == goal[corners[i]-size] && !isAlreadyIn(corners[i]-size, tiles) {
				tiles = append(tiles, corners[i]-size)
			}
		}
	}
	return len(tiles)*2 + <-c1

}

func CornerTiles(size int, state []int, goal []int) int {
	corners := []int{0, size - 1, size*size - size, size*size - 1}
	if size == 2 {
		return LinearConflict(size, state, goal)
	} else if size == 3 {
		return ctSize3(corners, size, state, goal)
	}
	counter := 0

	c1 := make(chan int)
	go func() {
		c1 <- LinearConflict(size, state, goal)
	}()

	for i := 0; i < len(corners); i++ {
		if state[corners[i]] == goal[corners[i]] {
			continue
		}
		//	p := FindPos(goal[corners[i]], state)
		//	if corners[i]%size == p%size || corners[i]/size == p/size {
		//		continue
		//	}
		if corners[i]%size == 0 && state[corners[i]+1] == goal[corners[i]+1] {
			counter++
		}
		if corners[i]%size == size-1 && state[corners[i]-1] == goal[corners[i]-1] {
			counter++
		}
		if corners[i]/size == 0 && state[corners[i]+size] == goal[corners[i]+size] {
			counter++
		}
		if corners[i]/size == size-1 && state[corners[i]-size] == goal[corners[i]-size] {
			counter++
		}
	}
	return counter*2 + <-c1
}
