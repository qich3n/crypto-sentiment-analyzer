document.addEventListener("DOMContentLoaded", () => {
  const analyzeBtn = document.getElementById("analyzeBtn");
  const coinInput = document.getElementById("coinInput");
  const resultDiv = document.getElementById("result");

  analyzeBtn.addEventListener("click", () => {
    const coin = coinInput.value.trim() || "bitcoin"; // default if empty
    fetch(`/api/analyze?coin=${coin}`)
      .then((res) => res.json())
      .then((data) => {
        if (data.error) {
          resultDiv.textContent = `Error: ${data.error}`;
          return;
        }
        resultDiv.textContent = 
          `Coin: ${data.coin} | Sentiment: ${data.sentimentDirection} (${data.sentimentPercent.toFixed(2)}%)`;
      })
      .catch((err) => {
        console.error(err);
        resultDiv.textContent = "An error occurred while analyzing sentiment.";
      });
  });
});
