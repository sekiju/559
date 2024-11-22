package speed_binb

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/rs/zerolog/log"
	"github.com/sekiju/htt"
	"github.com/sekiju/mdl/internal/renamer"
	"github.com/sekiju/mdl/internal/sdk/extractor/manga"
	"image"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"net/url"
	"strconv"
	"strings"
)

type Extractor struct {
	req *htt.Request
}

func (e *Extractor) FindChapterPages(chapter *manga.Chapter) ([]*manga.Page, error) {
	res, err := e.req.Get(chapter.URL)
	if err != nil || res.StatusCode != 200 {
		return nil, manga.ErrChapterNotFound
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	content := doc.Find("div#content.pages").First()
	apiURL, apiURLExists := content.Attr("data-ptbinb")

	parsedURL, err := url.Parse(chapter.URL)
	if err != nil {
		return nil, nil
	}

	query := parsedURL.Query()

	switch {
	case apiURLExists && strings.Contains(apiURL, "bibGetCntntInfo") && query.Has("u1"):
		return e.v016452(parsedURL, apiURL)
	default:
		return nil, manga.ErrMethodUnimplemented
	}
}

func (e *Extractor) v016452(parsedURL *url.URL, apiURL string) ([]*manga.Page, error) {
	log.Trace().Msg("SpeedBinb version: v016452")

	cid := parsedURL.Query().Get("cid")
	sharingKey := tt(cid)

	query := parsedURL.Query()
	query.Set("k", sharingKey)
	query.Del("rurl")

	parsedURL.Path = apiURL
	parsedURL.RawQuery = query.Encode()

	res, err := e.req.Get(parsedURL.String())
	if err != nil {
		return nil, err
	}

	var bibGetCntntInfoItems BibGetCntntInfo
	if err = res.JSON(&bibGetCntntInfoItems); err != nil {
		return nil, err
	}

	if bibGetCntntInfoItems.Result != 1 {
		return nil, errors.New("invalid bibGetCntntInfoItems result")
	}

	bibGetCntntInfo := bibGetCntntInfoItems.Items[0]

	if bibGetCntntInfo.ServerType != 0 {
		return nil, fmt.Errorf("unsupported speedbinb server type")
	}

	sbcGetCntntUrl, err := url.Parse(bibGetCntntInfo.ContentsServer + "/sbcGetCntnt.php")
	if err != nil {
		return nil, err
	}

	query.Del("k")
	query.Set("p", bibGetCntntInfo.P)
	query.Set("q", "1")
	query.Set("vm", strconv.Itoa(bibGetCntntInfo.ViewMode))
	query.Set("dmytime", bibGetCntntInfo.ContentDate)
	sbcGetCntntUrl.RawQuery = query.Encode()

	res, err = e.req.Get(sbcGetCntntUrl.String())
	if err != nil {
		return nil, err
	}

	var sbcGetCntn SbcGetCntnt
	if err = res.JSON(&sbcGetCntn); err != nil {
		return nil, err
	}

	tDoc, err := goquery.NewDocumentFromReader(strings.NewReader(sbcGetCntn.Ttx))
	if err != nil {
		return nil, err
	}

	ctbl := pt(cid, sharingKey, bibGetCntntInfo.Ctbl)
	ptbl := pt(cid, sharingKey, bibGetCntntInfo.Ptbl)

	sbcGetImgUrl := sbcGetCntntUrl
	sbcGetImgUrl.Path = strings.Replace(sbcGetCntntUrl.Path, "sbcGetCntnt", "sbcGetImg", 1)

	tImages := tDoc.Find("t-case:first-of-type t-img")

	pages := make([]*manga.Page, tImages.Length())
	indexNamer := renamer.New(tImages.Length())

	for i, el := range tImages.EachIter() {
		src, _ := el.Attr("src")

		query = sbcGetImgUrl.Query()
		query.Set("src", src)
		sbcGetImgUrl.RawQuery = query.Encode()

		pages[i] = &manga.Page{
			URL:      sbcGetImgUrl.String(),
			Filename: indexNamer.Name(i, ".png"),
			Index:    uint(i),
			Decode:   e.decode016130(src, ctbl, ptbl),
		}
	}

	return pages, nil
}

func (e *Extractor) decode016130(imgSrc string, ctbl, ptbl []string) manga.DecodeFunc {
	return func(b []byte) ([]byte, error) {
		prototype := lt(imgSrc, ctbl, ptbl)
		if prototype == nil || !prototype.vt() {
			return nil, fmt.Errorf("prototype.vt() dont exists")
		}

		img, _, err := image.Decode(bytes.NewReader(b))
		if err != nil {
			return nil, err
		}

		destRect := prototype.dt(img.Bounds())
		view := DescrambleView{Width: destRect.Dx(), Height: destRect.Dy(), Transfers: []DescrambleTransfer{{0, prototype.gt(img.Bounds())}}}

		descrambledImg := image.NewRGBA(image.Rect(0, 0, view.Width, view.Height))

		for _, part := range view.Transfers[0].Coords {
			wherePlaceRect := image.Rect(part.XDest, part.YDest, part.XDest+part.Width, part.YDest+part.Height)
			whereTakeRect := image.Rect(part.XSrc, part.YSrc, part.XSrc+part.Width, part.YSrc+part.Height)

			draw.Draw(descrambledImg, wherePlaceRect, img, whereTakeRect.Min, draw.Src)
		}

		var buf bytes.Buffer
		if err = png.Encode(&buf, descrambledImg); err != nil {
			return nil, err
		}

		return buf.Bytes(), nil
	}
}

func New(req *htt.Request) *Extractor {
	return &Extractor{req}
}
