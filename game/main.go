package main

import (
	"github.com/jakecoffman/fam2"
	"github.com/jakecoffman/fam2/scenes/start"
)

func main() {
	sceneManager := &fam2.SceneManager{}
	sceneManager.CurrentScene = &start.Scene{SceneManager: sceneManager}
	sceneManager.CurrentScene.Load()
	fam2.Run(sceneManager)
}
