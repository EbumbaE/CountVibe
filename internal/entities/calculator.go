package entities

type Portion struct{
	product Product
	amount float64
	calcportion Composition
}

type Meal struct{
	portions []Portion
	countOrder int
	calcmeal Composition
}

type DayMeals struct {
	meals []Meal
	date string
	calcdaymeals Composition
}

func divProductOnAmountUnit(product Product) Composition{

    amountunit := product.amountunit
    composition := product.unitcomposition
	
	composition.divCompositionOn(amountunit)
	return composition
}

func (p Portion)сalcPorcion(){
	
	product := p.product
	calcportion := divProductOnAmountUnit(product)

	amount := p.amount
	calcportion.multCompositionOn(amount)

	p.calcportion = calcportion
}

func (m Meal)сalcMeal(){

	calcmeal := Composition{0,0,0,0}
	for _, portion := range m.portions{
		c := portion.calcportion
		calcmeal.addComposition(c)
	}
	m.calcmeal = calcmeal
}

func (m DayMeals)сalcDayMeal(){
			
	calcdaymeals := Composition{0,0,0,0}
	for _, meel := range m.meals{
		c := meel.calcmeal
		calcdaymeals.addComposition(c)
	}
	m.calcdaymeals = calcdaymeals
}