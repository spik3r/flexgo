package flexgo

import (
	"errors"
	"fmt"
	"strings"
)

// Validate walks the tree and reports structural problems that would
// otherwise be resolved silently at render time: conflicting field
// settings, negative sizes, and mis-shaped nodes. Returns nil when the
// tree is well-formed.
//
// The checks cover the precedence rules documented on Node:
//   - View and Children set together (container wins, View is ignored).
//   - Padding and Paddings both non-zero (Paddings wins).
//   - Margin and Margins both non-zero (Margins wins).
//   - Negative Flex / Width / Height / Min* / Max* / Gap.
//   - MinWidth > MaxWidth (similarly for height) when both are set.
//
// Returns a multi-error joined with errors.Join; use errors.Is or
// errors.As to inspect individual issues.
func Validate(root *Node) error {
	if root == nil {
		return nil
	}
	var errs []error
	validateNode(root, "root", &errs)
	return errors.Join(errs...)
}

func validateNode(n *Node, path string, errs *[]error) {
	if n == nil {
		*errs = append(*errs, fmt.Errorf("%s: nil node", path))
		return
	}

	if n.View != nil && len(n.Children) > 0 {
		*errs = append(*errs, fmt.Errorf("%s: View and Children both set; View is ignored", path))
	}
	if n.Padding != 0 && n.Paddings != (Spacing{}) {
		*errs = append(*errs, fmt.Errorf("%s: Padding and Paddings both set; Paddings wins", path))
	}
	if n.Margin != 0 && n.Margins != (Spacing{}) {
		*errs = append(*errs, fmt.Errorf("%s: Margin and Margins both set; Margins wins", path))
	}
	if n.Flex < 0 {
		*errs = append(*errs, fmt.Errorf("%s: negative Flex (%d)", path, n.Flex))
	}
	if n.Width < 0 {
		*errs = append(*errs, fmt.Errorf("%s: negative Width (%d)", path, n.Width))
	}
	if n.Height < 0 {
		*errs = append(*errs, fmt.Errorf("%s: negative Height (%d)", path, n.Height))
	}
	if n.Gap < 0 {
		*errs = append(*errs, fmt.Errorf("%s: negative Gap (%d)", path, n.Gap))
	}
	if n.MinWidth > 0 && n.MaxWidth > 0 && n.MinWidth > n.MaxWidth {
		*errs = append(*errs, fmt.Errorf("%s: MinWidth (%d) > MaxWidth (%d)", path, n.MinWidth, n.MaxWidth))
	}
	if n.MinHeight > 0 && n.MaxHeight > 0 && n.MinHeight > n.MaxHeight {
		*errs = append(*errs, fmt.Errorf("%s: MinHeight (%d) > MaxHeight (%d)", path, n.MinHeight, n.MaxHeight))
	}

	for i, child := range n.Children {
		childPath := fmt.Sprintf("%s.Children[%d]", path, i)
		if child != nil && child.Name != "" {
			childPath = fmt.Sprintf("%s.%s", path, child.Name)
		}
		validateNode(child, childPath, errs)
	}
}

// Inspect returns a human-readable tree dump of root for debugging.
// Each line shows the node's shape (direction, flex/size, spacing) and
// its Name when set. Useful when a layout is misbehaving and you want
// to see what the tree actually looks like.
//
// Sizes shown are the structural values set on the Node, not the
// post-layout allocated sizes — that requires the two-phase layout
// (see X1 in todo.md).
func Inspect(root *Node) string {
	if root == nil {
		return "<nil>\n"
	}
	var b strings.Builder
	inspectNode(&b, root, "", true, true)
	return b.String()
}

func inspectNode(b *strings.Builder, n *Node, prefix string, isLast, isRoot bool) {
	if n == nil {
		b.WriteString(prefix)
		b.WriteString("<nil>\n")
		return
	}

	branch := "├─ "
	childPrefix := prefix + "│  "
	if isLast {
		branch = "└─ "
		childPrefix = prefix + "   "
	}
	if isRoot {
		branch = ""
		childPrefix = prefix
	}

	b.WriteString(prefix)
	b.WriteString(branch)
	b.WriteString(describeNode(n))
	b.WriteByte('\n')

	for i, child := range n.Children {
		inspectNode(b, child, childPrefix, i == len(n.Children)-1, false)
	}
}

func describeNode(n *Node) string {
	var parts []string

	kind := "Box"
	switch {
	case n.View != nil && len(n.Children) == 0:
		kind = "Leaf"
	case n.Dir == Row:
		kind = "HBox"
	case n.Dir == Col:
		kind = "VBox"
	}
	parts = append(parts, kind)

	if n.Name != "" {
		parts = append(parts, fmt.Sprintf("%q", n.Name))
	}

	var sizeParts []string
	if n.Flex > 0 {
		sizeParts = append(sizeParts, fmt.Sprintf("flex=%d", n.Flex))
	}
	if n.Width > 0 {
		sizeParts = append(sizeParts, fmt.Sprintf("w=%d", n.Width))
	}
	if n.Height > 0 {
		sizeParts = append(sizeParts, fmt.Sprintf("h=%d", n.Height))
	}
	if len(sizeParts) > 0 {
		parts = append(parts, strings.Join(sizeParts, " "))
	}

	if n.Gap > 0 {
		parts = append(parts, fmt.Sprintf("gap=%d", n.Gap))
	}
	if n.Padding > 0 || n.Paddings != (Spacing{}) {
		parts = append(parts, "pad="+describeSpacing(n.Padding, n.Paddings))
	}
	if n.Margin > 0 || n.Margins != (Spacing{}) {
		parts = append(parts, "mar="+describeSpacing(n.Margin, n.Margins))
	}
	if n.ShowBorder {
		parts = append(parts, "border")
	}
	if n.Debug {
		parts = append(parts, "debug")
	}

	return strings.Join(parts, " ")
}

func describeSpacing(shorthand int, sides Spacing) string {
	if sides != (Spacing{}) {
		return fmt.Sprintf("{%d,%d,%d,%d}", sides.Top, sides.Right, sides.Bottom, sides.Left)
	}
	return fmt.Sprintf("%d", shorthand)
}
