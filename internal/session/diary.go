package session

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"CountVibe/internal/entities"
	productSearch "CountVibe/internal/search"
	"CountVibe/internal/storage"
)

func getOrderMeal(button string) entities.OrderMeal {
	switch button {
	case "addBreakfastItems":
		return entities.Breakfast

	case "addLunchItems":
		return entities.Lunch

	case "addDinnerItems":
		return entities.Dinner

	case "addSnacksItems":
		return entities.Snacks
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
		return time.Now().Format("02.01.2006")
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

func (s *Session) diaryFormCheck(w http.ResponseWriter, r *http.Request, isLogin bool, userID int64) error {

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
					strAmount := r.FormValue("product_amount:" + strID)

					if strID != "" && strAmount != "" {
						strOD := strconv.FormatInt(int64(od), 10)
						strUserID := strconv.FormatInt(userID, 10)

						insertMap := map[string]string{
							"diary_id":   strUserID,
							"date":       date,
							"meal_order": strOD,
							"product_id": strID,
							"amount":     strAmount,
						}
						s.db.SetPortion(insertMap)
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

func parsePrtoduct(m map[string]string) (entities.Product, error) {
	amountUnit, err1 := strconv.ParseFloat(m["amount_unit"], 64)
	id, err2 := strconv.ParseInt(m["product_id"], 10, 64)
	if err1 != nil || err2 != nil {
		return entities.Product{}, fmt.Errorf("parse product from map", err1, err2)
	}

	strUnitComp := strings.Split(m["unit_composition"], "/")
	comp := []float64{0, 0, 0, 0}

	for i := 0; i < 4; i++ {
		comp[i], _ = strconv.ParseFloat(strUnitComp[i], 64)
	}
	return entities.Product{
		ID:   id,
		Name: m["name"],
		UnitComposition: entities.Composition{
			Calories:      comp[0],
			Proteins:      comp[1],
			Fats:          comp[2],
			Carbohydrates: comp[3],
		},
		Unit:       m["unit"],
		AmountUnit: amountUnit,
	}, nil
}

func getMeal(diaryID int64, date string, om entities.OrderMeal, db storage.Storage) (entities.Meal, error) {

	portions := entities.NewPorcions()

	strDiaryID := strconv.FormatInt(diaryID, 10)

	strOM := strconv.FormatInt(int64(om), 10)
	res, err := db.GetPortions(strDiaryID, date, strOM)
	if err != nil {
		return entities.Meal{}, err
	}

	for _, m := range *res {
		amount, err := strconv.ParseFloat(m["amount"], 64)
		if err != nil {
			return entities.Meal{}, err
		}

		mapProduct, err := db.GetProduct(m["product_id"])
		if err != nil {
			return entities.Meal{}, err
		}
		product, err := parsePrtoduct(*mapProduct)
		if err != nil {
			return entities.Meal{}, err
		}

		p := entities.SetPorcion(product, amount)
		portions = append(portions, p)
	}
	m := entities.SetMeal(portions, om)

	return m, nil
}

func getDayMeals(diaryID int64, date string, db storage.Storage) *entities.DayMeals {
	meals := entities.NewMeals(4)

	for order := entities.BeginMeal; order <= entities.EndMeal; order++ {
		getM, err := getMeal(diaryID, date, order, db)
		if err == nil {
			meals[order] = getM
		}
	}
	dm := entities.NewDayMeals(meals)
	return dm
}

func (s *Session) diaryHandler(w http.ResponseWriter, r *http.Request) {

	isLogin, err := s.compareLogin(r)
	if err != nil {
		s.Logger.Error(err, "Login verification")
	}

	var userID int64 = 0
	username := s.parseUsernameFromURL(r)
	if username != "" {
		strUserID, err := s.db.GetUserID(username)
		if err == nil {
			getUserID, err := strconv.ParseInt(strUserID, 10, 64)
			if err == nil {
				userID = getUserID
			}
		}
	}

	switch r.Method {
	case "GET":

		paths := []string{
			s.paths["diary"],
			s.paths["item"],
			s.paths["product"],
		}

		search := productSearch.GetSearch()

		date := getDateCookie(r)
		dm := getDayMeals(userID, date, s.db)

		diaryData := entities.NewViewDayData(date, isLogin, search.Products, dm)

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
		if err := s.diaryFormCheck(w, r, isLogin, userID); err != nil {
			s.Logger.Error(err, " Diary Form Check")
		}

		http.Redirect(w, r, r.URL.Path, http.StatusFound)
	}
}
