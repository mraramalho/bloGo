// Função para dividir o texto em caracteres individuais
function prepareTextForGlow() {
  const glowElements = document.querySelectorAll(".text-glow");

  glowElements.forEach((element) => {
    const text = element.innerText;
    element.innerHTML = "";

    // Divide o texto em caracteres individuais
    for (let i = 0; i < text.length; i++) {
      const charSpan = document.createElement("span");
      charSpan.classList.add("glow-char");

      // Preserva os espaços em branco usando &nbsp; para espaços
      if (text[i] === " ") {
        charSpan.innerHTML = "&nbsp;";
      } else {
        charSpan.innerText = text[i];
      }

      element.appendChild(charSpan);
    }
  });
}

// Função para aplicar o efeito de iluminação baseado na posição do mouse
function handleGlowEffect(e) {
  const chars = document.querySelectorAll(".glow-char");
  const mouseX = e.clientX;
  const mouseY = e.clientY;

  chars.forEach((char) => {
    const rect = char.getBoundingClientRect();
    const charX = rect.left + rect.width / 2;
    const charY = rect.top + rect.height / 2;

    // Calcula a distância entre o mouse e o caractere
    const distance = Math.sqrt(
      Math.pow(mouseX - charX, 2) + Math.pow(mouseY - charY, 2)
    );

    // Define um raio de efeito
    const maxDistance = 100;

    if (distance < maxDistance) {
      // Quanto mais próximo, mais forte o efeito
      const intensity = 1 - distance / maxDistance;
      char.classList.add("glow-active");
      char.style.color = `rgba(85, 209, 192, ${intensity})`;
      char.style.textShadow = `0 0 ${8 * intensity}px rgba(85, 209, 192, ${
        intensity * 0.8
      })`;
    } else {
      char.classList.remove("glow-active");
      char.style.color = "";
      char.style.textShadow = "";
    }
  });
}

// Inicializa o efeito quando o DOM estiver carregado
document.addEventListener("DOMContentLoaded", () => {
  prepareTextForGlow();
  document.addEventListener("mousemove", handleGlowEffect);
});
