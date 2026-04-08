package parsing

import (
	"fmt"
	"net/http"
	"os"
	"parser/database"
	"strings"
	"time"

	"parser/application"

	"github.com/PuerkitoBio/goquery"
)

func Parse(url string, app *application.App) error {
	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error: %w, url: %s", err, url)
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

	doc.Find(".tgme_widget_message").EachWithBreak(func(i int, s *goquery.Selection) bool {
		post_hash, exist := s.Attr("data-post")
		if !exist || last_hash == post_hash { //тут проверка на хэш, но сам он не хэш, мб докрутить потом
			return false
		}

		form := &database.Form{}
		form.Channel = channel
		form.Hash = post_hash
		form.Msg = strings.TrimSpace(s.Find(".tgme_widget_message_text").Text())
		time_str, exist := s.Find(".tgme_widget_message_date time").Attr("datetime")
		if !exist {
			fmt.Fprint(os.Stderr, "Something wrong with time_exist", post_hash)
			return true
		}
		form.Posted_at, err = time.Parse(time.RFC3339, time_str)
		if err != nil {
			fmt.Fprint(os.Stderr, "Something wrong with time_parse", post_hash, time_str)
			return true
		}

		forms = append(forms, form)
		return true
	})

	err = app.Model.PasteForms(forms)
	if err != nil {
		return err
	}

	return nil
}
