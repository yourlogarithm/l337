package openai

import (
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"

	"github.com/yourlogarithm/l337/internal/logging"
	"github.com/yourlogarithm/l337/providers"
)

var logger = logging.SetupLogger("provider.openai")

type openAIProvider struct {
	model  string
	client openai.Client
}

func NewModel(name string, opts ...option.RequestOption) *providers.Model {
	return &providers.Model{
		Name:     name,
		Provider: "openai",
		Impl:     &openAIProvider{model: name, client: openai.NewClient(opts...)},
	}
}
