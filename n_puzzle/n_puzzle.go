package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var Size int
var GoalState []int

func readFile(name string) ([]int, int, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, 0, errors.New("No such file or directory")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	scanner.Scan()

	for strings.Trim(scanner.Text(), "\n \t")[0] == '#' {
		scanner.Scan()
	}
	size, err := strconv.Atoi(strings.Trim(scanner.Text(), "\n \t"))
	if err != nil {
		return nil, 0, err
	}
	if size < 2 {
		return nil, 0, errors.New("Puzzle size must be at least 2")
	}
	initialState := make([]int, size*size)
	i := 0
	for scanner.Scan() {
		if strings.Trim(scanner.Text(), "\n \t")[0] == '#' {
			continue
		}
		if i == size {
			return nil, 0, errors.New("Too many rows")
		}

		parts := strings.Split(strings.Trim(scanner.Text(), "\n \t"), " ")
		for idx := 0; idx < len(parts); idx++ {
			if len(strings.Trim(parts[idx], " ")) == 0 {
				copy(parts[idx:], parts[idx+1:])
				parts = parts[:len(parts)-1]
				idx--
			}
		}
		if len(parts) < size {
			return nil, 0, errors.New("Row too short")
		}
		if len(parts) > size && strings.Trim(parts[size], "\n \t")[0] != '#' {
			return nil, 0, errors.New("Row too long")
		}

		for j := 0; j < size; j++ {
			tmp, err := strconv.Atoi(parts[j])
			if err != nil {
				return nil, 0, err
			}
			if tmp >= size*size || tmp < 0 {
				return nil, 0, errors.New("One of the values is too large")
			}
			initialState[i*size+j] = tmp
		}
		i++
	}
	if i < size {
		return nil, 0, errors.New("Not enough rows")
	}

	for i := 0; i < len(initialState); i++ {
		for j := i + 1; j < len(initialState); j++ {
			if initialState[i] == initialState[j] {
				return nil, 0, errors.New("Duplicate number")
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, 0, err
	}
	return initialState, size, nil
}

func main() {
	var input []int
	var size int
	var err error

	hFlag := flag.String("H", "m", "Pick one of the following heuristics:\nh	Hamming distance\nm	Manhattan distance\nl	Linear conflict + Manhattan\nc	Corner tiles + Linear conflict + Manhattan\nu\tUniform cost search (No heuristic)")
	gFlag := flag.String("g", "a", "Pick one of the following algorithms:\ng\tGreedy bfs search algorithm\na\tA* best search algorithm")
	fFlag := flag.String("f", "", "File to read input from")
	sFlag := flag.Int("s", 0, "Size of the map to generate")
	iFlag := flag.Int("i", 100, "Number of random swap done during map generation")
	flag.Parse()
	if len(*fFlag) != 0 && *sFlag != 0 {
		fmt.Printf("\033[1;31mYou can only use one of the s and f flags\033[m\n")
		return
	}
	if len(*fFlag) == 0 && *sFlag == 0 {
		fmt.Printf("\033[1;31mPlease specify a file or the size of the board to generate\033[m\n")
		return
	}
	if len(*fFlag) != 0 {
		input, size, err = readFile(*fFlag)
		if err != nil {
			fmt.Printf("\033[1;31m%s\033[m\n", err)
			return
		}
	} else if *sFlag < 2 {
		fmt.Printf("\033[1;31mSize must be at least 2\033[m\n")
		return
	} else {
		size = *sFlag
		input = Shuffle(size, *iFlag, Generator(*sFlag))
	}
	Resolve(size, input, *hFlag, *gFlag)
	return
}
