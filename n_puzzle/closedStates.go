package main

type node struct {
    child map[int]*node
}

func (n *node) insertState(state []int) {
    current := n
    for i := 0; i < len(state); i++ {
        if pNode := findNodeContainingValInChilds(state[i], current.child); pNode != nil {
            current = pNode
        } else {
            tmp := &node{
                child: make(map[int]*node, 0),
            }
            current.child[state[i]] = tmp
            current = tmp
        }
    }
}

func findNodeContainingValInChilds(val int, childs map[int]*node) *node {
    n, ok := childs[val]
    if ok {
        return n
    }

    return nil
}

func (n *node) stateAlreadyExist(state []int) bool {
    current := n
    for i := 0; i < len(state); i++ {
        if pNode := findNodeContainingValInChilds(state[i], current.child); pNode != nil {
            current = pNode
        } else {
            return false
        }
    }
   return true
}