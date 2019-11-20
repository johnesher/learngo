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
	fmt.Println("obeying", cmd.String(), avail)
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
		panic("unknown command in obey")
	}
	return retval
}

func (r Rect) Inside(pos Pos) bool {
	return true
}

func StartRobot(cmd chan Command, act chan Action) {
	for {
		what, ok := <-cmd
		// fmt.Println("StartRobot got", what, ok)
		if !ok {
			fmt.Println("command channel closing")
			close(act)
			break
		}
		fmt.Println("StartRobot sendng", what)
		to_send, ok := map[Command]Action{I: II, L: LL, A: AA, R: RR}[what]
		if !ok{
			panic("unrecognised command")
		}
		act <- to_send
	}
}

func Room(extent Rect, robot Step2Robot, act chan Action, rep chan Step2Robot) {
	fmt.Println("Room", extent, robot)
	for {
		avail := FreeSpaces{N: true, E: true, S: true, W: true}
		// complicated bit - establish if the robot can move
		avail[N] = extent.Max.Northing != robot.Pos.Northing
		avail[S] = extent.Min.Northing != robot.Pos.Northing
		avail[E] = extent.Max.Easting != robot.Pos.Easting
		avail[W] = extent.Min.Easting != robot.Pos.Easting
		fmt.Println("Romm avail", avail, extent, robot.Pos)
		what, ok := <-act
		fmt.Println("Room got", what, ok)
		if !ok {
			fmt.Println("action channel closing in Room", robot)
			rep <- robot
			break
		}
		cmd, ok := map[Action]Command{II: I, LL: L, AA: A, RR: R}[what]
		if !ok{
			panic("unrecognised action")
		}
		fmt.Printf("trying to obey %T, [%v]\n", cmd, cmd )
		robot.Obey(cmd, avail)
	}
}

type Action3 struct{
	action rune
	name string
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
		fmt.Println("sending", a)
		action <- Action3{a, name}
	}
	// close(action)
	action <- Action3{Q3, "quitter"}
}

func Room3(
	extent Rect,
	robots []Step3Robot,
	action chan Action3,
	report chan []Step3Robot,
	log chan string) {
	var names = make(map[string]int)
	for _, rob := range robots{
		if "" == rob.Name {
			log <-"Room3 robot without a name"
		}
		_, ok := names[rob.Name]
		fmt.Println("Room3 names", ok, rob.Name, names)
		if ok {
			log <- "duplicate name"
		}else{
			names[rob.Name] += 1
		}
	}
	robot := &robots[0]
	for {
		what, ok := <-action
		// if !ok{
		// 	report <-robots
		// }
		fmt.Println("Room3 got", what, ok)
		if !ok {
			report <-robots
			fmt.Println("action channel closed in Room3", robots)
			break
		}
		if Q3 == what.action {
			report <-robots
			fmt.Println("action channel got quit in Room3", robots)
			break
		}
		avail := FreeSpaces{N: true, E: true, S: true, W: true}
		avail[N] = extent.Max.Northing != robot.Pos.Northing
		avail[S] = extent.Min.Northing != robot.Pos.Northing
		avail[E] = extent.Max.Easting != robot.Pos.Easting
		avail[W] = extent.Min.Easting != robot.Pos.Easting
		fmt.Println("Romm avail", avail, extent, robot.Pos)
		cmd, ok := map[rune]Command{I3: I, L3: L, A3: A, R3: R}[what.action]
		if !ok{
			panic("unrecognised action")
		}
		fmt.Printf("trying to obey %T, [%v]\n", cmd, cmd )
		ok = robot.Obey(cmd, avail)
		if !ok{
			log <-"tried to hit wall"
		}
		fmt.Printf("after obey %v", robot)
	}
}
