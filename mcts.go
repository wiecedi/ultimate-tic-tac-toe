package main

import (
    "fmt"
    "time"
    "math"
)


type Head struct {
    root            *Node
    random          *Random
    limit           time.Time
}

type Node struct {
    parent          *Node
    causingAction   Action
 
    gameState       GameState
   
    children        []*Node
    untriedActions  []Action
    
    opening         *Node
    openinghash     uint64
    
    n               int
    q               int
    draws           int
    wins            int
    losses          int
}


func (node *Node) fullyExpand() {
    for i := 0; i < len(node.untriedActions); i++ {
        action := node.untriedActions[i]
        g := node.gameState
        g.play(action)
        child := createNode(node, g, action)
        node.children = append(node.children, child)
    }
    node.untriedActions = nil
}

func createHead() *Head {
    gameState, causingAction := initGame()
    root := createNode(&Node{}, gameState, causingAction)
    head := &Head { root: root, random: &Random{} }
    root.checkOpeningBook()
    root.fullyExpand()
    return head
}

func createNode(parent *Node, gameState GameState, causingAction Action) *Node {
    node := &Node {
        parent:            parent,
        causingAction:    causingAction,
        gameState:        gameState,
        children:         make([]*Node, 0),
        untriedActions:    gameState.getActions().ToList(),
    }
    node.checkOpeningBook()
    return node
}

func mcts(head *Head) {
    for time.Until(head.limit) > 0 {
        leaf := head.root.treePolicy(head)
        result := leaf.rollout(head.random)
        leaf.backpropagate(head, result)
    }
}

func printMove(head *Head) {
    var maxN int = -1
    var bestNode *Node
    
    for i := 0; i < len(head.root.children); i++ {
        if maxN < head.root.children[i].n {
            maxN = head.root.children[i].n
            bestNode = head.root.children[i]
        }
    }
    row, col := bestNode.causingAction.getRowCol()
    fmt.Println(row, col)

      bestNode.parent.parent = nil
    head.root = bestNode
}

func (node *Node) isTerminal() bool {
    return node.gameState.ended
}

func (node *Node) addChildNode(random *Random) *Node {
    i := random.Get(len(node.untriedActions))
    action := node.untriedActions[i]
    node.untriedActions = append(node.untriedActions[:i], node.untriedActions[i+1:]...)
    g := node.gameState
    g.play(action)
    child := createNode(node, g, action)
    node.children = append(node.children, child)
    return child
}

func (node *Node) getBestUCTChild(c float64) *Node {
    var maxChild *Node
    maxValue := -math.MaxFloat64
    var logN float64 = math.Log2(float64(node.n))
    
    for _, child := range node.children {
        if(child.n == 0) {
            return child
        }
        value := 0.1/float64(len(child.children) + len(child.untriedActions)) + 0.3*node.gameState.heuristic + float64(child.q) / float64(child.n) + c * math.Sqrt(logN / float64(child.n))
       
        if(value > maxValue) {
            maxValue = value
            maxChild = child
        }
    }
    return maxChild
}


func (node *Node) treePolicy(head *Head) *Node {
    for !node.isTerminal() {
        if len(node.untriedActions) != 0 {
            return node.addChildNode(head.random)
        }
        node = node.getBestUCTChild(1.2)
    }
    return node
}

func (node *Node) rollout(random *Random) int {
    clone := node.gameState
    return clone.simulate(random)
}

func (node *Node) backpropagate(head *Head, result int) {
    for node != nil && node != head.root.parent {
        node.n++
        node.q += result * (-node.gameState.nextPlayer)
        switch(result) {
            case node.gameState.nextPlayer:
                node.losses++
            case -node.gameState.nextPlayer:
                node.wins++
            default:
                node.draws++
        }
        node = node.parent
    }
}
