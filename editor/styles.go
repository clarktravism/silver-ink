package main

import (
	"silver-ink/data"
	"silver-ink/engine"
)

var boxTemplate = data.BoxStyle{
	PaddingTop:    1,
	PaddingBottom: 1,
	PaddingLeft:   1,
	PaddingRight:  1,
	//Background:        data.Cell{Char: 0, Color: data.NewColor(0, 1)},
	BorderTopRight:    data.Cell{Char: 184, Color: data.NewColor(15, 1)},
	BorderTopLeft:     data.Cell{Char: 213, Color: data.NewColor(15, 1)},
	BorderBottomRight: data.Cell{Char: 190, Color: data.NewColor(15, 1)},
	BorderBottomLeft:  data.Cell{Char: 212, Color: data.NewColor(15, 1)},
	BorderTop:         data.Cell{Char: 205, Color: data.NewColor(15, 1)},
	BorderBottom:      data.Cell{Char: 205, Color: data.NewColor(15, 1)},
	BorderLeft:        data.Cell{Char: 179, Color: data.NewColor(15, 1)},
	BorderRight:       data.Cell{Char: 179, Color: data.NewColor(15, 1)},
}

var red_button_style = engine.ButtonStyle{
	Background:      data.NewColor(data.ColorWHITE, data.ColorDKRED),
	HoverBackground: data.NewColor(data.ColorBLACK, data.ColorYELLOW),
}

var green_button_style = engine.ButtonStyle{
	Background:      data.NewColor(data.ColorWHITE, data.ColorDKGREEN),
	HoverBackground: data.NewColor(data.ColorBLACK, data.ColorYELLOW),
}

var dialog_box_style = data.BoxStyle{
	BorderTopRight:        data.Cell{Char: 191, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
	BorderTopLeft:         data.Cell{Char: 218, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
	BorderBottomRight:     data.Cell{Char: 217, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
	BorderBottomLeft:      data.Cell{Char: 192, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
	BorderTop:             data.Cell{Char: 196, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
	BorderBottom:          data.Cell{Char: 196, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
	BorderLeft:            data.Cell{Char: 179, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
	BorderRight:           data.Cell{Char: 179, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
	HorizontalDivide:      data.Cell{Char: 196, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
	HorizontalDivideLeft:  data.Cell{Char: 195, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
	HorizontalDivideRight: data.Cell{Char: 180, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
	Background:            data.Cell{Char: 0, Color: data.NewColor(data.ColorWHITE, data.ColorDKGRAY)},
	TitleBackground:       data.Cell{Char: 0, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
}

var inside_box_style = data.BoxStyle{
	BorderTopRight:    data.Cell{Char: 191, Color: data.NewColor(data.ColorBLACK, data.ColorDKGRAY)},
	BorderTopLeft:     data.Cell{Char: 218, Color: data.NewColor(data.ColorBLACK, data.ColorDKGRAY)},
	BorderBottomRight: data.Cell{Char: 217, Color: data.NewColor(data.ColorBLACK, data.ColorDKGRAY)},
	BorderBottomLeft:  data.Cell{Char: 192, Color: data.NewColor(data.ColorBLACK, data.ColorDKGRAY)},
	BorderTop:         data.Cell{Char: 196, Color: data.NewColor(data.ColorBLACK, data.ColorDKGRAY)},
	BorderBottom:      data.Cell{Char: 196, Color: data.NewColor(data.ColorBLACK, data.ColorDKGRAY)},
	BorderLeft:        data.Cell{Char: 179, Color: data.NewColor(data.ColorBLACK, data.ColorDKGRAY)},
	BorderRight:       data.Cell{Char: 179, Color: data.NewColor(data.ColorBLACK, data.ColorDKGRAY)},
	Background:        data.Cell{Char: 0, Color: data.NewColor(data.ColorWHITE, data.ColorDKGRAY)},
}

var text_prompt_style = data.BoxStyle{
	BorderTopRight:        data.Cell{Char: 191, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
	BorderTopLeft:         data.Cell{Char: 218, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
	BorderBottomRight:     data.Cell{Char: 217, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
	BorderBottomLeft:      data.Cell{Char: 192, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
	BorderTop:             data.Cell{Char: 196, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
	BorderBottom:          data.Cell{Char: 196, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
	BorderLeft:            data.Cell{Char: 179, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
	BorderRight:           data.Cell{Char: 179, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
	HorizontalDivide:      data.Cell{Char: 196, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
	HorizontalDivideLeft:  data.Cell{Char: 195, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
	HorizontalDivideRight: data.Cell{Char: 180, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
	Background:            data.Cell{Char: 0, Color: data.NewColor(data.ColorWHITE, data.ColorBLACK)},
	TitleBackground:       data.Cell{Char: 0, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
}

var scroll_box_style = data.ScrollBoxStyle{
	BoxStyle: data.BoxStyle{
		PaddingTop:            0,
		PaddingBottom:         0,
		PaddingLeft:           2,
		PaddingRight:          2,
		BorderTopRight:        data.Cell{Char: 191, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
		BorderTopLeft:         data.Cell{Char: 218, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
		BorderBottomRight:     data.Cell{Char: 217, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
		BorderBottomLeft:      data.Cell{Char: 192, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
		BorderTop:             data.Cell{Char: 196, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
		BorderBottom:          data.Cell{Char: 196, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
		BorderLeft:            data.Cell{Char: 179, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
		BorderRight:           data.Cell{Char: 179, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
		HorizontalDivide:      data.Cell{Char: 196, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
		HorizontalDivideLeft:  data.Cell{Char: 195, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
		HorizontalDivideRight: data.Cell{Char: 180, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
		Background:            data.Cell{Char: 0, Color: data.NewColor(data.ColorWHITE, data.ColorBLACK)},
		TitleBackground:       data.Cell{Char: 0, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
	},
	TextColor:      data.ColorGRAY,
	AltTextColor:   data.ColorWHITE,
	LinkColor:      data.ColorWHITE,
	LinkHoverColor: data.ColorYELLOW,
	LinkIcon:       data.Cell{Char: 16, Color: data.ColorYELLOW},
	CursorLeft:     data.Cell{Char: 62, Color: data.ColorRED},
	CursorRight:    data.Cell{Char: 60, Color: data.ColorRED},
}

var scroll_box_style_popup = data.ScrollBoxStyle{
	BoxStyle: data.BoxStyle{
		PaddingTop:            1,
		PaddingBottom:         1,
		PaddingLeft:           1,
		PaddingRight:          1,
		BorderTopRight:        data.Cell{Char: 191, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
		BorderTopLeft:         data.Cell{Char: 218, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
		BorderBottomRight:     data.Cell{Char: 217, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
		BorderBottomLeft:      data.Cell{Char: 192, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
		BorderTop:             data.Cell{Char: 196, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
		BorderBottom:          data.Cell{Char: 196, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
		BorderLeft:            data.Cell{Char: 179, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
		BorderRight:           data.Cell{Char: 179, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
		HorizontalDivide:      data.Cell{Char: 196, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
		HorizontalDivideLeft:  data.Cell{Char: 195, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
		HorizontalDivideRight: data.Cell{Char: 180, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
		Background:            data.Cell{Char: 0, Color: data.NewColor(data.ColorWHITE, data.ColorBLACK)},
		TitleBackground:       data.Cell{Char: 0, Color: data.NewColor(data.ColorGRAY, data.ColorBLACK)},
	},
	TextColor:      data.ColorGRAY,
	AltTextColor:   data.ColorWHITE,
	LinkColor:      data.ColorWHITE,
	LinkHoverColor: data.ColorYELLOW,
	LinkIcon:       data.Cell{Char: 16, Color: data.ColorYELLOW},
	CursorLeft:     data.Cell{Char: 0, Color: 0},
	CursorRight:    data.Cell{Char: 0, Color: 0},
}

var checkbox_style = engine.CheckBoxStyle{
	CheckedOn:  data.Cell{Char: 'Y', Color: data.ColorGREEN},
	CheckedOff: data.Cell{Char: 'N', Color: data.ColorRED},
}

var gray_button_style = engine.ButtonStyle{
	Background:      data.NewColor(data.ColorBLACK, data.ColorGRAY),
	HoverBackground: data.NewColor(data.ColorBLACK, data.ColorYELLOW),
}

func PushPopupMessage(text string) {
	size := data.Vector{X: len(text) + 6, Y: 5}
	sb := engine.ScrollBox{
		Rectangle: data.Rect{
			Position: engine.CenterPos(size),
			Size: size,
		},
		Style: scroll_box_style_popup,
	}
	sb.AppendText(text, false)
	sb.Push()
}