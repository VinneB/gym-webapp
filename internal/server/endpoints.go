package server

import (
	"log"
	"net/http"
	"slices"
	"strconv"
	"unicode"

	"github.com/VinneB/gym-webapp/internal/sql"
	"github.com/VinneB/gym-webapp/internal/structapi"
)

func ExercisesGetHandler(w http.ResponseWriter, r *http.Request) {
}

func ExercisesPostHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Exercises Post Handler")
	renderer := newTemplate()
	var data structapi.Data
	for {
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
			SendError(w, r, http.StatusInternalServerError, err.Error())
			break
		}
		log.Print(r.Form)
		log.Printf("detailName ya: %s\n", r.Form["detailName"])
		log.Println("POST to exercises")
		// Length Sanitation
		log.Println(len(r.Form["notreal"]))
		if len(r.Form["muscleName"]) < 1 || len(r.Form["muscleDetail"]) != len(r.Form["muscleName"]) || len(r.Form["exerciseName"]) != 1 {
			log.Println("Failed length sanitation")
			SendError(w, r, http.StatusUnprocessableEntity, "Failed length sanitation")
			break
		}
		// Make fractional muscle objects
		var muscles []structapi.MuscleFraction
		for index, value := range r.Form["muscleName"] {
			// Value sanitation
			muscleDetailFloat, err := strconv.ParseFloat(r.Form["muscleDetail"][index], 32)
			capitalizedValue := []rune(value)
			capitalizedValue[0] = unicode.ToUpper(capitalizedValue[0])
			if err != nil {
				log.Println(err)
				SendError(w, r, http.StatusUnprocessableEntity, err.Error())
				break
			}
			if !slices.Contains(structapi.Muscles, string(capitalizedValue)) {
				log.Println("Invalid muscle name")
				SendError(w, r, http.StatusUnprocessableEntity, "Invalid muscle name")
				break
			} else if muscleDetailFloat > 1.0 || muscleDetailFloat < 0.0 {
				log.Println("Muscle fractional detail out of range")
				SendError(w, r, http.StatusUnprocessableEntity, "Invalid muscle name")
				break
			}
			muscles = append(muscles, structapi.MuscleFraction{Name: string(capitalizedValue), Fraction: float32(muscleDetailFloat)})
		}
		exercise := structapi.Exercise{Name: r.Form["exerciseName"][0], MuscleFractions: muscles}
		err = sql.AddExercise(exercise)
		if err != nil {
			log.Println(err)
			SendError(w, r, http.StatusUnprocessableEntity, err.Error())
			break
		}
		// Render
		data, err = getData(r.URL.Path)
		if err != nil {
			log.Println(err)
			SendError(w, r, http.StatusUnprocessableEntity, "Error")
			break
		}
		renderer.Render(w, "add_exercise_list", data)
		break
	}
	renderer.Render(w, "add_exercise_form", data)
}
