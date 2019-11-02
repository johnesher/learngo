package robot

import (
	"fmt"
)

func Right() {
	Step1Robot.Dir.Turn(+1)
}

func Left() {
	Step1Robot.Dir.Turn(-1)
}

// advance 1 step in current direction
func Advance() {
	switch Step1Robot.Dir {
	case N:
		Step1Robot.Y += 1
	case E:
		Step1Robot.X += 1
	case S:
		Step1Robot.Y -= 1
	case W:
		Step1Robot.X -= 1
	default:
		fmt.Println("bad direction")
	}
}

const (
	N Dir = iota
	E     //exploits implicit repetition of the last non-empty expression list
	S
	W
)


// attach method to a type
func (dx Dir) String() string {
	m := make(map[Dir]string)
	m[N] = "north"
	m[S] = "south"
	m[E] = "east"
	m[W] = "west"
	return m[dx]
}

func (dx *Dir) Turn(way int) {
	var clockwise = map[Dir]Dir {N: E, E: S, S: W, W: N, }
	var anticlock = map[Dir]Dir {N: W, W: S, S: E, E: N, }
	switch {
	case way > 0:
		*dx = clockwise[*dx]
	case way < 0:
		*dx = anticlock[*dx]
	}
}