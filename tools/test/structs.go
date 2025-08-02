package test

type Bar struct {
	X string `json:"x"`
	Y int64  `json:"y"`
	Z bool   `json:"z"`
}

type Foo struct {
	Bar Bar `json:"bar"`
}
