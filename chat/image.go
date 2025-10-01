package chat

type ImageDetailLevel string

const (
	ImageDetailAuto   ImageDetailLevel = "auto"
	ImageDetailLow    ImageDetailLevel = "low"
	ImageDetailMedium ImageDetailLevel = "medium"
	ImageDetailHigh   ImageDetailLevel = "high"
)

type Image struct {
	Content []byte `json:"buffer,omitempty"`
	Url     string `json:"url,omitempty"`
	// OpenAI: The level of detail for the image description. Options are "auto", "low", "medium", and "high". Default is "auto".
	Detail ImageDetailLevel `json:"detail,omitempty"`
}
