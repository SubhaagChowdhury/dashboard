document.addEventListener('DOMContentLoaded', function() {
    // Assuming fetchCompanyData is a function that fetches your data
    fetchCompanyData().then(companies => {
        // Assuming companies is an array of company data
        populateTable(companies);
        createBarChart(companies);
    });
});

function fetchCompanyData() {
    return new Promise((resolve, reject) => {
        // Example AJAX call to fetch company data
        // Adjust URL/path as necessary
        fetch('/api/companies')
            .then(response => response.json())
            .then(data => resolve(data))
            .catch(error => reject(error));
    });
}

function populateTable(companies) {
    const tableBody = document.querySelector('#company-table tbody');
    tableBody.innerHTML = ''; // Clear existing rows
    companies.forEach(company => {
        const row = document.createElement('tr');
        row.innerHTML = `
            <td>${company.Name}</td>
            <td>${company.IShare.toFixed(2)}</td>
            <td>${company.CShare.toFixed(2)}</td>
        `;
        tableBody.appendChild(row);
    });
}

function createBarChart(companies) {
    const labels = companies.map(company => company.Name);
    const iShares = companies.map(company => company.IShare);
    const cShares = companies.map(company => company.CShare);

    const ctx = document.getElementById('company-chart').getContext('2d');
    const chart = new Chart(ctx, {
        type: 'bar',
        data: {
            labels: labels,
            datasets: [{
                label: 'IShare',
                data: iShares,
                backgroundColor: 'rgba(255, 99, 132, 0.2)',
                borderColor: 'rgba(255, 99, 132, 1)',
                borderWidth: 1
            }, {
                label: 'CShare',
                data: cShares,
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
            }
        }
    });
}
