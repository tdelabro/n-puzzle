package main

func isPermutationEven(state []int, goal []int) bool {
	var s = make([]int, len(state))
	var p int
	var tmp int
	isEven := true

	copy(s, state)

	for i := 0; i < len(s)-1; i++ {
		if s[i] != goal[i] {
			p = FindPos(goal[i], s)
			tmp = s[i]
			s[i] = s[p]
			s[p] = tmp
			isEven = !isEven
		}
	}
	return isEven
}

func isEmptyEven(size int, state []int, goal []int) bool {
	i := FindPos(0, state)
	p := FindPos(0, goal)
	n := abs(p/size-i/size) + abs(p%size-i%size)
	return n%2 == 0
}

func IsSolvable(size int, state []int, goal []int) bool {
	c1 := make(chan bool)
	c2 := make(chan bool)

	go func() {
		c1 <- isPermutationEven(state, goal)
	}()
	go func() {
		c2 <- isEmptyEven(size, state, goal)
	}()
	return <-c1 == <-c2
}
