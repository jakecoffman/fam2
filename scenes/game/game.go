package game

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/jakecoffman/cp"
	"github.com/jakecoffman/fam2"
)

var grabbableMaskBit uint = 1 << 31
var grabFilter = cp.ShapeFilter{
	cp.NO_GROUP, grabbableMaskBit, grabbableMaskBit,
}

type Scene struct {
	*fam2.SceneManager

	space      *cp.Space
	mouseBody  *cp.Body
	mousePos   cp.Vector
	mouseJoint *cp.Constraint

	accumulator, dt, alpha float32
}

func (s *Scene) Update(dt float32) {
	mousePos := rl.GetMousePosition()
	s.mousePos = cp.Vector{X: float64(mousePos.X/s.Camera.Zoom), Y: float64(mousePos.Y/s.Camera.Zoom)}
	newPoint := s.mouseBody.Position().Lerp(s.mousePos, 0.25)
	s.mouseBody.SetVelocityVector(newPoint.Sub(s.mouseBody.Position()).Mult(60.0))
	s.mouseBody.SetPosition(s.mousePos)

	// handle grabbing
	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		result := s.space.PointQueryNearest(s.mousePos, 5, grabFilter)
		if result.Shape != nil && result.Shape.Body().Mass() < cp.INFINITY {
			var nearest cp.Vector
			if result.Distance > 0 {
				nearest = result.Point
			} else {
				nearest = s.mousePos
			}

			// create a new constraint where the mousePos is to draw the body towards the mousePos
			body := result.Shape.Body()
			s.mouseJoint = cp.NewPivotJoint2(s.mouseBody, body, cp.Vector{}, body.WorldToLocal(nearest))
			s.mouseJoint.SetMaxForce(50000)
			s.mouseJoint.SetErrorBias(math.Pow(1.0-0.15, 60.0))
			s.space.AddConstraint(s.mouseJoint)
		}
	} else if rl.IsMouseButtonReleased(rl.MouseLeftButton) && s.mouseJoint != nil {
		s.space.RemoveConstraint(s.mouseJoint)
		s.mouseJoint = nil
	}

	// perform a fixed rate physics tick
	const physicsTickrate = 1.0 / 180.0
	s.accumulator += dt
	for s.accumulator >= physicsTickrate {
		s.space.Step(physicsTickrate)
		s.accumulator -= physicsTickrate
	}
	s.alpha = s.accumulator / dt
}

func (s *Scene) Draw() {
	// this is a generic way to iterate over the shapes in a space,
	// to avoid the type switch just keep a pointer to the shapes when they've been created
	s.space.EachShape(func(s *cp.Shape) {
		switch s.Class.(type) {
		case *cp.Segment:
			segment := s.Class.(*cp.Segment)
			a := segment.A()
			b := segment.B()
			rl.DrawLineEx(v(a), v(b), float32(segment.Radius()), rl.Black)
		case *cp.Circle:
			circle := s.Class.(*cp.Circle)
			pos := circle.Body().Position()
			rl.DrawCircleV(v(pos), float32(circle.Radius()), rl.Red)
		default:
			fmt.Println("unexpected shape", s.Class)
		}
	})
}

func (s *Scene) Load() {
	s.space = cp.NewSpace()
	body := s.space.AddBody(cp.NewBody(1, cp.MomentForCircle(1, 10, 10, cp.Vector{})))
	shape := body.AddShape(cp.NewCircle(body, 10, cp.Vector{0, 0}))
	s.space.AddShape(shape)
	body.SetPosition(cp.Vector{100, 100})
	s.space.SetGravity(cp.Vector{0, 100})

	const width, height = 800, 600
	s.space.AddShape(cp.NewSegment(s.space.StaticBody, cp.Vector{0, height}, cp.Vector{width, height}, 1))
	s.space.AddShape(cp.NewSegment(s.space.StaticBody, cp.Vector{width, 0}, cp.Vector{width, height}, 1))
	s.space.AddShape(cp.NewSegment(s.space.StaticBody, cp.Vector{width, 0}, cp.Vector{0, 0}, 1))
	s.space.AddShape(cp.NewSegment(s.space.StaticBody, cp.Vector{0, height}, cp.Vector{0, 0}, 1))

	s.mouseBody = cp.NewKinematicBody()
}

func (s *Scene) Unload() {
	s.space.EachShape(func(shape *cp.Shape) {
		s.space.RemoveShape(shape)
	})
	s.space.EachBody(func(body *cp.Body) {
		s.space.RemoveBody(body)
	})
}

func v(v cp.Vector) rl.Vector2 {
	return rl.Vector2{X: float32(v.X), Y: float32(v.Y)}
}
