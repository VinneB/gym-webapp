package structapi

import (
	"html/template"
)

type Data struct {
	Page                string
	Title               string
	ContentTemplate     string
	Exercises           []Exercise
	Muscles             []string
	MusclesJson         template.JS
	ChartDataSets       []ChartDataSet
	ChartDataSetsJson   template.JS
	ExercisesJson       template.JS
	CalcTypes           []string
	Errors              []string
	GraphDatasetOptions []string
	WorkoutSummary      []WorkoutSummary
	WorkoutPlans        []PlanWorkout
	WorkoutPlansJson    template.JS
	WorkoutInstanceType string
}

type ChartPoint struct {
	Label string
	Data  float32
}

type ChartDataSet struct {
	Label    string
	CalcType string
	Points   []ChartPoint
}

type ExerciseJsonFormat struct {
	Name            string
	MuscleFractions []MuscleFraction
}

func ExerciseToJsonFormat(exercises []Exercise) []ExerciseJsonFormat {
	jsonExercises := []ExerciseJsonFormat{}
	for _, exercise := range exercises {
		jsonExercises = append(jsonExercises, ExerciseJsonFormat{Name: exercise.Name, MuscleFractions: exercise.MuscleFractions})
	}
	return jsonExercises
}

type WorkoutSummary struct {
	StartTime      string
	Duration       string
	Sets           []Set
	ExerciseVolume []ChartPoint
	MuscleVolume   []ChartPoint
}

type SetCollection struct {
	ExerciseName string
	NumberOfSets int16
	Format       string
	RepRange     []struct {
		Upper int16
		Lower int16
	}
	RIRs []int16
	Type string
}

type SetCollectionFormat int

const (
	ConstantRepsAndConstantRIRSets SetCollectionFormat = iota
	ConstantReps
	ConstantRIR
	NeitherConstant
)

func (f SetCollectionFormat) String() string {
	if f == ConstantRepsAndConstantRIRSets {
		return "ConstantRepsAndConstantSets"
	} else if f == ConstantReps {
		return "ConstantReps"
	} else if f == ConstantRIR {
		return "ConstantRIR"
	} else {
		return "NeitherConstant"
	}
}

// Assumes all sets have the same exerciseName
func setsToSetCollection(sets []PlanSet) SetCollection {
	setCollection := SetCollection{NumberOfSets: int16(len(sets)), ExerciseName: sets[0].ExerciseName}
	isRIRConstant := true
	isRepsConstant := true
	prevRIR := sets[0].RIR
	prevRepUpperRange := sets[0].RepUpperRange
	prevRepLowerRange := sets[0].RepLowerRange
	for _, set := range sets {
		repRangeStruct := struct {
			Upper int16
			Lower int16
		}{Lower: set.RepLowerRange, Upper: set.RepUpperRange}
		setCollection.RIRs = append(setCollection.RIRs, set.RIR)
		setCollection.RepRange = append(setCollection.RepRange, repRangeStruct)
		if set.RepUpperRange != prevRepUpperRange || set.RepLowerRange != prevRepLowerRange {
			isRepsConstant = false
		}
		if set.RIR != prevRIR {
			isRIRConstant = false
		}
		prevRIR = set.RIR
		prevRepUpperRange = set.RepUpperRange
		prevRepLowerRange = set.RepLowerRange
	}
	if isRIRConstant && isRepsConstant {
		setCollection.Format = ConstantRepsAndConstantRIRSets.String()
	} else if isRIRConstant {
		setCollection.Format = ConstantRIR.String()
	} else if isRepsConstant {
		setCollection.Format = ConstantReps.String()
	} else {
		setCollection.Format = NeitherConstant.String()
	}
	return setCollection
}

func SetsToSetCollections(sets []PlanSet) []SetCollection {
	setsByExerciseNameMap := map[struct {
		name       string
		typestring string
	}][]PlanSet{}
	setCollectionArr := []SetCollection{}
	for _, set := range sets {
		exerciseNameTypePair := struct {
			name       string
			typestring string
		}{name: set.ExerciseName, typestring: set.Type}
		setsByExerciseNameMap[exerciseNameTypePair] = append(setsByExerciseNameMap[exerciseNameTypePair], set)
	}
	for _, setArr := range setsByExerciseNameMap {
		setCollectionArr = append(setCollectionArr, setsToSetCollection(setArr))
	}
	return setCollectionArr
}
