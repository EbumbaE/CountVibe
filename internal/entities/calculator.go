package entities

type Portion struct {
	Product     Product
	Amount      float64
	CalcPortion Composition
}

type Meal struct {
	Portions   []Portion
	CountOrder int
	CalcMeal   Composition
}

type DayMeals struct {
	Meals        []Meal
	Date         string
	CalcDayMeals Composition
}

func divProductOnAmountUnit(product Product) Composition {

	amountUnit := product.AmountUnit
	composition := product.UnitComposition

	composition.divCompositionOn(amountUnit)
	return composition
}

func (p Portion) сalcPorcion() {

	product := p.Product
	calcPortion := divProductOnAmountUnit(product)

	amount := p.Amount
	calcPortion.multCompositionOn(amount)

	p.CalcPortion = calcPortion
}

func (m Meal) сalcMeal() {

	calcMeal := Composition{0, 0, 0, 0}
	for _, portion := range m.Portions {
		c := portion.CalcPortion
		calcMeal.addComposition(c)
	}
	m.CalcMeal = calcMeal
}

func (m DayMeals) сalcDayMeal() {

	calcDayMeals := Composition{0, 0, 0, 0}
	for _, meel := range m.Meals {
		c := meel.CalcMeal
		calcDayMeals.addComposition(c)
	}
	m.CalcDayMeals = calcDayMeals
}
