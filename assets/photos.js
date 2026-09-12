// Paint each thumbnail over its white square once it has loaded
for (const thumb of document.querySelectorAll(".photo-thumb img")) {
	const reveal = () => thumb.classList.add("loaded");
	if (thumb.complete) {
		// Flush styles so the class change still runs the transition
		void thumb.offsetWidth;
		reveal();
	} else {
		thumb.addEventListener("load", reveal, { once: true });
		thumb.addEventListener("error", reveal, { once: true });
	}
}

const modal = document.querySelector(".photo-modal");
const image = modal.querySelector("img");
const counter = modal.querySelector(".photo-modal-counter");
const description = modal.querySelector(".photo-modal-description");
const buttons = {
	prev: modal.querySelector(".photo-modal-prev"),
	next: modal.querySelector(".photo-modal-next"),
	close: modal.querySelector(".photo-modal-close"),
};

let photos = [];
let index = 0;

function show(i) {
	index = (i + photos.length) % photos.length;
	image.src = photos[index];
	counter.textContent = `${index + 1} / ${photos.length}`;

	// Preload the neighbours so navigation feels instant
	if (photos.length > 1) {
		for (const offset of [-1, 1]) {
			new Image().src = photos[(index + offset + photos.length) % photos.length];
		}
	}
}

function open(item) {
	photos = [...item.querySelectorAll("a.photo-link")].map((link) => link.href);
	image.alt = item.querySelector(".photo-thumb img").alt;
	description.innerHTML = item.querySelector(".photo-description").innerHTML;

	const single = photos.length === 1;
	buttons.prev.hidden = single;
	buttons.next.hidden = single;
	counter.hidden = single;

	show(0);
	modal.showModal();
}

for (const item of document.querySelectorAll(".photo-item")) {
	item.querySelector(".photo-thumb").addEventListener("click", (event) => {
		event.preventDefault();
		open(item);
	});
}

buttons.prev.addEventListener("click", () => show(index - 1));
buttons.next.addEventListener("click", () => show(index + 1));
buttons.close.addEventListener("click", () => modal.close());

modal.addEventListener("keydown", (event) => {
	if (event.key === "ArrowLeft") show(index - 1);
	if (event.key === "ArrowRight") show(index + 1);
});

// Clicking anywhere but the photo, its caption or a button closes the modal
modal.addEventListener("click", (event) => {
	if (!event.target.closest("img, figcaption, button")) modal.close();
});
