package start

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/jakecoffman/fam2"
	"github.com/jakecoffman/fam2/scenes/game"
)

type Scene struct {
	*fam2.SceneManager
}

func (s *Scene) Update(float32) {
	if rl.IsKeyPressed(rl.KeyEnter) {
		s.CurrentScene = &game.Scene{SceneManager: s.SceneManager}
		s.CurrentScene.Load()
		s.Unload()
	}
}

func (s *Scene) Draw() {
	rl.DrawText("Press Enter", 200, 200, 10, rl.Black)
}

func (s *Scene) Load() {

}

func (s *Scene) Unload() {

}
