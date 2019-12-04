package fam2

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Scener interface {
	Update(float32)
	Draw()
	Load()
	Unload()
}

type SceneManager struct {
	Camera       *rl.Camera2D
	CurrentScene Scener
}
