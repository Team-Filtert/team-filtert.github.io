const cards = Array.from(document.querySelectorAll<HTMLElement>(".card"));

function setExpanded(card: HTMLElement, expanded: boolean) {
  const header = card.querySelector<HTMLElement>(".card-header");
  if (!header) {
    return;
  }

  card.classList.toggle("expanded", expanded);
  header.setAttribute("aria-expanded", String(expanded));
}

function expandCard(card: HTMLElement) {
  cards.forEach((other) => setExpanded(other, other === card));
}

function toggleCard(card: HTMLElement) {
  const header = card.querySelector<HTMLElement>(".card-header");
  const expanded = header?.getAttribute("aria-expanded") === "true";

  if (expanded) {
    setExpanded(card, false);
  } else {
    expandCard(card);
  }
}

cards.forEach((card) => {
  const header = card.querySelector<HTMLElement>(".card-header");
  if (!header) {
    return;
  }

  header.addEventListener("click", (event) => {
    const target = event.target as HTMLElement | null;
    if (target?.closest("a")) {
      return;
    }
    toggleCard(card);
  });

  header.addEventListener("keydown", (event) => {
    if (event.key !== "Enter" && event.key !== " ") {
      return;
    }

    const target = event.target as HTMLElement | null;
    if (target?.closest("a")) {
      return;
    }

    event.preventDefault();
    toggleCard(card);
  });
});

const hash = window.location.hash.slice(1);

if (hash) {
  const target = document.getElementById(hash);
  const card = target?.closest<HTMLElement>(".card");

  if (card) {
    expandCard(card);
  }
}

export {};
