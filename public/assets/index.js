(() => {
  // app/assets/index.js
  console.log("hello from index.js");
  document.addEventListener("updateLeaderboard", function() {
    const data = JSON.parse(document.getElementById("leaderboardData").textContent);
    const names = data.map((row) => row.Name);
    const points = data.map((row) => row.Points);
    const ctx = document.getElementById("leaderboardChart").getContext("2d");
    new Chart(ctx, {
      type: "bar",
      data: {
        labels: names,
        datasets: [{
          label: "Xp per team",
          data: points,
          backgroundColor: "rgba(75, 192, 192, 0.2)",
          borderColor: "rgba(75, 192, 192, 1)",
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
})();
