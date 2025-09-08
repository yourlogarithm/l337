package agent

import "fmt"

type ErrBuilderParams struct {
	Param string
	Msg   string
}

func (e ErrBuilderParams) Error() string {
	return fmt.Sprintf("error building agent: %s: %s", e.Param, e.Msg)
}

type ErrModelResponse struct {
	Msg string
}

func (e ErrModelResponse) Error() string {
	return fmt.Sprintf("error processing model response: %s", e.Msg)
}
