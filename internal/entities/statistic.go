package entities

type statistic interface{
	Day() Composition
	Week() Composition
	Month() Composition
}

type Week struct{
	total Composition
	goal Composition
}

type Month struct{
	total Composition
	goal Composition
}
