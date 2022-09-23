package entities

type OrderMeal int64

const (
	Breakfast OrderMeal = iota
	Lunch
	Dinner
	Snacks

	BeginMeal = Breakfast
	EndMeal   = Snacks
)

type ViewDayData struct {
	IsLogin bool

	Date  string
	Total Composition

	Search []Product

	Breakfast Meal
	Lunch     Meal
	Dinner    Meal
	Snacks    Meal
}

func NewViewDayData(date string, isLogin bool, search []Product, dm *DayMeals) ViewDayData {
	return ViewDayData{
		IsLogin: isLogin,
		Date:    date,

		Total: dm.CalcDayMeals,

		Search: search,

		Breakfast: dm.Meals[Breakfast],
		Lunch:     dm.Meals[Lunch],
		Dinner:    dm.Meals[Dinner],
		Snacks:    dm.Meals[Snacks],
	}
}

func plugNewViewDayData(date string, isLogin bool) ViewDayData {

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
	return ViewDayData{
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
