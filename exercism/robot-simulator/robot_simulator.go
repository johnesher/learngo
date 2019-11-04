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
	I Command = iota  // ignore
	L 
	A         //exploits implicit repetition of the last non-empty expression list
	R
)

type Action byte

const (
	// RL Action = iota // rotate left
	// F                // forward
	// RR               // rotate right
	AN Action = iota // advance North
	AE // advance East
	AS // advance South
	AW // advance West
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

func (cmd Command) String() string {
	m := make(map[Command]string)
	m[L] = "left"
	m[A] = "advance"
	m[R] = "right"
	return m[cmd]
}

func (dx *Dir) Turn(way Command) {
	var clockwise = map[Dir]Dir{N: E, E: S, S: W, W: N}
	var anticlock = map[Dir]Dir{N: W, W: S, S: E, E: N}
	switch {
	case way == R:
		*dx = clockwise[*dx]
	case way == L:
		*dx = anticlock[*dx]
	}
}


func (rob *Step2Robot)Obey(cmd Command){
	switch cmd{
	case A:
		fmt.Println("advance")
	default:  // L or R
		rob.Dir.Turn(cmd)
	}
}

func StartRobot(cmd chan Command, act chan Action) {
	for {
		what, ok := <-cmd
		fmt.Println("inpgot", what, ok)
		if !ok {
			fmt.Println("channel closing")
			break
		} else {
			fmt.Println("got", what, ok)
		}
	}
}

func Room(extent Rect, robot Step2Robot, act chan Action, rep chan Step2Robot) {
	fmt.Println(extent, robot)
	//var a2c = map[Action]Command {F:A, RR:R, RL: L}
	for {
		what, ok := <-act
		if !ok {
			fmt.Println("channel closing")
			break
		} else {
			fmt.Println("got", what, ok)
		}
	}
}
