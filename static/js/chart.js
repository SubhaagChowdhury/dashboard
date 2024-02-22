document.addEventListener('DOMContentLoaded', async (event) => {
    const ctx = document.getElementById('sharesChart').getContext('2d');

    // Fetch the company data from the Gin backend
    const response = await fetch('/api/company-share');
    const companies = await response.json();

    // Extract the data for the chart
    const labels = companies.map(company => company.Name);
    const initialShares = companies.map(company => company.IShare);
    const finalShares = companies.map(company => company.CShare);

    // Create the chart with the fetched data
    const sharesChart = new Chart(ctx, {
        type: 'bar',
        data: {
            labels: labels, // Company names from the fetched data
            datasets: [{
                label: 'Initial Share',
                data: initialShares, // Initial shares data from the fetched data
                backgroundColor: 'rgba(255, 99, 132, 0.2)',
                borderColor: 'rgba(255, 99, 132, 1)',
                borderWidth: 1
            }, {
                label: 'Final Share',
                data: finalShares, // Final shares data from the fetched data
                backgroundColor: 'rgba(54, 162, 235, 0.2)',
                borderColor: 'rgba(54, 162, 235, 1)',
                borderWidth: 1
            }]
        },
        options: {
            scales: {
                y: {
                    beginAtZero: true
                }
            },
            plugins: {
                legend: {
                    position: 'top',
                },
                title: {
                    display: true,
                    text: 'Company Shares Comparison'
                }
            }
        }
    });
});
