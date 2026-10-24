const cards: HTMLElement[] = Array.from(
  document.querySelectorAll<HTMLElement>(".card"),
);

function setExpanded(card: HTMLElement, expanded: boolean): void {
  const header: HTMLElement | null =
    card.querySelector<HTMLElement>(".card-header");
  if (header === null) {
    return;
  }

  card.classList.toggle("expanded", expanded);
  header.setAttribute("aria-expanded", String(expanded));
}

function expandCard(card: HTMLElement): void {
  cards.forEach((other: HTMLElement): void => {
    setExpanded(other, other === card);
  });
}

function toggleCard(card: HTMLElement): void {
  const header: HTMLElement | null =
    card.querySelector<HTMLElement>(".card-header");
  const expanded: boolean = header?.getAttribute("aria-expanded") === "true";

  if (expanded) {
    setExpanded(card, false);
  } else {
    expandCard(card);
  }
}

function isLink(target: EventTarget | null): boolean {
  return target instanceof Element && target.closest("a") !== null;
}

cards.forEach((card: HTMLElement): void => {
  const header: HTMLElement | null =
    card.querySelector<HTMLElement>(".card-header");
  if (header === null) {
    return;
  }

  header.addEventListener("click", (event: MouseEvent): void => {
    if (isLink(event.target)) {
      return;
    }
    toggleCard(card);
  });

  header.addEventListener("keydown", (event: KeyboardEvent): void => {
    if (event.key !== "Enter" && event.key !== " ") {
      return;
    }

    if (isLink(event.target)) {
      return;
    }

    event.preventDefault();
    toggleCard(card);
  });
});

const hash: string = window.location.hash.slice(1);

if (hash !== "") {
  const target: HTMLElement | null = document.getElementById(hash);
  const card: HTMLElement | null =
    target?.closest<HTMLElement>(".card") ?? null;

  if (card !== null) {
    expandCard(card);
  }
}

export {};
