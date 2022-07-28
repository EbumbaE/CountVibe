package entities

type Composition struct{
	calories float64
	protein float64
	fat float64
	carbohydrate float64
}

type Product struct{
	name string
	id int
	unitcomposition Composition
	unit string
	amountunit float64
}

func (c Composition) multCompositionOn(delta float64){
	c.calories *= delta
	c.protein *= delta
	c.fat *= delta
	c.carbohydrate *= delta
}

func (c Composition) divCompositionOn(delta float64){
	c.calories /= delta
	c.protein /= delta
	c.fat /= delta
	c.carbohydrate /= delta
}

func (c Composition) addComposition(composition Composition){
	c.calories += composition.calories
	c.protein += composition.protein
	c.fat += composition.fat
	c.carbohydrate += composition.carbohydrate
}