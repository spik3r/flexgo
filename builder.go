package flexgo

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

type NodeBuilder struct {
	node *Node
}

func NewNode() *NodeBuilder {
	return &NodeBuilder{node: &Node{}}
}

func (b *NodeBuilder) Dir(dir Direction) *NodeBuilder {
	b.node.Dir = dir
	return b
}

func (b *NodeBuilder) Flex(flex int) *NodeBuilder {
	b.node.Flex = flex
	return b
}

func (b *NodeBuilder) Width(width int) *NodeBuilder {
	b.node.Width = width
	return b
}

func (b *NodeBuilder) Height(height int) *NodeBuilder {
	b.node.Height = height
	return b
}

func (b *NodeBuilder) MinWidth(min int) *NodeBuilder {
	b.node.MinWidth = min
	return b
}

func (b *NodeBuilder) MaxWidth(max int) *NodeBuilder {
	b.node.MaxWidth = max
	return b
}

func (b *NodeBuilder) MinHeight(min int) *NodeBuilder {
	b.node.MinHeight = min
	return b
}

func (b *NodeBuilder) MaxHeight(max int) *NodeBuilder {
	b.node.MaxHeight = max
	return b
}

func (b *NodeBuilder) Gap(gap int) *NodeBuilder {
	b.node.Gap = gap
	return b
}

func (b *NodeBuilder) Padding(padding int) *NodeBuilder {
	b.node.Padding = padding
	return b
}

func (b *NodeBuilder) Paddings(padding Spacing) *NodeBuilder {
	b.node.Paddings = padding
	return b
}

func (b *NodeBuilder) Margin(margin int) *NodeBuilder {
	b.node.Margin = margin
	return b
}

func (b *NodeBuilder) Margins(margin Spacing) *NodeBuilder {
	b.node.Margins = margin
	return b
}

func (b *NodeBuilder) MarginTopAuto(v bool) *NodeBuilder {
	b.node.MarginTopAuto = v
	return b
}

func (b *NodeBuilder) MarginBottomAuto(v bool) *NodeBuilder {
	b.node.MarginBottomAuto = v
	return b
}

func (b *NodeBuilder) MarginLeftAuto(v bool) *NodeBuilder {
	b.node.MarginLeftAuto = v
	return b
}

func (b *NodeBuilder) MarginRightAuto(v bool) *NodeBuilder {
	b.node.MarginRightAuto = v
	return b
}

func (b *NodeBuilder) Justify(justify Justify) *NodeBuilder {
	b.node.Justify = justify
	return b
}

func (b *NodeBuilder) Align(align Align) *NodeBuilder {
	b.node.Align = align
	return b
}

func (b *NodeBuilder) Debug(debug bool) *NodeBuilder {
	b.node.Debug = debug
	return b
}

func (b *NodeBuilder) Name(name string) *NodeBuilder {
	b.node.Name = name
	return b
}

func (b *NodeBuilder) ShowBorder(v bool) *NodeBuilder {
	b.node.ShowBorder = v
	return b
}

func (b *NodeBuilder) Border(border lipgloss.Border) *NodeBuilder {
	b.node.ShowBorder = true
	b.node.Border = border
	return b
}

func (b *NodeBuilder) BorderForeground(c color.Color) *NodeBuilder {
	b.node.BorderForeground = c
	return b
}

func (b *NodeBuilder) BorderBackground(c color.Color) *NodeBuilder {
	b.node.BorderBackground = c
	return b
}

func (b *NodeBuilder) AlignSelf(align Align) *NodeBuilder {
	b.node.AlignSelf = &align
	return b
}

func (b *NodeBuilder) Background(c color.Color) *NodeBuilder {
	b.node.Background = c
	return b
}

func (b *NodeBuilder) Foreground(c color.Color) *NodeBuilder {
	b.node.Foreground = c
	return b
}

func (b *NodeBuilder) View(view func(w, h int) string) *NodeBuilder {
	b.node.View = view
	return b
}

func (b *NodeBuilder) Children(children ...*Node) *NodeBuilder {
	b.node.Children = children
	return b
}

func (b *NodeBuilder) AddChildren(children ...*Node) *NodeBuilder {
	b.node.Children = append(b.node.Children, children...)
	return b
}

func (b *NodeBuilder) Build() *Node {
	return b.node
}
