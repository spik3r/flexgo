package main

import (
	"fmt"

	"charm.land/lipgloss/v2"
	"github.com/spik3r/flexgo"
)

func box(label string) func(int, int) string {
	return func(w, h int) string {
		return lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			Width(w).
			Height(h).
			Align(lipgloss.Center, lipgloss.Center).
			Render(label)
	}
}

func main() {
	header := &flexgo.Node{
		Dir:     flexgo.Row,
		Height:  1,
		Justify: flexgo.JustifyCenter,
		View:    box("HEADER"),
	}

	content := &flexgo.Node{
		Dir:     flexgo.Row,
		Flex:    1,
		Justify: flexgo.JustifySpaceBetween,
		Children: []*flexgo.Node{
			{Flex: 3, View: box("LEFT")},
			{Flex: 7, View: box("RIGHT")},
		},
	}

	statusline := &flexgo.Node{
		Dir:     flexgo.Row,
		Flex:    1,
		Height:  1,
		Justify: flexgo.JustifySpaceBetween,
		Children: []*flexgo.Node{
			{Flex: 7, View: box("STATUS INFO")},
			{Flex: 2, View: box("1.0.0-alpha")},
		},
	}

	contentContainer := &flexgo.Node{
		Dir:     flexgo.Col,
		Flex:    1,
		Justify: flexgo.JustifyStart,
		Children: []*flexgo.Node{
			statusline,
			content,
		},
	}

	footer := &flexgo.Node{
		Dir:     flexgo.Row,
		Height:  2,
		Justify: flexgo.JustifyCenter,
		View:    box("Footer"),
	}

	page := &flexgo.Node{
		Dir:   flexgo.Col,
		Flex:  1,
		Debug: true,
		Children: []*flexgo.Node{
			header,
			contentContainer,
			footer,
		},
	}

	out := page.Render(80, 24)
	fmt.Println(out)
}
