document.addEventListener('DOMContentLoaded', () => {
		const navigation = [
		{ name: 'Product', href: '#' },
		{ name: 'Features', href: '#' },
		{ name: 'Marketplace', href: '#' },
		{ name: 'Company', href: '#' },
		];

		const desktopNav = document.getElementById('desktopNav');
		const mobileNav = document.getElementById('mobileNav');
		navigation.forEach(item => {
				const d = document.createElement('a');
				d.href = item.href;
				d.textContent = item.name;
				d.className = 'text-sm/6 font-semibold text-white';
				desktopNav.appendChild(d);

				const m = document.createElement('a');
				m.href = item.href;
				m.textContent = item.name;
				m.className =
				'-mx-3 block rounded-lg px-3 py-2 text-base/7 font-semibold text-white hover:bg-gray-800';
				mobileNav.appendChild(m);
				});

		document
			.getElementById('openMobileMenu')
			.addEventListener('click', () => document.getElementById('mobileMenu').classList.remove('hidden'));
		document
			.getElementById('closeMobileMenu')
			.addEventListener('click', () => document.getElementById('mobileMenu').classList.add('hidden'));
});
