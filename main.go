package main

import (
    "fmt"
    "time"
)


var bot = -1
var head *Head


func readInput() (int, int) {
    var opponentRow, opponentCol, validActionCount, tmp int
    fmt.Scan(&opponentRow, &opponentCol)
    fmt.Scan(&validActionCount)

    for i := 0; i < validActionCount; i++ {
        fmt.Scan(&tmp, &tmp)
    }
    if opponentRow == -1 && opponentCol == -1 {
        bot = 1
    }
    
    return opponentRow, opponentCol
}

func setInput(opponentRow, opponentCol int) bool {
    if opponentRow == -1 {
        return false
    }
    
    fail := true
    head.root.fullyExpand()
    action := toAction(opponentRow, opponentCol)
    for i := 0; i < len(head.root.children); i++ {
        ac := head.root.children[i].causingAction
        if ac == action {
            head.root = head.root.children[i]
            fail = false
            break
        }
    }
    return fail
}

func setTimer(turn int) {
    var limit = time.Now()
    if turn < 2 {
        limit = limit.Add(time.Millisecond * 1000)
    } else {
        limit = limit.Add(time.Millisecond * 95)
    }
    head.limit = limit
}

func main() {
     head = createHead()
     for turn := 1; ; turn += 2 {
         opponentRow, opponentCol := readInput()
         if(setInput(opponentRow, opponentCol)) {
             println("falscher Zug, deine Actions:")
             for _, c := range head.root.children {
                 row, col := c.causingAction.getRowCol()
                 println(fmt.Sprintf("- %d %d", row, col))
             }
             continue
         }
        setTimer(turn)
        head.root.fullyExpand()
        mcts(head)
        printMove(head)
    }
}
