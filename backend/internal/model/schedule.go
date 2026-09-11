package model

type Subgroup int

const (
	SubgroupAll Subgroup = 0
	SubgroupFirst Subgroup = 1
	SubgroupSecond Subgroup = 2
)

type Lesson struct {
	Discipline string `json:"disciplines"`
	Type string `json:"types"`
	TimeStart string `json:"timeStart"`
	TimeEnd   string `json:"timeEnd"`
	Number int `json:"number"`
	Auditorium string `json:"auditorium"`
	Corpus string `json:"corpus"`
	Teacher string `json:"teachers"`
	Subgroup Subgroup `json:"subgroup"`
}

type ScheduleDay struct {
	Date string `json:"date"`
	List []Lesson `json:"list"`
}