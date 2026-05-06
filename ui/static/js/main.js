var navLinks = document.querySelectorAll("nav a");
for (var i = 0; i < navLinks.length; i++) {
	var link = navLinks[i]
	if (link.getAttribute('href') == window.location.pathname) {
		link.classList.add("live");
		break;
	}
}

const dateElements = document.querySelectorAll('.local-date');

dateElements.forEach(el => {
	const isoDate = el.getAttribute('datetime');
	if(isoDate) {
		const date = new Date(isoDate);
		el.textContent = new Intl.DateTimeFormat('default', {
			year: 'numeric',
			month: 'short',
			day: 'numeric',
			hour: 'numeric',
			minute: 'numeric',
		}).format(date);
	}
});