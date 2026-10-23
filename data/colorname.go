package data

import (
	"strconv"
)

// loads the default values
func init() {
	ColorNames = make(map[Name]Color)
	ColorNames[NewName("BLACK")] = 0
	ColorNames[NewName("DKBLUE")] = 1
	ColorNames[NewName("DKGREEN")] = 2
	ColorNames[NewName("DKCYAN")] = 3
	ColorNames[NewName("DKRED")] = 4
	ColorNames[NewName("DKPURPLE")] = 5
	ColorNames[NewName("BROWN")] = 6
	ColorNames[NewName("GRAY")] = 7
	ColorNames[NewName("DKGRAY")] = 8
	ColorNames[NewName("BLUE")] = 9
	ColorNames[NewName("GREEN")] = 10
	ColorNames[NewName("CYAN")] = 11
	ColorNames[NewName("RED")] = 12
	ColorNames[NewName("PURPLE")] = 13
	ColorNames[NewName("YELLOW")] = 14
	ColorNames[NewName("WHITE")] = 15
}
var ColorNames map[Name]Color

func GetColor(name Name) (Color, bool) {
	color, ok := ColorNames[name]
	if !ok {
		str := name.String()
		if len(str) == 2 {
			v, err := strconv.ParseUint(str, 16, 8)
			if err != nil {
				return color, false
			}
			return Color(v), true
		}
	}
	return color, ok
}

var ColorBLACK = Color(0)
var ColorDKBLUE = Color(1)
var ColorDKGREEN = Color(2)
var ColorDKCYAN = Color(3)
var ColorDKRED = Color(4)
var ColorDKPURPLE = Color(5)
var ColorBROWN = Color(6)
var ColorGRAY = Color(7)
var ColorDKGRAY = Color(8)
var ColorBLUE = Color(9)
var ColorGREEN = Color(10)
var ColorCYAN = Color(11)
var ColorRED = Color(12)
var ColorPURPLE = Color(13)
var ColorYELLOW = Color(14)
var ColorWHITE = Color(15)
