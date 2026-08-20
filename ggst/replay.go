package ggst

import "time"

type Player struct {
	ID      string
	Name    string
	string1 string
	string2 string
	int1    int
}

type Replay struct {
	Floor       Floor
	Player1Char Character
	Player2Char Character
	Player1     Player
	Player2     Player
	Winner      int // 1 = player 1 won, 2 = player 2 won
	Date        time.Time
	Views       uint64
	Likes       uint64
	int1        uint64
	int2        int
	int7        int
	int8        int
}

func (r Replay) Player1Won() bool { return r.Winner == 1 }

func (r Replay) WinnerPlayer() (p Player, char Character, ok bool) {
	switch r.Winner {
	case 1:
		return r.Player1, r.Player1Char, true
	case 2:
		return r.Player2, r.Player2Char, true
	default:
		return Player{}, 0, false
	}
}
