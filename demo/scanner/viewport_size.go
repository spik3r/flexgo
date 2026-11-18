package main

import "charm.land/lipgloss/v2"

const (
	shellFrameSize      = 1
	shellBodyMargin     = 2
	shellBrandHeight    = 3
	shellLineHeight     = 1
	scanSidebarWidth    = 24
	historySidebarWidth = 28
	screenSplitGap      = 2
	tabRowHeight        = 1
	tabStripHeight      = 2 * tabRowHeight
)

var fileViewerBorder = lipgloss.RoundedBorder()

type viewportSize struct {
	width, height int
}

func pageBodySize(w, h int) viewportSize {
	return (viewportSize{width: w, height: h}).inset(
		2*shellFrameSize+2*shellBodyMargin,
		2*shellFrameSize+shellBrandHeight+5*shellLineHeight,
	)
}

func (size viewportSize) inset(w, h int) viewportSize {
	return viewportSize{
		width:  max(0, max(0, size.width)-w),
		height: max(0, max(0, size.height)-h),
	}
}

func (size viewportSize) positive() bool {
	return size.width > 0 && size.height > 0
}

func (size viewportSize) scroll(state *ViewportState, delta, contentLines int) {
	if size.positive() {
		state.ScrollBy(delta, contentLines, size.height)
	}
}

func (s ScanScreen) viewportSize() viewportSize {
	if s.openedPath != "" {
		return s.bodySize.inset(
			fileViewerBorder.GetLeftSize()+fileViewerBorder.GetRightSize(),
			fileViewerBorder.GetTopSize()+fileViewerBorder.GetBottomSize(),
		)
	}
	return s.bodySize.inset(scanSidebarWidth+screenSplitGap, tabStripHeight)
}

func (s ScanScreen) withBodySize(size viewportSize) ScanScreen {
	s.bodySize = size
	return s.reconcileViewport()
}

func (s ScanScreen) reconcileViewport() ScanScreen {
	s.viewportSize().scroll(&s.viewports[s.activeTab], 0, s.contentLen())
	return s
}

func (h HistoryScreen) viewportSize() viewportSize {
	return h.bodySize.inset(historySidebarWidth+screenSplitGap, 0)
}

func (h HistoryScreen) withBodySize(size viewportSize) HistoryScreen {
	h.bodySize = size
	h.viewportSize().scroll(&h.detail, 0, len(h.detailLines()))
	return h
}
