const text: string = "A Neuro-sama fan game presented by Team Filtered";
const reduceMotion: boolean = window.matchMedia(
  "(prefers-reduced-motion: reduce)",
).matches;
const typing: HTMLElement | null = document.getElementById("typing");

if (typing !== null) {
  if (reduceMotion) {
    typing.textContent = text;
  } else {
    let i: number = 0;

    const type: () => void = (): void => {
      if (i < text.length) {
        typing.textContent += text[i];
        i += 1;
        window.setTimeout(type, 40);
      }
    };

    type();
  }
}

const stars: HTMLElement | null = document.getElementById("stars");

if (stars !== null && !reduceMotion) {
  for (let i: number = 0; i < 200; i += 1) {
    const star: HTMLDivElement = document.createElement("div");
    const middle: HTMLDivElement = document.createElement("div");
    const inner: HTMLDivElement = document.createElement("div");
    star.classList.add("star");
    middle.classList.add("star-middle");
    inner.classList.add("star-inner");

    star.style.animationDelay = Math.random() * 9 + "s";
    star.style.animationDuration = 80 + Math.random() * 80 + "s";

    middle.style.animationDelay = Math.random() * 11 + "s";
    middle.style.animationDuration = 90 + Math.random() * 90 + "s";

    inner.style.left = Math.random() * 100 + "vw";
    inner.style.top = Math.random() * 100 + "vh";

    const size: number = Math.random() * 3 + 1;

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
