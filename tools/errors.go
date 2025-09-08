package tools

type ErrToolCreation struct {
	Message string
}

func (e ErrToolCreation) Error() string {
	return e.Message
}
