package entities

type statistic interface {
	Day()
	Week()
	Month()
}

type Week struct {
	total Composition
	goal  Composition
}

type Month struct {
	total Composition
	goal  Composition
}
