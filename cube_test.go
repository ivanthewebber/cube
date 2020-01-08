package cube

import (
	"crypto/rand"
	"testing"
)

// TestSpin tests the spinning move on a rubic's cube
func TestSpin(t *testing.T) {
	c := NewCube()

sides:
	for i := uint8(0); i < NumSides; i++ {
		for j := 0; j < 4; j++ {
			c.Spin(i)
		}
		if c != NewCube() {
			t.Logf("%s\nSide %d", "TestSpin: spinning disrupted state!", i)
			t.Fail()
			t.Log("\n" + c.String())
			c = NewCube()
			continue sides
		}
	}

}

// TestSolved tests Rubik's Cube state test
func TestSolved(t *testing.T) {
	c := NewCube()

	if !c.Solved() {
		t.FailNow()
	}

	c.Spin(0)

	if c.Solved() {
		t.FailNow()
	}
}

// TestSolveDepth tests rubics cube solving with limited move set
func TestSolveDepth(t *testing.T) {
	var cube = NewCube()

	var rands = make([]uint8, 15)
	rand.Read(rands)

	for _, r := range rands {
		cube.Spin(r % uint8(NumSides/2))
	}

	t.Log("\n" + cube.String())
	var solution = BFSolve(cube)

	for _, move := range solution {
		cube.Spin(move)
	}

	if !cube.Solved() {
		t.Log(solution)
		t.Log(cube.String())
		t.Fatal("Failed to solve cube!")
	}
}

func testSolver(t *testing.T, s Solver) {
	var cube = NewCube()

	cube.Shuffle()

	var solution = s(cube)

	for _, move := range solution {
		cube.Spin(move)
	}

	if !cube.Solved() {
		t.Fatal("Failed to solve cube!")
	}
}

// TestTwoCycleSolver tests it
func TestTwoCycleSolver(t *testing.T) {
	testSolver(t, TwoCycleSolve)
}
