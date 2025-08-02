package tools

type Toolkit []Tool

func (t *Toolkit) AddTool(fn any) error {
	tool, err := CreateToolFromFunc(fn)
	if err != nil {
		return err
	}
	*t = append(*t, tool)
	return nil
}
