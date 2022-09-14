package session

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"CountVibe/internal/entities"
)

type orderMeal int64

const (
	Breakfast orderMeal = iota
	Lunch
	Dinner
	Snacks
)

func getOrderMeal(button string) orderMeal {
	switch button {
	case "addBreakfastItems":
		return Breakfast

	case "addLunchItems":
		return Lunch

	case "addDinnerItems":
		return Dinner

	case "addSnacksItems":
		return Snacks
	}
	return -1
}

func correctDate(date string) (string, error) {
	parseDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return "", err
	}
	return parseDate.Format("02.01.2006"), nil
}

func getDateCookie(r *http.Request) string {
	c, err := r.Cookie("request_date")
	if err != nil {
		return time.Now().Format("01.02.2006")
	}
	return c.Value
}
func saveDateCookie(w http.ResponseWriter, date string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "request_date",
		Value:    date,
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
	})
}

func addItem(date string, id int64, amount float64, od orderMeal) {

}

func diaryFormCheck(w http.ResponseWriter, r *http.Request, isLogin bool) error {

	if isLogin {
		button := r.FormValue("button")

		if od := getOrderMeal(button); od != -1 {
			dateCookie := getDateCookie(r)

			date, err := correctDate(dateCookie)
			if err != nil {
				return err
			}

			for key, value := range r.Form {
				if key == "product_id" {

					strID := value[0]
					id, err := strconv.ParseInt(strID, 10, 64)
					if err != nil {
						return err
					}

					strAmount := r.FormValue("product_amount:" + strID)
					amount, err := strconv.ParseFloat(strAmount, 64)
					if err != nil {
						return err
					}

					if strID != "" && strAmount != "" {
						addItem(date, id, amount, od)
					}
				}
			}
		}
	}

	formDate := r.FormValue("calendar")
	if formDate != "" {
		if date, err := correctDate(formDate); err != nil {
			return err
		} else {
			saveDateCookie(w, date)
		}
	}

	return nil
}

func (s *Session) diaryHandler(w http.ResponseWriter, r *http.Request) {

	isLogin, err := s.compareLogin(r)
	if err != nil {
		s.Logger.Error(err, "Login verification")
	}

	switch r.Method {
	case "GET":

		date := getDateCookie(r)
		diaryData := entities.GetViewDiaryData(date, isLogin)

		paths := []string{
			s.paths["diary"],
			s.paths["item"],
		}
		if err := s.newTemplate(w, diaryData, paths); err != nil {
			s.Logger.Error(err, "new Template")
			http.Error(w, "error in create Template", http.StatusInternalServerError)
		}
	case "POST":
		err := r.ParseForm()
		if err != nil {
			fmt.Fprintf(w, "parse form err: %v", err)
			return
		}
		if err := diaryFormCheck(w, r, isLogin); err != nil {
			s.Logger.Error(err, " Diary Form Check")
		}

		http.Redirect(w, r, r.URL.Path, http.StatusFound)
	}
}
