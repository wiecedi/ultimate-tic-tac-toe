package main 


type Action int64

type GameState struct {
    nextPlayer          int
    
    grid_pieces         BitArray81
    grid_player1        BitArray81
    bigGrid_pieces      BitArray9
    bigGrid_player1     BitArray9
    bigGrid_player2     BitArray9
    
    lastAction          Action
    
    ended               bool
    result              int
    heuristic           float64
}

var miniBoards = [9]BitArray81 {}

var checkWinBigBoard = [8]BitArray9 { {0x54}, {0x111}, {0x124}, {0x92}, {0x49}, {0x1C0}, {0x38}, {0x7} }

var checkWinSubBoard = [9][8]BitArray81 {}


func toAction(row, col int) Action {
    if row < 0 {
        return Action(-1)
    }
    return Action(9 * row + col)
}

func (ac Action) getRowCol() (int, int) {
    if ac == -1 || ac == 0 {
        return int(ac), int(ac)
    }
    return int(ac / 9), int(ac % 9)
}

func (ac Action) getBigN() int {
    if ac == -1 || ac == 0 {
        return int(ac)
    }
    return int(3 * (ac / 27) + ac / 3 % 3)
}

func (ac Action) getNextN() int {
    if ac == -1 || ac == 0 {
        return int(ac)
    }
    return int(3 * (ac / 9 % 3) + ac % 3)
}



func initGame() (GameState, Action) {
    for i, zz := range[]Action {toAction(0, 0), toAction(0, 3), toAction(0, 6), toAction(3, 0), toAction(3, 3), toAction(3, 6), toAction(6, 0), toAction(6, 3), toAction(6, 6)} {
        for r := 0; r < 3; r++ {
            for c := 0; c < 3; c++ {
                miniBoards[i].SetBit(uint64(zz + toAction(r, c)))
            }
        }
    }
    
    for i, zz := range[]Action {toAction(0, 0), toAction(0, 3), toAction(0, 6), toAction(3, 0), toAction(3, 3), toAction(3, 6), toAction(6, 0), toAction(6, 3), toAction(6, 6)} {
        for j, wins := range[][]Action {
            {toAction(0, 0), toAction(0, 1), toAction(0, 2)},
            {toAction(1, 0), toAction(1, 1), toAction(1, 2)},
            {toAction(2, 0), toAction(2, 1), toAction(2, 2)},
            {toAction(0, 0), toAction(1, 0), toAction(2, 0)},
            {toAction(0, 1), toAction(1, 1), toAction(2, 1)},
            {toAction(0, 2), toAction(1, 2), toAction(2, 2)},
            {toAction(0, 0), toAction(1, 1), toAction(2, 2)},
            {toAction(0, 2), toAction(1, 1), toAction(2, 0)},
        } {
            for _, a := range wins {
                 checkWinSubBoard[i][j].SetBit(uint64(zz + a))
            }
        }
    }
    
    return GameState {
        nextPlayer: 1,
        lastAction: -1,
    }, -1
}

func (gameState *GameState) getHeuristicScore() float64 {
    return float64(-gameState.nextPlayer * (gameState.bigGrid_player1.OnesCount() - gameState.bigGrid_player2.OnesCount())) / 6.0
}

func (gameState *GameState) play(action Action) {
    gameState.lastAction = action
    gameState.grid_pieces.SetBit(uint64(action))
    if gameState.nextPlayer == 1 {
        gameState.grid_player1.SetBit(uint64(action))
    }
    
    gameState.evaluateGame()
    gameState.nextPlayer *= -1
    gameState.heuristic = gameState.getHeuristicScore()
}

func (gameState *GameState) evaluateGame() {
    sub_winner := gameState.checkWinSubGrid()
    if(sub_winner != 0) {
        n := gameState.lastAction.getBigN()
        gameState.bigGrid_pieces.SetBit(uint64(n))
        if sub_winner == 1 {
            gameState.bigGrid_player1.SetBit(uint64(n))
        } else if sub_winner == -1 {
            gameState.bigGrid_player2.SetBit(uint64(n))
        }
        
        winner := gameState.checkWinBigGrid()
        if(winner != 0) {
            gameState.ended = true
            if(winner != 2) {
                gameState.result = winner
            }
        }
    }
}

func (gameState *GameState) getActions() BitArray81 {
    if gameState.lastAction != -1 {
        n := gameState.lastAction.getNextN()
        if !gameState.bigGrid_pieces.GetBit(uint64(n)) {
            return miniBoards[n].AndNot(gameState.grid_pieces)
        }
    }
    b := gameState.grid_pieces.Not()
    for i := 0; i < 9; i++ {
        if gameState.bigGrid_pieces.GetBit(uint64(i)) {
            b = b.AndNot(miniBoards[i])
        }
    }
    return b
}

func (gameState *GameState) getRandomAction(random *Random) Action {
    actions := gameState.getActions()
    r := random.Get(actions.OnesCount())
    return Action(actions.getNthPos(r))
}

func (gameState *GameState) simulate(random *Random) int {
    for !gameState.ended {
        action := gameState.getRandomAction(random)
        gameState.lastAction = action
        
        gameState.grid_pieces.SetBit(uint64(action))
        if gameState.nextPlayer == 1 {
            gameState.grid_player1.SetBit(uint64(action))
        }
    
        gameState.evaluateGame()
        gameState.nextPlayer *= -1
    }
    return gameState.result
}

func (s *GameState) checkWinBigGrid() (winner int) {
    var b BitArray9
    if s.nextPlayer == 1 {
        b = s.bigGrid_player1
    } else {
        b = s.bigGrid_player2
    }
    for i := 0; i < 8; i++ {
        if b.And(checkWinBigBoard[i]).Equals(checkWinBigBoard[i]) {
            return s.nextPlayer
        }
    }
    
    if s.bigGrid_pieces.IsFull() {
        if s.bigGrid_player1.OnesCount() > s.bigGrid_player2.OnesCount() {
            return 1
        }
        
        if s.bigGrid_player1.OnesCount() < s.bigGrid_player2.OnesCount() {
            return -1
        }
        
        return 2
    }
    return 0
}

func (s *GameState) checkWinSubGrid() (winner int) {
    var b BitArray81
    if s.nextPlayer == 1 {
        b = s.grid_player1
    } else {
        b = s.grid_pieces.AndNot(s.grid_player1)
    }
    n := s.lastAction.getBigN()
    for i := 0; i < 8; i++ {
        if b.And(checkWinSubBoard[n][i]).Equals(checkWinSubBoard[n][i]) {
            return s.nextPlayer
        }
    }
    if s.grid_pieces.And(miniBoards[n]).Equals(miniBoards[n]) {
        return 2
    }
    return 0
}
