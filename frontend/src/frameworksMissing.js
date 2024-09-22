
document.getElementById('getMissingFrameworksButton').addEventListener('click',function() {
    window.go.main.App.GetMissingFramework()
        .then(function(records) {
            displayMissingRecords(records);
        })
        .catch(function(err) {
            console.error('Error fetching records:', err);
            alert('Failed to retrieve records.');
        });
});

function displayMissingRecords(records) {
    var modal = document.getElementById('recordsModal');
    var recordsContainer = document.getElementById('recordsContainer');


    var content = '<h3>Frameworks Missing From Lookup Table</h3>';
    content += '<div id="selectedRecordLabel"></div>';
    content += '<div id="tableContainer"><table><thead><tr>';

    content += '<th>Framework</th>';
    content += '</tr></thead><tbody>';

    records.forEach(function(record, index) {
        content += '<tr data-index="' + index + '">';
        content += '<td>' + record + '</td>';
        content += '</tr>';    
    });
    content += '</tbody></table></div>';
    recordsContainer.innerHTML = content;
    modal.style.display = 'block';

    addMissingRowEventListeners(records);
}

function addMissingRowEventListeners(records) {
    var tableRows = document.querySelectorAll('#recordsContainer tbody tr');
    tableRows.forEach(function(row, index) {
        row.addEventListener('click', function() {
            tableRows.forEach(function(r) {
                r.classList.remove('selected');
            });
            row.classList.add('selected');
            var selectedRecord = records[index];
            displaySelectedFrameworkDetails(selectedRecord);
        });
    }); 
}

function displaySelectedFrameworkDetails(record) {
    // var fields = record.fields;
    // var name = fields['Name'] || 'N/A';
    // var uatStage = fields['UAT_Stage'] || 'N/A';
    // var prodNumber = fields['Production Framework Number'] || 'N/A';
    // var stageNumber = fields['Stage Framework Number'] || 'N/A';

    var label = document.getElementById('selectedRecordLabel');
    label.innerHTML = `
        <strong>Selected Framework</strong>
        <br><strong>Name:</strong> ${record}
    `;
}