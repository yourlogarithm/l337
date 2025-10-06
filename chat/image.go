package chat

type ImageDetailLevel string

const (
	ImageDetailAuto   ImageDetailLevel = "auto"
	ImageDetailLow    ImageDetailLevel = "low"
	ImageDetailMedium ImageDetailLevel = "medium"
	ImageDetailHigh   ImageDetailLevel = "high"
)

type ImageFormat string

const (
	ImageFormatPNG  ImageFormat = "png"
	ImageFormatJPG  ImageFormat = "jpg"
	ImageFormatWEBP ImageFormat = "webp"
)

type ImageData struct {
	Base64 string      `json:"base64,omitempty"`
	Format ImageFormat `json:"format,omitempty"`
}

type Image struct {
	ImageData *ImageData `json:"image_data,omitempty"`
	Url       string     `json:"url,omitempty"`
	// OpenAI: The level of detail for the image description. Options are "auto", "low", "medium", and "high". Default is "auto".
	Detail ImageDetailLevel `json:"detail,omitempty"`
}
