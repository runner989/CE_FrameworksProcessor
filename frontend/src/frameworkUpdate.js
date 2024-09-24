document.getElementById('addFrameworkButton').addEventListener('click',function() {
    window.go.main.App.GetUniqueFrameworks()
        .then(function(records) {
            displayUniqueFrameworks(records);
        })
        .catch(function(err) {
            console.error('Error fetching records:', err);
            alert('Failed to retrieve records.');
        });
});

function displayUniqueFrameworks(frameworks) {
    var modal = document.getElementById('recordsModal');
    var recordsContainer = document.getElementById('recordsContainer');


    var content = '<h3>Frameworks from Lookup Table</h3>';
    content += '<div id="selectedRecordLabel"></div>';
    content += '<div id="tableContainer"><table><thead><tr>';

    content += '<th>Framework</th>';
    content += '</tr></thead><tbody>';

    frameworks.forEach(function(framework, index) {
        content += '<tr data-index="' + index + '">';
        content += '<td>' + framework + '</td>';
        content += '</tr>';
    });
    content += '</tbody></table></div>';
    recordsContainer.innerHTML = content;
    modal.style.display = 'block';

    recordsContainer.innerHTML = content;
    modal.style.display = 'block';

    var tableRows = document.querySelectorAll('#recordsContainer tbody tr');
    tableRows.forEach(function(row, index) {
        row.addEventListener('click', function() {
            tableRows.forEach(function(r) {
                r.classList.remove('selected');
            });
            row.classList.add('selected');
            var selectedRecord = frameworks[index];
            fetchFrameworkDetails(selectedRecord);
        });
    });
}

function fetchFrameworkDetails(selectedFramework) {
    window.selectedFramework = selectedFramework;
    var label = document.getElementById('selectedRecordLabel');
    label.innerHTML = `
        <strong>Selected Framework</strong>
        <br><strong>Name:</strong> ${selectedFramework}
        <div id="selectedFrameworkDetails"></div>
    `;
    window.go.main.App.GetFrameworkDetails(selectedFramework)
        .then(function(frameworkDetails) {
            console.log(frameworkDetails);
            if(!frameworkDetails.EvidenceLibraryMappedName) {
                window.go.main.App.GetMappedFrameworks()
                    .then(function(records) {
                        displayMappedFrameworkRecords(records);
                        fetchFrameworkDetails(selectedFramework);
                    })
                    .catch(function(err) {
                        console.error('Error fetching records:', err);
                        alert('Failed to retrieve records.');
                    });
            } else if (!frameworkDetails.AirtableFramework || !frameworkDetails.AirtableView) {
                var fields = frameworkDetails.fields;
                console.log("fields", frameworkDetails.fields)
                var name = frameworkDetails.CEName || 'N/A';
                var mappedName = frameworkDetails.EvidenceLibraryMappedName || 'N/A';
                var uatStage = frameworkDetails.FrameworkId_UAT || 'N/A';
                var prodNumber = frameworkDetails.FrameworkId_Prod || 'N/A';
                var stageNumber = frameworkDetails.FrameworkId_Staging || 'N/A';
                window.selectedFrameworkDetails = {
                    name: name,
                    mappedName: mappedName,
                    uatStage: uatStage,
                    stageNumber: stageNumber,
                    prodNumber: prodNumber
                }
                window.selectedMissingFramework = selectedFramework;
                fetchAirtableTablesandViews();
            } else {
                console.log('fetching framework from Airtable!')
                fetchFrameworkFromAirtable(selectedFramework);
            }
        })
        .catch(function(err) {
            console.error('Error fetching framework details:', err);
            alert('Failed to retrieve framework details.');
        });
}

function displayMappedFrameworkRecords(records) {
    var modal = document.getElementById('recordsModal');
    var recordsContainer = document.getElementById('recordsContainer');


    var content = '<h3>Frameworks From Mapping Table</h3>';
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

function fetchFrameworkFromAirtable(selectedFramework) {
    console.log('Fetching framework from Airtable for: ' + selectedFramework);
    // Placeholder for fetching actual data from Airtable
    var modal = document.getElementById('recordsModal');
    modal.style.display = 'none';
}