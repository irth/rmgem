package main

import (
	"git.sr.ht/~adnano/go-gemini"
	"github.com/irth/go-simple"
	ui "github.com/irth/go-simple"
)

type BrowserScene struct {
	r       *RMGem
	url     string
	gemtext gemini.Text
}

func NewBrowserScene(r *RMGem, url string) simple.Scene {
	b := &BrowserScene{
		r, url, nil,
	}
	return b
}

func (b *BrowserScene) Render() (ui.Widget, error) {
	return ui.WidgetList{
		ui.Justify(ui.Left),
		ui.FontSize(42),
		ui.TextInput(
			"address",
			ui.Pos(ui.Abs(100), ui.Abs(100), ui.Abs(b.r.simple.ScreenWidth()-250-5), ui.Abs(55)),
			b.url,
			func(a *ui.App, t *ui.TextInputWidget, value string) error {
				b.url = value
				return nil
			},
		),
		ui.Button(
			"go",
			ui.Pos(ui.Abs(b.r.simple.ScreenWidth()-150), ui.Same, ui.Abs(50), ui.Abs(55)),
			"Go",
			func(a *ui.App, button *ui.ButtonWidget) error {
				err := b.fetch()
				if err != nil {
					panic(err)
				}
				println(b.gemtext.String())
				return nil
			},
		),
		// Displaying user generated content using paragraph is unsafe, because AFAIK ] cannot be escaped.
		// Newlines should not be a problem, as each newline in Gemini generates a new widget.
		b.renderSite(),
	}, nil
}

func (b *BrowserScene) fetch() error {
	var err error
	b.gemtext, err = Fetch(b.url)
	return err
}

func (b *BrowserScene) renderSite() simple.Widget {
	pos := ui.Pos(ui.Abs(100), ui.Step, ui.Abs(b.r.simple.ScreenWidth()-200), ui.Abs(25))
	wl := ui.WidgetList{
		ui.FontSize(32),
		ui.Label(pos, " "),
	}
	for _, line := range b.gemtext {
		var widget simple.Widget
		switch line := line.(type) {
		case gemini.LineText:
			widget = b.textWidget(pos, line)
		default:
			continue
		}
		wl = append(wl, widget)
	}

	return wl
}

func (b *BrowserScene) textWidget(pos ui.Position, l gemini.LineText) simple.Widget {
	return ui.Paragraph(
		pos,
		l.String(),
	)
}
