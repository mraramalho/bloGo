const spotlight = document.querySelector(".spotlight");

document.addEventListener("mousemove", (e) => {
  const { clientX: x, clientY: y } = e;
  spotlight.style.background = `radial-gradient(circle at ${x}px ${y}px, rgb(55, 135, 132, 0.1) 100px, rgb(15, 23, 42, 0.95) 400px)`;
});
