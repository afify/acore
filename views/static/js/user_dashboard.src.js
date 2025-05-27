			document.addEventListener('DOMContentLoaded', () => {
				// Data
				const navigation = [
					{ name: 'Home', href: '/home', icon: 'FolderIcon', current: false },
					{ name: 'Settings', href: '/settings', icon: 'Cog6ToothIcon', current: true },
				];
				const teams = [
					{ id: 1, name: 'Planetaria', href: '#', initial: 'P', current: false },
					{ id: 2, name: 'Protocol', href: '#', initial: 'P', current: false },
					{ id: 3, name: 'Tailwind Labs', href: '#', initial: 'T', current: false },
				];
				const secondaryNavigation = [
					{ name: 'Account', href: '#', current: true },
					{ name: 'Security', href: '#', current: false },
					{ name: 'Integrations', href: '#', current: false },
					{ name: 'Billing', href: '#', current: false },
					{ name: 'Notifications', href: '#', current: false },
				];

				// State (using a simple object to manage sidebar visibility)
				let state = {
					sidebarOpen: false,
				};

				// DOM Elements
				const mobileSidebar = document.getElementById('mobile-sidebar');
				const mobileSidebarPanel = document.getElementById('mobile-sidebar-panel');
				const openSidebarButton = document.getElementById('open-sidebar-button');
				const closeSidebarButton = document.getElementById('close-sidebar-button');
				const mobileMainNavigation = document.getElementById('mobile-main-navigation');
				const mobileTeamNavigation = document.getElementById('mobile-team-navigation');
				const desktopMainNavigation = document.getElementById('desktop-main-navigation');
				const desktopTeamNavigation = document.getElementById('desktop-team-navigation');
				const secondaryNavigationElement = document.getElementById('secondary-navigation');

				// Helper function to get icon SVG (you'd need a more robust system for all icons)
				function getIconSvg(iconName) {
					// This is a simplified example. In a real app, you'd likely
					// have a map of icon names to their SVG paths or load them dynamically.
						const icons = {
							FolderIcon: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-6 shrink-0" aria-hidden="true"><path stroke-linecap="round" stroke-linejoin="round" d="M2.25 12.75V12A2.25 2.25 0 0 1 4.5 9.75h15A2.25 2.25 0 0 1 21.75 12v.75m-4.5-7.5a2.25 2.25 0 0 0-2.25 2.25v2.25m-12.75 0V7.5A2.25 2.25 0 0 1 4.5 5.25h15A2.25 2.25 0 0 1 21.75 7.5v2.25m-18 0V12m0-2.25A2.25 2.25 0 0 0 4.5 12h15a2.25 2.25 0 0 0 2.25-2.25m-18 0a2.25 2.25 0 0 1 2.25-2.25h15M6.75 12v3M17.25 12v3m-10.5-3v3m10.5-3v3M6.75 15v3M17.25 15v3m-10.5-3v3m10.5-3v3" /></svg>`,
							ServerIcon: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-6 shrink-0" aria-hidden="true"><path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15M12 4.5a3 3 0 0 1 3 3v9a3 3 0 0 1-3 3M12 4.5a3 3 0 0 0-3 3v9a3 3 0 0 0 3 3M19.5 7.5h-15m15 0a3 3 0 0 0 3 3v9a3 3 0 0 0-3 3M4.5 7.5h-15m15 0a3 3 0 0 1-3-3v9a3 3 0 0 1 3 3M4.5 7.5h-15m15 0a3 3 0 0 0 3 3v9a3 3 0 0 0-3 3" /></svg>`,
							SignalIcon: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-6 shrink-0" aria-hidden="true"><path stroke-linecap="round" stroke-linejoin="round" d="m11.25 11.25.041-.02A.75.75 0 0 1 12 11.25c2.25 0 4.288 1.134 5.517 2.836A6 6 0 0 0 11.25 9.75v.003ZM12 12.75a.75.75 0 0 1-.75-.75.75.75 0 0 1 .75-.75c2.25 0 4.288 1.134 5.517 2.836A6 6 0 0 0 11.25 9.75v.003ZM12 12.75a.75.75 0 0 1-.75-.75.75.75 0 0 1 .75-.75c2.25 0 4.288 1.134 5.517 2.836A6 6 0 0 0 11.25 9.75v.003Z" /></svg>`,
							GlobeAltIcon: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-6 shrink-0" aria-hidden="true"><path stroke-linecap="round" stroke-linejoin="round" d="m12 21.75 4.5-4.5c.39-.39.39-1.023 0-1.413L12 12l-4.5 4.5c-.39.39-.39 1.023 0 1.413L12 21.75Zm0-18L7.5 7.5c-.39.39-.39 1.023 0 1.413L12 12l4.5-4.5c.39-.39.39-1.023 0-1.413L12 3.75Zm0-18L7.5 7.5c-.39.39-.39 1.023 0 1.413L12 12l4.5-4.5c.39-.39.39-1.023 0-1.413L12 3.75Z" /></svg>`,
							ChartBarSquareIcon: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-6 shrink-0" aria-hidden="true"><path stroke-linecap="round" stroke-linejoin="round" d="M12 21V3m0 0a3 3 0 0 1 3 3v15a3 3 0 0 1-3 3M12 3a3 3 0 0 0-3 3v15a3 3 0 0 0 3 3" /></svg>`,
							Cog6ToothIcon: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-6 shrink-0" aria-hidden="true"><path stroke-linecap="round" stroke-linejoin="round" d="M10.375 2.125a6 6 0 0 0-6 6c0 3.314 2.686 6 6 6h1.5a6 6 0 0 0 6-6c0-3.314-2.686-6-6-6h-1.5Zm0 0a6 6 0 0 1 6 6c0 3.314-2.686 6-6 6h-1.5a6 6 0 0 1-6-6c0-3.314-2.686-6-6-6h-1.5Z" /></svg>`,
						};
					return icons[iconName] || ''; // Return empty string if icon not found
				}

				// Render functions
				function renderNavigation(container, items) {
					container.innerHTML = `<ul role="list" class="-mx-2 space-y-1"></ul>`;
					const ul = container.querySelector('ul');

					items.forEach(item => {
						const li = document.createElement('li');
						const a = document.createElement('a');
						a.href = item.href;
						a.className = item.current ? 'bg-gray-800 text-white' : 'text-gray-400 hover:bg-gray-800 hover:text-white';
						a.className += ' group flex gap-x-3 rounded-md p-2 text-sm/6 font-semibold';
						a.innerHTML = `${getIconSvg(item.icon)} ${item.name}`;
						li.appendChild(a);
						ul.appendChild(li);
					});
				}

				function renderTeams(container, items) {
					container.innerHTML = ''; // Clear previous content
					items.forEach(team => {
						const li = document.createElement('li');
						const a = document.createElement('a');
						a.href = team.href;
						a.className = team.current ? 'bg-gray-800 text-white' : 'text-gray-400 hover:bg-gray-800 hover:text-white';
						a.className += ' group flex gap-x-3 rounded-md p-2 text-sm/6 font-semibold';
						a.innerHTML = `
							<span class="flex size-6 shrink-0 items-center justify-center rounded-lg border border-gray-700 bg-gray-800 text-[0.625rem] font-medium text-gray-400 group-hover:text-white">${team.initial}</span>
							<span class="truncate">${team.name}</span>
							`;
						li.appendChild(a);
						container.appendChild(li);
					});
				}

				function renderSecondaryNavigation(container, items) {
					container.innerHTML = ''; // Clear previous content
					items.forEach(item => {
						const li = document.createElement('li');
						const a = document.createElement('a');
						a.href = item.href;
						a.className = item.current ? 'text-indigo-400' : '';
						a.textContent = item.name;
						li.appendChild(a);
						container.appendChild(li);
					});
				}

				// Update UI based on state
				function updateSidebarVisibility() {
					if (state.sidebarOpen) {
						mobileSidebar.classList.remove('hidden');
						// Add transition classes after a short delay to allow for the initial render
						setTimeout(() => {
							mobileSidebarPanel.classList.remove('-translate-x-full');
							mobileSidebarPanel.classList.add('translate-x-0');
							mobileSidebar.classList.remove('opacity-0');
							mobileSidebar.classList.add('opacity-100');
						}, 10); // A small delay
					} else {
						mobileSidebarPanel.classList.remove('translate-x-0');
						mobileSidebarPanel.classList.add('-translate-x-full');
						mobileSidebar.classList.remove('opacity-100');
						mobileSidebar.classList.add('opacity-0');
						// Hide the sidebar completely after the transition
						mobileSidebarPanel.addEventListener('transitionend', function handler() {
							if (!state.sidebarOpen) {
								mobileSidebar.classList.add('hidden');
							}
							mobileSidebarPanel.removeEventListener('transitionend', handler);
						});
					}
				}

				// Event Listeners
				openSidebarButton.addEventListener('click', () => {
					state.sidebarOpen = true;
					updateSidebarVisibility();
				});

				closeSidebarButton.addEventListener('click', () => {
					state.sidebarOpen = false;
					updateSidebarVisibility();
				});

				// Initial render
				renderNavigation(mobileMainNavigation, navigation);
				renderNavigation(desktopMainNavigation, navigation);
				renderTeams(mobileTeamNavigation, teams);
				renderTeams(desktopTeamNavigation, teams);
				renderSecondaryNavigation(secondaryNavigationElement, secondaryNavigation);
			});
