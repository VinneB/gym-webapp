export function initSubAccordian(root) {

  const containers = root.querySelectorAll('[data-feature="sub-accordian"]')
  console.log("sub-accordian count: " + containers.length)
  if (containers.length === 0) return
  let elements = Array.from(containers)
  for (let i = 0; i < elements.length; i++) {
    let element = elements[i]
    let accordian_header = element.querySelector('.sub-accordian-header')
    accordian_header.addEventListener("click", (event) => {
      let currentAccordian = event.currentTarget.closest(".sub-accordian");
      currentAccordian.classList.toggle("active");

      let accordians = Array.from(event.currentTarget.closest(".sub-accordian-container").querySelectorAll(".sub-accordian"));
      accordians.forEach((item) => {
        if (item != currentAccordian) {
          item.classList.remove("active")
        }
      });
    })
  }
}
