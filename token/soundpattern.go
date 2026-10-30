package token

const minOctave = uint8(0)
const maxOctave = uint8(7)
const maxTone = (maxOctave + 1) * 12

func (bcs *ByteCodeScanner) putSoundPatternParam(str string, buf []byte) (w int, err error) {
	octave := uint8(4)
	duration := uint8(6)
	buf[0] = ArgSoundPattern
	p := 0
	for _, r := range str {
		switch r {
		case 'T':
			duration = 6
		case 'S':
			duration = 12
		case 'I':
			duration = 24
		case 'Q':
			duration = 48
		case 'H':
			duration = 96
		case 'W':
			duration = 192
		case '.':
			duration = uint8((int(duration) * 3) / 2)
		case '3':
			duration = duration / 3
		case '+':
			if octave != maxOctave {
				octave++
			}
		case '-':
			if octave != minOctave {
				octave--
			}
		case 'A', 'B', 'C', 'D', 'E', 'F', 'G':
			buf[p+2] = noteValue(octave, r)
			buf[p+3] = duration
			p += 2
		case '!':
			if p == 0 {
				return 0, ErrBadParam
			}
			note := buf[p] //p+2-2
			if note != 0 && note <= maxTone {
				buf[p] = note - 1
			}
		case '#':
			if p == 0 {
				return 0, ErrBadParam
			}
			note := buf[p] //p+2-2
			if note < maxTone {
				buf[p] = note + 1
			}
		case '0', '1', '2', '4', '5', '6', '7', '8', '9':
			buf[p+2] = drumNum[r]
			buf[p+3] = duration
			p += 2
		case 'X':
			buf[p+2] = SilentNote
			buf[p+3] = duration
			p += 2
		case '%':
			buf[p+2] = WaveFormSquare
			buf[p+3] = 0
			p += 2
		case '~':
			buf[p+2] = WaveFormSine
			buf[p+3] = 0
			p += 2
		case '^':
			buf[p+2] = WaveFormTriangle
			buf[p+3] = 0
			p += 2
		case '*':
			buf[p+2] = WaveFormSaw
			buf[p+3] = 0
			p += 2
		default:
			return 0, ErrBadParam
		}
	}
	buf[1] = uint8(p / 2) // number of notes
	return p + 2, err
}

func noteValue(octave uint8, note rune) uint8 {
	return (octave * 12) + noteNum[note]
}

var noteNum = map[rune]uint8{
	'C': 0,
	'D': 2,
	'E': 4,
	'F': 5,
	'G': 7,
	'A': 9,
	'B': 11,
}

const (
	SilentNote       byte = 255
	WaveFormSquare   byte = 254
	WaveFormSine     byte = 253
	WaveFormTriangle byte = 252
	WaveFormSaw      byte = 251

	DrumTick     byte = 240
	DrumTweet    byte = 241
	DrumCowbell  byte = 242
	DrumHiSnare  byte = 243
	DrumHiWood   byte = 244
	DrumLowSnare byte = 245
	DrumLowTom   byte = 246
	DrumLowWood  byte = 247
	DrumBassDrum byte = 248
)

var drumNum = map[rune]byte{
	'0': DrumTick,
	'1': DrumTweet,
	'2': DrumCowbell,
	'4': DrumHiSnare,
	'5': DrumHiWood,
	'6': DrumLowSnare,
	'7': DrumLowTom,
	'8': DrumLowWood,
	'9': DrumBassDrum,
}