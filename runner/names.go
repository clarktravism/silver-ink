package main

import (
	"silver-ink/data"
)

// commands
var CHAR = data.NewName("CHAR")
var COLOR = data.NewName("COLOR")
var END = data.NewName("END")
var RESTART = data.NewName("RESTART")

// object find
var NONE = data.NewName("NONE")
var SELF = data.NewName("SELF")
var ALL = data.NewName("ALL")
var OTHERS = data.NewName("OTHERS")
var DIR = data.NewName("DIR")

// misc keywords / counters
var FGCOLOR = data.NewName("FGCOLOR")
var BGCOLOR = data.NewName("BGCOLOR")
var STEPX = data.NewName("STEPX")
var STEPY = data.NewName("STEPY")
var THISX = data.NewName("THISX")
var THISY = data.NewName("THISY")
var PRIORITY = data.NewName("PRIORITY")
var FOLLOW = data.NewName("FOLLOW")
var ENTERX = data.NewName("ENTERX")
var ENTERY = data.NewName("ENTERY")
var PLAYERX = data.NewName("PLAYERX")
var PLAYERY = data.NewName("PLAYERY")
var BOARDID = data.NewName("BOARDID")

// built-in events
var START_GAME = data.NewName("START_GAME")
var SAVE_GAME = data.NewName("SAVE_GAME")
var LOAD_GAME = data.NewName("LOAD_GAME")
var TICK = data.NewName("TICK")
var THUD = data.NewName("THUD")
var ENTER = data.NewName("ENTER")
var READY = data.NewName("READY")
var PUSHED = data.NewName("PUSHED")
var STEP = data.NewName("STEP")

// actions
var INPUT_UP = data.NewName("INPUT_UP")
var INPUT_DOWN = data.NewName("INPUT_DOWN")
var INPUT_LEFT = data.NewName("INPUT_LEFT")
var INPUT_RIGHT = data.NewName("INPUT_RIGHT")
var INPUT_ACCEPT = data.NewName("INPUT_ACCEPT")
var INPUT_CANCEL = data.NewName("INPUT_CANCEL")

// hotkeys
var HK_PLAY = data.NewName("HK_PLAY")
var HK_LOAD = data.NewName("HK_LOAD")
var HK_QUIT = data.NewName("HK_QUIT")
var HK_SAVE = data.NewName("HK_SAVE")