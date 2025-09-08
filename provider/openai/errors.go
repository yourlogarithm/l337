package openai

type ErrParams struct {
	Msg   string
	Param string
}

func (e ErrParams) Error() string {
	if e.Msg != "" {
		return "invalid parameter " + e.Param + ": " + e.Msg
	}
	return "invalid parameter: " + e.Param
}
