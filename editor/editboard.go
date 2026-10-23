package main

import (
	"log"
	"silver-ink/data"
	"silver-ink/engine"
	"silver-ink/render"
)

type EditBoard struct {
	// tool vars
	tool    Tool
	brush   Brush
	prefab  *data.Prefab
	element *data.Element

	// controls
	container   engine.Container
	sidePanel   sidePanel
	scrollBoard ScrollBoard
}

func (eb *EditBoard) Push() {
	eb.init()
	engine.PushStack(eb)
}

func (eb *EditBoard) init() {
	// container
	eb.container = engine.Container{
		Rectangle: data.NewRect(0, 0, 80, 60),
		//Background: data.Cell{Char: 0, Color: data.ColorBLACK},
	}
	eb.container.Init()

	// scroll board
	eb.scrollBoard = ScrollBoard{
		Rectangle: data.NewRect(0, 0, 62, 60),
	}
	eb.scrollBoard.Init()

	// side panel
	eb.sidePanel = sidePanel{}

	eb.container.AddChild(&eb.scrollBoard)
	eb.container.AddChild(&eb.sidePanel)

	// signal hooks
	eb.scrollBoard.BoardMouseChange.Listen(eb.onBoardMouseChange)
	eb.sidePanel.bufferTool.ItemClicked.Listen(eb.onBufferToolClicked)
	eb.sidePanel.bufferTool.ItemHover.Listen(eb.onBufferToolHover)
	eb.sidePanel.colorSelect.ColorChanged.Listen(eb.onColorChanged)
	eb.sidePanel.charSelect.CharChanged.Listen(eb.onCharChanged)

	// button hooks
	eb.sidePanel.objectButton.Click.Listen(eb.onObjectButtonClick)
	eb.sidePanel.elementButton.Click.Listen(eb.onElementButtonClick)
	eb.sidePanel.drawButton.Click.Listen(eb.onDrawClick)
	eb.sidePanel.remobjButton.Click.Listen(eb.onRemoveObjClick)
	eb.sidePanel.textButton.Click.Listen(eb.onTextClick)
	eb.sidePanel.fillButton.Click.Listen(eb.onFillClick)
	eb.sidePanel.quitButton.Click.Listen(eb.onQuitClick)
	eb.sidePanel.saveButton.Click.Listen(eb.onSaveClick)
	eb.sidePanel.loadButton.Click.Listen(eb.onLoadClick)
	eb.sidePanel.boardsButton.Click.Listen(eb.onBoardsClick)
	eb.sidePanel.boardinfoButton.Click.Listen(eb.onBoardInfoClick)
	eb.sidePanel.worldinfoButton.Click.Listen(eb.onWorldInfoClick)
	eb.sidePanel.viewButton.Click.Listen(eb.onViewButtonClick)
	eb.sidePanel.prefabButton.Click.Listen(eb.onPrefabButtonClick)

	// tool and brush
	eb.toolDraw()
	eb.SetBrush(Brush{
		ElementID: 0,
		Cell:      data.Cell{Char: 0, Color: data.ColorWHITE},
	})
}

// brush
func (eb *EditBoard) SetBrush(value Brush) {
	if eb.brush == value {
		return
	}
	eb.brush = value
	eb.tool.SetBrush(value)
	eb.sidePanel.SetBrush(value)
}

func (eb *EditBoard) onColorChanged(args ColorChangedArgs) {
	if args.Color != eb.brush.Cell.Color {
		brush := eb.brush
		brush.Cell.Color = args.Color
		eb.SetBrush(brush)
	}
}

func (eb *EditBoard) onCharChanged(args CharChangedArgs) {
	if args.Char != eb.brush.Cell.Char {
		brush := eb.brush
		brush.Cell.Char = args.Char
		eb.SetBrush(brush)
	}
}


// edit object
func (eb *EditBoard) onObjectButtonClick(engine.MouseState) {
	eb.selectObject()
}

func (eb *EditBoard) selectObject() {
	if eb.brush.Prefab == nil {
		prefab := defaultPrefab
		prefab.Cell = eb.brush.Cell
		eb.brush.Prefab = &prefab
	}

	eo := EditObject{
		Callback: eb.onEditObjectCallback,
		Prefab:   eb.brush.Prefab,
	}
	eo.Push()
}

func (eb *EditBoard) onEditObjectCallback(callback EditObjectCallback) {
	if callback.Accepted {
		brush := Brush{
			Prefab: &callback.Data,
			Cell:   callback.Data.Cell,
		}
		eb.SetBrush(brush)
	}
}

var defaultPrefab = data.Prefab{
	Cycle: 1,
	Collision: data.Collision1 | data.Collision2,
}

// select element
func (eb *EditBoard) onElementButtonClick(engine.MouseState) {
	eb.selectElement()
}

func (eb *EditBoard) selectElement() {
	se := SelectElement{
		Callback: eb.onSelectElementCallback,
		Color:    eb.brush.Cell.Color,
	}
	se.Push()
}

func (eb *EditBoard) onSelectElementCallback(callback SelectElementCallback) {
	if callback.Accepted {
		brush := Brush{
			ElementID: callback.Element.ID,
			Cell:      data.Cell{Char: callback.Element.Char, Color: eb.brush.Cell.Color},
		}
		eb.SetBrush(brush)
	}
}

// select prefab
func (eb *EditBoard) onPrefabButtonClick(engine.MouseState) {
	eb.selectPrefab()
}

func (eb *EditBoard) selectPrefab() {
	se := SelectPrefab{
		Callback: eb.onSelectPrefabCallback,
	}
	se.Push()
}

func (eb *EditBoard) onSelectPrefabCallback(callback SelectPrefabCallback) {
	if callback.Accepted {
		brush := Brush{
			Prefab: &callback.Prefab,
			Cell:   callback.Prefab.Cell,
		}
		eb.SetBrush(brush)
	}
}

// quit
func (eb *EditBoard) onQuitClick(engine.MouseState) {
	onQuit()
}

// save board
func (eb *EditBoard) onSaveClick(engine.MouseState) {
	eb.saveBoard()
}

func (eb *EditBoard) saveBoard() {
	SaveBoard()
	PushPopupMessage("Board Saved")
}

// select board
func (eb *EditBoard) onBoardsClick(engine.MouseState) {
	eb.selectBoard()
}

func (eb *EditBoard) selectBoard() {
	sb := SelectBoard{
		Callback: eb.onSelectBoardCallback,
	}
	sb.Push()
}

func (eb *EditBoard) onSelectBoardCallback(callback SelectBoardCallback) {
	if callback.Accepted {

		// auto-save for now...
		engine.WorldPath.SaveBoard(engine.Board)

		board, err := engine.WorldPath.LoadBoard(callback.Name)
		if err != nil {
			log.Printf("error while loading board: %v\r\n", err)
			return
		}
		engine.SetBoard(board)
		eb.sidePanel.boardLabel.SetText(engine.Board.Name.String())
		engine.Dirty()
	}
}

// edit board info
func (eb *EditBoard) onBoardInfoClick(engine.MouseState) {
	eb.openBoardInfo()
}

func (eb *EditBoard) openBoardInfo() {
	ebi := EditBoardInfo{
		Callback: func() {
			eb.sidePanel.boardLabel.SetText(engine.Board.Name.String())
		},
	}
	ebi.Push()
}

// edit world info
func (eb *EditBoard) onWorldInfoClick(engine.MouseState) {
	eb.openWorldInfo()
}

func (eb *EditBoard) openWorldInfo() {
	ewi := EditWorldInfo{}
	ewi.Push()
}

// view button
func (eb *EditBoard) onViewButtonClick(engine.MouseState) {
	eb.changeView()
}

func (eb *EditBoard) changeView() {
	mode := (eb.scrollBoard.BoardViewMode + 1) % 3
	eb.scrollBoard.SetViewMode(mode)
	switch mode {
	case BoardViewModeAll:
		eb.sidePanel.viewButton.RightText = " View: All "
	case BoardViewModeTiles:
		eb.sidePanel.viewButton.RightText = " View: Tiles "
	case BoardViewModeObjects:
		eb.sidePanel.viewButton.RightText = " View: Objs "
	}
	engine.Dirty()
}

// tools
func (eb *EditBoard) SetTool(tool Tool) {
	if eb.tool != nil {
		eb.tool.Close()
	}
	eb.tool = tool
	eb.tool.Init()
	eb.tool.SetBrush(eb.brush)
}

func (eb *EditBoard) onDrawClick(engine.MouseState) {
	eb.toolDraw()
}

func (eb *EditBoard) toolDraw() {
	eb.SetTool(&DrawTool{ScrollBoard: &eb.scrollBoard})
	eb.sidePanel.toolLabelBox.SetText(" Tool:Draw ")
	engine.Dirty()
}

func (eb *EditBoard) onRemoveObjClick(engine.MouseState) {
	eb.toolRemoveObj()
}

func (eb *EditBoard) toolRemoveObj() {
	eb.SetTool(&RemoveObjTool{ScrollBoard: &eb.scrollBoard})
	eb.sidePanel.toolLabelBox.SetText(" Tool:Remove ")
	engine.Dirty()
}

func (eb *EditBoard) onTextClick(engine.MouseState) {
	eb.toolText()
}

func (eb *EditBoard) toolText() {
	eb.SetTool(&TextTool{ScrollBoard: &eb.scrollBoard})
	eb.sidePanel.toolLabelBox.SetText(" Tool:Text ")
	engine.Dirty()
}

func (eb *EditBoard) onFillClick(engine.MouseState) {
	eb.toolFill()
}

func (eb *EditBoard) toolFill() {
	eb.SetTool(&FillTool{ScrollBoard: &eb.scrollBoard})
	eb.sidePanel.toolLabelBox.SetText(" Tool:Fill ")
	engine.Dirty()
}

func (eb *EditBoard) onBufferToolClicked(args BufferItemArgs) {
	if args.LeftButtonClick {
		// save brush to buffer
		args.Brush.Prefab = eb.brush.Prefab
		args.Brush.ElementID = eb.brush.ElementID
		args.Brush.Cell = eb.brush.Cell
	} else if args.RightButtonClick {
		// copy buffer to brush
		brush := *args.Brush
		eb.SetBrush(brush)
	}
	engine.Dirty()
}

func (eb *EditBoard) onBufferToolHover(args BufferItemArgs) {
	if args.Brush != nil {
		eb.sidePanel.infoBrushView.SetBrush(*args.Brush)
	}
	engine.Dirty()
}

// board mouse
func (eb *EditBoard) onBoardMouseChange(arg BoardMouseArgs) {

	// get brush at location
	brush, brushOk := eb.brushAt(arg.BoardPosition)
	if brushOk {
		eb.sidePanel.infoBrushView.SetBrush(brush)
	}
	eb.sidePanel.setInfoCoords(arg.BoardPosition)

	// special right click
	if arg.RightButtonClick {
		if engine.IsCtrl(arg.Modifier) {
			eb.inspectObject(arg.BoardPosition)
		} else if brushOk {
			eb.SetBrush(brush)
		}
		return
	}

	eb.tool.Apply(arg.BoardPosition, arg.MouseState)
	engine.Dirty()
}

func (eb *EditBoard) inspectObject(pos data.Vector) {
	if eb.scrollBoard.BoardViewMode == BoardViewModeTiles {
		return
	}

	i := engine.Board.ObjectAt(pos)
	if i == -1 {
		return
	}

	object := &engine.Board.Objects[i]
	pf := object.Prefab()

	if engine.Board.Follow == i {
		pf.Follow = true
	}
	if engine.Board.Player == i {
		pf.Player = true
	}

	eo := EditObject{
		Prefab: &pf,
		Callback: func(callback EditObjectCallback) {
			if callback.Accepted {
				object.Name = callback.Data.Name
				object.Group = callback.Data.Group
				object.Cycle = callback.Data.Cycle
				object.ScriptName = callback.Data.ScriptName
				object.Target = callback.Data.Target
				object.Collision = callback.Data.Collision
				object.DrawFlags = callback.Data.DrawFlags
				object.PushFlags = callback.Data.PushFlags
				object.Cell = callback.Data.Cell
				object.Counters = callback.Data.Counters

				if callback.Data.Player && engine.Board.Player != i {
					engine.Board.Player = i
				} else if !callback.Data.Player && engine.Board.Player == i {
					engine.Board.Player = -1
				}

				if callback.Data.Follow && engine.Board.Follow != i {
					engine.Board.Follow = i
				} else if !callback.Data.Follow && engine.Board.Follow == i {
					engine.Board.Follow = -1
				}

				engine.Dirty()
			}
		},
	}
	eo.Push()
}

func (eb *EditBoard) brushAt(pos data.Vector) (Brush, bool) {
	if pos.X < 0 || pos.Y < 0 || pos.X >= data.BoardWidth || pos.Y >= data.BoardHeight {
		return Brush{}, false
	}

	if eb.scrollBoard.BoardViewMode != BoardViewModeTiles {
		if i := engine.Board.ObjectAt(pos); i != -1 {
			object := &engine.Board.Objects[i]
			pf := object.Prefab()
			pf.Player = engine.Board.Player == i
			pf.Follow = engine.Board.Follow == i
			brush := Brush{
				Prefab: &pf,
				Cell:   pf.Cell,
			}
			return brush, true
		}
	}
	if eb.scrollBoard.BoardViewMode != BoardViewModeObjects {
		tid := data.TileID(pos)
		tile := engine.Board.Tiles[tid]
		element := engine.ElementsFile.GetElement(tile.ID)
		if element != nil {
			brush := Brush{
				ElementID: element.ID,
				Cell:      tile.Cell,
			}
			return brush, true
		}
	}
	return Brush{}, false
}

// world management
func (eb *EditBoard) onLoadClick(engine.MouseState) {
	eb.loadWorld()
}

func (eb *EditBoard) loadWorld() {
	SaveBoard()
	sw := SelectWorld{
		Force:    false,
		Callback: eb.callbackSelectWorld,
	}
	sw.Push()
}

func (eb *EditBoard) callbackSelectWorld(callback SelectWorldCallback) {
	if callback.Accepted {
		eb.sidePanel.worldLabel.SetText(callback.WorldName)
		LoadWorld(callback.WorldPath, callback.WorldName)
	}
}

// debug
func (eb *EditBoard) debugGame() {
	SaveBoard()
	PushDebug()
	/*path, err := exec.LookPath("./runner")
	if err != nil {
		log.Println(err)
		return
	}
	cmd := exec.Command(path, fmt.Sprintf("-w=%v", engine.WorldPath.Path), fmt.Sprintf("-b=%v", engine.Board.Name.String()))
	if err := cmd.Run(); err != nil {
		log.Println(err)
	}*/
}

// for stack
func (eb *EditBoard) Render(buffer *render.Buffer, top bool) {
	eb.container.Render(buffer, data.VectorZero)
}

func (eb *EditBoard) Active() {
	engine.FPS(60)
	eb.sidePanel.worldLabel.SetText(engine.WorldFile.Name)
	eb.sidePanel.boardLabel.SetText(engine.Board.Name.String())
}

func (eb *EditBoard) Process() {
	for _, ms := range engine.PullMouseState() {
		eb.container.ProcessMouse(ms)
	}

	if !engine.CycleTick(8) {
		return
	}

	if action := engine.GetAction(); action != data.DefaultName {
		switch action {
		case HK_WORLD_INFO:
			eb.openWorldInfo()
		case HK_BOARDS:
			eb.selectBoard()
		case HK_BOARD_INFO:
			eb.openBoardInfo()
		case HK_VIEW:
			eb.changeView()
		case HK_OBJECT:
			eb.selectObject()
		case HK_PREFAB:
			eb.selectPrefab()
		case HK_ELEMENT:
			eb.selectElement()
		case HK_DRAW:
			eb.toolDraw()
		case HK_TEXT:
			eb.toolText()
		case HK_FILL:
			eb.toolFill()
		case HK_REMOVE_OBJ:
			eb.toolRemoveObj()
		case HK_SAVE:
			eb.saveBoard()
		case HK_LOAD:
			eb.loadWorld()
		case HK_QUIT:
			onQuit()
		case engine.INPUT_UP:
			eb.scrollBoard.Move(data.Vector{X: 0, Y: -4})
		case engine.INPUT_DOWN:
			eb.scrollBoard.Move(data.Vector{X: 0, Y: 4})
		case engine.INPUT_LEFT:
			eb.scrollBoard.Move(data.Vector{X: -4, Y: 0})
		case engine.INPUT_RIGHT:
			eb.scrollBoard.Move(data.Vector{X: 4, Y: 0})
		case HK_DEBUG:
			eb.debugGame()
		}
	}
}

// hotkey names
var (
	HK_WORLD_INFO = data.NewName("HK_WORLD_INFO")
	HK_BOARDS     = data.NewName("HK_BOARDS")
	HK_BOARD_INFO = data.NewName("HK_BOARD_INFO")
	HK_VIEW       = data.NewName("HK_VIEW")
	HK_OBJECT     = data.NewName("HK_OBJECT")
	HK_PREFAB     = data.NewName("HK_PREFAB")
	HK_ELEMENT    = data.NewName("HK_ELEMENT")
	HK_REMOVE_OBJ = data.NewName("HK_REMOVE_OBJ")
	HK_DRAW       = data.NewName("HK_DRAW")
	HK_TEXT       = data.NewName("HK_TEXT")
	HK_FILL       = data.NewName("HK_FILL")
	HK_LOAD       = data.NewName("HK_LOAD")
	HK_SAVE       = data.NewName("HK_SAVE")
	HK_QUIT       = data.NewName("HK_QUIT")
	HK_DEBUG = data.NewName("HK_DEBUG")
)
