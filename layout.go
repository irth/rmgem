package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"git.sr.ht/~adnano/go-gemini"
	ui "github.com/irth/go-simple"
	"github.com/mitchellh/go-wordwrap"
)

type Dimensions struct {
	LineHeight  int
	FontSize    int
	Padding     Padding
	LineSpacing int
}

func (d Dimensions) Total(lines int) int {
	betweenLines := lines - 1
	if betweenLines < 0 {
		betweenLines = 0
	}
	return d.Padding.Top + lines*d.LineHeight + betweenLines*d.LineSpacing + d.Padding.Bottom
}

type Padding struct {
	Top    int
	Bottom int
}

func (p Padding) Total() int {
	return p.Top + p.Bottom
}

type LayoutEngine struct {
	PageHeight int

	Text     Dimensions
	Link     Dimensions
	Heading1 Dimensions
	Heading2 Dimensions
	Heading3 Dimensions
}

func NewLayoutEngine(pageHeight int) LayoutEngine {
	return LayoutEngine{
		PageHeight: pageHeight,
		Text: Dimensions{
			LineHeight:  28,
			FontSize:    32,
			Padding:     Padding{Top: 8, Bottom: 8},
			LineSpacing: 8,
		},

		Link: Dimensions{
			LineHeight: 29 + 15,
			FontSize:   32,
			Padding:    Padding{Top: 4, Bottom: 4},
		},

		Heading1: Dimensions{
			LineHeight: 56,
			FontSize:   64,
			Padding:    Padding{Top: 32, Bottom: 16},
		},
		Heading2: Dimensions{
			LineHeight: 42,
			FontSize:   48,
			Padding:    Padding{Top: 24, Bottom: 8},
		},
		Heading3: Dimensions{
			LineHeight: 33,
			FontSize:   38,
			Padding:    Padding{Top: 16, Bottom: 4},
		},
	}
}

func (l *LayoutEngine) getDimensions(line gemini.Line) Dimensions {
	go func() {
		for {
			time.Sleep(1 * time.Second)
		}
	}()
	switch line.(type) {
	case gemini.LineText, gemini.LineListItem:
		return l.Text

	case gemini.LineHeading1:
		// TODO: wordwrap headings and links
		return l.Heading1

	case gemini.LineHeading2:
		return l.Heading2

	case gemini.LineHeading3:
		return l.Heading3

	case gemini.LineLink:
		return l.Link

	default:
		// log.Printf("layout: getDimensions: unknown line type: %T", line)
		return Dimensions{}
	}

}

func (l LayoutEngine) wrapLines(text string) []string {
	// 71 - line wrap length, checked experimentally
	return strings.Split(wordwrap.WrapString(text, 71), "\n")
}

func (l LayoutEngine) estimateLines(line gemini.Line) int {
	switch line := line.(type) {
	case gemini.LineText, gemini.LineListItem:
		if line.String() == "" {
			return 0
		}
		wrapped := l.wrapLines(line.String())
		noOfLines := len(wrapped)
		return noOfLines

	case gemini.LineHeading1, gemini.LineHeading2, gemini.LineHeading3, gemini.LineLink:
		// TODO: wordwrap headings and links
		return 1

	default:
		// log.Printf("layout: estimateLines: unknown line type: %T", line)
		return 0
	}
}

func (l LayoutEngine) getWidget(pos ui.Position, line gemini.Line, idx int, followUrl func(url string)) ui.Widget {
	switch line := line.(type) {
	case gemini.LineText, gemini.LineListItem:
		return l.textWidget(pos, line.String())

	case gemini.LineHeading1, gemini.LineHeading2, gemini.LineHeading3:
		// TODO: wordwrap headings and links
		return l.headingWidget(pos, line.String())

	case gemini.LineLink:
		return l.linkWidget(pos, line, idx, followUrl)

	default:
		return nil
	}

}

func (l LayoutEngine) textWidget(pos ui.Position, line string) ui.Widget {
	wrapped := l.wrapLines(line)

	wl := ui.WidgetList{}

	x := pos.X
	y := (pos.Y.(ui.Abs))
	w := pos.Width
	h := pos.Height

	for _, textLine := range wrapped {
		wl = append(wl,
			ui.Paragraph(
				ui.Pos(x, y, w, h),
				textLine,
			),
		)
		y += ui.Abs(l.Text.LineHeight + l.Text.LineSpacing)
	}

	return wl
}

func (l LayoutEngine) linkWidget(pos ui.Position, line gemini.LineLink, idx int, followUrl func(url string)) ui.Widget {
	text := line.Name
	if text == "" {
		text = line.URL
	}
	return ui.Button(
		fmt.Sprintf("link_%d", idx),
		pos,
		fmt.Sprintf("=> %s", text),
		func(a *ui.App, b *ui.ButtonWidget) error {
			followUrl(line.URL)
			return nil
		},
	)
}

func (l LayoutEngine) headingWidget(pos ui.Position, text string) ui.Widget {
	return ui.Paragraph(pos, text)
}

func (l LayoutEngine) newPage(gemtext gemini.Text) Page {
	return Page{LayoutEngine: l, Gemtext: gemtext}
}

func (l LayoutEngine) trySplitting(line gemini.Line, remainingSpace int) (gemini.Line, gemini.Line, bool) {
	switch line := line.(type) {
	case gemini.LineText:
		wrapped := l.wrapLines(line.String())
		lineCount := len(wrapped)
		dim := l.getDimensions(line)
		splitBoundary := (remainingSpace - dim.Padding.Total() + lineCount*dim.LineSpacing) /
			(lineCount * (dim.LineSpacing + dim.LineHeight))
		if splitBoundary <= 0 {
			return nil, nil, false
		}

		if splitBoundary > lineCount {
			log.Println("splitBoundary > lineCount: this shouldn't ever happen")
			return nil, nil, false
		}

		firstText := wrapped[:splitBoundary]
		secondText := wrapped[splitBoundary:]

		return gemini.LineText(strings.Join(firstText, " ")),
			gemini.LineText(strings.Join(secondText, " ")),
			true

	default:
		return nil, nil, false
	}

}

func (l LayoutEngine) splitPages(gemtext gemini.Text) []Page {
	pages := []Page{}

	page := l.newPage(nil)

	height := 0
	for _, line := range gemtext {
		lineCount := l.estimateLines(line)
		lineHeight := l.getDimensions(line).Total(lineCount)
		if height+lineHeight > l.PageHeight {
			log.Println("exceeded max page length at height", height)
			remainingSpace := l.PageHeight - height
			first, second, ok := l.trySplitting(line, remainingSpace)

			if ok {
				log.Println("solved by splitting the line")
				page.Gemtext = append(page.Gemtext, first)
				pages = append(pages, page)
				page = l.newPage(gemini.Text{second})
				height = l.getDimensions(second).Total(l.estimateLines(second))
				fmt.Println(first.String(), "///", second.String())
				continue
			}

			// if the element won't fit, start a new page
			pages = append(pages, page)
			page = l.newPage(nil)
			height = 0

		}
		if lineHeight > l.PageHeight {
			// if the element is too big to fit on a single page,
			// give it a dedicated page and hope for the best
			// TODO: split the element in two!
			pages = append(pages, l.newPage(gemini.Text{line}))
			continue
		}
		height += lineHeight
		page.Gemtext = append(page.Gemtext, line)
	}

	pages = append(pages, page)
	for n, ll := range pages {
		fmt.Println("")
		fmt.Println("page", n)
		fmt.Println("first:", ll.Gemtext[0].String())
		fmt.Println("last:", ll.Gemtext[len(ll.Gemtext)-1].String())
		fmt.Println("")
	}
	return pages
}

type Page struct {
	Gemtext      gemini.Text
	LayoutEngine LayoutEngine
}

func (p Page) AtPosition(Y int, followUrl func(url string)) ui.Widget {
	x := 100
	y := Y

	wl := ui.WidgetList{}

	totalHeight := 0

	for idx, line := range p.Gemtext {
		lineCount := p.LayoutEngine.estimateLines(line)
		dim := p.LayoutEngine.getDimensions(line)
		pos := ui.Pos(
			ui.Abs(x), ui.Abs(y+dim.Padding.Top),
			ui.Abs(1304), ui.Abs(dim.LineHeight*lineCount))
		y += dim.Total(lineCount)

		wl = append(wl, ui.WidgetList{
			ui.FontSize(dim.FontSize),
			p.LayoutEngine.getWidget(pos, line, idx, followUrl),
		})

		totalHeight += dim.Total(lineCount)
	}

	log.Println("drawing page with estimated height: ", totalHeight)
	return wl
}
