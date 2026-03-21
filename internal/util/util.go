package util

import (
	"slices"
	"time"

	"github.com/VinneB/gym-webapp/internal/structapi"
)

func RemoveDuplicates[T comparable](arr []T) []T {
	encountered := map[T]bool{}
	result := []T{}
	for i := range arr {
		if !(encountered[arr[i]] == true) {
			encountered[arr[i]] = true
			result = append(result, arr[i])
		}
	}
	return result
}

// Pair represents a key-value pair.
type pair[K comparable, V any] struct {
	Key   K
	Value V
}

// SortMapByValues converts a map to a slice of Pairs and sorts it by the values
// using the provided custom comparator function.
func SortMapByValuesReturnKeys[K comparable, V any](
	m map[K]V,
	less func(a, b V) bool,
) []K {
	var pairs []pair[K, V]
	for k, v := range m {
		pairs = append(pairs, pair[K, V]{Key: k, Value: v})
	}

	// Sort the slice using slices.SortFunc and the custom comparator
	slices.SortFunc(pairs, func(a, b pair[K, V]) int {
		if less(a.Value, b.Value) {
			return -1 // a comes before b
		}
		if less(b.Value, a.Value) {
			return 1 // b comes before a
		}
		return 0 // a and b are equal
	})

	keys := []K{}
	for _, pair := range pairs {
		keys = append(keys, pair.Key)
	}

	return keys
}

func ConcatArraysNoDuplicates[T comparable](arrs ...[]T) []T {
	result := []T{}
	for _, arr := range arrs {
		for _, element := range arr {
			if !slices.Contains(arr, element) {
				result = append(result, element)
			}
		}
	}
	return result
}

func TimesToDateTimeString(times []time.Time) []string {
	result := []string{}
	for _, time := range times {
		result = append(result, time.Format("2006-01-02 03:04 PM"))
	}
	return result
}

func SeparateSetsIntoMapByWorkoutId(sets []structapi.Set) map[int][]structapi.Set {
	returnArr := map[int][]structapi.Set{}
	for _, set := range sets {
		setArr, setArr_ok := returnArr[set.WorkoutId]
		if !setArr_ok {
			setArr = []structapi.Set{}
		}
		setArr = append(setArr, set)
		returnArr[set.WorkoutId] = setArr
	}
	return returnArr
}

func SortWorkoutsByDate(workouts []structapi.WorkoutInstance) []structapi.WorkoutInstance {
	slices.SortFunc(workouts, func(a, b structapi.WorkoutInstance) int {
		if a.StartTime.Before(b.StartTime) {
			return -1 // a comes before b
		}
		if b.StartTime.Before(b.StartTime) {
			return 1 // b comes before a
		}
		return 0 // a and b are equal
	})
	return workouts
}

func WorkoutDuration(workout structapi.WorkoutInstance) string {
	return workout.EndTime.Sub(workout.StartTime).String()
}
