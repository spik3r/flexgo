package main

import (
	"fmt"
	"os"

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
	if os.Getenv("FLEXGO_GOLDEN") == "1" {
		fmt.Print(render())
		return
	}

	fmt.Println(render())
}

func render() string {
	header := &flexgo.Node{
		Dir:     flexgo.Row,
		Height:  3,
		Justify: flexgo.JustifyCenter,
		View:    box("HEADER"),
	}

	content := &flexgo.Node{
		Dir:        flexgo.Row,
		Flex:       1,
		Justify:    flexgo.JustifySpaceBetween,
		Background: lipgloss.Color("237"),
		Children: []*flexgo.Node{
			{Flex: 3, Background: lipgloss.Color("61"), View: box("LEFT")},
			{Flex: 7, Background: lipgloss.Color("61"), View: box("RIGHT")},
		},
	}

	statusline := &flexgo.Node{
		Dir:     flexgo.Row,
		Height:  3,
		Justify: flexgo.JustifySpaceBetween,
		Children: []*flexgo.Node{
			{Width: 56, View: box("STATUS INFO")},
			{Width: 16, View: box("1.0.0-alpha")},
		},
	}

	contentContainer := &flexgo.Node{
		Dir:        flexgo.Col,
		Flex:       1,
		Justify:    flexgo.JustifyStart,
		Background: lipgloss.Color("237"),
		Children: []*flexgo.Node{
			statusline,
			content,
		},
	}

	footer := &flexgo.Node{
		Dir:     flexgo.Row,
		Height:  3,
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

	return page.Render(80, 24)
}
