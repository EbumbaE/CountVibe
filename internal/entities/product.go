package entities

type Composition struct {
	Calories      float64
	Proteins      float64
	Fats          float64
	Carbohydrates float64
}

type Product struct {
	ID              int
	Name            string
	UnitComposition Composition
	Unit            string
	AmountUnit      float64
}

func (c Composition) multCompositionOn(delta float64) {
	c.Calories *= delta
	c.Proteins *= delta
	c.Fats *= delta
	c.Carbohydrates *= delta
}

func (c Composition) divCompositionOn(delta float64) {
	c.Calories /= delta
	c.Proteins /= delta
	c.Fats /= delta
	c.Carbohydrates /= delta
}

func (c Composition) addComposition(composition Composition) {
	c.Calories += composition.Calories
	c.Proteins += composition.Proteins
	c.Fats += composition.Fats
	c.Carbohydrates += composition.Carbohydrates
}
