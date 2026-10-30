package engine

import (
	"math"
	"math/rand"
	"silver-ink/token"
	"sync"
	"time"
)

const Pi2 = math.Pi * 2
const Pih = math.Pi / 2

var _audio_rand *rand.Rand

func init() {
	source := rand.NewSource(time.Now().Unix())
	_audio_rand = rand.New(source)
}

type Note struct {
	Tone     uint8
	Duration uint8
}

func (n Note) Freq() float64 {
	return FreqTable[n.Tone]
}

func (n Note) Empty() bool {
	return n.Duration == 0
}

type AudioBuffer struct {
	buffer     [256]Note
	head, tail uint8 // head is read pos, tail is write pos
	mutex      sync.Mutex
	note       Note
	notePos    int
	noteLen    int
	noteTime   float64
	noteDelta  float64
	waveForm   byte
}

func (ab *AudioBuffer) Clear() {
	ab.mutex.Lock()
	ab.head = 0
	ab.tail = 0
	ab.note = Note{}
	ab.notePos = 0
	ab.noteLen = 0
	ab.noteTime = 0
	ab.noteDelta = 0
	ab.mutex.Unlock()
}

func (ab *AudioBuffer) Stop() {
	ab.mutex.Lock()
	ab.head = 0
	ab.tail = 0
	if ab.waveForm == 0 {
		ab.note = Note{}
		ab.notePos = 0
		ab.noteLen = 0
		ab.noteTime = 0
		ab.noteDelta = 0
	} else if ab.noteDelta != 0 && ab.noteLen-ab.notePos > _note_fade {
		ab.notePos = ab.noteLen - _note_fade
	}
	ab.mutex.Unlock()
}

func (ab *AudioBuffer) IsPlaying() bool {
	return ab.noteLen != 0
}

func (ab *AudioBuffer) Push(note Note) bool {
	ab.mutex.Lock()
	defer ab.mutex.Unlock()
	if ab.tail+1 == ab.head {
		return false
	}
	ab.buffer[ab.tail] = note
	ab.tail++
	return true
}

func (ab *AudioBuffer) Pop() bool {
	ab.mutex.Lock()
	defer ab.mutex.Unlock()

NEXT:
	ab.notePos = 0
	ab.noteLen = 0
	ab.noteTime = 0
	if ab.Empty() {
		ab.note = Note{}
		return false
	}

	note := ab.buffer[ab.head]
	ab.head++

	// switch wave form
	switch note.Tone {
	case token.WaveFormSquare:
		ab.waveForm = 0
		goto NEXT
	case token.WaveFormSine:
		ab.waveForm = 1
		goto NEXT
	case token.WaveFormTriangle:
		ab.waveForm = 2
		goto NEXT
	case token.WaveFormSaw:
		ab.waveForm = 3
		goto NEXT
	}

	ab.note = note
	ab.noteLen = _samples_per_dur * int(ab.note.Duration)
	if ab.IsNote() {
		freq := FreqTable[ab.note.Tone]
		ab.noteDelta = (freq * Pi2) / _frequency
	} else if ab.IsDrum() {
		ab.noteDelta = 0
		if ab.waveForm != 0 {
			ab.noteTime += Pih
		}
	}

	return true
}

func (ab *AudioBuffer) IsNote() bool {
	return ab.note.Tone < 96
}

func (ab *AudioBuffer) IsDrum() bool {
	return ab.note.Tone >= token.DrumTick && ab.note.Tone <= token.DrumBassDrum
}

// inc time within the wave
func (ab *AudioBuffer) IncNoteTime() {
	ab.noteTime += ab.noteDelta
	if ab.noteTime > Pi2 {
		ab.noteTime = math.Mod(ab.noteTime, Pi2)
	}
}

// inc sample pos
func (ab *AudioBuffer) IncNotePos() {
	ab.notePos++
	if ab.notePos == ab.noteLen-1 {
		ab.note = Note{}
		return
	}

	// move drum to next freq
	// drums change freq every 16 samples
	if ab.IsDrum() && ab.notePos%48 == 0 {
		len := DrumLen[ab.note.Tone]
		i := ab.notePos / 48 //div 16
		if i < len {
			freq := DrumFreq(ab.note.Tone, i)
			ab.noteDelta = freq * Pi2 / _frequency
		} else {
			ab.noteDelta = 0
			ab.noteTime = 0
		}
	}
}

func (ab *AudioBuffer) Empty() bool {
	return ab.head == ab.tail
}

// octaves 0-7
var FreqTable = [96]float64{
	16.35, //C0
	17.32,
	18.35,
	19.45,
	20.60,
	21.83,
	23.12,
	24.50,
	25.96,
	27.50,
	29.14,
	30.87,
	32.70, //C1
	34.65,
	36.71,
	38.89,
	41.20,
	43.65,
	46.25,
	49.00,
	51.91,
	55.00,
	58.27,
	61.74,
	65.41, //C2
	69.30,
	73.42,
	77.78,
	82.41,
	87.31,
	92.50,
	98.00,
	103.83,
	110.00,
	116.54,
	123.47,
	130.81, //C3
	138.59,
	146.83,
	155.56,
	164.81,
	174.61,
	185.00,
	196.00,
	207.65,
	220.00,
	233.08,
	246.94,
	261.63, //C4
	277.18,
	293.66,
	311.13,
	329.63,
	349.23,
	369.99,
	392.00,
	415.30,
	440.00,
	466.16,
	493.88,
	523.25, //C5
	554.37,
	587.33,
	622.25,
	659.25,
	698.46,
	739.99,
	783.99,
	830.61,
	880.00,
	932.33,
	987.77,
	1046.50, //C6
	1108.73,
	1174.66,
	1244.51,
	1318.51,
	1396.91,
	1479.98,
	1567.98,
	1661.22,
	1760.00,
	1864.66,
	1975.53,
	2093.00, //C7
	2217.46,
	2349.32,
	2489.02,
	2637.02,
	2793.83,
	2959.96,
	3135.96,
	3322.44,
	3520.00,
	3729.31,
	3951.07,
}

var DrumMap = map[byte]byte{
	token.DrumTick:     0,
	token.DrumTweet:    1,
	token.DrumCowbell:  2,
	token.DrumHiSnare:  3,
	token.DrumHiWood:   4,
	token.DrumLowSnare: 5,
	token.DrumLowTom:   6,
	token.DrumLowWood:  7,
	token.DrumBassDrum: 8,
}

var DrumLen = map[byte]int{
	token.DrumTick:     2,
	token.DrumTweet:    15,
	token.DrumCowbell:  15,
	token.DrumHiSnare:  15,
	token.DrumHiWood:   15,
	token.DrumLowSnare: 15,
	token.DrumLowTom:   15,
	token.DrumLowWood:  15,
	token.DrumBassDrum: 15,
}

func DrumFreq(drum byte, i int) float64 {
	switch drum {
	case token.DrumTick:
		return 3200.00
	case token.DrumTweet:
		return DrumTable[0][i-1]
	case token.DrumCowbell:
		return DrumTable[1][i-1]
	case token.DrumHiSnare:
		return float64(_audio_rand.Intn(5000) + 500)
	case token.DrumHiWood:
		if i%2 == 0 {
			return float64(_audio_rand.Intn(1600) + 800)
		}
		return 1600
	case token.DrumLowSnare:
		return DrumTable[2][i-1]
	case token.DrumLowTom:
		return DrumTable[3][i-1]
	case token.DrumLowWood:
		return float64((i*20 + 1200) - _audio_rand.Intn(i*40))
	case token.DrumBassDrum:
		return float64(_audio_rand.Intn(440) + 220)
	}
	return 0
}

var DrumTable = [4][14]float64{
	{1100, 1200, 1300, 1400, 1500, 1600, 1700, 1800, 1900, 2000, 2100, 2200, 2300, 2400}, // tweet
	{4800, 4800, 8000, 1600, 4800, 4800, 8000, 1600, 4800, 4800, 8000, 1600, 4800, 4800}, // cowbell
	{2200, 1760, 1760, 1320, 2640, 880, 2200, 1760, 1760, 1320, 2640, 880, 2200, 1760},   // low snare
	{688, 676, 664, 652, 640, 628, 616, 604, 592, 580, 568, 556, 544, 532},               // low tom
}
