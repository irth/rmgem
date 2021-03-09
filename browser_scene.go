package main

import (
	"fmt"

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
	pos := ui.Pos(ui.Abs(100), ui.Step, ui.Abs(b.r.simple.ScreenWidth()-200), ui.Abs(60))
	wl := ui.WidgetList{
		ui.FontSize(32),
		ui.Label(pos, " "),
	}
	wasButton := false
	isButton := false
	for idx, line := range b.gemtext {
		wasButton = isButton
		isButton = false
		var widget simple.Widget
		switch line := line.(type) {
		case gemini.LineText:
			widget = b.textWidget(pos, line)
		case gemini.LineLink:
			widget = b.linkWidget(pos, line, idx)
			isButton = true
		case gemini.LineHeading1:
			widget = b.headingWidget(pos, 1, line.String())
		case gemini.LineHeading2:
			widget = b.headingWidget(pos, 2, line.String())
		case gemini.LineHeading3:
			widget = b.headingWidget(pos, 3, line.String())
		default:
			continue
		}
		if wasButton && !isButton {
			// insert padding
			wl = append(wl, ui.FontSize(16), ui.Label(pos, " "), ui.FontSize(32))
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

func (b *BrowserScene) linkWidget(pos ui.Position, l gemini.LineLink, idx int) simple.Widget {
	text := l.Name
	if text == "" {
		text = l.URL
	}
	return ui.Button(
		fmt.Sprintf("link_%d", idx),
		pos,
		fmt.Sprintf("=> %s", text),
		nil, // TODO: follow link onClick
	)
}

func (b *BrowserScene) headingWidget(pos ui.Position, level int, text string) simple.Widget {
	size := 32
	switch level {
	case 3:
		size = 38
	case 2:
		size = 48
	case 1:
		size = 64

	}
	return simple.WidgetList{
		ui.FontSize(size),
		ui.Paragraph(pos, text),
		ui.FontSize(32),
	}
}
