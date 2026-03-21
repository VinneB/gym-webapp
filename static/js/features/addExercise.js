export function initAddExercise(root) {

  const container = root.querySelector('[data-feature="add-exercise"]')
  if (!container) return

  const jsonDataElement = root.querySelector('#addexercise-data')
  if (!jsonDataElement) return

  const { detailNames } = JSON.parse(jsonDataElement.textContent)

  const button = root.querySelector('#addexercise-add-input-button')
  if (!button || button._initalized) return
  button._initalized = true

  button.addEventListener("click", () => {
    createExerciseInput(container, detailNames)
  })

  createExerciseInput(container, detailNames)
}

function createExerciseInput(container, detailNames) {
  const exerciseDetails = container.querySelector("#addexercise-detail-collection")
  const newInputGroup = document.createElement("div");
  newInputGroup.classList.add("exercise-inputs");

  // Create a label for the select field (Detail Name)
  const newNameLabel = document.createElement("label");
  newNameLabel.textContent = "Detail Name";
  newNameLabel.setAttribute("for", "muscleName");

  // Create new select element with static options
  const newNameSelect = document.createElement("select");
  newNameSelect.setAttribute("name", "muscleName");
  newNameSelect.setAttribute("id", "muscleName");
  newNameSelect.required = true;

  // Add default "Select Detail" option
  const defaultOption = document.createElement("option");
  defaultOption.value = "";
  defaultOption.textContent = "Select Detail";
  newNameSelect.appendChild(defaultOption);

  // Add static options to the select dropdown (no randomness)
  for (let i = 0; i < detailNames.length; i++) {
    const option = document.createElement("option");
    option.value = detailNames[i];
    option.textContent = detailNames[i];
    newNameSelect.appendChild(option);
  }

  // Create a label for the input field (Detail Value)
  const newValueLabel = document.createElement("label");
  newValueLabel.textContent = "Detail Value";
  newValueLabel.setAttribute("for", "muscleDetail");

  // Create a new input field for the value (float)
  const newValueInput = document.createElement("input");
  newValueInput.setAttribute("type", "number");
  newValueInput.setAttribute("name", "muscleDetail");
  newValueInput.setAttribute("id", "muscleDetail");
  newValueInput.setAttribute("value", "1");
  newValueInput.setAttribute("min", "0.1");
  newValueInput.setAttribute("max", "1");
  newValueInput.setAttribute("placeholder", "Value (float)");
  newValueInput.setAttribute("step", "0.1");
  newValueInput.required = true;

  // Create a remove button
  const removeButton = document.createElement("button");
  removeButton.type = "button";
  removeButton.textContent = "Remove";
  removeButton.classList.add("remove-btn");

  // Add remove functionality
  removeButton.addEventListener("click", function() {
    newInputGroup.remove();
  });

  // Append the labels, select, input fields, and remove button to the new group
  newInputGroup.appendChild(newNameLabel);
  newInputGroup.appendChild(newNameSelect);
  newInputGroup.appendChild(newValueLabel);
  newInputGroup.appendChild(newValueInput);
  newInputGroup.appendChild(removeButton);

  // Add the new group of inputs to the exerciseDetails container
  exerciseDetails.appendChild(newInputGroup);
}
