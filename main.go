package main

import (
	ui "github.com/irth/go-simple"
)

type RMGem struct {
	simple     *ui.App
	sceneStack *ui.SceneStack
}

func main() {
	app := &RMGem{}
	app.sceneStack = &ui.SceneStack{NewBrowserScene(app, "gemini://irth.pl")}
	app.simple = ui.NewApp(app.sceneStack)
	app.simple.RunForever()
}
