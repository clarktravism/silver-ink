package main

import (
	"silver-ink/data"
	"silver-ink/engine"
)

func processObjects() {

	engine.Board.FreeObjects()

	if engine.Board.ObjectID >= len(engine.Board.Objects) || engine.Board.ObjectID < 0 {
		engine.Board.ObjectID = 0
	}

	if engine.Board.Player != -1 {
		pos := engine.Board.Objects[engine.Board.Player].Position
		EventAt(pos, STEP, 0)
	}

	for engine.Board.ObjectID < len(engine.Board.Objects) {
		i := engine.Board.ObjectID
		engine.Board.ObjectID++
		obj := &engine.Board.Objects[i]
		if !obj.Free && canTick(obj.Cycle, engine.Board.Tick, i) {
			// run script
			if ctx := NewObjectContext(i); ctx != nil {
				RunScript(ctx)
			}
			// walk
			WalkObject(i, obj)
		}
	}
	engine.Board.ObjectID = 0
}

func WalkObject(i int, obj *data.Object) {
	if obj.Free || obj.Step == data.VectorZero {
		return
	}
	pos := obj.Position.Add(obj.Step)
	if data.OutsideBoard(pos) {
		ObjectEvent(i, THUD, 0)
		return
	}
	if CollisionAt(pos, obj.Collision, i) {
		// attempt to push
		if obj.PushFlags.Has(data.CanPush) {
			if Push(pos, obj.Step, obj.Collision, 0) {
				obj.Position = pos
				UpdateCamera(i, engine.Board)
				return
			}
		}
		ObjectEvent(i, THUD, 0)
		return
	}
	obj.Position = pos
	UpdateCamera(i, engine.Board)
}

func CollisionAt(v data.Vector, collision data.BitField, id int) bool {
	if data.OutsideBoard(v) {
		return true
	}
	ele := engine.ElementsFile.GetElement(engine.Board.Tiles[data.TileID(v)].ID)
	if ele != nil && ele.Collision&collision != 0 {
		return true
	}
	objects := engine.Board.Objects
	for i := range objects {
		if i != id && !objects[i].Free && objects[i].Position == v && objects[i].Collision.Collide(collision) {
			return true
		}
	}
	return false
}

func canTick(cycle, tick uint8, id int) bool {
	if cycle == 0 {
		return false
	}
	if cycle == 1 {
		return true
	}
	return tick%cycle == uint8(id)%cycle
}

func DirToPushFlag(dir data.Vector) data.BitField {
	if dir.X == 0 && dir.Y == 1 {
		return data.PushableSouth
	}
	if dir.X == 0 && dir.Y == -1 {
		return data.PushableNorth
	}
	if dir.X == 1 && dir.Y == 0 {
		return data.PushableEast
	}
	if dir.X == -1 && dir.Y == 0 {
		return data.PushableWest
	}
	return 0
}

// return true on success, or if no collision
func Push(pos data.Vector, dir data.Vector, col data.BitField, push data.BitField) bool {
	if push == 0 {
		push = DirToPushFlag(dir)
		if push == 0 {
			return false
		}
	}

	if ok := pushElement(pos, dir, col, push); !ok {
		return false
	}
	if ok := pushObject(pos, dir, col, push); !ok {
		return false
	}

	return true
}

// return true on success, or if no collision
func pushElement(pos data.Vector, dir data.Vector, col data.BitField, push data.BitField) bool {
	if i := data.TileID(pos); i != -1 {
		tile := engine.Board.Tiles[i]
		element := engine.ElementsFile.GetElement(tile.ID)
		if element == nil {
			return false
		}
		if element.Collision.Collide(col) {
			if element.PushFlags&push != push {
				return false // cant push in this dir
			}
			nPos := pos.Add(dir)
			if CollisionAt(pos.Add(dir), element.Collision, -1) {
				if ok := Push(nPos, dir, element.Collision, push); !ok {
					// push failed
					return false
				}
			}
			// moving element
			moveElement(pos, dir)
			return true
		} else {
			// no collision
			return true
		}
	}
	// out of bounds
	return false
}

func moveElement(pos data.Vector, dir data.Vector) {
	i := data.TileID(pos)
	tile := engine.Board.Tiles[i]
	j := data.TileID(pos.Add(dir))
	engine.Board.Tiles[j] = tile
	engine.Board.Tiles[i] = data.Tile{}
	TileEvent(j, PUSHED, 0)
}

func pushObject(pos data.Vector, dir data.Vector, col data.BitField, push data.BitField) bool {
	objects := engine.Board.Objects
	for i, obj := range objects {
		// eligible object
		if !obj.Free && obj.Position == pos && obj.Collision.Collide(col) {
			// can't push in this direction
			if obj.PushFlags&push != push {
				return false
			}
			nPos := pos.Add(dir)
			if CollisionAt(nPos, obj.Collision, -1) {
				if ok := Push(nPos, dir, obj.Collision, push); !ok {
					// push failed
					return false
				}
			}
			//move object
			obj.Position = nPos
			ObjectEvent(i, PUSHED, 0)
			UpdateCamera(i, engine.Board)
		}
	}
	return true
}
