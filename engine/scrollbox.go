package engine

import (
	"silver-ink/data"
	"silver-ink/render"
)

// Scrolls displays text in a scrollable frame
// can return links if chosen
// implements Processor
type ScrollBox struct {
	Rectangle   data.Rect
	Title       string
	Style       data.ScrollBoxStyle
	SelectIndex int // which line the cursor is on
	ForceSelect bool

	lines []ScrollBoxLine

	hoverIndex int // where the mouse is hovering

	Callback func(ScrollBoxCallback)

	// box
	box    Box
	inRect data.Rect
	inPos  data.Vector
	inSize data.Vector

	// helper vars
	selectOffset     int // lines the 'cursor' is from the top
	leftCursorIndex  int
	rightCursorIndex int
}

type ScrollBoxCallback struct {
	Accepted bool
	Link     data.Name
	Ext      any
}

func (sb *ScrollBox) Push() {
	sb.box.Rectangle = sb.Rectangle
	sb.box.Style = sb.Style.BoxStyle
	sb.box.Title = sb.Title
	sb.hoverIndex = -1

	// helper vars
	sb.inPos = sb.box.InsidePostion()
	sb.inSize = sb.box.InsideSize()
	sb.inRect = data.Rect{
		Position: sb.inPos,
		Size:     sb.inSize,
	}
	sb.selectOffset = int(sb.inSize.Y / 2)
	sb.leftCursorIndex = _buffer.Index(sb.inPos.X-1, sb.inPos.Y+sb.selectOffset)
	sb.rightCursorIndex = _buffer.Index(sb.inPos.X+sb.inSize.X, sb.inPos.Y+sb.selectOffset)

	PushStack(sb)
}

type ScrollBoxLine struct {
	Text     string
	Link     data.Name
	Ext      any
	Centered bool
	Icon     data.Cell
}

func (sbox *ScrollBox) ScrollUp() {
	if sbox.SelectIndex > 0 {
		sbox.SelectIndex--
		sbox.setMouseHover(PeekMouseState())
		Dirty()
	}
}

func (sbox *ScrollBox) ScrollDown() {
	if sbox.SelectIndex < len(sbox.lines)-1 {
		sbox.SelectIndex++
		sbox.setMouseHover(PeekMouseState())
		Dirty()
	}
}

func (sbox *ScrollBox) actionCancel() {
	if sbox.ForceSelect {
		return
	}

	PopStack()

	if sbox.Callback != nil {
		sbox.Callback(ScrollBoxCallback{
			Accepted: false,
			Link:     data.DefaultName,
		})
	}
}

func (sbox *ScrollBox) actionAccept() {
	// validations
	if sbox.SelectIndex < 0 || sbox.SelectIndex >= len(sbox.lines) {
		return
	}

	line := sbox.lines[sbox.SelectIndex]
	if sbox.hoverIndex >= 0 && sbox.hoverIndex < len(sbox.lines) {
		line = sbox.lines[sbox.hoverIndex]
	}

	if line.Link.Empty() && sbox.ForceSelect {
		return
	}

	PopStack()

	if sbox.Callback != nil {
		sbox.Callback(ScrollBoxCallback{
			Accepted: true,
			Link:     line.Link,
			Ext:      line.Ext,
		})
	}
}

func (sbox *ScrollBox) AppendLine(line ScrollBoxLine) {
	sbox.lines = append(sbox.lines, line)
}

func (sbox *ScrollBox) AppendText(text string, centered bool) {
	sbox.lines = append(sbox.lines, ScrollBoxLine{
		Text:     text,
		Centered: centered,
	})
}

func (sbox *ScrollBox) AppendLink(text string, link data.Name, ext any) {
	sbox.lines = append(sbox.lines, ScrollBoxLine{
		Text: text,
		Link: link,
		Ext:  ext,
	})
}

func (sbox *ScrollBox) Rect() data.Rect {
	return sbox.Rectangle
}

// processor
func (sbox *ScrollBox) Render(buffer *render.Buffer, top bool) {
	sbox.box.Render(buffer, data.VectorZero)

	// draw cursor
	if sbox.Style.CursorLeft.Char != 0 {
		buffer.SetCell(sbox.leftCursorIndex, sbox.Style.CursorLeft)
	}
	if sbox.Style.CursorRight.Char != 0 {
		buffer.SetCell(sbox.rightCursorIndex, sbox.Style.CursorRight)
	}

	// draw lines
	indexOffset := sbox.SelectIndex - sbox.selectOffset
	index := buffer.Index(sbox.inPos.X+1, sbox.inPos.Y)
	for i := 0; i < sbox.inSize.Y; i++ {
		if j := i + indexOffset; j >= 0 && j < len(sbox.lines) {
			sbox.drawLine(j, index, buffer)
		}
		index += buffer.Stride()
	}
}

func (sbox *ScrollBox) drawLine(line int, index int, buffer *render.Buffer) {
	l := sbox.lines[line]
	color := sbox.Style.TextColor

	if !l.Link.Empty() {
		if l.Icon.Char == 0 && l.Icon.Color == 0 {
			buffer.SetCell(index, sbox.Style.LinkIcon)
		} else {
			buffer.SetCell(index, l.Icon)
		}
		index += 2
		color = sbox.Style.LinkColor
		if sbox.hoverIndex == line || (sbox.hoverIndex == -1 && sbox.SelectIndex == line) {
			color = sbox.Style.LinkHoverColor
		}
	} else if l.Centered {
		index += (int(sbox.inSize.X) - len(l.Text)) / 2
		color = sbox.Style.AltTextColor
	}

	buffer.WriteText(index, color, l.Text)
}

func (sbox *ScrollBox) Active() {
	FPS(15)
}

func (sbox *ScrollBox) Process() {
	sbox.processAction()
	sbox.processMouse()
}

func (sbox *ScrollBox) processAction() {
	action := GetAction()
	if action == data.DefaultName {
		return
	}
	switch action {
	case INPUT_UP:
		sbox.ScrollUp()
	case INPUT_DOWN:
		sbox.ScrollDown()
	case INPUT_MOUSE_ACCEPT:
		fallthrough
	case INPUT_ACCEPT:
		sbox.actionAccept()
	case INPUT_CANCEL:
		sbox.actionCancel()
	case INPUT_SCROLL_UP:
		sbox.ScrollUp()
	case INPUT_SCROLL_DOWN:
		sbox.ScrollDown()
	}
}

func (sbox *ScrollBox) processMouse() {
	for _, ms := range PullMouseState() {
		sbox.setMouseHover(ms)
		Dirty()
	}
}

func (sbox *ScrollBox) setMouseHover(ms MouseState) {
	if !sbox.inRect.IsInside(ms.Position) {
		sbox.hoverIndex = -1
		return
	}

	ms.Position = ms.Position.Sub(sbox.inPos)

	indexOffset := sbox.SelectIndex - sbox.selectOffset
	hoverindex := indexOffset + ms.Position.Y
	if hoverindex >= 0 && hoverindex < len(sbox.lines) {
		sbox.hoverIndex = hoverindex
	} else {
		sbox.hoverIndex = -1
	}
}
