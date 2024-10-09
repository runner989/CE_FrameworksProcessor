document.getElementById("mappingCountsButton").addEventListener("click", function() {
    window.go.main.App.GetMappingCounts()
        .then(function(records) {
            // console.log('Records:', records)
            displayMappingCountsRecords(records);
        })
        .catch(function(err) {
            console.error('Error fetching records:', err);
            alert('Failed to retrieve records.');
        });
})

function displayMappingCountsRecords(records) {
    let modal = document.getElementById('recordsModal');
    let container = document.getElementById('recordsContainer');
    let content = '<h3>Framework Mapping Counts</h3>';
    content += `<table><thead><tr><th>Framework</th><th>Count of EvIDs</th></tr></thead><tbody>`;
    records.forEach(function(record) {
        content += `<tr><td>${record.Framework}</td><td>${record.Count}</td></tr>`;
    });

    content += `</tbody></table>`;
    container.innerHTML = content;
    modal.style.display = 'block';
}