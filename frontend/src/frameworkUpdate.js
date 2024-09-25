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
            if(!frameworkDetails.EvidenceLibraryMappedName) {
                window.go.main.App.GetMappedFrameworks()
                    .then(function(records) {
                        displayMappedFrameworkRecords(records);
                       // fetchFrameworkDetails(selectedFramework);
                    })
                    .catch(function(err) {
                        console.error('Error fetching records:', err);
                        alert('Failed to retrieve records.');
                    });
            } else if (!frameworkDetails.AirtableFramework || !frameworkDetails.AirtableView) {
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
                fetchFrameworkFromAirtable(data);
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

    addFrameworkRowEventListeners(records);
}

function addFrameworkRowEventListeners(records) {
    var tableRows = document.querySelectorAll('#recordsContainer tbody tr');
    tableRows.forEach(function(row, index) {
        row.addEventListener('click', function() {
            tableRows.forEach(function(r) {
                r.classList.remove('selected');
            });
            row.classList.add('selected');
            var selectedRecord = records[index];
            displaySelectedFrameworkInfo(selectedRecord);
        });
    }); 
}

function displaySelectedFrameworkInfo(record) {
    window.selectedMissingFramework = record;
    // openFrameworkBuildListModal();
    var label = document.getElementById('selectedRecordLabel');
    label.innerHTML = `
        <strong>Selected Framework</strong>
        <br><strong>Name:</strong> ${record}
        <div id="selectedFrameworkDetails"></div>
    `;
    openFrameworkBuildListSelectionModal();
}

function openFrameworkBuildListSelectionModal() {
    window.go.main.App.GetFrameworkLookup()
        .then(function(records) {
            displayFrameworkBuildListSelection(records, function(selectedFramework) {
                closeFrameworkBuildListModal();
                displaySelectedFrameworkFromBuildListSelection(selectedFramework);
            });
        })
        .catch(function(err) {
            console.error('Error fetching framework build list:', err);
            alert('Failed to retrieve framework build list.');
        });
}

function displayFrameworkBuildListSelection(records, onFrameworkSelected) {
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

    addFrameworkBuildSelectionRowEventListeners(records, onFrameworkSelected);
}

function addFrameworkBuildSelectionRowEventListeners(records, onFrameworkSelected) {
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


function displaySelectedFrameworkFromBuildListSelection(selectedFramework) {
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
        fetchAirtableTablesandViewsUpdate();
    });
}

function fetchAirtableTablesandViewsUpdate() {
    window.go.main.App.GetAirtableBaseTables()
        .then(function(response) {
            displayUpdateFrameworkTablesModal(response.tables);
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

function displayUpdateFrameworkTablesModal(tables) {
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
            handleUpdateTableViewSelector(tableName, viewName);
        });
    });
}


function handleUpdateTableViewSelector(tableName, viewName) {
    window.selectedFrameworkTable = tableName;
    window.selectedFrameworkView = viewName;

    var modal = document.getElementById('frameworkTablesModal');
    modal.style.display = 'none';

    updateSelectFrameworkModalWithTableView();
}


function updateSelectFrameworkModalWithTableView() {
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
        updateSelectedFrameworkLookup();
    });
}

function updateSelectedFrameworkLookup() {
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
    window.go.main.App.UpdateFrameworkLookup(data)
        .then(function (response) {
            alert('Framework Lookup updated successfully.');
            if (window.selectedFrameworkDetails.name != "N/A" && window.selectedFrameworkTable != "" && window.selectedFrameworkView != "") {
                fetchFrameworkFromAirtable(data)
            }
            // Close the Missing Frameworks modal
            var modal = document.getElementById('recordsModal');
            modal.style.display = 'none';
        })
        .catch(function(err) {
            console.error('Error updating Framework Lookup', err);
            alert('Failed to update Framework Lookup.');
        });
}

function fetchFrameworkFromAirtable(data) {
    console.log(data)

    console.log('Fetching framework from Airtable for: ' + selectedFramework);
    // Placeholder for fetching actual data from Airtable
    var modal = document.getElementById('recordsModal');
    modal.style.display = 'none';
}