package main

import (
	"fmt"
	"os"

	"charm.land/lipgloss/v2"
	"github.com/spik3r/flexgo"
)

func main() {
	if os.Getenv("FLEXGO_GOLDEN") == "1" {
		fmt.Print(render())
		return
	}

	fmt.Println(render())
}

func render() string {
	bg := lipgloss.Color("234")
	headerBg := lipgloss.Color("61")
	contentBg := lipgloss.Color("237")
	footerBg := lipgloss.Color("66")

	header := flexgo.NewNode().
		Dir(flexgo.Row).
		Height(3).
		Background(headerBg).
		Foreground(lipgloss.Color("255")).
		Justify(flexgo.JustifyCenter).
		View(func(w, h int) string {
			return lipgloss.NewStyle().
				Width(w).
				Height(h).
				Bold(true).
				Align(lipgloss.Center, lipgloss.Center).
				Render("FLEXGO FOREGROUND")
		}).Build()

	leftContent := flexgo.NewNode().
		Flex(1).
		Background(contentBg).
		Foreground(lipgloss.Color("250")).
		View(func(w, h int) string {
			return lipgloss.NewStyle().
				Width(w).
				Height(h).
				Align(lipgloss.Center, lipgloss.Center).
				Render("Foreground\ncolor\napplied")
		}).Build()

	rightContent := flexgo.NewNode().
		Flex(1).
		Background(contentBg).
		Foreground(lipgloss.Color("61")).
		View(func(w, h int) string {
			return lipgloss.NewStyle().
				Width(w).
				Height(h).
				Align(lipgloss.Center, lipgloss.Center).
				Render("Different\nforeground\ncolors")
		}).Build()

	content := flexgo.NewNode().
		Dir(flexgo.Row).
		Flex(1).
		Background(contentBg).
		Children(leftContent, rightContent).Build()

	footer := flexgo.NewNode().
		Dir(flexgo.Row).
		Height(3).
		Background(footerBg).
		Foreground(lipgloss.Color("255")).
		Justify(flexgo.JustifyCenter).
		View(func(w, h int) string {
			return lipgloss.NewStyle().
				Width(w).
				Height(h).
				Bold(true).
				Align(lipgloss.Center, lipgloss.Center).
				Render("Footer with foreground color")
		}).Build()

	page := flexgo.NewNode().
		Dir(flexgo.Col).
		Flex(1).
		Background(bg).
		Children(header, content, footer).Build()

	return page.Render(80, 24)
}
