const buttons: NodeListOf<HTMLButtonElement> =
  document.querySelectorAll<HTMLButtonElement>("[data-credit-filter]");
const cards: NodeListOf<HTMLElement> =
  document.querySelectorAll<HTMLElement>("[data-tags]");
const sections: NodeListOf<HTMLElement> =
  document.querySelectorAll<HTMLElement>("[data-credit-member]");
const empty: HTMLElement | null =
  document.querySelector<HTMLElement>("[data-credit-empty]");

function readTags(card: HTMLElement): string[] {
  return (card.dataset.tags ?? "").split(",").filter(Boolean);
}

function apply(tag: string): void {
  let visible: number = 0;

  cards.forEach((card: HTMLElement): void => {
    const cardTags: string[] = readTags(card);
    const show: boolean = tag === "" || cardTags.includes(tag);
    card.hidden = !show;
    if (show) {
      visible += 1;
    }
  });

  sections.forEach((section: HTMLElement): void => {
    const shown: boolean =
      section.querySelector("[data-tags]:not([hidden])") !== null;
    section.hidden = !shown;
  });

  if (empty !== null) {
    empty.hidden = visible > 0;
  }

  buttons.forEach((button: HTMLButtonElement): void => {
    const active: boolean = (button.dataset.creditFilter ?? "") === tag;
    button.classList.toggle("active", active);
    button.setAttribute("aria-pressed", String(active));
  });
}

buttons.forEach((button: HTMLButtonElement): void => {
  button.addEventListener("click", (): void => {
    apply(button.dataset.creditFilter ?? "");
  });
});

const knownTags: Set<string> = new Set<string>(
  Array.from(
    buttons,
    (button: HTMLButtonElement): string => button.dataset.creditFilter ?? "",
  ),
);

// The footer links to /credits/#<tag>, so honour the fragment on load and when
// it changes.
function applyFromHash(): void {
  const tag: string = decodeURIComponent(window.location.hash.slice(1));
  apply(knownTags.has(tag) ? tag : "");
}

window.addEventListener("hashchange", applyFromHash);
applyFromHash();

export {};
