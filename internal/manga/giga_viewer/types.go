package giga_viewer

type EpisodeResult struct {
	ReadableProduct struct {
		FinishReadingNotificationUri interface{} `json:"finishReadingNotificationUri"`
		HasPurchased                 bool        `json:"hasPurchased"`
		Id                           string      `json:"id"`
		IsPublic                     bool        `json:"isPublic"`
		NextReadableProductUri       *string     `json:"nextReadableProductUri"`
		Number                       int         `json:"number"`
		PageStructure                *struct {
			Pages []EpisodeResultPage `json:"pages"`
		} `json:"pageStructure"`
		Permalink              string  `json:"permalink"`
		PrevReadableProductUri *string `json:"prevReadableProductUri"`
		Title                  string  `json:"title"`
	} `json:"readableProduct"`
}

type EpisodeResultPage struct {
	Type string `json:"type"`
	Src  string `json:"src,omitempty"`
}
