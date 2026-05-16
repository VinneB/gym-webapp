export function initAddWorkoutCustom(root) {
  const container = root.querySelector('[data-feature="add-workout-custom"]')
  if (!container) return
  const jsonDataElement = root.querySelector('#addworkout-data')
  if (!jsonDataElement) return
  const { detailNames } = JSON.parse(jsonDataElement.textContent)
  const detailNames_justName = detailNames.map(p => p.Name)

  const button = root.querySelector('#addworkout-add-input-button')
  if (!button || button._initalized) return
  button._initalized = true

  button.addEventListener("click", () => {
    createWorkoutInput(container, detailNames_justName)
  })

  // Time stuff
  const startTimeInput = root.querySelector('#addworkout-start-time')
  if (!startTimeInput || startTimeInput._initalized) return
  startTimeInput._initalized = true

  const endTimeInput = root.querySelector('#addworkout-end-time')
  if (!endTimeInput || endTimeInput._initalized) return
  endTimeInput._initalized = true

  startTimeInput.addEventListener("change", () => {
    const startTimeInput = root.querySelector('#addworkout-start-time')
    const endTimeInput = root.querySelector('#addworkout-end-time')
    // The end date can't be before the start
    endTimeInput.min = startTimeInput.value
    // Set the send time to one hour after start time
    const startDate = new Date(startTimeInput.value)
    startDate.setHours(startDate.getHours() + 1)
    endTimeInput.value = formatDatetimeLocal(startDate);
  })

  createWorkoutInput(container, detailNames_justName)
}

// Function to create a new exercise input set (dropdown and input)
function createWorkoutInput(container, detailNames) {
  const workoutDetails = container.querySelector("#addworkout-detail-collection")
  const newInputGroup = document.createElement("div");

  const newInputGroupInputs = document.createElement("div")
  newInputGroup.classList.add("mb-5");

  // Tailwind on input div
  newInputGroupInputs.classList.add("border");
  newInputGroupInputs.classList.add("rounded-sm");
  newInputGroupInputs.classList.add("p-2");

  // // Create a label for the select field (Detail Name)
  // const newNameLabel = document.createElement("label");
  // newNameLabel.textContent = "Exercise";
  // newNameLabel.setAttribute("for", "exerciseName");

  // Create new select element with static options
  const newNameSelect = document.createElement("select");
  newNameSelect.setAttribute("name", "exerciseName");
  newNameSelect.setAttribute("id", "exerciseName");
  newNameSelect.required = true;
  newNameSelect.classList.add("mr-4")
  newNameSelect.classList.add("border")
  newNameSelect.classList.add("rounded-sm")
  newNameSelect.classList.add("bg-secondarydark")

  // Add default "Select Detail" option
  const defaultOption = document.createElement("option");
  defaultOption.value = "";
  defaultOption.textContent = "Select Exercise";
  newNameSelect.appendChild(defaultOption);

  // Add static options to the select dropdown (no randomness)
  for (let i = 0; i < detailNames.length; i++) {
    const option = document.createElement("option");
    option.value = detailNames[i];
    option.textContent = detailNames[i];
    newNameSelect.appendChild(option);
  }

  // Rep Field

  const newRepLabel = document.createElement("label");
  newRepLabel.textContent = "Reps:";
  newRepLabel.setAttribute("for", "repAmount");
  newRepLabel.classList.add("mr-1")

  const newRepInput = document.createElement("input");
  newRepInput.setAttribute("type", "number");
  newRepInput.setAttribute("name", "repAmount");
  newRepInput.setAttribute("id", "repAmount");
  newRepInput.setAttribute("value", "1");
  newRepInput.setAttribute("min", "1");
  newRepInput.setAttribute("max", "99");
  newRepInput.setAttribute("placeholder", "# Reps");
  newRepInput.setAttribute("step", "1");
  newRepInput.required = true;
  newRepInput.classList.add("mr-4")
  newRepInput.classList.add("rounded-sm")
  newRepInput.classList.add("border")

  // Weight Field

  const newWeightLabel = document.createElement("label");
  newWeightLabel.textContent = "Weight: ";
  newWeightLabel.setAttribute("for", "weightAmount");
  newWeightLabel.classList.add("mr-1")

  const newWeightInput = document.createElement("input");
  newWeightInput.setAttribute("type", "number");
  newWeightInput.setAttribute("name", "weightAmount");
  newWeightInput.setAttribute("id", "weightAmount");
  newWeightInput.setAttribute("value", "1");
  newWeightInput.setAttribute("min", "1");
  newWeightInput.setAttribute("max", "999");
  newWeightInput.setAttribute("placeholder", "Weight");
  newWeightInput.setAttribute("step", "0.05");
  newWeightInput.required = true;
  newWeightInput.classList.add("mr-4")
  newWeightInput.classList.add("border")
  newWeightInput.classList.add("rounded-sm")

  // Partial Reps

  const newPartialRepLabel = document.createElement("label");
  newPartialRepLabel.textContent = "Partial Reps:";
  newPartialRepLabel.setAttribute("for", "partialRepAmount");
  newPartialRepLabel.classList.add("mr-1")

  const newPartialRepInput = document.createElement("input");
  newPartialRepInput.setAttribute("type", "number");
  newPartialRepInput.setAttribute("name", "partialRepAmount");
  newPartialRepInput.setAttribute("id", "partialRepAmount");
  newPartialRepInput.setAttribute("value", "1");
  newPartialRepInput.setAttribute("min", "1");
  newPartialRepInput.setAttribute("max", "999");
  newPartialRepInput.setAttribute("placeholder", "# Partial Reps");
  newPartialRepInput.setAttribute("step", "1");
  newPartialRepInput.required = true;
  newPartialRepInput.classList.add("mr-4")
  newPartialRepInput.classList.add("rounded-sm")
  newPartialRepInput.classList.add("border")


  // Number of Sets

  const newSetCountLabel = document.createElement("label");
  newSetCountLabel.textContent = "Number of Sets:";
  newSetCountLabel.setAttribute("for", "setAmount");
  newSetCountLabel.classList.add("mr-1")

  const newSetCountInput = document.createElement("input");
  newSetCountInput.setAttribute("type", "number");
  newSetCountInput.setAttribute("name", "setAmount");
  newSetCountInput.setAttribute("id", "setAmount");
  newSetCountInput.setAttribute("value", "1");
  newSetCountInput.setAttribute("min", "1");
  newSetCountInput.setAttribute("max", "999");
  newSetCountInput.setAttribute("placeholder", "# Sets");
  newSetCountInput.setAttribute("step", "1");
  newSetCountInput.required = true;
  newSetCountInput.classList.add("mr-4")
  newSetCountInput.classList.add("border")
  newSetCountInput.classList.add("rounded-sm")

  // Workout Start time

  const startTimeInput = container.querySelector('#addworkout-start-time')
  const startTimeString = startTimeInput.value.slice(11, 16)

  const newTimeLabel = document.createElement("label");
  newTimeLabel.textContent = "Start Time: ";
  newTimeLabel.setAttribute("for", "startTime");
  newTimeLabel.classList.add("mr-1")

  const newTimeInput = document.createElement("input");
  newTimeInput.setAttribute("type", "time");
  newTimeInput.setAttribute("value", startTimeString);
  newTimeInput.setAttribute("name", "startTime");
  newTimeInput.setAttribute("id", "startTime");
  newTimeInput.required = true;
  newTimeInput.classList.add("mr-4")
  newTimeInput.classList.add("border")
  newTimeInput.classList.add("rounded-sm")

  // Create a remove button
  const removeButton = document.createElement("button");
  removeButton.type = "button";
  removeButton.textContent = "Remove";
  removeButton.classList.add("bg-deletebutton");
  removeButton.classList.add("box-border");
  removeButton.classList.add("border");
  removeButton.classList.add("hover:bg-deletebutton-strong");
  removeButton.classList.add("focus:ring-1");
  removeButton.classList.add("focus:ring-black");
  removeButton.classList.add("shadow-xs");
  removeButton.classList.add("shadow-xs");
  removeButton.classList.add("leading-5");
  removeButton.classList.add("rounded-sm");
  removeButton.classList.add("text-sm");
  removeButton.classList.add("px-4");
  removeButton.classList.add("py-1");
  removeButton.classList.add("focus:outline-none");


  // Add remove functionality
  removeButton.addEventListener("click", function() {
    newInputGroup.remove();
  });

  // Append the labels, select, input fields, and remove button to the new group
  //newInputGroup.appendChild(newNameLabel);
  newInputGroup.appendChild(newInputGroupInputs);
  newInputGroupInputs.appendChild(newNameSelect);
  newInputGroupInputs.appendChild(newTimeLabel);
  newInputGroupInputs.appendChild(newTimeInput);
  newInputGroupInputs.appendChild(newWeightLabel);
  newInputGroupInputs.appendChild(newWeightInput);
  newInputGroupInputs.appendChild(newRepLabel);
  newInputGroupInputs.appendChild(newRepInput);
  newInputGroupInputs.appendChild(newPartialRepLabel);
  newInputGroupInputs.appendChild(newPartialRepInput);
  newInputGroupInputs.appendChild(newSetCountLabel);
  newInputGroupInputs.appendChild(newSetCountInput);

  newInputGroup.appendChild(removeButton);

  // Add the new group of inputs to the exerciseDetails container
  workoutDetails.appendChild(newInputGroup);
}

function formatDatetimeLocal(date) {
  const year = date.getFullYear();
  const month = (date.getMonth() + 1).toString().padStart(2, '0'); // Months are 0-indexed
  const day = date.getDate().toString().padStart(2, '0');
  const hours = date.getHours().toString().padStart(2, '0');
  const minutes = date.getMinutes().toString().padStart(2, '0');

  return `${year}-${month}-${day}T${hours}:${minutes}`;
}
