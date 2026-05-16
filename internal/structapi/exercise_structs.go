package structapi

import "time"

type MuscleFraction struct {
	Name     string  `json:"name"`
	Fraction float32 `json:"fraction"`
}

type Exercise struct {
	Id              int              `db:"id"`
	Name            string           `db:"name"`
	MuscleFractions []MuscleFraction `db:"data"`
}

type ExerciseSqlForm struct {
	Id              int    `db:"id"`
	Name            string `db:"name"`
	MuscleFractions string `db:"data"`
}

type Set struct {
	Id           int       `db:"id"`
	ExerciseName string    `db:"exercise_name"`
	Reps         int16     `db:"reps"`
	PartialReps  int16     `db:"partial_reps"`
	Weight       float32   `db:"weight"`
	WorkoutId    int       `db:"workout_id"`
	Time         time.Time `db:"time"`
	Type         string    `db:"type"`
	UserEmail    string    `db:"user_email"`
}

type PlanSet struct {
	ExerciseName  string  `json:"name"`
	RepUpperRange int16   `json:"upper_rep"`
	RepLowerRange int16   `json:"lower_rep"`
	Weight        float32 `json:"weight"`
	RIR           int16   `json:"rir"`
	Type          string  `json:"type"`
}

type PlanWorkout struct {
	Id             int       `db:"id"`
	Name           string    `db:"name"`
	UserEmail      string    `db:"user_email"`
	Sets           []PlanSet `json:"data"`
	SetCollections []SetCollection
}

type PlanWorkoutSqlForm struct {
	Id        int    `db:"id"`
	Name      string `db:"name"`
	UserEmail string `db:"user_email"`
	Sets      string `db:"data"`
}

type WorkoutInstance struct {
	Id        int       `db:"id"`
	UserEmail string    `db:"user_email"`
	StartTime time.Time `db:"start_time"`
	EndTime   time.Time `db:"end_time"`
	PlanName  string    `db:"plan_name"`
}

var Muscles = []string{"Abs", "Biceps", "Calves", "Forearms", "Front Delts", "Glutes", "Hamstrings", "Lats", "Lower Back", "Pectorals", "Quads", "Rear Delts", "Rhomboids", "Side Delts", "Traps", "Triceps"}

func ExerciseToMap(exercises []Exercise) map[string]Exercise {
	exerciseMap := map[string]Exercise{}
	for _, exercise := range exercises {
		exerciseMap[exercise.Name] = exercise
	}
	return exerciseMap
}

func MuscleFractionsToMap(muscleFractions []MuscleFraction) map[string]MuscleFraction {
	muscleFractionMap := map[string]MuscleFraction{}
	for _, element := range muscleFractions {
		muscleFractionMap[element.Name] = element
	}
	return muscleFractionMap
}
