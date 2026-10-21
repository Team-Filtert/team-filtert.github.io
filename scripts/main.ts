const text = "A Neuro-sama fan game presented by Team Filtered";
const reduceMotion = window.matchMedia("(prefers-reduced-motion: reduce)").matches;
const typing = document.getElementById("typing");

if (typing) {
  if (reduceMotion) {
    typing.textContent = text;
  } else {
    let i = 0;

    const type = () => {
      if (i < text.length) {
        typing.textContent += text[i];
        i++;
        setTimeout(type, 40);
      }
    };

    type();
  }
}

const stars = document.getElementById("stars");

if (stars && !reduceMotion) {
  for (let i = 0; i < 200; i++) {
    const star = document.createElement("div");
    const middle = document.createElement("div");
    const inner = document.createElement("div");
    star.classList.add("star");
    middle.classList.add("star-middle");
    inner.classList.add("star-inner");

    star.style.animationDelay = Math.random() * 9 + "s";
    star.style.animationDuration = 80 + Math.random() * 80 + "s";

    middle.style.animationDelay = Math.random() * 11 + "s";
    middle.style.animationDuration = 90 + Math.random() * 90 + "s";

    inner.style.left = Math.random() * 100 + "vw";
    inner.style.top = Math.random() * 100 + "vh";

    const size = Math.random() * 3 + 1;

    inner.style.width = size + "px";
    inner.style.height = size + "px";

    inner.style.animationDelay = Math.random() * 5 + "s";
    inner.style.animationDuration = 5 + Math.random() * 10 + "s";

    middle.appendChild(inner);
    star.appendChild(middle);
    stars.appendChild(star);
  }
}

export {};
