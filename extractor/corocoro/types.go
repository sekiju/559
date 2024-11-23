package corocoro

type pagesResult []struct {
	Src       string  `json:"src"`
	Type      string  `json:"type"`
	SizeRatio float64 `json:"sizeRatio"`
	Crypto    struct {
		Method string `json:"method"`
		Key    string `json:"key"`
		Iv     string `json:"iv"`
	} `json:"crypto"`
	WidthRatio string `json:"widthRatio"`
}
