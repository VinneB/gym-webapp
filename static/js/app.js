import { initAddExercise } from "./features/addExercise.js"
import { initAddWorkoutCustom } from "./features/addWorkoutCustom.js"
import { initRenderGraph } from "./features/renderGraph.js"
import { initAddWorkoutPlan } from "./features/addWorkoutPlan.js"
import { initAccordian } from "./features/accordian.js"
import { initSubAccordian } from "./features/sub-accordian.js"

const features = [
  initAddExercise,
  initAddWorkoutCustom,
  initRenderGraph,
  initAddWorkoutPlan,
  initAccordian,
  initSubAccordian,
]

function initFeatures(root) {
  console.log("init features")
  for (const init of features) {
    init(root)
  }
}


document.body.addEventListener("htmx:beforeSwap", function(evt) {
  if (evt.detail.xhr.status === 409 || evt.detail.xhr.status === 422) {
    evt.detail.shouldSwap = true;
    evt.detail.isError = false;
  }
});

document.body.addEventListener("htmx:afterSwap", function(evt) {
  // Init data-features in returned htmx fragment
  console.log(evt.target)
  initFeatures(evt.target)
  /*
  if (evt.target.id === "add_exercise_form_div") {
    createExerciseInput() 
  }
  if (evt.target.id === "add_workout_form_div") {
    createWorkoutInput() 
  }
  */
});

document.addEventListener("DOMContentLoaded", () => {
  // Init data-feature in base document on page load
  initFeatures(document)
})
