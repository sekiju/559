package giga_viewer

type episodeResult struct {
	ReadableProduct struct {
		FinishReadingNotificationUri interface{} `json:"finishReadingNotificationUri"`
		HasPurchased                 bool        `json:"hasPurchased"`
		Id                           string      `json:"id"`
		IsPublic                     bool        `json:"isPublic"`
		NextReadableProductUri       *string     `json:"nextReadableProductUri"`
		Number                       int         `json:"number"`
		PageStructure                *struct {
			Pages []episodeResultPage `json:"pages"`
		} `json:"pageStructure"`
		Permalink              string  `json:"permalink"`
		PrevReadableProductUri *string `json:"prevReadableProductUri"`
		Title                  string  `json:"title"`
	} `json:"readableProduct"`
}

type episodeResultPage struct {
	Type string `json:"type"`
	Src  string `json:"src,omitempty"`
}
