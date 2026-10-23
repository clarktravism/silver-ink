package token

const (
	ArgInteger    byte = iota + 1 // int(2)
	ArgScope                      // scope(1)
	ArgScopeName                  // name(2)
	ArgCmd                        // 
	ArgElement                    // name(2)
	ArgPrefab                     // name(2)
	ArgTarget                     // 
	ArgCond                       // byte(1) args..
	ArgMath                       // byte(1) args..
	ArgDir                        // byte(1) args..
	ArgPos                        // byte(1) args..
	ArgObject                     // byte(1)
	ArgObjectName                 // name(2)
	ArgLabel                      // name(2)
	ArgCounter                    // name(2)
	ArgName                       // name(2)
	ArgFlag                       // byte(1)
	ArgText                       // len(1) text..
)
