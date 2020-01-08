package cube

import (
	"crypto/rand"
	"fmt"
	"strings"
)

// Cube represents a standard rubik's cube
type Cube [NumSides][NumTiles]uint8

// NumSides is the number of sides of a cube
const NumSides = 6

// NumTiles is the number of mobile tiles on a Rubik's Cube
const NumTiles = 8

/* An enumeration of the edges of a Rubic's Cube */
const (
	Front = iota + 1
	Top
	Right
	Back
	Bottom
	Left
)

// Solution is the sequence of moves to solve a Rubik's Cube
type Solution []uint8

// Solver A method for solving a Rubik's Cube
type Solver func(Cube) Solution

// NewCube initilizes a solved Cube
func NewCube() Cube {
	var c Cube
	for s := range c {
		for i := range c[s] {
			c[s][i] = uint8(s)
		}
	}

	return c
}

func (c *Cube) String() string {
	var str [4 * 3]string
	//   1
	// 5 0 2
	//   4
	//   3

	// Space for each row but 5's
	for i := range str {
		if i >= 3 && i <= 5 {
			continue
		}

		str[i] += "      "
	}

	for _, i := range []int{5, 1, 0, 2, 4, 3} {
		var offset int
		switch i {
		case 1:
			offset = 0
		case 5:
			fallthrough
		case 0:
			fallthrough
		case 2:
			offset = 3
		case 4:
			offset = 6
		case 3:
			offset = 9
		}

		if i < 3 {
			str[0+offset] += fmt.Sprint(c[i][0], c[i][1], c[i][2], " ")
			str[1+offset] += fmt.Sprint(c[i][7], "   ", c[i][3], " ")
			str[2+offset] += fmt.Sprint(c[i][6], c[i][5], c[i][4], " ")
		} else {
			str[0+offset] += fmt.Sprint(c[i][4], c[i][5], c[i][6], " ")
			str[1+offset] += fmt.Sprint(c[i][3], "   ", c[i][7], " ")
			str[2+offset] += fmt.Sprint(c[i][2], c[i][1], c[i][0], " ")
		}
	}

	return strings.Join(str[:], "\n")
}

// front, right, top, back, left, bottom
// front: top, right, left, bottom :back
// left: front, top, back, bottom :right
// top: front, right, left, back :bottom

// <-> right <-> back <-> left <-> front
// <-> front <-> bottom <-> back <-> top
// <-> right <-> top <-> left <-> bottom

// we can move the middle by moving each side. We can reverse a move by moving the other two in the same direction.
func swap(left []uint8, right []uint8) {
	var temp [3]uint8
	for i, l := range left {
		temp[i] = l
	}
	for i, r := range right {
		left[i] = r
	}
	for i, t := range temp {
		right[i] = t
	}
}

// Spin spins the cube once
// the spin is clockwise on positive axis 0, 1, 2 = z, y, x
// the spin is symetric about the origin (counter) for 3, 4, 5, = -z, -y, -x
// 3 moves on a single face reverses a move
// a move of the middle is equivelent to one move on one side and 3 on the other
func (c *Cube) Spin(side uint8) {
	if side < 0 || side >= 6 {
		panic("Cube.Spin(side): Invalid size!")
	}

	// 0 1 2 \< 2 3 4
	// 7   3 |^ 1   5
	// 6 5 4 /> 0 7 6
	fa, ce := c[side][:2], c[side][2:NumTiles]
	for i, elem := range append(ce, fa...) {
		c[side][i] = elem
	}

	// turning faces affects related sides
	var next [3]uint8
	for i := range c {
		if ((int(side) + NumSides - i) % (NumSides / 2)) == 0 {
			continue
		}

		var temp [3]uint8
		switch side {
		case 0:
			switch i {
			case 1:
				fallthrough
			case 4:
				// 6 5 4
				swap(c[i][4:7], next[:])
			case 2:
				fallthrough
			case 5:
				// 6 7 0
				swap(temp[:], append(c[i][6:], c[i][0]))

				c[i][6] = next[0]
				c[i][7] = next[1]
				c[i][0] = next[2]

				swap(next[:], temp[:])
			}
		case 1:
			switch i {
			case 0:
				fallthrough
			case 2:
				fallthrough
			case 3:
				// 0 1 2
				swap(c[i][:3], next[:])
			case 5:
				// 4 5 6
				swap(c[i][4:7], next[:])
			}
		case 2:
			switch i {
			case 0:
				fallthrough
			case 1:
				// 2 3 4
				swap(c[i][2:5], next[:])
			case 3:
				fallthrough
			case 4:
				// 6 7 0
				swap(temp[:], append(c[i][6:], c[i][0]))

				c[i][6] = next[0]
				c[i][7] = next[1]
				c[i][0] = next[2]

				swap(next[:], temp[:])
			}
		case 3:
			switch i {
			case 1:
				fallthrough
			case 4:
				// 0 1 2
				swap(c[i][:3], next[:])
			case 2:
				fallthrough
			case 5:
				// 2 3 4
				swap(c[i][2:5], next[:])
			}
		case 4:
			switch i {
			case 0:
				fallthrough
			case 2:
				fallthrough
			case 3:
				// 6 5 4
				swap(c[i][4:7], next[:])
			case 5:
				// 0 1 2
				swap(c[i][:3], next[:])
			}
		case 5:
			switch i {
			case 0:
				fallthrough
			case 1:
				// 6 7 0
				swap(temp[:], append(c[i][6:], c[i][0]))

				c[i][6] = next[0]
				c[i][7] = next[1]
				c[i][0] = next[2]

				swap(next[:], temp[:])
			case 3:
				fallthrough
			case 4:
				// 2 3 4
				swap(c[i][2:5], next[:])
			}
		}
	}

	i := 0
	for ((int(side) + NumSides - i) % (NumSides / 2)) == 0 {

		i++
	}

	// one last time, slide next into place.
	switch side {
	case 0:
		switch i {
		case 1:
			fallthrough
		case 4:
			// 6 5 4
			swap(c[i][4:7], next[:])
		case 2:
			fallthrough
		case 5:
			// 6 7 0
			c[i][6] = next[0]
			c[i][7] = next[1]
			c[i][0] = next[2]
		}
	case 1:
		switch i {
		case 0:
			fallthrough
		case 2:
			fallthrough
		case 3:
			// 0 1 2
			swap(c[i][:3], next[:])
		case 5:
			// 4 5 6
			swap(c[i][4:7], next[:])
		}
	case 2:
		switch i {
		case 0:
			fallthrough
		case 1:
			// 2 3 4
			swap(c[i][2:5], next[:])
		case 3:
			fallthrough
		case 4:
			// 6 7 0
			c[i][6] = next[0]
			c[i][7] = next[1]
			c[i][0] = next[2]
		}
	case 3:
		switch i {
		case 1:
			fallthrough
		case 4:
			// 0 1 2
			swap(c[i][:3], next[:])
		case 2:
			fallthrough
		case 5:
			// 2 3 4
			swap(c[i][2:5], next[:])
		}
	case 4:
		switch i {
		case 0:
			fallthrough
		case 2:
			fallthrough
		case 3:
			// 6 5 4
			swap(c[i][4:7], next[:])
		case 5:
			// 0 1 2
			swap(c[i][:3], next[:])
		}
	case 5:
		switch i {
		case 0:
			fallthrough
		case 1:
			// 6 7 0
			c[i][6] = next[0]
			c[i][7] = next[1]
			c[i][0] = next[2]
		case 3:
			fallthrough
		case 4:
			// 2 3 4
			swap(c[i][2:5], next[:])
		}
	}
}

// Shuffle shuffles the Rubik's Cube
func (c *Cube) Shuffle() {
	var rands = make([]uint8, 1000)
	rand.Read(rands)

	for _, r := range rands {
		c.Spin(r % uint8(NumSides))
	}

}

// Solved iff each side is only a single color
func (c *Cube) Solved() bool {
	// either 2+ sides are wrong or none
	for s := 1; s < NumSides; s++ {
		for i := range c[s] {
			if c[s][i] != uint8(s) {
				return false
			}
		}
	}
	return true
}

// Cross iff the center and middle edges are equivilent
func (c *Cube) Cross(si uint8) bool {
	for _, i := range []uint8{2, 3, 5, 7} {
		if si != c[si][i] {
			return false
		}
	}
	return true
}

// Daisy iff middle edges are equivilent
func (c *Cube) Daisy(si uint8) bool {
	for _, i := range []uint8{3, 5, 7} {
		if c[si][2] != c[si][i] {
			return false
		}
	}
	return true
}

// RAlg manipulates corners
func (c *Cube) RAlg(s1, s2 uint8) {
	// X'3 == X
	c.Spin(s1)
	c.Spin(s1)
	c.Spin(s1)

	c.Spin(s2)
	c.Spin(s2)
	c.Spin(s2)

	c.Spin(s1)

	c.Spin(s2)
}

// BFSolve returns the shortest solution to the rubik's cube
func BFSolve(c Cube) Solution {
	// Uses Breadth-First Search
	var fringe, known = make(map[Cube]Solution), make(map[Cube]bool)

	// initialize fringe to hold source
	fringe[c] = nil
	known[c] = true

	for {
		var newFringe = make(map[Cube]Solution)
		for cube, path := range fringe {
			for i := uint8(0); i < NumSides; i++ {
				var next = cube
				next.Spin(i)

				sofar := append(append(Solution(nil), path...), i)
				if next.Solved() {
					return reduce(sofar)
				} else if !known[next] {
					known[next] = true
					newFringe[next] = sofar
				}
			}
			delete(fringe, cube)
		}
		fringe = newFringe
	}
}

// Reduce reduces the solution to higher-level moves
func reduce(s Solution) Solution {
	return s
}

// TwoCycleSolve is probs okay
func TwoCycleSolve(c Cube) Solution {
	return nil
}

// TODO
func (s *Solution) String() string {
	return ""
}
