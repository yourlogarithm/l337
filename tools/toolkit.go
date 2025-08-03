package tools

type Toolkit []Tool

func (t *Toolkit) AddTool(tool Tool) {
	*t = append(*t, tool)
}
