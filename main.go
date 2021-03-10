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
	bs := NewBrowserScene(app, "gemini://gemini.circumlunar.space/docs/specification.gmi")
	bs.(*BrowserScene).fetch("gemini://gemini.circumlunar.space/docs/specification.gmi")
	app.sceneStack = &ui.SceneStack{bs}

	app.simple = ui.NewApp(app.sceneStack)
	app.simple.RunForever()
}
