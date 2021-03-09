package main

import (
	"fmt"
	"strings"

	"git.sr.ht/~adnano/go-gemini"
	"github.com/irth/go-simple"
	ui "github.com/irth/go-simple"
	"github.com/mitchellh/go-wordwrap"
)

type BrowserScene struct {
	r           *RMGem
	url         string
	gemtext     gemini.Text
	pages       []gemini.Text
	currentPage int
	pageHeight  int
}

func NewBrowserScene(r *RMGem, url string) simple.Scene {
	b := &BrowserScene{
		r:           r,
		url:         url,
		gemtext:     nil,
		pages:       nil,
		currentPage: 0,
		pageHeight:  30,
	}
	return b
}

func (b *BrowserScene) Render() (ui.Widget, error) {
	var prev simple.Widget
	var next simple.Widget
	if b.currentPage > 0 {
		prev = ui.Button(
			"prev",
			ui.Pos(ui.Abs(50), ui.Abs(b.r.simple.ScreenHeight()-90), ui.Abs(150), ui.Abs(100)),
			"[prev]",
			func(a *ui.App, btn *ui.ButtonWidget) error { b.previousPage(); return nil },
		)
	}
	if b.currentPage < len(b.pages)-1 {
		next = ui.Button(
			"next",
			ui.Pos(ui.Abs(b.r.simple.ScreenWidth()-50-120), ui.Abs(b.r.simple.ScreenHeight()-90), ui.Abs(150), ui.Abs(100)),
			"[next]",
			func(a *ui.App, btn *ui.ButtonWidget) error { b.nextPage(); return nil },
		)
	}

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
		ui.Label(
			ui.Pos(ui.Abs(0), ui.Abs(b.r.simple.ScreenHeight()-140), ui.Abs(b.r.simple.ScreenWidth()), ui.Abs(55)),
			"______________________________________________________________________________________",
		),
		ui.FontSize(40),
		prev,
		next,
	}, nil
}

func (b *BrowserScene) fetch() error {
	var err error
	b.currentPage = 0
	b.gemtext, err = Fetch(b.url)
	if err != nil {
		b.gemtext = nil
		b.pages = nil
		return err
	}
	b.pages = b.splitSite(b.gemtext)
	return err
}

func (b *BrowserScene) nextPage() {
	if b.currentPage < len(b.pages)-1 {
		b.currentPage++
	}
}

func (b *BrowserScene) previousPage() {
	if b.currentPage > 0 {
		b.currentPage--
	}
}

func (b *BrowserScene) estimateHeight(l gemini.Line) int {
	switch l := l.(type) {
	case gemini.LineText:
		if l.String() == "" {
			return 0
		}
		wrapped := wordwrap.WrapString(l.String(), 71) // 71 - line wrap length, checked experimentally
		return len(strings.Split(wrapped, "\n"))
	case gemini.LineHeading1, gemini.LineHeading2, gemini.LineHeading3:
		return 2
	default:
		return 1
	}
}

func (b *BrowserScene) splitSite(site gemini.Text) []gemini.Text {
	pages := []gemini.Text{}
	var page gemini.Text
	height := 0
	for _, line := range site {
		lineHeight := b.estimateHeight(line)
		height += lineHeight
		if height > b.pageHeight {
			// if the element won't fit, start a new page
			pages = append(pages, page)
			page = gemini.Text{}
			height = 0
		}
		if lineHeight > b.pageHeight {
			// if the element is too big to fit on a single page,
			// give it a dedicated page and hope for the best
			// TODO: solve this better
			pages = append(pages, gemini.Text{line})
			continue
		}
		page = append(page, line)
	}
	pages = append(pages, page)
	return pages
}

func (b *BrowserScene) getCurrentPage() gemini.Text {
	if b.pages == nil {
		return nil
	}
	return b.pages[b.currentPage]
}

func (b *BrowserScene) renderSite() simple.Widget {
	pos := ui.Pos(ui.Abs(100), ui.Step, ui.Abs(b.r.simple.ScreenWidth()-200), ui.Abs(55))
	wl := ui.WidgetList{
		ui.FontSize(32),
		ui.Label(pos, " "),
	}

	wasButton := false
	isButton := false
	for idx, line := range b.getCurrentPage() {
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
