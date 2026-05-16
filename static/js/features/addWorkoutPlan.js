export function initAddWorkoutPlan(root) {
  const container = root.querySelector('[data-feature="add-workoutplan"]')
  if (!container) return
  console.log("found data feature")
  const jsonDataElement = root.querySelector('#addworkoutplan-data')
  if (!jsonDataElement) return
  const { detailNames } = JSON.parse(jsonDataElement.textContent)
  const detailNames_justName = detailNames.map(p => p.Name)
  console.log("stuff")

  const button = root.querySelector('#workoutplans-new-form-input-button')
  if (!button || button._initalized) return
  button._initalized = true

  button.addEventListener("click", () => {
    createWorkoutPlanInput(container, detailNames_justName)
  })

  console.log("stuff")
  createWorkoutPlanInput(container, detailNames_justName)
}

function createWorkoutPlanInput(container, detailNames) {
  const workoutDetails = container.querySelector("#workoutplans-new-form-set-collection")
  const newInputGroup = document.createElement("div");
  newInputGroup.classList.add("exercise-inputs");
  newInputGroup.classList.add("mb-4");

  // Create a label for the select field (Detail Name)
  const newNameLabel = document.createElement("label");
  newNameLabel.textContent = "Exercise:";
  newNameLabel.setAttribute("for", "exercise-name");
  newNameLabel.classList.add("pr-4")

  // Create new select element with static options
  const newNameSelect = document.createElement("select");
  newNameSelect.setAttribute("name", "exercise-name");
  newNameSelect.setAttribute("id", "exercise-name");
  newNameSelect.required = true;

  // Add default "Select Detail" option
  const defaultOption = document.createElement("option");
  defaultOption.value = "";
  defaultOption.textContent = "Select Exercise";
  newNameSelect.appendChild(defaultOption);

  newNameSelect.classList.add("mr-4")
  newNameSelect.classList.add("rounded-sm")
  newNameSelect.classList.add("border")
  newNameSelect.classList.add("bg-secondarydark")

  // Add static options to the select dropdown (no randomness)
  for (let i = 0; i < detailNames.length; i++) {
    const option = document.createElement("option");
    option.value = detailNames[i];
    option.textContent = detailNames[i];
    newNameSelect.appendChild(option);
  }

  // Upper Rep Limit Field

  const newUpperRepLabel = document.createElement("label");
  newUpperRepLabel.textContent = "Upper Rep Limit:";
  newUpperRepLabel.setAttribute("for", "upper-rep-limit");
  newUpperRepLabel.classList.add("pr-4")

  const newUpperRepInput = document.createElement("input");
  newUpperRepInput.setAttribute("type", "number");
  newUpperRepInput.setAttribute("name", "upper-rep-limit");
  newUpperRepInput.setAttribute("id", "upper-rep-limit");
  newUpperRepInput.setAttribute("value", "10");
  newUpperRepInput.setAttribute("min", "1");
  newUpperRepInput.setAttribute("max", "99");
  newUpperRepInput.setAttribute("placeholder", "# Reps");
  newUpperRepInput.setAttribute("step", "1");
  newUpperRepInput.required = true;

  newUpperRepInput.classList.add("mr-4")
  newUpperRepInput.classList.add("rounded-sm")
  newUpperRepInput.classList.add("border")

  // Lower Rep Limit Field

  const newLowerRepLabel = document.createElement("label");
  newLowerRepLabel.textContent = "Lower Rep Limit:";
  newLowerRepLabel.setAttribute("for", "lower-rep-limit");
  newLowerRepLabel.classList.add("pr-4")

  const newLowerRepInput = document.createElement("input");
  newLowerRepInput.setAttribute("type", "number");
  newLowerRepInput.setAttribute("name", "lower-rep-limit");
  newLowerRepInput.setAttribute("id", "lower-rep-limit");
  newLowerRepInput.setAttribute("value", "10");
  newLowerRepInput.setAttribute("min", "1");
  newLowerRepInput.setAttribute("max", "99");
  newLowerRepInput.setAttribute("placeholder", "# Reps");
  newLowerRepInput.setAttribute("step", "1");
  newLowerRepInput.required = true;

  newLowerRepInput.classList.add("mr-4")
  newLowerRepInput.classList.add("rounded-sm")
  newLowerRepInput.classList.add("border")

  // RIR

  const newRIRLabel = document.createElement("label");
  newRIRLabel.textContent = "RIR:";
  newRIRLabel.setAttribute("for", "reps-in-reserve");
  newRIRLabel.classList.add("pr-4")

  const newRIRInput = document.createElement("input");
  newRIRInput.setAttribute("type", "number");
  newRIRInput.setAttribute("name", "reps-in-reserve");
  newRIRInput.setAttribute("id", "reps-in-reserve");
  newRIRInput.setAttribute("value", "0");
  newRIRInput.setAttribute("min", "0");
  newRIRInput.setAttribute("max", "10");
  newRIRInput.setAttribute("placeholder", "# Reps");
  newRIRInput.setAttribute("step", "1");
  newRIRInput.required = true;

  newRIRInput.classList.add("mr-4")
  newRIRInput.classList.add("rounded-sm")
  newRIRInput.classList.add("border")

  // Number of Sets

  const newSetCountLabel = document.createElement("label");
  newSetCountLabel.textContent = "Number of Sets:";
  newSetCountLabel.setAttribute("for", "set-amount");
  newSetCountLabel.classList.add("pr-4")

  const newSetCountInput = document.createElement("input");
  newSetCountInput.setAttribute("type", "number");
  newSetCountInput.setAttribute("name", "set-amount");
  newSetCountInput.setAttribute("id", "set-amount");
  newSetCountInput.setAttribute("value", "1");
  newSetCountInput.setAttribute("min", "1");
  newSetCountInput.setAttribute("max", "999");
  newSetCountInput.setAttribute("placeholder", "# Sets");
  newSetCountInput.setAttribute("step", "1");
  newSetCountInput.required = true;

  newSetCountInput.classList.add("mr-4")
  newSetCountInput.classList.add("rounded-sm")
  newSetCountInput.classList.add("border")

  // Create a remove button
  const removeButton = document.createElement("button");
  removeButton.type = "button";
  removeButton.textContent = "Remove";
  removeButton.classList.add("remove-btn");
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
  newInputGroup.appendChild(newNameLabel);
  newInputGroup.appendChild(newNameSelect);
  newInputGroup.appendChild(newUpperRepLabel);
  newInputGroup.appendChild(newUpperRepInput);
  newInputGroup.appendChild(newLowerRepLabel);
  newInputGroup.appendChild(newLowerRepInput);
  newInputGroup.appendChild(newRIRLabel);
  newInputGroup.appendChild(newRIRInput);
  newInputGroup.appendChild(newSetCountLabel);
  newInputGroup.appendChild(newSetCountInput);

  newInputGroup.appendChild(removeButton);

  // Add the new group of inputs to the exerciseDetails container
  workoutDetails.appendChild(newInputGroup);
}
