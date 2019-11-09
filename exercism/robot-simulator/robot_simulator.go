package robot

import (
	"fmt"
)

func Right() {
	Step1Robot.Dir.Turn(R)
}

func Left() {
	Step1Robot.Dir.Turn(L)
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

const (
	I Command = iota // ignore
	L
	A //exploits implicit repetition of the last non-empty expression list
	R
)

type Action byte

const (
	II Action = iota // rotate left
	LL               // rotate left
	AA               // forward
	RR               // rotate right
)

type FreeSpaces map[Dir]bool

// attach method to a type
func (dx Dir) String() string {
	return map[Dir]string{N: "north", S: "south", E: "east", W: "west"}[dx]
}

func (cmd Command) String() string {
	return map[Command]string{L: "left", A: "advance", R: "right"}[cmd]
}

func (dx *Dir) Turn(way Command) {
	clockwise := map[Dir]Dir{N: E, E: S, S: W, W: N}
	anticlock := map[Dir]Dir{N: W, W: S, S: E, E: N}
	switch {
	case way == R:
		*dx = clockwise[*dx]
	case way == L:
		*dx = anticlock[*dx]
	}
}

// advance 1 step in current direction
func (rob *Step2Robot) Advance(extent FreeSpaces) {
	switch rob.Dir {
	case N:
		rob.Pos.Northing += 1
	case E:
		rob.Pos.Easting += 1
	case S:
		rob.Pos.Northing -= 1
	case W:
		rob.Pos.Easting -= 1
	default:
		fmt.Println("bad direction")
	}
}

func (rob *Step2Robot) Turn(cmd Command) {
	fmt.Println("gothere")
	rob.Dir.Turn(cmd)
}

func (rob *Step2Robot) Obey(cmd Command, extent FreeSpaces) {
	switch cmd {
	case A:
		rob.Advance(extent)
	case I:
		// ignore it
	default: // L or R
		rob.Turn(cmd)
	}
}

func (r Rect)Inside(pos Pos)bool{
	return true
}

func StartRobot(cmd chan Command, act chan Action) {
	for {
		what, ok := <-cmd
		fmt.Println("inpgot", what, ok)
		if !ok {
			fmt.Println("channel closing")
			break
		}
		fmt.Println("got", what, ok)
		act <- map[Command]Action{I: II, L: LL, A: AA, R: RR}[what]
	}
}

func Room(extent Rect, robot Step2Robot, act chan Action, rep chan Step2Robot) {
	fmt.Println(extent, robot)
	avail := FreeSpaces{N:true, E:true, S:true, W:true}
	for {
		what, ok := <-act
		if !ok {
			fmt.Println("channel closing")
			rep <- robot
			break
		}
		fmt.Println("got", what, ok)
		robot.Obey(map[Action]Command{II: I, LL: L, AA: A, RR: R}[what], avail)
	}
}
