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
				for i := 1; i <= 2; i++ {
					moveToAdd := Move{To: Coordinates{X: x, Y: y + i*modifier}, From: Coordinates{X: x, Y: y}}
					if canAddMove(&moveToAdd, board) {
						piece.legalMoves = append(piece.legalMoves, moveToAdd)
					}
				}
				// Calculate legal moves for pawn
				// Add the moves to the piece's legalMoves slice
				// ...
			case Rook:
				// Calculate legal moves for rook
				// Add the moves to the piece's legalMoves slice
				// ...
			case Knight:
				// Calculate legal moves for knight
				// Add the moves to the piece's legalMoves slice
				// ...
			case Bishop:
				// Calculate legal moves for bishop
				// Add the moves to the piece's legalMoves slice
				// ...
			case Queen:
				// Calculate legal moves for queen
				// Add the moves to the piece's legalMoves slice
				// ...
			case King:
				// Calculate legal moves for king
				// Add the moves to the piece's legalMoves slice
				// ...
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
