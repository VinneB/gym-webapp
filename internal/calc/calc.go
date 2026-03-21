package calc

import (
	"errors"
	"log"
	"time"

	"github.com/VinneB/gym-webapp/internal/sql"
	"github.com/VinneB/gym-webapp/internal/structapi"
	"github.com/VinneB/gym-webapp/internal/util"
)

var CalcTypes = []string{"Exercise Volume", "Muscle Volume", "Average Exercise Set Volume"}

var DefaultCalcType string = CalcTypes[0]

func GetOptions(calcType string) ([]string, error) {
	var options []string
	switch calcType {
	case CalcTypes[0]:
		optionsRaw, err := sql.GetExercises(0)
		if err != nil {
			return nil, err
		}
		for _, option := range optionsRaw {
			options = append(options, option.Name)
		}
	case CalcTypes[1]:
		options = structapi.Muscles
	case CalcTypes[2]:
		optionsRaw, err := sql.GetExercises(0)
		if err != nil {
			return nil, err
		}
		for _, option := range optionsRaw {
			options = append(options, option.Name)
		}
	}
	return options, nil
}

func Calculate(CalcType string, params ...string) ([]structapi.ChartPoint, []time.Time, error) {
	switch CalcType {
	case CalcTypes[0]:
		return exerciseVolume(params[0], params[1])
	case CalcTypes[1]:
		return muscleVolume(params[0], params[1])
	case CalcTypes[2]:
		return avgSetVolume(params[0], params[1])
	default:
		return nil, nil, errors.New("Invalid Calc Type")
	}
}

func exerciseVolume(exerciseName string, userName string) ([]structapi.ChartPoint, []time.Time, error) {
	sets, err := sql.GetAllSetsOfExercise(exerciseName, userName)
	if err != nil {
		return nil, nil, err
	}
	volumes, _, dates := setsIntoExerciseVolumeAndDates(sets)
	sortedWorkoutIds := util.SortMapByValuesReturnKeys(dates, func(a, b time.Time) bool {
		return a.Before(b)
	})
	points := []structapi.ChartPoint{}
	times := []time.Time{}
	for _, id := range sortedWorkoutIds {
		points = append(points, structapi.ChartPoint{Label: dates[id].Format("2006-01-02 03:04 PM"), Data: volumes[id]})
		times = append(times, dates[id])
	}

	//for key, val := range volumes {
	//	points = append(points, structapi.ChartPoint{Label: dates[key].Format("2006-01-02 03:04 PM"), Data: val})
	//}
	return points, times, nil
}

func avgSetVolume(exerciseName string, userName string) ([]structapi.ChartPoint, []time.Time, error) {
	sets, err := sql.GetAllSetsOfExercise(exerciseName, userName)
	if err != nil {
		return nil, nil, err
	}
	volumes, numberOfSets, dates := setsIntoExerciseVolumeAndDates(sets)
	avgSetVolume := map[int]float32{}
	for workoutId, volume := range volumes {
		avgSetVolume[workoutId] = volume / float32(numberOfSets[workoutId])
	}
	sortedWorkoutIds := util.SortMapByValuesReturnKeys(dates, func(a, b time.Time) bool {
		return a.Before(b)
	})

	points := []structapi.ChartPoint{}
	times := []time.Time{}
	for _, id := range sortedWorkoutIds {
		points = append(points, structapi.ChartPoint{Label: dates[id].Format("2006-01-02 03:04 PM"), Data: avgSetVolume[id]})
		times = append(times, dates[id])
	}
	return points, times, nil
}

func setsIntoExerciseVolumeAndDates(sets []structapi.Set) (map[int]float32, map[int]int, map[int]time.Time) {
	volumes := map[int]float32{}
	numberOfSets := map[int]int{}
	dates := map[int]time.Time{}
	for i := range sets {
		_, ok := volumes[sets[i].WorkoutId]
		if !ok {
			volumes[sets[i].WorkoutId] = float32(sets[i].Weight) * float32(sets[i].Reps)
			numberOfSets[sets[i].WorkoutId] = 1
			dates[sets[i].WorkoutId] = sets[i].Time
		} else {
			volumes[sets[i].WorkoutId] += float32(sets[i].Weight) * float32(sets[i].Reps)
			numberOfSets[sets[i].WorkoutId] += 1
			if sets[i].Time.Before(dates[i]) {
				dates[sets[i].WorkoutId] = sets[i].Time
			}
		}
	}
	return volumes, numberOfSets, dates
}

// Assumes that 'sets' is a slice of sets in which 'exercises' contains sets.ExerciseName
func muscleVolume(muscleName string, userName string) ([]structapi.ChartPoint, []time.Time, error) {
	log.Println("calculating muscle volume")
	sets, exercises, err := sql.GetAllSetsThatTargetAMuscle(muscleName, userName)
	if err != nil {
		return nil, nil, err
	}
	log.Println(exercises)
	log.Println(sets)
	volumes := map[int]float32{}
	dates := map[int]time.Time{}
	exerciseMap := structapi.ExerciseToMap(exercises)
	exerciseNameToFractionMap := map[string]float32{}
	for i := range sets {
		fraction, exerciseNameToFractionMap_ok := exerciseNameToFractionMap[sets[i].ExerciseName]
		if !exerciseNameToFractionMap_ok {
			exercise, exercise_ok := exerciseMap[sets[i].ExerciseName]
			if !exercise_ok {
				continue
			}
			res := structapi.MuscleFractionsToMap(exercise.MuscleFractions)[muscleName].Fraction
			exerciseNameToFractionMap[sets[i].ExerciseName] = res
			fraction = res
		}
		_, volumes_ok := volumes[sets[i].WorkoutId]
		if !volumes_ok {
			volumes[sets[i].WorkoutId] = float32(sets[i].Weight) * float32(sets[i].Reps) * float32(fraction)
			dates[sets[i].WorkoutId] = sets[i].Time
		} else {
			volumes[sets[i].WorkoutId] += float32(sets[i].Weight) * float32(sets[i].Reps) * float32(fraction)
			if sets[i].Time.Before(dates[i]) {
				dates[sets[i].WorkoutId] = sets[i].Time
			}
		}
	}
	sortedWorkoutIds := util.SortMapByValuesReturnKeys(dates, func(a, b time.Time) bool {
		return a.Before(b)
	})

	points := []structapi.ChartPoint{}
	times := []time.Time{}
	for _, id := range sortedWorkoutIds {
		log.Println(volumes[id])
		points = append(points, structapi.ChartPoint{Label: dates[id].Format("2006-01-02 03:04 PM"), Data: volumes[id]})
		times = append(times, dates[id])
	}
	return points, times, nil
}

func ExerciseVolumeForSets(sets []structapi.Set) map[int][]structapi.ChartPoint {
	outerMap := map[int]map[string]structapi.ChartPoint{}
	returnArr := map[int][]structapi.ChartPoint{}
	for _, set := range sets {
		workoutId := set.WorkoutId
		exerciseName := set.ExerciseName
		innerMap, innerMap_ok := outerMap[workoutId]
		if !innerMap_ok {
			innerMap = map[string]structapi.ChartPoint{}
		}
		volume, volume_ok := innerMap[exerciseName]
		if !volume_ok {
			volume = structapi.ChartPoint{Label: exerciseName, Data: float32(0.0)}
		}
		volume.Data += float32(set.Reps) * float32(set.Weight)
		innerMap[exerciseName] = volume
		outerMap[workoutId] = innerMap
	}
	for workoutId, innerMap := range outerMap {
		innerArr := []structapi.ChartPoint{}
		for _, volume := range innerMap {
			innerArr = append(innerArr, volume)
		}
		returnArr[workoutId] = innerArr
	}
	return returnArr
}

func MuscleVolumeForSets(sets []structapi.Set) map[int][]structapi.ChartPoint {
	outerMap := map[int]map[string]structapi.ChartPoint{}
	returnArr := map[int][]structapi.ChartPoint{}
	exerciseMap := map[string][]structapi.MuscleFraction{}
	exercises, err := sql.GetExercises(0)
	if err != nil {
		return nil
	}
	for _, exercise := range exercises {
		exerciseMap[exercise.Name] = exercise.MuscleFractions
	}
	for _, set := range sets {
		workoutId := set.WorkoutId
		fractions := exerciseMap[set.ExerciseName]
		innerMap, innerMap_ok := outerMap[workoutId]
		if !innerMap_ok {
			innerMap = map[string]structapi.ChartPoint{}
		}
		for _, fraction := range fractions {
			volume, volume_ok := innerMap[fraction.Name]
			if !volume_ok {
				volume = structapi.ChartPoint{Label: fraction.Name, Data: 0}
			}
			volume.Data += float32(set.Reps) * float32(set.Weight) * float32(fraction.Fraction)
			innerMap[fraction.Name] = volume
			outerMap[workoutId] = innerMap
		}
	}
	for workoutId, innerMap := range outerMap {
		innerArr := []structapi.ChartPoint{}
		for _, volume := range innerMap {
			innerArr = append(innerArr, volume)
		}
		returnArr[workoutId] = innerArr
	}
	return returnArr
}

func GreatestExerciseVolume(exerciseName string, userName string) {
	volumes, _, _ := exerciseVolume(exerciseName, userName)
	maxVolume := float32(0)
	for _, volumes := range volumes {
		if volumes.Data > maxVolume {
			maxVolume = volumes.Data
		}
	}
}

func GreatestAverageExerciseVolume(exerciseName string, userName string) {
	volumes, _, _ := avgSetVolume(exerciseName, userName)
	maxVolume := float32(0)
	for _, volumes := range volumes {
		if volumes.Data > maxVolume {
			maxVolume = volumes.Data
		}
	}
}
