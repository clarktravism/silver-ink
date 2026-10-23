package data

type BoxStyle struct {
	PaddingTop            int  `json:"PaddingTop"`
	PaddingBottom         int  `json:"PaddingBottom"`
	PaddingLeft           int  `json:"PaddingLeft"`
	PaddingRight          int  `json:"PaddingRight"`
	BorderTopRight        Cell `json:"BorderTopRight"`
	BorderTopLeft         Cell `json:"BorderTopLeft"`
	BorderBottomRight     Cell `json:"BorderBottomRight"`
	BorderBottomLeft      Cell `json:"BorderBottomLeft"`
	BorderTop             Cell `json:"BorderTop"`
	BorderBottom          Cell `json:"BorderBottom"`
	BorderLeft            Cell `json:"BorderLeft"`
	BorderRight           Cell `json:"BorderRight"`
	HorizontalDivide      Cell `json:"HorizontalDivide"`
	HorizontalDivideLeft  Cell `json:"HorizontalDivideLeft"`
	HorizontalDivideRight Cell `json:"HorizontalDivideRight"`
	Background            Cell `json:"Background"`
	TitleBackground       Cell `json:"TitleBackground"`
}
