package main

import (
	"fmt"
	"log"
	"net/url"

	"git.sr.ht/~adnano/go-gemini"
	"github.com/irth/go-simple"
	ui "github.com/irth/go-simple"
)

type BrowserScene struct {
	r           *RMGem
	url         string
	gemtext     gemini.Text
	pages       []Page
	layout      LayoutEngine
	currentPage int
	pageHeight  int
}

func NewBrowserScene(r *RMGem, url string) simple.Scene {
	b := &BrowserScene{
		r:           r,
		url:         url,
		gemtext:     nil,
		pages:       nil,
		layout:      NewLayoutEngine(r.simple.ScreenHeight() - 300),
		currentPage: 0,
		pageHeight:  20,
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
		b.hr(80),
		ui.FontSize(42),
		ui.Button(
			"exit",
			ui.Pos(ui.Abs(20), ui.Abs(20), ui.Abs(70), ui.Abs(70)),
			"[X]",
			func(a *ui.App, button *ui.ButtonWidget) error {
				// this scene will clear screen and call os.Exit(0)
				b.r.sceneStack.Replace(ExitScene{})
				return nil
			},
		),
		ui.TextInput(
			"address",
			ui.Pos(ui.Abs(100), ui.Abs(30), ui.Abs(b.r.simple.ScreenWidth()-250-5), ui.Abs(55)),
			b.url,
			func(a *ui.App, t *ui.TextInputWidget, value string) error {
				b.url = value
				return nil
			},
		),
		ui.Button(
			"go",
			ui.Pos(ui.Abs(b.r.simple.ScreenWidth()-130), ui.Abs(25), ui.Abs(90), ui.Abs(70)),
			"[Go]",
			func(a *ui.App, button *ui.ButtonWidget) error {
				err := b.fetch(b.url)
				if err != nil {
					panic(err)
				}
				return nil
			},
		),
		b.getCurrentPage().AtPosition(150, func(url string) {
			err := b.fetch(url)
			if err != nil {
				log.Println(err)
			}
		}),
		b.hr(b.r.simple.ScreenHeight() - 140),
		ui.FontSize(40),
		prev,
		next,
	}, nil
}

func (b *BrowserScene) fetch(newUrl string) error {
	u := newUrl
	if b.url != "" {
		new, err := url.Parse(newUrl)
		if err != nil {
			return fmt.Errorf("while parsing new url: %w", err)
		}
		base, err := url.Parse(b.url)
		if err != nil {
			return fmt.Errorf("while parsing base url: %w", err)
		}
		u = base.ResolveReference(new).String()
	}

	b.url = u
	var err error
	b.currentPage = 0
	b.gemtext, err = Fetch(b.url)
	if err != nil {
		b.gemtext = nil
		b.pages = nil
		if err != nil {
			log.Println("Fetch error:", err.Error())
		}
		return err
	}
	b.pages = b.layout.splitPages(b.gemtext)
	return nil
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

func (b *BrowserScene) hr(y int) simple.Widget {
	return ui.WidgetList{
		ui.FontSize(32),
		ui.Label(
			ui.Pos(ui.Abs(0), ui.Abs(y), ui.Abs(b.r.simple.ScreenWidth()), ui.Abs(55)),
			"______________________________________________________________________________________",
		),
	}
}

func (b *BrowserScene) getCurrentPage() Page {
	if b.pages == nil {
		return Page{nil, b.layout}
	}
	return b.pages[b.currentPage]
}
