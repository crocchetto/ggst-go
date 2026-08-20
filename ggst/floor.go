package ggst

import "fmt"

type Floor int

// known floors
const (
	Floor1    Floor = 1
	Floor2    Floor = 2
	Floor3    Floor = 3
	Floor4    Floor = 4
	Floor5    Floor = 5
	Floor6    Floor = 6
	Floor7    Floor = 7
	Floor8    Floor = 8
	Floor9    Floor = 9
	Floor10   Floor = 10
	Celestial Floor = 99
)

const FloorAll Floor = 0

func (f Floor) IsValid() bool {
	return (f >= Floor1 && f <= Floor10) || f == Celestial
}

func (f Floor) Name() string {
	switch {
	case f == Celestial:
		return "Celestial"
	case f >= Floor1 && f <= Floor10:
		return fmt.Sprintf("Floor %d", int(f))
	default:
		return fmt.Sprintf("Unknown(%d)", int(f))
	}
}

func (f Floor) String() string {
	return f.Name()
}
