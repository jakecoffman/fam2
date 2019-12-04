package fam2

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func Run(sceneManager *SceneManager) {
	const width, height = 800, 450

	rl.SetConfigFlags(rl.FlagVsyncHint)
	rl.InitWindow(width, height, "raylib [physics] example - chipmunk")

	camera := rl.NewCamera2D(rl.Vector2{}, rl.Vector2{}, 0, 1)
	sceneManager.Camera = &camera

	lastTime := rl.GetTime()

	rl.SetMouseScale(1, 1) // fix for high dpi getting wrong mousePos position

	var dt float32
	for !rl.WindowShouldClose() {
		// calculate dt
		now := rl.GetTime()
		dt = now - lastTime
		lastTime = now

		sceneManager.CurrentScene.Update(dt)

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		rl.BeginMode2D(camera)

		sceneManager.CurrentScene.Draw()

		rl.EndMode2D()
		rl.DrawFPS(0, 0)
		rl.EndDrawing()
	}

	sceneManager.CurrentScene.Unload()
	rl.CloseWindow()
}

