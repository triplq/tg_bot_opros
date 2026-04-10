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

	lastPD, err := app.Model.LastHash(channel)
	if err != nil {
		return err
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}

	forms := []*database.Form{}

	for true { //убрать потом этот вонючий цикл и сделать maxAttmpts
		selector := fmt.Sprintf(`[data-post="%s"]`, lastPD)
		found := doc.Find(selector)
		if found.Length() == 0 {
			splitPD := strings.Split(lastPD, "/")
			intPD, err := strconv.Atoi(strings.TrimSpace(splitPD[1]))
			if err != nil {
				return err
			}
			intPD++
			lastPD = fmt.Sprintf("%s/%d", splitPD[0], intPD)
		} else {
			break
		}
	}

	posts := doc.Find(".tgme_widget_message")
	for i := posts.Length() - 1; i >= 0; i-- {
		s := posts.Eq(i)
		postPD, exist := s.Attr("data-post")
		if !exist || lastPD == postPD {
			break
		}

		form := &database.Form{}
		form.Channel = channel
		form.PostData = postPD
		form.Msg = strings.TrimSpace(s.Find(".tgme_widget_message_text").Text())
		timeStr, exist := s.Find(".tgme_widget_message_date time").Attr("datetime")
		if !exist {
			fmt.Fprint(os.Stderr, "Something wrong with time_exist", postPD)
		}
		form.Posted_at, err = time.Parse(time.RFC3339, timeStr)
		if err != nil {
			fmt.Fprint(os.Stderr, "Something wrong with time_parse", postPD, timeStr)
		}

		forms = append(forms, form)
	}

	err = app.Model.PasteForms(forms)
	if err != nil {
		return err
	}

	return nil
}
