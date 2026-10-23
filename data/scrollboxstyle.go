package data

type ScrollBoxStyle struct {
	BoxStyle
	CursorLeft   Cell  `json:"CursorLeft"`
	CursorRight  Cell  `json:"CursorRight"`
	TextColor    Color `json:"TextColor"`
	AltTextColor Color `json:"AltTextColor"`
	LinkIcon     Cell  `json:"LinkIcon"`
	LinkColor      Color `json:"LinkColor"`
	LinkHoverColor Color `json:"LinkHoverColor"`
}