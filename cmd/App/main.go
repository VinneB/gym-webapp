package main

import (
	"fmt"
	"os"
	"slices"

	//	"time"

	"github.com/VinneB/gym-webapp/internal/server"
	//"github.com/VinneB/gym-webapp/internal/sql"
	//"github.com/VinneB/gym-webapp/internal/structapi"
)

var cliString string = "--cli"
var serverString string = "--server"
var defaultBehavior string = serverString

func main() {
	if slices.Contains(os.Args[1:], cliString) && !slices.Contains(os.Args[1:], serverString) {
		//		db := sql.Connect()
		//		result := sql.GetAllUserWorkouts(db, "lj@lorenzojones@gmail.com")
		//		for _, value := range result {
		//			fmt.Println(value)
		//		}
		//		sql.Close(db)
		//		fmt.Println("done")
	} else if slices.Contains(os.Args[1:], serverString) && !slices.Contains(os.Args[1:], cliString) {
		//		db := sql.Connect()
		//		time1 := time.Now()
		//		time2 := time.Now().Add(1000)
		//		var workout structapi.WorkoutInstance = structapi.WorkoutInstance{StartTime: time1, EndTime: time2, Exercises: []structapi.SetCollection{
		//			{"workout3", []structapi.Set{{11, 1}, {34, 1}}}}}
		//		sql.AddWorkoutInstance(db, workout, "lj@lorenzojones@gmail.com")
		//		sql.Close(db)
		//	db := sql.Connect()
		//	var exercise structapi.Exercise = structapi.Exercise{Name: "exercise1", MuscleFractions: []structapi.MuscleFraction{{"bicep", 0.1}, {"quads", 0.5}}}
		//	var exercise2 structapi.Exercise = structapi.Exercise{Name: "exercise2", MuscleFractions: []structapi.MuscleFraction{{"biep", 0.2}, {"quads", 0.5}}}
		//	sql.AddExercise(db, exercise)
		//	sql.AddExercise(db, exercise2)
		//	sql.Close(db)
		fmt.Println("Start Server")
		server.StartServer()
	} else if defaultBehavior == cliString {

	} else {
		server.StartServer()
	}

}
