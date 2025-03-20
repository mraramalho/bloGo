document.addEventListener("DOMContentLoaded", () => {
  document.querySelectorAll("pre").forEach((pre) => {
    // Criar um container para o bloco de código
    const container = document.createElement("div");
    container.classList.add("col", "code-container");

    // Criar um botão de copiar
    const button = document.createElement("button");
    button.innerText = "Copiar";
    button.classList.add("copy-btn", "btn", "btn-tag");

    // Adicionar evento de clique no botão
    button.addEventListener("click", () => {
      const code = pre.querySelector("code").innerText;
      navigator.clipboard.writeText(code).then(() => {
        button.innerText = "Copiado!";
        setTimeout(() => (button.innerText = "Copiar"), 2000);
      });
    });

    // Adicionar o botão e o bloco de código ao container
    pre.parentNode.replaceChild(container, pre);
    container.appendChild(pre);
    container.appendChild(button);

    // Garantir que o pre tenha posição relativa para posicionamento do botão
    pre.style.position = "relative";

    // Mover o botão para dentro do pre
    pre.appendChild(button);
  });
});
