const buttons = document.querySelectorAll<HTMLButtonElement>("[data-credit-filter]");
const cards = document.querySelectorAll<HTMLElement>("[data-tags]");
const sections = document.querySelectorAll<HTMLElement>("[data-credit-member]");
const empty = document.querySelector<HTMLElement>("[data-credit-empty]");

function apply(tag: string) {
  let visible = 0;

  cards.forEach((card) => {
    const tags = (card.dataset.tags ?? "").split(",").filter(Boolean);
    const show = tag === "" || tags.includes(tag);
    card.hidden = !show;
    if (show) {
      visible += 1;
    }
  });

  sections.forEach((section) => {
    const shown =
      section.querySelector<HTMLElement>("[data-tags]:not([hidden])") !== null;
    section.hidden = !shown;
  });

  if (empty) {
    empty.hidden = visible > 0;
  }

  buttons.forEach((button) => {
    const active = (button.dataset.creditFilter ?? "") === tag;
    button.classList.toggle("active", active);
    button.setAttribute("aria-pressed", String(active));
  });
}

buttons.forEach((button) => {
  button.addEventListener("click", () => apply(button.dataset.creditFilter ?? ""));
});

apply("");

export {};
