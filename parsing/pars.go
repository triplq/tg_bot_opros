package parsing

import (
	"fmt"
	"net/http"
	"os"
	"parser/database"
	"strconv"
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

	urlParts := strings.Split(url, "/")
	channel := urlParts[len(urlParts)-1]

	lastHash, err := app.Model.LastHash(channel)
	if err != nil {
		return err
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}

	forms := []*database.Form{}

	for true { //убрать потом этот вонючий цикл и сделать maxAttmpts
		selector := fmt.Sprintf(`[data-post="%s"]`, lastHash)
		found := doc.Find(selector)
		if found.Length() == 0 {
			hashSplit := strings.Split(lastHash, "/")
			hashInt, err := strconv.Atoi(hashSplit[1])
			if err != nil {
				return err
			}
			hashInt++
			lastHash = fmt.Sprintf("%s/%d", hashSplit[0], hashInt)
		} else {
			break
		}
	}

	doc.Find(".tgme_widget_message").EachWithBreak(func(i int, s *goquery.Selection) bool {
		postHash, exist := s.Attr("data-post")
		fmt.Println(postHash, lastHash)
		if !exist || lastHash == postHash { //тут проверка на хэш, но сам он не хэш, мб докрутить потом
			return false
		}

		form := &database.Form{}
		form.Channel = channel
		form.Hash = postHash
		form.Msg = strings.TrimSpace(s.Find(".tgme_widget_message_text").Text())
		timeStr, exist := s.Find(".tgme_widget_message_date time").Attr("datetime")
		if !exist {
			fmt.Fprint(os.Stderr, "Something wrong with time_exist", postHash)
			return true
		}
		form.Posted_at, err = time.Parse(time.RFC3339, timeStr)
		if err != nil {
			fmt.Fprint(os.Stderr, "Something wrong with time_parse", postHash, timeStr)
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
