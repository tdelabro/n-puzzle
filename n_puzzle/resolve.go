package main

import (
	"fmt"
)

type position struct {
	state     []int
	cost      int
	heuristic int
	prev      *position
}

var algo string

func getHeuristic(size int, heuristic string, state []int, goalState []int) int {
	if heuristic == "h" {
		return Hamming(size, state, goalState)
	} else if heuristic == "m" {
		return Manhattan(size, state, goalState)
	} else if heuristic == "l" {
		return LinearConflict(size, state, goalState)
	} else if heuristic == "c" {
		return CornerTiles(size, state, goalState)
	} else if heuristic == "u" {
		return 0
	} else {
		return Manhattan(size, state, goalState)
	}
}

func getZeroIndex(state []int) int {
	zeroIndex := 0
	for zeroIndex < len(state) {
		if state[zeroIndex] == 0 {
			break
		}
		zeroIndex++
	}
	return zeroIndex
}

func isSameState(s1 []int, s2 []int) bool {
	for i := 0; i < len(s1); i++ {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

func insertInOpenList(openList []position, closeList node, pos position) []position {
    if closeList.stateAlreadyExist(pos.state) {
    	return openList
	}
	l := -1
	for i := len(openList) - 1; i >= 0; i-- {
		if l == -1 && pos.heuristic < openList[i].heuristic {
			l = i
		}
		if isSameState(pos.state, openList[i].state) {
			if openList[i].cost <= pos.cost {
				return openList
			}
			copy(openList[i:], openList[i+1:])
			openList = openList[:len(openList)-1]
		}
	}
	if l != len(openList) {
		l++
		return append(openList[:l], append([]position{pos}, openList[l:]...)...)
	} else {
		return append(openList, pos)
	}
}

func visitPosition(size int, heuristic string, goalState []int, pos position, openList []position, closedList node) []position {
	zeroIndex := getZeroIndex(pos.state)

	state := make([]int, len(pos.state))
	var newPos position
	if zeroIndex/size > 0 {
		copy(state, pos.state)
		state[zeroIndex] = state[zeroIndex-size]
		state[zeroIndex-size] = 0
		newPos = createPosition(size, heuristic, state, goalState, pos.cost+1, &pos)
		openList = insertInOpenList(openList, closedList, newPos)
	}
	if zeroIndex/size < size-1 {
		copy(state, pos.state)
		state[zeroIndex] = state[zeroIndex+size]
		state[zeroIndex+size] = 0
		newPos = createPosition(size, heuristic, state, goalState, pos.cost+1, &pos)
		openList = insertInOpenList(openList, closedList, newPos)
	}
	if zeroIndex%size > 0 {
		copy(state, pos.state)
		state[zeroIndex] = state[zeroIndex-1]
		state[zeroIndex-1] = 0
		newPos = createPosition(size, heuristic, state, goalState, pos.cost+1, &pos)
		openList = insertInOpenList(openList, closedList, newPos)
	}
	if zeroIndex%size < size-1 {
		copy(state, pos.state)
		state[zeroIndex] = state[zeroIndex+1]
		state[zeroIndex+1] = 0
		newPos = createPosition(size, heuristic, state, goalState, pos.cost+1, &pos)
		openList = insertInOpenList(openList, closedList, newPos)
	}
	return openList
}

func createPosition(size int, heuristic string, state []int, goalState []int, cost int, prev *position) position {
	var tmp position
	tmp.state = make([]int, len(state))
	copy(tmp.state, state)
	if (algo == "g") {
		tmp.cost = 0
		tmp.heuristic = getHeuristic(size, heuristic, state, goalState)
	} else {
		tmp.cost = cost
		tmp.heuristic = getHeuristic(size, heuristic, state, goalState) + cost
	}
	tmp.prev = prev
	return tmp
}

func rewind(size int, pos position, closedList node) int {
	var ret int
	if pos.prev != nil {
		ret = rewind(size, *pos.prev, closedList) + 1
	} else {
		ret = 0
	}
	printTaquin(size, pos.state)
	return ret
}

func printTaquin(size int, state []int) {
	for i := 0; i < len(state); i++ {
		if i != 0 && i%size == 0 {
			fmt.Printf("\n")
		}
		fmt.Printf("%2d ", state[i])
	}
	fmt.Printf("\n\n")
}

func Resolve(size int, initialState []int, heuristic string, algorithm string) {
	closedList := node{
		child: make(map[int]*node, size*size),
	}
	openList := make([]position, 0, 1024)
	goalState := Generator(size)
	var pos position
	algo = algorithm

	if !IsSolvable(size, initialState, goalState) {
		fmt.Println("Unsolvable puzzle")
		return
	}

	closedCounter := 0
	start := createPosition(size, heuristic, initialState, goalState, 0, nil)
	openList = append(openList, start)
	for len(openList) != 0 {
		pos = openList[len(openList)-1]
		openList = openList[:len(openList)-1]
		closedList.insertState(pos.state)
		closedCounter++
		if isSameState(pos.state, goalState) {
			n_moves := rewind(size, pos, closedList)
			fmt.Printf("Time complexity: %d\n", closedCounter)
			fmt.Printf("Size complexity: %d\n", closedCounter+len(openList))
			fmt.Printf("Number of moves: %d\n", n_moves)
			return
		}
		openList = visitPosition(size, heuristic, goalState, pos, openList, closedList)
	}
	fmt.Println("Unsolvable puzzle")
}
