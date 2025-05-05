let current = 0;
const sections = document.querySelectorAll("section");

window.addEventListener(
  "wheel",
  (e) => {
    e.preventDefault(); // impede scroll normal
    if (e.deltaY > 0 && current < sections.length - 1) {
      current++;
    } else if (e.deltaY < 0 && current > 0) {
      current--;
    }
    sections[current].scrollIntoView({ behavior: "smooth" });
  },
  { passive: false }
); // necessário pro preventDefault funcionar no wheel
