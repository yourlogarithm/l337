package chat

type ContentType string

type Content struct {
	Text  string `json:"text,omitempty"`
	Image *Image `json:"image,omitempty"`
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

func NewImageContent(imageContent []byte) Content {
	return Content{
		Image: &Image{
			Content: imageContent,
		},
	}
}
