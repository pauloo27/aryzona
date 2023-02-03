package dice

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

// reference: https://en.wikipedia.org/wiki/Dice_notation

type DiceNotation struct {
	Dices, Faces int
}

var (
	reNotation = regexp.MustCompile(`^(\d+)$|^(\d*)d(\d*)$`)

	ErrInvalidNotation = errors.New("invalid notation")
	DefaultDice        = &DiceNotation{Dices: 1, Faces: 6}
)

func parseMatch(match string, defaultValue int) (int, error) {
	if match == "" {
		return defaultValue, nil
	}
	return strconv.Atoi(match)
}

func (d *DiceNotation) String() string {
	return fmt.Sprintf("%dd%d", d.Dices, d.Faces)
}

func ParseNotation(str string) (*DiceNotation, error) {
	matches := reNotation.FindStringSubmatch(str)
	if matches == nil {
		return nil, ErrInvalidNotation
	}

	facesIndex := 3

	if matches[1] != "" {
		facesIndex = 1
	}
	dicesIndex := 2

	dices, err := parseMatch(matches[dicesIndex], 1)
	if err != nil {
		return nil, ErrInvalidNotation
	}

	faces, err := parseMatch(matches[facesIndex], 6)
	if err != nil {
		return nil, ErrInvalidNotation
	}

	if faces <= 0 || dices <= 0 {
		return nil, ErrInvalidNotation
	}

	return &DiceNotation{
		Dices: dices,
		Faces: faces,
	}, nil
}
