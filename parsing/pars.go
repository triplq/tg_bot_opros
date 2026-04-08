package parsing

import (
	"net/http"
	"parser/database"
	"strings"
	"time"

	"parser/application"

	"github.com/PuerkitoBio/goquery"
)

func Parse(url string, app application.App) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	url_parts := strings.Split(url, "/")
	channel := url_parts[len(url_parts)-1]

	last_hash, err := app.Model.LastHash(channel)
	if err != nil {
		return err
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}

	forms := []*database.Form{}

	doc.Find("tgme_widget_message").EachWithBreak(func(i int, s *goquery.Selection) bool {
		post_hash, exist := s.Attr("data-post")
		if exist == false || last_hash == post_hash { //тут проверка на хэш, но сам он не хэш, мб докрутить потом
			return false
		}

		form := &database.Form{}
		form.Channel = channel
		form.Msg = s.Find("tgme_widget_message_text").Text()
		time_str := s.Find("tgme_widget_message_date").Text()
		form.Posted_at, err = time.Parse(time.RFC3339, time_str)
		if err != nil {
			return false
		}
	})

}
