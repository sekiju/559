package comic_walker

import "time"

type WorkResult struct {
	Work struct {
		Code       string `json:"code"`
		Id         string `json:"id"`
		Thumbnail  string `json:"thumbnail"`
		BookCover  string `json:"bookCover"`
		Title      string `json:"title"`
		IsOriginal bool   `json:"isOriginal"`
		LabelInfo  struct {
			Name         string `json:"name"`
			IconImageUrl string `json:"iconImageUrl"`
			Code         string `json:"code"`
		} `json:"labelInfo"`
		Language string `json:"language"`
		Internal struct {
			DepartmentCode string        `json:"departmentCode"`
			ScrollType     string        `json:"scrollType"`
			LabelNames     []string      `json:"labelNames"`
			FairIds        []interface{} `json:"fairIds"`
		} `json:"internal"`
		Summary string `json:"summary"`
		Genre   struct {
			Code string `json:"code"`
			Id   string `json:"id"`
			Name string `json:"name"`
		} `json:"genre"`
		SubGenre struct {
			Code string `json:"code"`
			Id   string `json:"id"`
			Name string `json:"name"`
		} `json:"subGenre"`
		Tags []struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		} `json:"tags"`
		Authors []struct {
			Id   string `json:"id"`
			Name string `json:"name"`
			Role string `json:"role"`
		} `json:"authors"`
		FollowerCount       int    `json:"followerCount"`
		IsNew               bool   `json:"isNew"`
		NextUpdateDateText  string `json:"nextUpdateDateText"`
		IsOneShot           bool   `json:"isOneShot"`
		SerializationStatus string `json:"serializationStatus"`
		RatingLevel         string `json:"ratingLevel"`
	} `json:"work"`
	FirstComic struct {
		Id        string `json:"id"`
		Title     string `json:"title"`
		Thumbnail string `json:"thumbnail"`
		Release   string `json:"release"`
		Episodes  []struct {
			Id             string        `json:"id"`
			Code           string        `json:"code"`
			Title          string        `json:"title"`
			SubTitle       string        `json:"subTitle"`
			UpdateDate     time.Time     `json:"updateDate"`
			DeliveryPeriod time.Time     `json:"deliveryPeriod"`
			IsNew          bool          `json:"isNew"`
			HasRead        bool          `json:"hasRead"`
			Stores         []interface{} `json:"stores"`
			ServiceId      string        `json:"serviceId"`
			Internal       struct {
				EpisodeNo   int    `json:"episodeNo"`
				PageCount   int    `json:"pageCount"`
				Episodetype string `json:"episodetype"`
			} `json:"internal"`
			Type     string `json:"type"`
			IsActive bool   `json:"isActive"`
		} `json:"episodes"`
		Stores []struct {
			Code  string `json:"code"`
			Name  string `json:"name"`
			Url   string `json:"url"`
			Image struct {
				Src         string `json:"src"`
				Height      int    `json:"height"`
				Width       int    `json:"width"`
				BlurDataURL string `json:"blurDataURL"`
				BlurWidth   int    `json:"blurWidth"`
				BlurHeight  int    `json:"blurHeight"`
			} `json:"image"`
		} `json:"stores"`
	} `json:"firstComic"`
	LatestComic struct {
		Id        string `json:"id"`
		Title     string `json:"title"`
		Thumbnail string `json:"thumbnail"`
		Release   string `json:"release"`
		Episodes  []struct {
			Id             string        `json:"id"`
			Code           string        `json:"code"`
			Title          string        `json:"title"`
			SubTitle       string        `json:"subTitle"`
			UpdateDate     time.Time     `json:"updateDate"`
			DeliveryPeriod time.Time     `json:"deliveryPeriod"`
			IsNew          bool          `json:"isNew"`
			HasRead        bool          `json:"hasRead"`
			Stores         []interface{} `json:"stores"`
			ServiceId      string        `json:"serviceId"`
			Internal       struct {
				EpisodeNo   int    `json:"episodeNo"`
				PageCount   int    `json:"pageCount"`
				Episodetype string `json:"episodetype"`
			} `json:"internal"`
			Type     string `json:"type"`
			IsActive bool   `json:"isActive"`
		} `json:"episodes"`
		Stores []struct {
			Code  string `json:"code"`
			Name  string `json:"name"`
			Url   string `json:"url"`
			Image struct {
				Src         string `json:"src"`
				Height      int    `json:"height"`
				Width       int    `json:"width"`
				BlurDataURL string `json:"blurDataURL"`
				BlurWidth   int    `json:"blurWidth"`
				BlurHeight  int    `json:"blurHeight"`
			} `json:"image"`
		} `json:"stores"`
	} `json:"latestComic"`
	FirstEpisodes struct {
		Total  int `json:"total"`
		Result []struct {
			Id             string        `json:"id"`
			Code           string        `json:"code"`
			Title          string        `json:"title"`
			SubTitle       string        `json:"subTitle"`
			Thumbnail      string        `json:"thumbnail"`
			UpdateDate     time.Time     `json:"updateDate"`
			DeliveryPeriod time.Time     `json:"deliveryPeriod"`
			IsNew          bool          `json:"isNew"`
			HasRead        bool          `json:"hasRead"`
			Stores         []interface{} `json:"stores"`
			ServiceId      string        `json:"serviceId"`
			Internal       struct {
				EpisodeNo   int    `json:"episodeNo"`
				PageCount   int    `json:"pageCount"`
				Episodetype string `json:"episodetype"`
			} `json:"internal"`
			Type     string `json:"type"`
			IsActive bool   `json:"isActive"`
		} `json:"result"`
	} `json:"firstEpisodes"`
	LatestEpisodes struct {
		Total  int `json:"total"`
		Result []struct {
			Id             string        `json:"id"`
			Code           string        `json:"code"`
			Title          string        `json:"title"`
			SubTitle       string        `json:"subTitle"`
			Thumbnail      string        `json:"thumbnail"`
			UpdateDate     time.Time     `json:"updateDate"`
			DeliveryPeriod time.Time     `json:"deliveryPeriod"`
			IsNew          bool          `json:"isNew"`
			HasRead        bool          `json:"hasRead"`
			Stores         []interface{} `json:"stores"`
			ServiceId      string        `json:"serviceId"`
			Internal       struct {
				EpisodeNo   int    `json:"episodeNo"`
				PageCount   int    `json:"pageCount"`
				Episodetype string `json:"episodetype"`
			} `json:"internal"`
			Type     string `json:"type"`
			IsActive bool   `json:"isActive"`
		} `json:"result"`
	} `json:"latestEpisodes"`
	Comics struct {
		Total  int `json:"total"`
		Result []struct {
			Id        string `json:"id"`
			Title     string `json:"title"`
			Thumbnail string `json:"thumbnail"`
			Release   string `json:"release"`
			Episodes  []struct {
				Id             string        `json:"id"`
				Code           string        `json:"code"`
				Title          string        `json:"title"`
				SubTitle       string        `json:"subTitle"`
				UpdateDate     time.Time     `json:"updateDate"`
				DeliveryPeriod time.Time     `json:"deliveryPeriod"`
				IsNew          bool          `json:"isNew"`
				HasRead        bool          `json:"hasRead"`
				Stores         []interface{} `json:"stores"`
				ServiceId      string        `json:"serviceId"`
				Internal       struct {
					EpisodeNo   int    `json:"episodeNo"`
					PageCount   int    `json:"pageCount"`
					Episodetype string `json:"episodetype"`
				} `json:"internal"`
				Type     string `json:"type"`
				IsActive bool   `json:"isActive"`
			} `json:"episodes"`
			Stores []struct {
				Code  string `json:"code"`
				Name  string `json:"name"`
				Url   string `json:"url"`
				Image struct {
					Src         string `json:"src"`
					Height      int    `json:"height"`
					Width       int    `json:"width"`
					BlurDataURL string `json:"blurDataURL"`
					BlurWidth   int    `json:"blurWidth"`
					BlurHeight  int    `json:"blurHeight"`
				} `json:"image"`
			} `json:"stores"`
		} `json:"result"`
	} `json:"comics"`
	Promotions []struct {
		Type string `json:"type"`
		Url  string `json:"url"`
	} `json:"promotions"`
	RelatedBooks struct {
		TotalCount int           `json:"totalCount"`
		Resources  []interface{} `json:"resources"`
	} `json:"relatedBooks"`
	Label struct {
		Id            string `json:"id"`
		Name          string `json:"name"`
		Code          string `json:"code"`
		Color         string `json:"color"`
		IconImageUrl  string `json:"iconImageUrl"`
		Description   string `json:"description"`
		CoverImageUrl string `json:"coverImageUrl"`
		LogoImageUrl  string `json:"logoImageUrl"`
	} `json:"label"`
	Labels []struct {
		Id            string `json:"id"`
		Name          string `json:"name"`
		Code          string `json:"code"`
		Color         string `json:"color"`
		IconImageUrl  string `json:"iconImageUrl"`
		Description   string `json:"description"`
		CoverImageUrl string `json:"coverImageUrl"`
		LogoImageUrl  string `json:"logoImageUrl"`
	} `json:"labels"`
	LabelWorks []struct {
		Code       string `json:"code"`
		Id         string `json:"id"`
		Thumbnail  string `json:"thumbnail"`
		BookCover  string `json:"bookCover"`
		Title      string `json:"title"`
		IsOriginal bool   `json:"isOriginal"`
		Language   string `json:"language"`
		Internal   struct {
			LabelNames []string `json:"labelNames"`
		} `json:"internal"`
		IsNew   bool `json:"isNew"`
		Episode struct {
			Type  string `json:"type"`
			Code  string `json:"code"`
			Title string `json:"title"`
		} `json:"episode"`
	} `json:"labelWorks"`
	LatestEpisodeId string `json:"latestEpisodeId"`
}

type EpisodeResult struct {
	Episode struct {
		Id        string `json:"id"`
		Code      string `json:"code"`
		Title     string `json:"title"`
		ServiceId string `json:"serviceId"`
		Thumbnail string `json:"thumbnail"`
		Internal  struct {
			EpisodeNo   int    `json:"episodeNo"`
			PageCount   int    `json:"pageCount"`
			Episodetype string `json:"episodetype"`
			IsLatest    bool   `json:"isLatest"`
		} `json:"internal"`
		UpdateDate time.Time `json:"updateDate"`
		IsActive   bool      `json:"isActive"`
	} `json:"episode"`
}

type ViewerResult struct {
	PromotionsEnd   []interface{} `json:"promotionsEnd"`
	LabelLogo       string        `json:"labelLogo"`
	ScrollDirection string        `json:"scrollDirection"`
	ExpiresAt       time.Time     `json:"expiresAt"`
	StartPosition   string        `json:"startPosition"`
	DisplayAds      bool          `json:"displayAds"`
	Manuscripts     []struct {
		DrmMode     string `json:"drmMode"`
		DrmHash     string `json:"drmHash"`
		DrmImageUrl string `json:"drmImageUrl"`
		Page        int    `json:"page"`
		Width       int    `json:"width"`
		Height      int    `json:"height"`
	} `json:"manuscripts"`
}

type ViewerJumpForwardResult struct {
	Type    *string `json:"type"`
	Episode *struct {
		Id             string        `json:"id"`
		Code           string        `json:"code"`
		Title          string        `json:"title"`
		SubTitle       string        `json:"subTitle"`
		Thumbnail      string        `json:"thumbnail"`
		UpdateDate     time.Time     `json:"updateDate"`
		DeliveryPeriod time.Time     `json:"deliveryPeriod"`
		IsNew          bool          `json:"isNew"`
		HasRead        bool          `json:"hasRead"`
		Stores         []interface{} `json:"stores"`
		ServiceId      string        `json:"serviceId"`
		Internal       struct {
			EpisodeNo   int    `json:"episodeNo"`
			PageCount   int    `json:"pageCount"`
			Episodetype string `json:"episodetype"`
		} `json:"internal"`
		Type     string `json:"type"`
		IsActive bool   `json:"isActive"`
	} `json:"episode"`
}
