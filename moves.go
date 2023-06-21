package main

import "fmt"

func calculateLegalMoves(board *Board) {
	// Iterate over each tile on the board
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			tile := &board.pos[y][x]
			piece := &tile.piece

			switch piece.pieceType {
			case Pawn:
				modifier := 1
				if piece.color == "white" {
					modifier = -1
				}

				// Forward move
				moveToAdd := Move{To: Coordinates{X: x, Y: y + modifier}, From: Coordinates{X: x, Y: y}}
				if canAddMove(&moveToAdd, board) {
					piece.legalMoves = append(piece.legalMoves, moveToAdd)
				}

				// Double forward move from starting position
				if (y == 1 && piece.color == "black") || (y == 6 && piece.color == "white") {
					moveToAdd := Move{To: Coordinates{X: x, Y: y + 2*modifier}, From: Coordinates{X: x, Y: y}}
					if canAddMove(&moveToAdd, board) {
						piece.legalMoves = append(piece.legalMoves, moveToAdd)
					}
				}

				// Capture moves
				captureMoves := []Coordinates{{X: x - 1, Y: y + modifier}, {X: x + 1, Y: y + modifier}}
				for _, captureMove := range captureMoves {
					if isValidPosition(captureMove.X, captureMove.Y) {
						moveToAdd := Move{To: captureMove, From: Coordinates{X: x, Y: y}}
						if canCaptureMove(&moveToAdd, board) {
							piece.legalMoves = append(piece.legalMoves, moveToAdd)
						}
					}
				}

			case Rook:
				// Calculate legal moves for rook in horizontal and vertical directions
				addHorizontalVerticalMoves(piece, x, y, board)

			case Knight:
				// Calculate legal moves for knight
				knightMoves := []Coordinates{
					{X: x + 1, Y: y + 2}, {X: x + 2, Y: y + 1},
					{X: x + 2, Y: y - 1}, {X: x + 1, Y: y - 2},
					{X: x - 1, Y: y - 2}, {X: x - 2, Y: y - 1},
					{X: x - 2, Y: y + 1}, {X: x - 1, Y: y + 2},
				}
				for _, knightMove := range knightMoves {
					if isValidPosition(knightMove.X, knightMove.Y) {
						moveToAdd := Move{To: knightMove, From: Coordinates{X: x, Y: y}}
						if canCaptureMove(&moveToAdd, board) {
							piece.legalMoves = append(piece.legalMoves, moveToAdd)
						}
					}
				}

			case Bishop:
				// Calculate legal moves for bishop in diagonal directions
				addDiagonalMoves(piece, x, y, board)

			case Queen:
				// Calculate legal moves for queen in all directions
				addHorizontalVerticalMoves(piece, x, y, board)
				addDiagonalMoves(piece, x, y, board)

			case King:
				// Calculate legal moves for king
				kingMoves := []Coordinates{
					{X: x + 1, Y: y}, {X: x + 1, Y: y + 1},
					{X: x, Y: y + 1}, {X: x - 1, Y: y + 1},
					{X: x - 1, Y: y}, {X: x - 1, Y: y - 1},
					{X: x, Y: y - 1}, {X: x + 1, Y: y - 1},
				}
				for _, kingMove := range kingMoves {
					if isValidPosition(kingMove.X, kingMove.Y) {
						moveToAdd := Move{To: kingMove, From: Coordinates{X: x, Y: y}}
						if canCaptureMove(&moveToAdd, board) {
							piece.legalMoves = append(piece.legalMoves, moveToAdd)
						}
					}
				}
			}
		}
	}
}

func movePiece(move *Move, board *Board) {
	from := move.From
	to := move.To
	tile := board.pos[from.Y][from.X]
	pieceToMove := tile.piece

	// Check if the move is valid for the piece being moved
	if isValidMove(move, &pieceToMove, board) {
		// Perform the move
		board.pos[from.Y][from.X].piece = Piece{pieceType: Empty, color: "white"}
		board.pos[to.Y][to.X].piece = pieceToMove

		// Print the updated board
		printBoard(board)
	} else {
		fmt.Println("Invalid move!")
	}
}

func isValidMove(move *Move, piece *Piece, board *Board) bool {
	// Check if the move is present in the piece's legalMoves slice
	for _, legalMove := range piece.legalMoves {
		if legalMove.To.X == move.To.X && legalMove.To.Y == move.To.Y {
			return true
		}
	}
	return false
}

func getRandomLegalMove(board *Board) *Move {
	for i := 0; i < len(board.pos); i++ {
		for j := 0; j < len(board.pos[i]); j++ {
			potentialLegalMoves := board.pos[i][j].piece.legalMoves
			if len(potentialLegalMoves) > 0 {
				return &potentialLegalMoves[0]
			}
		}
	}
	panic("Could not find any legal moves.")
}

func addHorizontalVerticalMoves(piece *Piece, x, y int, board *Board) {
	// Add legal moves in horizontal direction
	for i := x + 1; i < 8; i++ {
		moveToAdd := Move{To: Coordinates{X: i, Y: y}, From: Coordinates{X: x, Y: y}}
		if canAddMove(&moveToAdd, board) {
			piece.legalMoves = append(piece.legalMoves, moveToAdd)
		}
		if board.pos[y][i].piece.pieceType != Empty {
			break
		}
	}

	for i := x - 1; i >= 0; i-- {
		moveToAdd := Move{To: Coordinates{X: i, Y: y}, From: Coordinates{X: x, Y: y}}
		if canAddMove(&moveToAdd, board) {
			piece.legalMoves = append(piece.legalMoves, moveToAdd)
		}
		if board.pos[y][i].piece.pieceType != Empty {
			break
		}
	}

	// Add legal moves in vertical direction
	for i := y + 1; i < 8; i++ {
		moveToAdd := Move{To: Coordinates{X: x, Y: i}, From: Coordinates{X: x, Y: y}}
		if canAddMove(&moveToAdd, board) {
			piece.legalMoves = append(piece.legalMoves, moveToAdd)
		}
		if board.pos[i][x].piece.pieceType != Empty {
			break
		}
	}

	for i := y - 1; i >= 0; i-- {
		moveToAdd := Move{To: Coordinates{X: x, Y: i}, From: Coordinates{X: x, Y: y}}
		if canAddMove(&moveToAdd, board) {
			piece.legalMoves = append(piece.legalMoves, moveToAdd)
		}
		if board.pos[i][x].piece.pieceType != Empty {
			break
		}
	}
}

func addDiagonalMoves(piece *Piece, x, y int, board *Board) {
	// Add legal moves in diagonal directions
	for i, j := x+1, y+1; i < 8 && j < 8; i, j = i+1, j+1 {
		moveToAdd := Move{To: Coordinates{X: i, Y: j}, From: Coordinates{X: x, Y: y}}
		if canAddMove(&moveToAdd, board) {
			piece.legalMoves = append(piece.legalMoves, moveToAdd)
		}
		if board.pos[j][i].piece.pieceType != Empty {
			break
		}
	}

	for i, j := x-1, y+1; i >= 0 && j < 8; i, j = i-1, j+1 {
		moveToAdd := Move{To: Coordinates{X: i, Y: j}, From: Coordinates{X: x, Y: y}}
		if canAddMove(&moveToAdd, board) {
			piece.legalMoves = append(piece.legalMoves, moveToAdd)
		}
		if board.pos[j][i].piece.pieceType != Empty {
			break
		}
	}

	for i, j := x+1, y-1; i < 8 && j >= 0; i, j = i+1, j-1 {
		moveToAdd := Move{To: Coordinates{X: i, Y: j}, From: Coordinates{X: x, Y: y}}
		if canAddMove(&moveToAdd, board) {
			piece.legalMoves = append(piece.legalMoves, moveToAdd)
		}
		if board.pos[j][i].piece.pieceType != Empty {
			break
		}
	}

	for i, j := x-1, y-1; i >= 0 && j >= 0; i, j = i-1, j-1 {
		moveToAdd := Move{To: Coordinates{X: i, Y: j}, From: Coordinates{X: x, Y: y}}
		if canAddMove(&moveToAdd, board) {
			piece.legalMoves = append(piece.legalMoves, moveToAdd)
		}
		if board.pos[j][i].piece.pieceType != Empty {
			break
		}
	}
}

func canAddMove(move *Move, board *Board) bool {
	return coordinatesIsInBound(&move.To) && isEmptyTile(getTileByCoordinates(&move.To, board))
}

func coordinatesIsInBound(coordinates *Coordinates) bool {
	return coordinates.X >= 0 && coordinates.X < 8 && coordinates.Y >= 0 && coordinates.Y < 8
}

func isEmptyTile(tile *Tile) bool {
	return tile.piece.pieceType == Empty
}

func getTileByCoordinates(coordinates *Coordinates, board *Board) *Tile {
	return &board.pos[coordinates.X][coordinates.Y]
}
