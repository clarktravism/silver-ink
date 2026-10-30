package engine

import (
	"math"
	"unsafe"

	"github.com/Zyko0/go-sdl3/sdl"
)

const _frequency = 48000
const _note_delay = _frequency / 1000
const _note_fade = _frequency / 50
const _amp = 1600.0
const _clip_high = float64(math.MaxInt16)
const _clip_low = float64(math.MinInt16)

var _samples_per_dur = 440

const _audio_device = sdl.AUDIO_DEVICE_DEFAULT_PLAYBACK

var audioStream *sdl.AudioStream

var _audio_spec = sdl.AudioSpec{
	Format:   sdl.AUDIO_S16LE,
	Freq:     _frequency,
	Channels: 1,
}

var AudioBuffers = [4]AudioBuffer{}

func OpenAudio() {
	sdl.Init(sdl.INIT_AUDIO)
	callback := sdl.NewAudioStreamCallback(audioCallBack)
	audioStream = _audio_device.OpenAudioDeviceStream(&_audio_spec, callback)
	audioStream.ResumeDevice()
}

func CloseAudio() {
	if audioStream != nil {
		audioStream.Destroy()
	}
}

var _audioSampleBuffer = [512]int16{}
var _audioSampleBufferBytes = unsafe.Slice((*byte)(unsafe.Pointer(unsafe.SliceData(_audioSampleBuffer[:]))), len(_audioSampleBuffer)*2) // magic ..

func audioCallBack(stream *sdl.AudioStream, additionalAmount, totalAmount int32) {

	w := 0 // samples written

	// max samples we may write
	max := int(additionalAmount) >> 1
	if max > len(_audioSampleBuffer) {
		max = len(_audioSampleBuffer)
	}

	for w < max { // iter sample writes

		sample := float64(0)

		for i := 0; i < len(AudioBuffers); i++ {

			ab := &AudioBuffers[i]

			// get next note
			if ab.note.Empty() {
				if !ab.Pop() {
					continue
				}
			}

			if ab.noteDelta != 0 {
				if !ab.IsNote() || ab.notePos >= _note_delay {
					amp := _amp
					if ab.IsNote() && ab.waveForm != 0 { // to prevent clicks and pops
						if ab.notePos < _note_fade {
							amp = amp * (float64(ab.notePos) / float64(_note_fade))
						} else if ab.noteLen-ab.notePos < _note_fade {
							amp = amp * (float64(ab.noteLen-ab.notePos) / float64(_note_fade))
						}
					}
					ab.IncNoteTime()
					sample += waveForm(ab.waveForm, ab.noteTime, amp)
				}
			}
			ab.IncNotePos()
		}

		if sample > _clip_high {
			sample = _clip_high
		} else if sample < _clip_low {
			sample = _clip_low
		}

		_audioSampleBuffer[w] = int16(math.Floor(sample))
		w++
	}

	// write to audio stream
	if w != 0 {
		stream.PutData(_audioSampleBufferBytes[:w<<1])
	}
}

func waveForm(wave byte, time float64, amp float64) float64 {
	switch wave {
	case 0: // square
		if time < math.Pi {
			return -amp
		} else {
			return amp
		}
	case 1: // sine
		return (amp * math.Sin(time)) * 4
	case 2: // triangle
		return ((amp * 2) / math.Pi) * math.Asin(math.Sin(time)) * 4
	case 3: // saw
		return ((amp * 2) / math.Pi) * math.Atan(math.Tan(time/2))
	default:
		return 0
	}
}
