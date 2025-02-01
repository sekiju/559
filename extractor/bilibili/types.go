package bilibili

type getEpisodeResult struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Title      string `json:"title"`
		ComicId    int    `json:"comic_id"`
		ShortTitle string `json:"short_title"`
		ComicTitle string `json:"comic_title"`
	} `json:"data"`
}

type getImageIndexResult struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Path   string `json:"path"`
		Images []struct {
			Path      string `json:"path"`
			X         int    `json:"x"`
			Y         int    `json:"y"`
			VideoPath string `json:"video_path"`
			VideoSize string `json:"video_size"`
		} `json:"images"`
		LastModified string `json:"last_modified"`
		Host         string `json:"host"`
		Video        struct {
			Svid      string        `json:"svid"`
			Filename  string        `json:"filename"`
			Route     string        `json:"route"`
			Resource  []interface{} `json:"resource"`
			RawWidth  string        `json:"raw_width"`
			RawHeight string        `json:"raw_height"`
			RawRotate string        `json:"raw_rotate"`
			ImgUrls   []interface{} `json:"img_urls"`
			BinUrl    string        `json:"bin_url"`
			ImgXLen   int           `json:"img_x_len"`
			ImgXSize  int           `json:"img_x_size"`
			ImgYLen   int           `json:"img_y_len"`
			ImgYSize  int           `json:"img_y_size"`
		} `json:"video"`
		Cpx string `json:"cpx"`
	} `json:"data"`
}

type imageTokenResult struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		Url         string `json:"url"`
		Token       string `json:"token"`
		CompleteUrl string `json:"complete_url"`
		HitEncrpyt  bool   `json:"hit_encrpyt"`
	} `json:"data"`
}
