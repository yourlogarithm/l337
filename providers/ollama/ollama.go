package ollama

import (
	"net/http"
	"net/url"

	"github.com/ollama/ollama/api"
	"github.com/yourlogarithm/l337/internal/logging"
	"github.com/yourlogarithm/l337/providers"
)

var logger = logging.SetupLogger("provider.ollama")

type ollamaProvider struct {
	model  string
	client *api.Client
}

func NewModel(name string, baseUrl string, http *http.Client) (*providers.Model, error) {
	baseUrlParsed, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}
	client := api.NewClient(baseUrlParsed, http)
	return &providers.Model{
		Name:     name,
		Provider: "ollama",
		Impl:     &ollamaProvider{model: name, client: client},
	}, nil
}
