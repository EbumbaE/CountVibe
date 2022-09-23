package entities

type Calculator interface {
	NewPorcion(product Product, amount float64) Portion
	SumPortions(portions ...Portion) Composition
	SumDayMeals(meals ...Meal) Composition
}

type Portion struct {
	Product     Product
	Amount      float64
	CalcPortion Composition
}

type Meal struct {
	Portions   []Portion
	CountOrder OrderMeal
	CalcMeal   Composition
}

type DayMeals struct {
	Meals        []Meal
	CalcDayMeals Composition
}

func SetPorcion(product Product, amount float64) Portion {

	calcPortion := product.divProductOnUnit()
	calcPortion.multCompositionOn(amount)
	calcPortion.roundComposition()

	return Portion{
		Product:     product,
		Amount:      amount,
		CalcPortion: calcPortion,
	}
}

func SetMeal(portions []Portion, om OrderMeal) Meal {
	m := Meal{}
	m.Portions = portions
	m.SumPortions()
	m.CountOrder = om
	return m
}
func NewPorcions() []Portion {
	return make([]Portion, 0)
}

func NewMeals(amount int64) []Meal {
	return make([]Meal, amount)
}

func NewDayMeals(meals []Meal) *DayMeals {
	dm := &DayMeals{}
	dm.Meals = meals
	dm.SumDayMeals()
	return dm
}

func (m *Meal) SumPortions() {

	sumPortions := Composition{0, 0, 0, 0}
	for _, portion := range m.Portions {
		c := portion.CalcPortion
		sumPortions.addComposition(c)
	}

	m.CalcMeal = sumPortions
}

func (d *DayMeals) SumDayMeals() {

	sumDayMeals := Composition{0, 0, 0, 0}
	for _, meel := range d.Meals {
		c := meel.CalcMeal
		sumDayMeals.addComposition(c)
	}
	d.CalcDayMeals = sumDayMeals
}
