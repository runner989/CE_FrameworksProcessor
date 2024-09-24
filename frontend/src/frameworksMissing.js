
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
        <div id="selectedFrameworkDetails"></div>
    `;
    openFrameworkBuildListModal();
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
        <br><button id="selectTableViewButton">Select Framework Table and View</button>
    `;

    window.selectedFrameworkDetails = {
        name: name,
        uatStage: uatStage,
        stageNumber: stageNumber,
        prodNumber: prodNumber
    }

    document.getElementById('selectTableViewButton').addEventListener('click', function() {
        fetchAirtableTablesandViews();
    });
}

function fetchAirtableTablesandViews() {
    window.go.main.App.GetAirtableBaseTables()
        .then(function(response) {
            displayFrameworkTablesModal(response.tables);
        })
        .catch(function(err) {
            console.error('Error fetching Airtable tables and views.', err);
            alert('Failed to fetch Airtable tables and views.');
        })
}



function closeFrameworkTablesModalX() {
    var modal = document.getElementById('frameworkTablesModal');
    modal.style.display = 'none';
}

function displayFrameworkTablesModal(tables) {
    document.getElementById('closeFrameworkTablesModal').addEventListener('click',function() {
        closeFrameworkTablesModalX();
    });

    var modal = document.getElementById('frameworkTablesModal');
    var container = document.getElementById('frameworkTablesContainer');

    var content = '<h4>Select a Framework Table and View</h4><ul>';

    tables.forEach(function (table, index) {
        content += `<li class="table-item" data-index="${index}">
            <span class="table-name">${table.name}</span>
            <ul class="views-list" id="views-${index}" style="display: none;">`;

        table.views.forEach(function(view) {
            content += `<li class="view-item" data-table-name="${table.name}" data-view-name="${view.name}">
                ${view.name}
            </li>`;
        });
        content += `</ul></li>`;
    });

    content += '</ul>';

    container.innerHTML = content;
    modal.style.display = 'block';

    var tableItems = container.querySelectorAll('.table-item');
    tableItems.forEach(function(item) {
        var index = item.getAttribute('data-index');
        var viewsList = document.getElementById(`views-${index}`);

        item.querySelector('.table-name').addEventListener('click', function() {
            if (viewsList.style.display === 'none') {
                viewsList.style.display = 'block';
            } else {
                viewsList.style.display = 'none';
            }
        });
    });

    var viewItems = container.querySelectorAll('.view-item');
    viewItems.forEach(function(item) {
        item.addEventListener('click', function() {
            var tableName = item.getAttribute('data-table-name');
            var viewName = item.getAttribute('data-view-name');
            handleTableViewSelector(tableName, viewName);
        });
    });
}

function handleTableViewSelector(tableName, viewName) {
    window.selectedFrameworkTable = tableName;
    window.selectedFrameworkView = viewName;

    var modal = document.getElementById('frameworkTablesModal');
    modal.style.display = 'none';

    updateMissingFrameworkModalWithTableView();
}

function updateMissingFrameworkModalWithTableView() {
    var detailsDiv = document.getElementById('selectedFrameworkDetails');

    var tableViewInfo = `
    <br><strong>Selected Framework Table and View</strong>
    <br><strong>Table Name:</strong> ${window.selectedFrameworkTable}
    <br><strong>View Name:</strong> ${window.selectedFrameworkView}
    <br><button id="okButton">OK</button>
    `;

    var existingInfo = detailsDiv.querySelector('.table-view-info');

    if (existingInfo) {
        existingInfo.innerHTML = tableViewInfo;
    } else {
        var div = document.createElement('div');
        div.classList.add('table-view-info');
        div.innerHTML = tableViewInfo;
        detailsDiv.appendChild(div);
    }

    document.getElementById('okButton').addEventListener('click', function() {
        updateFrameworkLookup();
    });
}

function updateFrameworkLookup() {
    var data = {
        missingFrameworkName: window.selectedMissingFramework,
        selectedFrameworkDetails: {
            cename: window.selectedFrameworkDetails.name,
            uatStage: window.selectedFrameworkDetails.uatStage,
            prodNumber: window.selectedFrameworkDetails.prodNumber,
            stageNumber: window.selectedFrameworkDetails.stageNumber,
            tableName: window.selectedFrameworkTable,
            viewName: window.selectedFrameworkView,
        }
    };
    console.log(data)
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
