
# Silver Ink Guide

This guide is still in progress

### World

When a new world is made in the editor, a folder is created under /worlds, is the world folder. Inside the folder, several files and folders are included to allow the runner to function.

##### /boards folder
This is where board files (.sib) are placed by the editor. You may copy, delete and rename as you like.

See: [Boards](#boards)

##### /charsets folder
This is where charsets for your world are kept. They must be .png format, 128x128. Black color is 'background', any other color is 'foreground'.

##### /palettes folder
This is where palettes for your world are kept. They must be .png format and exactly 16 pixels total. Reading left to right, top to bottom will be the colors and indexes for that palette.

##### /saves folder
Where the save files are kept, for saving and loading game states.

##### /scripts folder
Where the script files (.obj) are kept. Scripts are managed outside the editor, using your preferred text editor.

See: [INK-OOP](#ink-oop)

##### elements.json
This file defines your world's elements. Elements are what controls the behavior of tiles.

See: [Tiles and Elements](#tiles-and-elements)

##### keybinds.json
Keybinds for the world, non-hotkey bindings sent to the player as an event.

See: [Player](#player), [Events](#events)

##### prefabs.json
This file defines your world's prefabs. Prefabs are predefined objects that are ready to use.

See: [Objects and Prefabs](#objects-and-prefabs)
##### textboxes.json
This file defines the textboxes look and positioning. Textboxes are used for dialogs and menus within the game.

See: [TextBoxes](#text-boxes)
##### views.json
This file defines the views. Views are used to set the camera and display information.

See: [Views](#views)
##### world.json
This file set world values:
TitleBoard: Name of which board to use as the title screen.
StartBoard: Name of which board start on when the game begins.
FPS: Target frames per second to run at.
InputCycle: How often keybinds are checked for input. (1 = Every frame, 2 = Every other, etc)

Counters: Initial values for the world counters. See: [Counters](#counters)

### Editor
Quick Editor Tips:
F5 to launch the game and load the board you are editing.
The editor will auto-save when you switch boards.
Use the arrow keys to move the visible board area.
Right Click the board area to copy the element or object info.
Ctrl + Right Click to inspect/edit an object.
With 'Draw Tool' active, Shift + Left Click drag a straight line.
Left Click in the buffer to hold the element or object for later.
Right Click in the buffer to copy it back.

### Boards
Boards are 128 x 96 tiles in size, with 0,0 being the top-left corner and 127,95 the bottom-right. They can contain up to 256 objects.
 
##### Charset
Which charset this board will use.
##### Palette
Which palette this board will use.
##### Persist
If false, the board will reset when the player re-enters the board.
##### North,South,East,West
Connections to other boards when the player touches the edge.
##### Camera
Position of the top-left corner of the camera.
##### View
Id of the view the board is using. See: [Views](#views)


### Tiles and Elements

Elements are defined in elements.json in the world folder. They define how tiles should behave. When a new world is made, elements.json will be prefilled with many familiar elements ready to use.

Here is a sample for the Gem element.
```
{
	"ID" : 12,
	"Name" : "Gem",
	"Char" : 4,
	"Color" : "CYAN",
	"Collision" : "00000011",
	"Script" : "GEM",
	"PushFlags" : "PUSHABLE_ALL"
}
```

###### ID: Internal ID
0-255, must be unique.
##### Name
Name of the element, must be unique among elements and prefabs.
##### Char
Default char for the element, 0-255
##### Color
Default color for the element. Can be defined 0-255 or a color name.
##### Collision
A value of "00000011" registers as collides with water and walls, but not empties.
You may also use the byte value, which in this case would be 3.
##### Script 
Script file name, without extension. Note that scripts attached to elements have limited functionality. To summarize:
They do not run like an object, but instead respond to events.
They do not have their own local counters, but may use the board or world counters.
They cannot use the following commands:
\#FOLLOW
\#GROUP
\#LOCK
\#NAME
\#PLAYER
\#SCRIPT
\#RESTORE
\#UNLOCK
\#WALK
\#ZAP
  
##### PushFlags
How the tile to being pushed. No value indicates it cannot be pushed.
PUSHABLE_N
PUSHABLE_S
PUSHABLE_E
PUSHABLE_W
PUSHABLE_NS
PUSHABLE_EW
PUSHABLE_ALL

### Objects and Prefabs
Objects are entities on the board that exist on a separate layer than tiles. Unlike tiles, their scripts can potentially run every frame, in addition to responding to events. Maximum 256 objects are allowed on the board at a time.

Prefabs are simply predefined objects to make it easy to use within the editor or spawn during gameplay. \#BECOME, \#PUT and \#SHOOT all allow prefabs as an argument. Once a prefab is placed on the board, it becomes no different than any other object. Thus, if a prefab definition is changed, it won't retroactively change objects that were made with that prefab.

### Collision Flags

Every element and object has a set of collision flags. If one entity tries to move into a space and any of the flags match, then they collide. Technically, it allows multiple objects to occupy the same space, if they don't collide.

Flag Names:
COLLISION_1
COLLISION_2
COLLISION_3
COLLISION_4
COLLISION_5
COLLISION_6
COLLISION_7
COLLISION_8
WALLS, alias for COLLISION_1
WATER, alias for COLLISION_2
FLOOR, alias for COLLISION_3
In .json files you can set the "Collision" value like so:

```
"Collision" : "00000011"
```
Which will set COLLISION_1 and COLLISION_2. Left most digit is 8 and right most digit is 1.
  

Alternately you may set it like:
```
"Collision" : "WALLS,WATER,COLLISION_7"
```
Which is equivalent to:
```
"Collision" : "01000011"
```

Generally, walls will have COLLISION_1, water will have COLLISION_2, and floor is COLLISION_3.

"00000011" is a good collision for most solid entities, that cannot traverse water or walls.
"00000001" for projectiles that need to travel over water.
"00000101" for aquatic creatures that must stay on dry land.

Of course these aren't hard rules, you may use collision flags 8-1 however you like.

### Push Flags

Every element and object has defined set of push flags. They determine if they can push or be pushed by other entities.

Flag Names:
CAN_PUSH
PUSHABLE_N
PUSHABLE_S
PUSHABLE_E
PUSHABLE_W
PUSHABLE_NS
PUSHABLE_EW
PUSHABLE_ALL

In .json files you can set the "PushFlags" value like so:
```
"PushFlags" : "CAN_PUSH,PUSHABLE_ALL"
```
Allows it to push, but also be pushed in any direction
```
"PushFlags" : "PUSHABLE_N,PUSHABLE_E"
```
Allows it to only be pushed the north and east directions.
You can also use #SET and #CLEAR to update the flags for objects only.

### Camera and Follow
The camera doesn't strictly follow the player object, rather it follows the objected marked as 'Follow'.
The follow flag can be set in the editor, only one object can have it set.

### Counters

Counters are a signed 16 bit integers that scripts can update or read. You can use them like:
```
#SET AMMO 50
#TAKE SELF:HEALTH 10
#IF REPUTATION > 10
``` 
Or you can use them like flags:
```
#SET OLDKEY
#CLEAR SELF:POISONED
#IF BOARD:ALARM
```
Board counters are addressed with "BOARD:" prefix. They can be accessed by any object or element inside the board.

Objects counters are addressed with a "SELF:" or "THIS:" prefix. These are personal counters only seen by that object.

World counters span across boards, and do not require a prefix.
Currently, there is no hard limit on the number of counters.

### Built-in Values

  

### Events
An event is when code is executed out of turn. In other words, when an event is received it is executed immediately until it reaches a command that halts the cycle (ex: \#END \#GO \#IDLE). A script can signal its own events, even inside another event. Labels are used to specify where the event starts. If the label isn't found it does nothing. Events honor \#ZAP and \#LOCK.

#### Built-in Events
:ENTER
Sent to every element and object every time the board is entered.

:THUD
When the object is walking and collides with something else.

:PUSHED
When the element or object has been pushed.

:TICK
Called every frame for elements if they are on camera.

### Player
The player isn't a built-in element, rather it is a pointer to an object that is tracked on the board. There can only be one player at time, when a premade is adding with the player flag or when the \#PLAYER command is used, the old player is dropped. In the editor, you can set the player by editing the object.

While the player is set:
ALIGNED and CONTACT will use the player's position to determine if true.
SEEK will give a direction towards the player object.
The position PLAYER will return the player object position.
PLAYERX and PLAYERY will return the x and y of the player respectively.
If the object tries to move to the board edge, it will transition to the next board, if available.
Player scope 'PLAYER:LABEL' will send to the player for \#SEND and \#EVENT
Keybind actions will be sent to the player as an event.
 
### INK-OOP

##### Summary of Changes from ZZT-OOP
All scripts are stored in your world subfolder /scripts. Every object that requires code must reference a script file. On the startup, all scripts are evaluated and checked for (most) syntax errors. If any are found, they are written out and the game immediately closes.

\#PLAY is not yet implemented

\#BIND is no longer needed, as all objects that reference the same script share the code internally. \#ZAP and \#RESTORE still function as normal, however you cannot pre-zap labels.

Counters have replaced flags completely. \#SET, \#CLEAR and \#IF still function as before.

Hyperlinking external files is not supported.

@Name is not implemented, objects are named in the editor or can be renamed with \#NAME command.

Math operators can be used in <value> parameters:
ADD, SUB, MUL, DIV, MOD, RANDOM, ABS, DIST, DIFF

\#BECOME and \#PUT commands and ANY condition allows prefabs in addition to elements. In addition, \#PUT can fail if prefab has collision with the destination.

\#SHOOT now requires a prefab to be specified.

\#IF has value compare conditions

IS, EQUAL, ==, !=, <>, <, >, >=, <=, BETWEEN

\#SEND allows an optional [pos|dir] parameter

Direction RND is equal chance in all directions

#### Commands
Commands with Element: Yes can be used in an element script.
Parameters that are inside [brackets] are optional.

##### \#AND \<cmd>
Element: Yes
End Cycle: No

Executes <cmd> argument the last \#IF returned true.

##### \#BECOME \[color] \<element|prefab>
Element: Yes
End Cycle: Yes

Transforms the object or tile into the element or prefab.

##### \#BOARD <board|target>  \[x] \[y]
Element: Yes
End Cycle: Yes

Switches to another board.
You can name a specific board:
```
#BOARD Board1
```
Or use the object's Target value:
```
#BOARD Target
```

If \[x] and \[y] are given, will move the Player object to that position.
Otherwise, moves player to the position of an object that matches name and group of the calling object.

##### \#CAMERA \<x>  \<y>
Element: Yes
End Cycle: No

Moves the camera so that top-left corner of the camera is at <x> and <y>. It also clears the \#FOLLOW flag if any object was being followed.

##### \#CHANGE \[color] <element>  \[color] <element>
Element: Yes
End Cycle: No

Transforms tiles with the color and element to the new color and element.

##### \#CHAR \<value>
Element: Yes
End Cycle: No

Changes the char of the tile or object.
##### \#CLEAR \<flag>
Element: No
End Cycle: No

##### \#CLEAR \[scope:]\<counter>
Element: Yes
End Cycle: No

Sets the counter to zero.
##### \#COLOR \<color>
Element: Yes
End Cycle: No

Changes the color of the tile or object.
##### \#CYCLE \<value>
Element: No
End Cycle: No

Changes the how often the object's script will execute.

0 = Stop
1 = Every Frame
2 = Every Second Frame
3 = Every Third
etc..

##### \#DIE
Element: Yes
End Cycle: Yes

Removes if its an object. If its a tile, it becomes an empty tile.

##### \#DUPLICATE \<dir|pos>  \<dir|pos>
Element: Yes
End Cycle: No

Copies the object to the target position. All properties, counters, etc will copy, except for script position and label zaps.

##### \#ELSE \<cmd>
Element: Yes
End Cycle: No

Executes <cmd> argument the last \#IF returned false.

##### \#END
Element: Yes
End Cycle: Yes

Ends the current execution. For events, it will end the event and the caller will resume.

##### \#ENDGAME
Element: Yes
End Cycle: Yes

Opens a special textbox which will allow the user to load a previous save or return to the title screen. Any text messages prior to this will show as well.
```
$You have died!
Next time be a bit quicker..
#ENDGAME
```

##### \#EVENT \<dir|pos>  \<label>  \[skip]
Element: Yes
End Cycle: No

Sends an event at the position or direction. Similiar to \#SEND, however the code is executed immediately. The event will be sent to the element at the tile location and any objects present there. Once the event(s) are finish the script will resume.

The option skip value will skip the number of labels as if they were zapped.

##### \#EVENT [dir|pos] \[scope:]\<label> \[skip value]
Element: Yes
End Cycle: No

Sends an event. Similar to \#SEND, however the code is executed immediately. If no scope is provided, then it sends within the same script. Once the event(s) are finish the script will resume.

If a scope is provided, then it will send to all objects with the matching name or group value
```
#EVENT FOO:BAR
```
Will send the 'BAR' event to all objects named 'FOO' or in the group of 'FOO'.

Special Scopes:
If one of the following is provided, it will behave differently:
SELF - Sends to self only.
PLAYER - Sends to the player object.
OTHERS - Sends to every object, except this object. If the object has a group, will send to all objects outside the group. If \[dir|pos] is given, send to the tile as well.
GROUP - Sends to every object with the same group of this object.
ALL - Sends to all.

The option skip value will skip the number of labels as if they were zapped.

In addition, if the \[dir|pos] parameter is provided, it will only send to objects/tiles in that position.
```
#EVENT N TOUCH
```
Will send 'TOUCH' to the tiles and any objects north of the caller.
```
#EVENT SEEK OTHERS:ATTACK
```
Will send 'ATTACK' to objects in the seek direction, if they don't share a group with the caller.

See [Events](#events)

##### \#FOLLOW \<object>
Element: Yes
End Cycle: No

Sets the camera follow object. When the object moves, the camera will move with it.

##### \#GIVE \[scope:]\<counter> \<value>
Element: Yes
End Cycle: No

Increases the <counter> by <value>. If the counter does not exist, it will be added.

\[scope:] determines the scope of the counter

SELF:COUNTER will stay local to the object
BOARD:COUNTER is available to everything in the board
No scope is a world counter, and will persist across boards

##### \#GO \<dir|pos>
Element: Yes
End Cycle: Yes

Long form of /<dir|pos>. Moves the object or tile to the position or direction. Repeats if blocked.
##### \#GOTO \<label> \[skip value]
Element: Yes
End Cycle: No

Long form of \#<label>. Jumps execution to that label. [skip value] will skip over that number of labels, as if they were zapped.
##### \#GROUP <name>
Element: No
End Cycle: No

Changes the group name of the object.

##### \#IDLE
Element: Yes
End Cycle: Yes

Long form of /i. Stops execution until next frame.
##### \#IF \<cond>  \[then] \<cmd>
Element: Yes
End Cycle: No

Checks the condition, if passes it will execute \<cmd>.

Conditions:
BLOCKED \<dir|pos>
Checks own collision against the destination
 
CONTACT
Checks in the player object is adjacent to self.

NOT \<cond>
Reverses the condition
ALIGNED

Checks in the player object is aligned to self.
ANY \[color] \<element>

Returns true if there are any color and element combination on the board.

DETECT \<dir|pos>  \[color] <element|name|group>
Returns true if there are is a color and element combination at the destination.
  
ACCEPTED
Returns true if the last \#EVENT command was successfully executed by an object or element.

Math Conditions:
\<value> IS \<value>
\<value> EQUAL \<value>
\<value> == \<value>
\<value> != \<value>
\<value> <> \<value>
\<value> > \<value>
\<value> >= \<value>
\<value> < \<value>
\<value> <= \<value>
\<value> BETWEEN \<value> \<value>

##### \#LOCK
Element: No
End Cycle: No

Blocks all events and send to the object.

##### \#LOG \<text>
Element: Yes
End Cycle: No

Writes \<text> to runner.log, useful for debugging.

##### \#NAME \<name>
Element: No
End Cycle: No

Sets the object's name value.

##### \#PUT <dir|pos>  \[color] \<element|prefab>
Element: Yes
End Cycle: No

Places the element or prefab at the destination.
##### \#REPEAT \<value>  \<cmd>
Element: Yes
End Cycle: No

Will execute the <cmd>, <value> number of times.
```
#REPEAT 10 IDLE
```
Will execute idle 10 times.

If called during an event, it is much more limited, because the repeat count is lost when the event ends. The sample about will only \#IDLE once.

##### \#RESTART
Element: Yes
End Cycle: No

Sets the script position to the start of the script.
##### \#RESTORE <label>
Element: No
End Cycle: No

Undo a zapped label

##### \#SCRIPT <script>
Element: No
End Cycle: No

Changes the object's script.

##### \#SEND [dir|pos] \[scope:]\<label>  \[skip]
Element: Yes
End Cycle: No

Jumps script instruction to the label.

Note: You can use \#SEND SELF:LABEL to move the object's own script position while inside an event. While #LABEL will only move within the event.

##### \#SET \<flag>
Element: No
End Cycle: No

Available Flags:
COLLISION_1
COLLISION_2
COLLISION_3
COLLISION_4
COLLISION_5
COLLISION_6
COLLISION_7
COLLISION_8
COLLISION_WALLS same as COLLISION_1
COLLISION_WATER same as COLLISION_2
COLLISION_FLOOR same as COLLISION_3

CAN_PUSH
PUSHABLE_N
PUSHABLE_S
PUSHABLE_E
PUSHABLE_W
PUSHABLE_NS

PUSHABLE_EW
PUSHABLE_ALL

PRIORITY
SHOW_BG

##### \#SET \[scope:]\<counter> \[value]
Element: Yes
End Cycle: No

Sets the counter to the value. If the counter doesn't exist, it will add it.

BOARD: scope will set to the board counters.
SELF: scope will set to the object counters.
No scope will set to the world counters.

##### \#VIEW \<id>
Element: Yes
End Cycle: No

Changes the view to \<id>. This only applies to the current board.

See [Views](#views)

##### \#SHOOT \<dir>  \[color] \<prefab>
Element: Yes
End Cycle: Yes

Places a prefab in the \<dir> relative to the caller, and forces a #WALK \<dir> to that object.

##### \#TAKE \[scope:]\<counter> \<value>
Element: Yes
End Cycle: No

Decreases the counter by the value
##### \#TARGET <name>
Element: No
End Cycle: No

Sets the object's target name.
Note: Target has limited functionality, more to come.

##### \#TEXTBOX \<id>  \[title]
Element: Yes
End Cycle: No

Sets the next textbox id.
##### \#TRY <dir|pos>  \[cmd]
Element: Yes
End Cycle: Yes

Long form of ?<dir|pos>. Moves the object or tile to the position or direction. If blocked, executes \<cmd>.

##### \#UNLOCK
Element: No
End Cycle: No

Removes \#LOCK
##### \#WALK \<dir>
Element: No
End Cycle: No

Moves the object in the direction every frame. If blocked the :THUD event is sent to the object. #WALK IDLE stops movement.

##### \#ZAP <label>
Element: No
End Cycle: No

Prevents the label from being used in event or send. If there are multiple labels of the same name, the next one will be used.

##### Directions
I, IDLE
N, NORTH
S, SOUTH
E, EAST
W, WEST
CW \<dir>
CCW \<dir>
OPP \<dir>
FLOW
SEEK
RND
RNDNS
RNDNE
RNDDP \<dir>

#### Positions
##### SELF
This object or element position.

##### AT \<x>  \<y>
Position at \<x> and \<y>

##### WHERE \<name|group>
Finds the position of the closest object with the matching name or group.

##### BY \<x> \<y>
Adds x and y to the current position.

##### TO \<pos|dir> \<pos|dir>
Adds the two positions or directions together

##### PLAYER
Location of the player object.

##### FOLLOW
Location of the object the camera is following.
