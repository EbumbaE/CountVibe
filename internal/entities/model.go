package entities

type ViewDiaryData struct {
	IsLogin bool

	Date  string
	Total Composition

	Search []Product

	Breakfast Meal
	Lunch     Meal
	Dinner    Meal
	Snacks    Meal
}

func GetViewDiaryData(date string, isLogin bool) ViewDiaryData {

	search := []Product{
		Product{
			ID:   1,
			Name: "kuraga",
			UnitComposition: Composition{
				Calories:      365,
				Proteins:      1,
				Fats:          40,
				Carbohydrates: 12,
			},
			Unit:       "gr",
			AmountUnit: 100,
		},
	}

	m := Meal{
		Portions: []Portion{
			Portion{
				Product: Product{
					ID:   1,
					Name: "kuraga",
					UnitComposition: Composition{
						Calories:      365,
						Proteins:      1,
						Fats:          40,
						Carbohydrates: 12,
					},
					Unit:       "gr",
					AmountUnit: 100,
				},
				Amount: 0,
				CalcPortion: Composition{
					Calories:      365,
					Proteins:      1,
					Fats:          40,
					Carbohydrates: 12,
				},
			},
			Portion{
				Product: Product{
					ID:   1,
					Name: "kuraga",
					UnitComposition: Composition{
						Calories:      365,
						Proteins:      1,
						Fats:          40,
						Carbohydrates: 12,
					},
					Unit:       "gr",
					AmountUnit: 100,
				},
				Amount: 0,
				CalcPortion: Composition{
					Calories:      365,
					Proteins:      1,
					Fats:          40,
					Carbohydrates: 12,
				},
			},
		},
		CountOrder: 0,
		CalcMeal: Composition{
			Calories:      0,
			Proteins:      0,
			Fats:          0,
			Carbohydrates: 0,
		},
	}
	return ViewDiaryData{
		IsLogin: isLogin,
		Date:    date,

		Total: Composition{
			Calories:      0,
			Proteins:      0,
			Fats:          0,
			Carbohydrates: 0,
		},

		Search: search,

		Breakfast: m,
		Lunch:     m,
		Dinner:    m,
		Snacks:    m,
	}
}
