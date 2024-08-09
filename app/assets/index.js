function calculateSettingAsThemeString({ localStorageTheme, systemSettingDark }) {
  if (localStorageTheme !== null) {
    return localStorageTheme;
  }

  return (systemSettingDark.matches) ? "dark" : "light";
}

const localStorageTheme = localStorage.getItem("theme");
const systemSettingDark = window.matchMedia("(prefers-color-scheme: dark)");

let currentThemeSetting = calculateSettingAsThemeString({ localStorageTheme, systemSettingDark });
document.querySelector("html").setAttribute("data-theme", currentThemeSetting);

const checkbox = document.getElementById("themeToggle");
checkbox.checked = currentThemeSetting == "dark"
checkbox.addEventListener('change', event => {
  const newTheme = currentThemeSetting === "dark" ? "light" : "dark";

  checkbox.checked = newTheme == "dark"
  // update theme attribute on HTML to switch theme in CSS
  document.querySelector("html").setAttribute("data-theme", newTheme);

  // update in local storage
  localStorage.setItem("theme", newTheme);

  // update the currentThemeSetting in memory
  currentThemeSetting = newTheme;
})

document.addEventListener("updateLeaderboard", function() {
	let leaderboardData = document.getElementById('leaderboardData')
	if (!leaderboardData || !leaderboardData.textContent) {
		return;
	}
	const data = JSON.parse(document.getElementById('leaderboardData').textContent);
	if (!data) {
		return;
	}

	let noLeaderboardText = document.getElementById('noSubmissionText');
	const ctxElement = document.getElementById('leaderboardChart')
	if (leaderboardData.length == 0) {
		noLeaderboardText.classList.remove('hidden')
		ctxElement.classList.add('hidden')
		return;
	}

	const ctx = ctxElement.getContext('2d');

	const names = data.map(row => row.Name);
	const points = data.map(row => row.Points);

	ctxElement.classList.remove('hidden')
	noLeaderboardText.classList.add('hidden')
	new Chart(ctx, {
		type: 'bar',
		data: {
			labels: names,
			datasets: [{
				label: 'Xp per team',
				data: points,
				backgroundColor: 'rgba(75, 192, 192, 0.2)',
				borderColor: 'rgba(75, 192, 192, 1)',
				borderWidth: 1
			}]
		},
		options: {
			scales: {
				y: {
					beginAtZero: true
				}
			}
		}
	});
});
