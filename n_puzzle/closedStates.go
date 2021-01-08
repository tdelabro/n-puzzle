package main

type node struct {
    value int
    child []*node
}

func (n *node) insertState(state []int) {
    current := n
    for i := 0; i < len(state); i++ {
        if pNode := findNodeContainingValInChilds(state[i], current.child); pNode != nil {
            current = pNode
        } else {
            tmp := &node{
                value: state[i],
                child: make([]*node, 0),
            }
            current.child = append(current.child, tmp)
            current = tmp
        }
    }
}

func findNodeContainingValInChilds(val int, childs []*node) *node {
    for i := 0; i < len(childs); i++ {
        if (*childs[i]).value == val {
            return childs[i]
        }
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