package structapi

import "time"

type MuscleFraction struct {
	Name     string  `json:"name"`
	Fraction float32 `json:"fraction"`
}

type Exercise struct {
	Name            string `json:"name"`
	MuscleFractions []MuscleFraction
}

type Set struct {
	Reps   int16 `json:"reps"`
	Weight int16 `json:"weight"`
}

type SetCollection struct {
	Exercise string `json:"name"`
	Sets     []Set
}

type WorkoutInstance struct {
	StartTime time.Time `json:"start"`
	EndTime   time.Time `json:"end"`
	Exercises []SetCollection
}
