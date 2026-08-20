package ggst

import "fmt"

type Character int

// known character IDs
const (
	Sol        Character = 0
	Ky         Character = 1
	May        Character = 2
	Axl        Character = 3
	Chipp      Character = 4
	Potemkin   Character = 5
	Faust      Character = 6
	Millia     Character = 7
	Zato       Character = 8
	Ramlethal  Character = 9
	Leo        Character = 10
	Nagoriyuki Character = 11
	Giovanna   Character = 12
	Anji       Character = 13
	Ino        Character = 14
	Goldlewis  Character = 15
	JackO      Character = 16
	HappyChaos Character = 17
	Baiken     Character = 18
	Testament  Character = 19
	Bridget    Character = 20
	Sin        Character = 21
	Bedman     Character = 22
	Asuka      Character = 23
	Johnny     Character = 24
	Elphelt    Character = 25
	ABA        Character = 26
	Slayer     Character = 27
	Dizzy      Character = 28
	Venom      Character = 29
	Unika      Character = 30
	Lucy       Character = 31
	Jam        Character = 32
	RoboKy     Character = 33
)

var characterNames = []string{
	"Sol",
	"Ky",
	"May",
	"Axl",
	"Chipp",
	"Potemkin",
	"Faust",
	"Millia",
	"Zato-1",
	"Ramlethal",
	"Leo",
	"Nagoriyuki",
	"Giovanna",
	"Anji",
	"I-No",
	"Goldlewis",
	"Jack-O'",
	"Happy Chaos",
	"Baiken",
	"Testament",
	"Bridget",
	"Sin",
	"Bedman?",
	"Asuka",
	"Johnny",
	"Elphelt",
	"A.B.A.",
	"Slayer",
	"Dizzy",
	"Venom",
	"Unika",
	"Lucy",
	"Jam",
	"Robo-Ky",
}

var characterStatPrefixes = []string{
	"SOL", "KYK", "MAY", "AXL", "CHP", "POT", "FAU", "MLL", "ZAT", "RAM",
	"LEO", "NAG", "GIO", "ANJ", "INO", "GLD", "JKO", "COS", "BKN", "TST",
	"BGT", "SIN", "BED", "ASK", "JHN", "ELP", "ABA", "SLY", "DZY", "VEN",
	"UNI", "LUC", "USG", "ABS",
}

func (c Character) statPrefix() string {
	if !c.IsValid() {
		return ""
	}
	return characterStatPrefixes[c]
}

func (c Character) IsValid() bool {
	return c >= 0 && int(c) < len(characterNames)
}

func (c Character) Name() string {
	if !c.IsValid() {
		return fmt.Sprintf("Unknown(%d)", int(c))
	}
	return characterNames[c]
}

func (c Character) String() string {
	return c.Name()
}
