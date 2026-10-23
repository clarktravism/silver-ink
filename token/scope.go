package token

const (
	ScopeNone byte = iota + 1
	ScopeSelf
	ScopeAll
	ScopeOthers
	ScopePlayer
	ScopeName
	ScopeBoard
	ScopeGroup
)