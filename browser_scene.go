package main

import (
	ui "github.com/irth/go-simple"
)

type BrowserScene struct{ r *RMGem }

func (b *BrowserScene) Render() (ui.Widget, error) {
	return ui.WidgetList{
		ui.Justify(ui.Left),
		ui.FontSize(42),
		ui.TextInput(
			"address",
			ui.Pos(ui.Abs(100), ui.Abs(100), ui.Abs(b.r.simple.ScreenWidth()-250-5), ui.Abs(55)),
			"gemini://irth.pl",
			nil,
		),
		ui.Button(
			"go",
			ui.Pos(ui.Abs(b.r.simple.ScreenWidth()-150), ui.Same, ui.Abs(50), ui.Abs(55)),
			"Go",
			nil,
		),
		// Displaying user generated content using paragraph is unsafe, because AFAIK ] cannot be escaped.
		// Newlines should not be a problem, as each newline in Gemini generates a new widget.
		ui.Paragraph(
			ui.Pos(ui.Abs(100), ui.Abs(200), ui.Abs(b.r.simple.ScreenWidth()-200), ui.Abs(100)),
			"gemini site content placeholder",
		),
	}, nil
}
