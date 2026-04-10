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

func searchClosestPD(curLastPD string, doc *goquery.Document, MaxAttempts int) (string, error) {
	for i := 0; i != MaxAttempts; i++ {
		selector := fmt.Sprintf(`[data-post="%s"]`, curLastPD)
		found := doc.Find(selector)
		if found.Length() == 0 {
			splitPD := strings.Split(curLastPD, "/")
			intPD, err := strconv.Atoi(strings.TrimSpace(splitPD[1]))
			if err != nil {
				return "", err
			}
			intPD++
			curLastPD = fmt.Sprintf("%s/%d", splitPD[0], intPD)
		} else {
			return curLastPD, nil
		}
	}
	//  верхний пост в HTML - ближайший если переходить по ?after=ID
	// докрутить это далее
	return "", ErrNoClosest
}

func Parse(baseUrl string, app *application.App) ([]*database.Form, error) {
	urlParts := strings.Split(baseUrl, "/")
	channel := urlParts[len(urlParts)-1]
	forms := []*database.Form{}

	lastPD, err := app.Model.LastHash(channel)
	curLastPD := lastPD
	if err != nil {
		return nil, err
	}

	for {
		url := baseUrl + "?after=" + strings.TrimSpace(strings.Split(curLastPD, "/")[1])

		res, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("error: %w, url: %s", err, url)
		}

		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			return nil, err
		}
		res.Body.Close()

		curLastPD, err = searchClosestPD(curLastPD, doc, 50)
		if err != nil {
			return forms, err
		} //случай если новейший на старнице

		endOfPosts := doc.Find(".tme_no_messages_found")
		if endOfPosts.Length() != 0 {
			return forms, nil
		} //случай если нет на странице
		lastHTMLPD, exist := doc.Find(".tgme_widget_message").Eq(0).Attr("data-post")
		if !exist {
			return nil, ErrNotExist
		}

		posts := doc.Find(".tgme_widget_message")
		for i := posts.Length() - 1; i >= 0; i-- {
			s := posts.Eq(i)
			postPD, exist := s.Attr("data-post")
			if !exist {
				continue
			}
			if postPD == curLastPD && lastHTMLPD != curLastPD {
				return forms, nil
			} // случай если последняя страница и пост находится в середине
			if postPD == curLastPD && lastHTMLPD == curLastPD {
				continue
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
		curLastPD = lastHTMLPD
	}
}
