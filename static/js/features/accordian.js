export function initAccordian(root) {

  const containers = root.querySelectorAll('[data-feature="accordian"]')
  if (containers.length === 0) return
  let elements = Array.from(containers)
  for (let i = 0; i < elements.length; i++) {
    let element = elements[i]
    let accordian_header = element.querySelector('.accordian-header')
    accordian_header.addEventListener("click", (event) => {
      let currentAccordian = event.currentTarget.closest(".accordian");
      currentAccordian.classList.toggle("active");

      let accordians = Array.from(event.currentTarget.closest(".accordian-container").querySelectorAll(".accordian"));
      accordians.forEach((item) => {
        if (item != currentAccordian) {
          item.classList.remove("active")
        }
      });
    })
  }
}
