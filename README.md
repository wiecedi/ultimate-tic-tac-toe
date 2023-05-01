# Ultimate Tic Tac Toe AI using Monte Carlo Tree Search Algorithm
This is an AI written in golang, that competes with other AIs on CodinGame, using the Monte Carlo Tree Search algorithm to play Ultimate Tic Tac Toe. The AI is currently in the highest league, the Legend league, and is ranked in the top 5% overall.

## What is Ultimate Tic Tac Toe?
Ultimate Tic Tac Toe is a variation of the classic game Tic Tac Toe. The game is played on a 9x9 board, with each square containing a smaller 3x3 board. Players take turns placing their marker on one of the smaller boards, and the goal is to win three smaller boards, which are aligned to each other. However, the catch is that the board on which the player can place their marker is determined by the opponent's last move.

## Monte Carlo Tree Search Algorithm
The Monte Carlo Tree Search (MCTS) algorithm is a heuristic search algorithm used in decision-making processes. In the context of game-playing AIs, MCTS can be used to select the best move to make based on the current state of the game.

The algorithm works by simulating a large number of games from the current state, making random moves until a final state is reached. The results of these simulations are then used to determine the best move to make.
