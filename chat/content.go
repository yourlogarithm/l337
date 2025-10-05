package chat

import "encoding/base64"

type ContentType string

type Content struct {
	Text  string `json:"text,omitempty"`
	Image *Image `json:"image,omitempty"`
	Audio *Audio `json:"audio,omitempty"`
}

func (c Content) AsSlice() []Content {
	return []Content{c}
}

func NewTextContent(text string) Content {
	return Content{
		Text: text,
	}
}

func NewImageUrlContent(url string, detail ...ImageDetailLevel) Content {
	if len(detail) > 0 {
		return Content{
			Image: &Image{
				Url:    url,
				Detail: detail[0],
			},
		}
	}
	return Content{
		Image: &Image{
			Url: url,
		},
	}
}

func NewImageContent(imageContent []byte, format ImageFormat) Content {
	return Content{
		Image: &Image{
			ImageData: &ImageData{
				Base64: base64.StdEncoding.EncodeToString(imageContent),
				Format: format,
			},
		},
	}
}

// NewAudioContent creates a new Content instance containing audio data.
// It takes the audio data as a byte slice and the audio format as an AudioFormat.
// The returned Content object will have its Audio field populated with the provided data and format.
func NewAudioContent(data []byte, format AudioFormat) Content {
	base64Data := base64.StdEncoding.EncodeToString(data)

	return Content{
		Audio: &Audio{
			Base64: base64Data,
			Format: format,
		},
	}
}
