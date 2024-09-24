
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
    window.selectedMissingFramework = record;
    // openFrameworkBuildListModal();
    var label = document.getElementById('selectedRecordLabel');
    label.innerHTML = `
        <strong>Selected Framework</strong>
        <br><strong>Name:</strong> ${record}
<!--        <br><button id="chooseFrameworkButton">Choose Framework from Build List</button>-->
        <div id="selectedFrameworkDetails"></div>
    `;
    openFrameworkBuildListModal();

    // // from here, open the Framework Build List
    // document.getElementById('chooseFrameworkButton').addEventListener('click', function() {
    //     window.selectedMissingFramework = record;
    //     openFrameworkBuildListModal();
    // });
}

function openFrameworkBuildListModal() {
    window.go.main.App.GetFrameworkLookup()
        .then(function(records) {
            displayFrameworkBuildList(records, function(selectedFramework) {
                closeFrameworkBuildListModal();
                displaySelectedFrameworkFromBuildList(selectedFramework);
            });
        })
        .catch(function(err) {
            console.error('Error fetching framework build list:', err);
            alert('Failed to retrieve framework build list.');
        });
}

function displayFrameworkBuildList(records, onFrameworkSelected) {
    var modal = document.getElementById('frameworkBuildModal');
    var recordsContainer = document.getElementById('frameworkBuildContainer');

    var content = '<h3>Frameworks Build List</h3>';
    content += '<div id="tableContainer"><table><thead><tr>';

    orderedFields.forEach(function(field){
        content += '<th>' + field + '</th>';
    });
    content += '</tr></thead><tbody>';

    records.forEach(function(record, index) {
        var fields = record.fields;
        content += '<tr data-index="' + index + '">';
        orderedFields.forEach(function(field) {
            var value = fields[field];
            if (Array.isArray(value)) {
                value = value.join(', ');
            } else if (typeof value == 'object' && value !== null) {
                value = JSON.stringify(value);
            } else if (value == undefined || value == null) {
                value = '';
            }
            content += '<td>' + value + '</td>';
        });
        content += '</tr>';
    });
    content += '</tbody></table></div>';
    recordsContainer.innerHTML = content;
    modal.style.display = 'block';

    addFrameworkBuildRowEventListeners(records, onFrameworkSelected);
}

function addFrameworkBuildRowEventListeners(records, onFrameworkSelected) {
    var tableRows = document.querySelectorAll('#frameworkBuildContainer tbody tr');
    tableRows.forEach(function(row, index) {
        row.addEventListener('click', function() {
            tableRows.forEach(function(r) {
                r.classList.remove('selected');
            });
            row.classList.add('selected');
            var selectedRecord = records[index];
            onFrameworkSelected(selectedRecord);
        });
    });
}

function closeFrameworkBuildListModal() {
    var modal = document.getElementById('frameworkBuildModal');
    modal.style.display = 'none';
}

function displaySelectedFrameworkFromBuildList(selectedFramework) {
    var fields = selectedFramework.fields;
    var name = fields['Name'] || 'N/A';
    var uatStage = fields['UAT_Stage'] || 'N/A';
    var prodNumber = fields['Production Framework Number'] || 'N/A';
    var stageNumber = fields['Stage Framework Number'] || 'N/A';

    var detailsDiv = document.getElementById('selectedFrameworkDetails');
    detailsDiv.innerHTML = `
        <strong>Selected Framework from Build List</strong>
        <br><strong>Name:</strong> ${name}
        <br><strong>UAT Stage:</strong> ${uatStage}
        <br><strong>Staging Number:</strong> ${stageNumber}
        <br><strong>Production Framework Number:</strong> ${prodNumber}
        <br><button id="okButton">OK</button>
    `;

    window.selectedFrameworkDetails = {
        name: name,
        uatStage: uatStage,
        stageNumber: stageNumber,
        prodNumber: prodNumber
    }

    document.getElementById('okButton').addEventListener('click', function() {
        updateFrameworkLookup();
    });
}

function updateFrameworkLookup() {
    var data = {
        missingFrameworkName: window.selectedMissingFramework,
        selectedFrameworkDetails: window.selectedFrameworkDetails
    };

    window.go.main.App.UpdateFrameworkLookup(data)
        .then(function (response) {
            alert('Framework Lookup updated successfully.');
            // Close the Missing Frameworks modal
            var modal = document.getElementById('recordsModal');
            modal.style.display = 'none';
        })
        .catch(function(err) {
            console.error('Error updating Framework Lookup', err);
            alert('Failed to update Framework Lookup.');
        });
}
