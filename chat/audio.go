package chat

type AudioFormat string

const (
	AudioFormatWAV   AudioFormat = "wav"
	AudioFormatMP3   AudioFormat = "mp3"
	AudioFormatAAC   AudioFormat = "aac"
	AudioFormatFLAC  AudioFormat = "flac"
	AudioFormatOpus  AudioFormat = "opus"
	AudioFormatPcm16 AudioFormat = "pcm16"
)

type Audio struct {
	Base64 string      `json:"base64,omitempty"`
	Format AudioFormat `json:"format,omitempty"`
}
