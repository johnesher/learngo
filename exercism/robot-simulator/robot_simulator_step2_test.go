// +build step2 !step1,!step3

package robot

import (
	"testing"
	// "reflect"
)

// For step 1 you implemented robot movements, but it's not much of a simulation.
// For example where in the source code is "the robot"?  Where is "the grid"?
// Where are the computations that turn robot actions into grid positions,
// in the robot, or in the grid?  The physical world is different.
//
// Step 2 introduces a "room."  It seems a small addition, but we'll make
// big changes to clarify the roles of "room", "robot", and "test program"
// and begin to clarify the physics of the simulation.  You will define Room
// and Robot as functions which the test program "brings into existence" by
// launching them as goroutines.  Information moves between test program,
// robot, and room over Go channels.
//
// Think of Room as a "physics engine," something that models and simulates
// a physical room with walls and a robot.  It should somehow model the
// coordinate space of the room, the location of the robot and the walls,
// and ensure for example that the robot doesn't walk through walls.
// We want Robot to be an agent that performs actions, but we want Room to
// maintain a coherent truth.
//
// Step 2 API:
//
// StartRobot(chan Command, chan Action)
// Room(extent Rect, robot Step2Robot, act chan Action, rep chan Step2Robot)
//
// You get to define Action; see defs.go for other definitions.
//
// The test program creates the channels and starts both Room and Robot.
// The test program then sends commands to Robot.  When it is done sending
// commands, it closes the command channel.

// Robot must accept commands and inform Room of actions it is attempting.

// When it senses the command channel closing, it must shut itself down.

// The room must interpret the physical
// consequences of the robot actions.  When it senses the robot shutting down,
// it sends a final report back to the test program, telling the robot's final
// position and direction.

var test2 = []struct {
	Command
	Step2Robot
}{
	0:  {' ', Step2Robot{N, Pos{1, 1}}}, // no command, this is the start DirAt
	1:  {'A', Step2Robot{N, Pos{1, 2}}},
	2:  {'R', Step2Robot{E, Pos{1, 2}}},
	3:  {'A', Step2Robot{E, Pos{2, 2}}},
	4:  {'L', Step2Robot{N, Pos{2, 2}}},
	5:  {'L', Step2Robot{W, Pos{2, 2}}},
	6:  {'L', Step2Robot{S, Pos{2, 2}}},
	7:  {'A', Step2Robot{S, Pos{2, 1}}},
	8:  {'R', Step2Robot{W, Pos{2, 1}}},
	9:  {'A', Step2Robot{W, Pos{1, 1}}},
	10: {'A', Step2Robot{W, Pos{1, 1}}}, // would be 0,1 but bump W wall
	11: {'L', Step2Robot{S, Pos{1, 1}}},
	12: {'A', Step2Robot{S, Pos{1, 1}}}, // bump S wall
	13: {'L', Step2Robot{E, Pos{1, 1}}},
	14: {'A', Step2Robot{E, Pos{2, 1}}},
	15: {'A', Step2Robot{E, Pos{2, 1}}}, // bump E wall
	16: {'L', Step2Robot{N, Pos{2, 1}}},
	17: {'A', Step2Robot{N, Pos{2, 2}}},
	18: {'A', Step2Robot{N, Pos{2, 2}}}, // bump N wall
}

func TestStep2(t *testing.T) {
	// run incrementally longer tests
	for i := 1; i <= len(test2); i++ {
		cmd := make(chan Command)
		act := make(chan Action)
		rep := make(chan Step2Robot)
		go StartRobot(cmd, act)
		go Room(Rect{Pos{1, 1}, Pos{2, 2}}, test2[0].Step2Robot, act, rep)
		for j := 1; j < i; j++ {
			// fmt.Println("TestStep2", test2[j].Command)
			// t.Log("TestStep2", test2[j].Command)
			cmd <- test2[j].Command
		}
		// fmt.Println("looped")
		close(cmd)
		da := <-rep
		last := i - 1
		want := test2[last].Step2Robot
		if da.Dir != want.Dir {
			t.Fatalf("Command #%d, Dir = %v, want %v", last, da.Dir, want.Dir)
		}
		if da.Pos != want.Pos {
			t.Fatalf("Command #%d, Pos = %v, want %v", last, da.Pos, want.Pos)
		}
	}
}

func TestStartRobot(t *testing.T) {
	cmd := make(chan Command)
	act := make(chan Action)
	cases := []struct {
		in   Command
		want Action
	}{
		{A, AA},
		{R, RR},
		{L, LL},
		{I, II},
	}
	defer close(cmd)
	//defer close(act)
	go StartRobot(cmd, act)
	for _, c := range cases {
		cmd <- c.in
		// must send to the act channel
		resp := <-act
		if resp != c.want {
			t.Errorf("Sent %v, got %v, wanted %v", c.in, resp, c.want)
		}
	}
}

func TestRoomNoLimits(t *testing.T) {
	cases := []struct {
		in   []Action
		want Pos
	}{
		{[]Action{AA}, Pos{1, 2}},
		{[]Action{RR, AA}, Pos{2, 1}},
		{[]Action{RR, RR, AA}, Pos{1, 0}},
		{[]Action{RR, RR, RR, AA}, Pos{0, 1}},
		{[]Action{LL, AA}, Pos{0, 1}},
		{[]Action{RR, II}, Pos{1, 1}},
	}
	// Note order of closing as close of act sends to rpt
	//defer close(rpt)
	//defer close(act)
	for _, c := range cases {
		rpt := make(chan Step2Robot)
		act := make(chan Action)
		go Room(Rect{Pos{0, 0}, Pos{2, 2}}, Step2Robot{N, Pos{1, 1}}, act, rpt)
		for _, action := range c.in {
			act <- action
		}
		close(act)  // triggers rpt
		final := <- rpt
		if final.Pos != c.want{
			t.Errorf("Sent %v, got %v, wanted %v", c.in, final.Pos, c.want)
		}
	}
}

func TestTurnRobot(t *testing.T) {
	cases := []struct {
		in   Command
		want Dir
	}{
		{R, E}, {R, S}, {R, W}, {R, N}, {L, W},
	}
	ref_pos := Pos{1, 1}
	test_rob := Step2Robot{N, ref_pos}
	for _, c := range cases {
		test_rob.Turn(c.in)
		// These two both fail - why?
		//if reflect.DeepEqual(test_rob.Pos, Pos{1,1}) {
		//if test_rob.Pos != Pos{1,1}) {
		if test_rob.Pos != ref_pos {
			t.Errorf("Sent %v, Pos = %v, want %v", c.in, test_rob.Pos, ref_pos)
		}
		if test_rob.Dir != c.want {
			t.Errorf("Sent %v, Pos = %v, want %v", c.in, test_rob.Dir, c.want)
		}
	}
}

func TestAdvanceRobot2Singles(t *testing.T) {
	// first test single advances
	cases := []struct {
		in   Dir
		want Pos
	}{{N, Pos{0, 1}}, {E, Pos{1, 0}}, {S, Pos{0, -1}}, {W, Pos{-1, 0}}}
	ref_pos := Pos{0, 0}
	for _, c := range cases {
		test_rob := Step2Robot{c.in, ref_pos}
		test_rob.Advance(FreeSpaces{N: true, E: true, S: true, W: true})
		if test_rob.Pos != c.want {
			t.Errorf("Sent %v, Pos = %v, want %v", c.in, test_rob.Pos, c.want)
		}
	}
}

func TestAdvanceRobot2Multiples(t *testing.T) {
	// test mulitple advances
	cases := []struct {
		in   Dir
		want Pos
	}{{N, Pos{0, 1}}, {E, Pos{1, 0}}, {S, Pos{0, -1}}, {W, Pos{-1, 0}}}
	ref_pos := Pos{0, 0}
	const repeats = 5
	for _, c := range cases {
		wanted := Pos{c.want.Easting * repeats, c.want.Northing * repeats}
		test_rob := Step2Robot{c.in, ref_pos}
		for i := 0; i < repeats; i++ {
			test_rob.Advance(FreeSpaces{N: true, E: true, S: true, W: true})
		}
		if test_rob.Pos != wanted {
			t.Errorf("Sent %v, Pos = %v, want %v", c.in, test_rob.Pos, wanted)
		}
	}
}

func TestAdvanceRobot2LimitsWithMovement(t *testing.T) {
	// test advance with limit so should not move
	const lv = 3
	cases := []struct {
		in   Dir
		want Pos
	}{{N, Pos{0, 1}}, {E, Pos{1, 0}}, {S, Pos{0, -1}}, {W, Pos{-1, 0}}}
	ref_pos := Pos{0, 0}
	for _, c := range cases {
		test_rob := Step2Robot{c.in, ref_pos}
		test_rob.Advance(FreeSpaces{N: true, E: true, S: true, W: true})
		if test_rob.Pos != c.want {
			t.Errorf("Sent %v, Pos = %v, avail %v", c.in, test_rob.Pos, c.want)
		}
	}
}

func TestAdvanceRobot2LimitsNoMovement(t *testing.T) {
	// test advance with limit so should move
	const lv = 3
	cases := []struct {
		in    Dir
		avail FreeSpaces
	}{
		{N, FreeSpaces{N: false, E: true, S: true, W: true}},
		{E, FreeSpaces{N: true, E: false, S: true, W: true}},
		{S, FreeSpaces{N: true, E: true, S: false, W: true}},
		{W, FreeSpaces{N: true, E: true, S: true, W: false}},
		{W, FreeSpaces{N: true, E: true, S: true}},
		{W, FreeSpaces{W: false}},
	}
	ref_pos := Pos{0, 0}
	for _, c := range cases {
		test_rob := Step2Robot{c.in, ref_pos}
		test_rob.Advance(c.avail)
		if test_rob.Pos != ref_pos {
			t.Errorf("Sent %v, Pos = %v, avail %v", c.in, test_rob.Pos, c.avail)
		}
	}
}

func not_used_TestInsideRect(t *testing.T) {
	// is pos inside rect?
	const lv = 3
	cases := []struct {
		in   Pos
		want bool
	}{
		{Pos{0, lv}, true},
		{Pos{lv, 0}, true},
		{Pos{0, -lv}, true},
		{Pos{-lv, 0}, true},
		{Pos{0, lv * 2}, false},
		{Pos{lv * 2, 0}, false},
		{Pos{0, -lv * 2}, false},
		{Pos{-lv * 2, 0}, false},
	}
	const repeats = 5
	limit_pos := Pos{lv, lv}
	test_rect := Rect{Pos{0, 0}, limit_pos}
	for _, c := range cases {
		wanted := c.want
		if test_rect.Inside(c.in) != wanted {
			t.Errorf("Sent %v, wanted %v", c.in, wanted)
		}
	}
}
