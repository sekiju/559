package giga_viewer

import "time"

type EpisodeResult struct {
	ReadableProduct struct {
		FinishReadingNotificationUri interface{} `json:"finishReadingNotificationUri"`
		HasPurchased                 bool        `json:"hasPurchased"`
		Id                           string      `json:"id"`
		ImageUrisDigest              *string     `json:"imageUrisDigest"`
		IsPublic                     bool        `json:"isPublic"`
		NextReadableProductUri       *string     `json:"nextReadableProductUri"`
		Number                       int         `json:"number"`
		PageStructure                *struct {
			ChoJuGiga        string              `json:"choJuGiga"`
			Pages            []EpisodeResultPage `json:"pages"`
			ReadingDirection string              `json:"readingDirection"`
			StartPosition    string              `json:"startPosition"`
		} `json:"pageStructure"`
		Permalink                               string      `json:"permalink"`
		PointGettableEpisodeWhenCompleteReading interface{} `json:"pointGettableEpisodeWhenCompleteReading"`
		PrevReadableProductUri                  *string     `json:"prevReadableProductUri"`
		PublishedAt                             time.Time   `json:"publishedAt"`
		// Series eq nil only for magazines.
		Series *struct {
			Id           string `json:"id"`
			ThumbnailUri string `json:"thumbnailUri"`
			Title        string `json:"title"`
		} `json:"series"`
		Title string `json:"title"`
		// Toc is Table Of Contents, used only for magazines.
		Toc *struct {
			Items []struct {
				StartAt int    `json:"startAt"`
				Title   string `json:"title"`
			} `json:"items"`
			NextId       string `json:"next_id"`
			NextUrl      string `json:"next_url"`
			PrevId       string `json:"prev_id"`
			PrevUrl      string `json:"prev_url"`
			ThumbnailUrl string `json:"thumbnail_url"`
			Title        string `json:"title"`
		} `json:"toc"`
		// TypeName eq "episode" or "magazine"
		TypeName string `json:"typeName"`
	} `json:"readableProduct"`
}

type EpisodeResultPage struct {
	LinkPosition string `json:"linkPosition,omitempty"`
	Type         string `json:"type"`
	Buttons      []struct {
		Type string `json:"type"`
		Uri  string `json:"uri"`
	} `json:"buttons,omitempty"`
	ContentStart string `json:"contentStart,omitempty"`
	Height       int    `json:"height,omitempty"`
	Src          string `json:"src,omitempty"`
	Width        int    `json:"width,omitempty"`
	ContentEnd   string `json:"contentEnd,omitempty"`
}
