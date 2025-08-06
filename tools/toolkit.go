package tools

type Toolkit map[string]Tool

func (t *Toolkit) AddTool(tool Tool) {
	if *t == nil {
		*t = make(Toolkit)
	}
	(*t)[tool.Name] = tool
}

func (t *Toolkit) Get(name string) (Tool, bool) {
	tool, exists := (*t)[name]
	return tool, exists
}
