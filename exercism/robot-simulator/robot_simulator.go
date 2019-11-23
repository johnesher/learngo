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
	I Command = ' ' // ignore
	L = 'L'
	A = 'A'
	R = 'R'
)

type Action byte

const (
	II Action = iota // rotate left
	LL               // rotate left
	AA               // advance
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


func (pos Pos) String() string {
	return fmt.Sprintf("E:%d, N:%d", pos.Easting, pos.Northing)
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
func (rob *Step2Robot) Advance(avail FreeSpaces) bool {
	type Movement struct{ dx, dy RU }
	moves := map[Dir]Movement{N: {0, 1}, E: {1, 0}, S: {0, -1}, W: {-1, 0}}
	// allowed, ok := avail[rob.Dir]
	// if ok && allowed{
	if avail[rob.Dir] {
		move := moves[rob.Dir]
		rob.Pos.Northing += move.dy
		rob.Pos.Easting += move.dx
		return true
	}
	return false
}

func (rob *Step2Robot) Turn(cmd Command) {
	rob.Dir.Turn(cmd)
}


func (rob Step2Robot) String()string {
	return fmt.Sprintf("Step2Robot at %s, facing %s", rob.Pos.String(), rob.Dir.String())
}

func (rob *Step2Robot) Obey(cmd Command, avail FreeSpaces) bool {
	var retval bool = true
	switch cmd {
	case A:
		retval = rob.Advance(avail)
	case R:
		rob.Turn(cmd)
	case L:
		rob.Turn(cmd)
	case I:
		// ignore it
	default:
		// cannot panic as need to log it
		// panic("unknown command in obey")
	}
	return retval
}

func (r Rect) Inside(pos Pos) bool {
	return true
}

func StartRobot(cmd chan Command, act chan Action) {
	for {
		what, ok := <-cmd
		if !ok {
			close(act)
			break
		}
		to_send, ok := map[Command]Action{I: II, L: LL, A: AA, R: RR}[what]
		if !ok{
			panic("unrecognised command")
		}
		act <- to_send
	}
}

func Room(extent Rect, robot Step2Robot, act chan Action, rep chan Step2Robot) {
	for {
		avail := FreeSpaces{N: true, E: true, S: true, W: true}
		// complicated bit - establish if the robot can move
		avail[N] = extent.Max.Northing != robot.Pos.Northing
		avail[S] = extent.Min.Northing != robot.Pos.Northing
		avail[E] = extent.Max.Easting != robot.Pos.Easting
		avail[W] = extent.Min.Easting != robot.Pos.Easting
		what, ok := <-act
		if !ok {
			rep <- robot
			break
		}
		cmd, ok := map[Action]Command{II: I, LL: L, AA: A, RR: R}[what]
		if !ok{
			panic("unrecognised action in Room")
		}
		robot.Obey(cmd, avail)
	}
}

type Action3 struct{
	action rune
	name string
}	

func (act Action3) String() string {
	return fmt.Sprintf("[name:%s, action:%v]", act.name, string(act.action))
}

const (
	I3 rune = ' ' // ignore
	L3 = 'L'
	A3 = 'A'
	R3 = 'R'
	Q3 = 'Q'  // quit
)

func StartRobot3(name, script string, action chan Action3, log chan string){
	for _, a := range script{
		action <- Action3{a, name}
	}
	action <- Action3{Q3, name}
}

func Room3(
	extent Rect,
	robots []Step3Robot,
	action chan Action3,
	report chan []Step3Robot,
	log chan string) {
	defer close(report)
	actions := map[rune]bool{I3:true, L3:true, A3:true, R3:true, Q3:true}
	// map from name to robot
	var names = make(map[string]*Step3Robot)
	// map to use as sets to check for uniqueness
	var initial_positions = make(map[Pos]int)
	for i, rob := range robots{
		if "" == rob.Name {
			log <-"Room3 robot without a name"
		}
		_, ok := names[rob.Name]
		if ok {
			log <- "duplicate name"
		}else{
			names[rob.Name] = &robots[i]
		}
		_, ok = initial_positions[rob.Step2Robot.Pos]
		if ok {
			log <- "initial positions"
		}else{
			initial_positions[rob.Step2Robot.Pos] += 1
		}
		outside_room := false ||
			rob.Pos.Northing > extent.Max.Northing ||
			rob.Pos.Northing < extent.Min.Northing ||
			rob.Pos.Easting > extent.Max.Easting ||
			rob.Pos.Easting < extent.Min.Easting
		if outside_room {
			log <-"outside room"
		}
	}
	for {
		what, ok := <-action
		if !ok {
			report <-robots
			break
		}
		_, ok = actions[what.action]
		if !ok {
			log <- "unknown action in Room3"
			continue
		}
		robot, ok := names[what.name]
		if !ok{
			log <-"unknown robot in room3"
			break
		}
		cmd, ok := map[rune]Command{I3: I, L3: L, A3: A, R3: R, Q3: I}[what.action]
		if !ok{
			log <-"unrecognised action in Room3"
			continue
		}
		if Q3 == what.action {
			delete(names, what.name)
			if 0 == len(names){
				report <-robots
				break  // should return now
			}
		}
		avail := FreeSpaces{N: true, E: true, S: true, W: true}
		avail[N] = extent.Max.Northing != robot.Pos.Northing
		avail[S] = extent.Min.Northing != robot.Pos.Northing
		avail[E] = extent.Max.Easting != robot.Pos.Easting
		avail[W] = extent.Min.Easting != robot.Pos.Easting
		ok = robot.Obey(cmd, avail)
		if !ok{
			log <-"tried to hit wall"
		}
		
	}
}
