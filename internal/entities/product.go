package entities

type Composition struct {
	Calories      float64
	Proteins      float64
	Fats          float64
	Carbohydrates float64
}

type Product struct {
	ID              int64
	Name            string
	UnitComposition Composition
	Unit            string
	AmountUnit      float64
}

func round(x float64) float64 {
	return float64(int64(x*100)) / 100
}

func (c *Composition) roundComposition() {
	c.Calories = round(c.Calories)
	c.Proteins = round(c.Proteins)
	c.Fats = round(c.Fats)
	c.Carbohydrates = round(c.Carbohydrates)
}

func (p *Product) divProductOnUnit() Composition {

	amountUnit := p.AmountUnit
	comp := p.UnitComposition

	comp.divCompositionOn(amountUnit)
	return comp
}

func (c *Composition) multCompositionOn(delta float64) {
	c.Calories *= delta
	c.Proteins *= delta
	c.Fats *= delta
	c.Carbohydrates *= delta
}

func (c *Composition) divCompositionOn(delta float64) {
	c.Calories /= delta
	c.Proteins /= delta
	c.Fats /= delta
	c.Carbohydrates /= delta
}

func (c *Composition) addComposition(composition Composition) {
	c.Calories += composition.Calories
	c.Proteins += composition.Proteins
	c.Fats += composition.Fats
	c.Carbohydrates += composition.Carbohydrates
}
