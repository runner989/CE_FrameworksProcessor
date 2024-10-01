// Define the specific fields to display in the desired order
var orderedFields = [
    'Name',
    'UAT_Stage',
    'Stage Framework Number',
    'Production Framework Number',
    'Notes',
    'Category',
    'Status',
];

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
    let modal = document.getElementById('recordsModal');
    let recordsContainer = document.getElementById('recordsContainer');


    let content = '<h3>Frameworks from Lookup Table</h3>';
    content += '<p>This list is the Framework Build list that was imported.</p>';
    content += '<strong>NOTE: </strong>Not all frameworks in this list exist in Airtable. If you do not see it in the next selection, it is not yet ready.</p>';
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

    let tableRows = document.querySelectorAll('#recordsContainer tbody tr');
    tableRows.forEach(function(row, index) {
        row.addEventListener('click', function() {
            tableRows.forEach(function(r) {
                r.classList.remove('selected');
            });
            row.classList.add('selected');
            let selectedRecord = frameworks[index];
            fetchFrameworkDetails(selectedRecord);
        });
    });
}

function fetchFrameworkDetails(selectedFramework) {
    window.selectedFramework = selectedFramework;
    let label = document.getElementById('selectedRecordLabel');
    label.innerHTML = `
        <strong>Selected Framework</strong>
        <br><strong>Name:</strong> ${selectedFramework}
        <div id="selectedFrameworkDetails"></div>
    `;
    window.go.main.App.GetFrameworkDetails(selectedFramework)
        .then(function(frameworkDetails) {
            if(!frameworkDetails.EvidenceLibraryMappedName) {
                // console.log("Missing EvidenceLibraryMappedName")
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
                // console.log("Missing AirtableFramework and/or AirtableView")
                let name = frameworkDetails.CEName || 'N/A';
                let mappedName = frameworkDetails.EvidenceLibraryMappedName || 'N/A';
                let uatStage = frameworkDetails.FrameworkId_UAT || 'N/A';
                let prodNumber = frameworkDetails.FrameworkId_Prod || 'N/A';
                let stageNumber = frameworkDetails.FrameworkId_Staging || 'N/A';
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
                let data = {
                    name: frameworkDetails.CEName,
                    mappedName: frameworkDetails.EvidenceLibraryMappedName,
                    tableID: frameworkDetails.AirtableTableID,
                    tableName: frameworkDetails.AirtableFramework,
                    tableView: frameworkDetails.AirtableView,
                    uatStage: frameworkDetails.FrameworkId_UAT,
                    stageNumber: frameworkDetails.FrameworkId_Staging,
                    prodNumber: frameworkDetails.FrameworkId_Prod,
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
    let modal = document.getElementById('recordsModal');
    let recordsContainer = document.getElementById('recordsContainer');


    let content = '<h3>Frameworks From Mapping Table</h3>';
    content += '<p>Framework is missing the Framework name from the Mapping table.</p>'
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
    let tableRows = document.querySelectorAll('#recordsContainer tbody tr');
    tableRows.forEach(function(row, index) {
        row.addEventListener('click', function() {
            tableRows.forEach(function(r) {
                r.classList.remove('selected');
            });
            row.classList.add('selected');
            let selectedRecord = records[index];
            displaySelectedFrameworkInfo(selectedRecord);
        });
    }); 
}

function displaySelectedFrameworkInfo(record) {
    window.selectedMissingFramework = record;
    // openFrameworkBuildListModal();
    let label = document.getElementById('selectedRecordLabel');
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

function closeFrameworkBuildListModal() {
    let modal = document.getElementById('frameworkBuildModal');
    modal.style.display = 'none';
}

function displayFrameworkBuildListSelection(records, onFrameworkSelected) {
    let modal = document.getElementById('frameworkBuildModal');
    let recordsContainer = document.getElementById('frameworkBuildContainer');

    let content = '<h3>Frameworks Build List</h3>';
    content += '<div id="tableContainer"><table><thead><tr>';

    orderedFields.forEach(function(field){
        content += '<th>' + field + '</th>';
    });
    content += '</tr></thead><tbody>';

    records.forEach(function(record, index) {
        let fields = record.fields;
        content += '<tr data-index="' + index + '">';
        orderedFields.forEach(function(field) {
            let value = fields[field];
            if (Array.isArray(value)) {
                value = value.join(', ');
            } else if (typeof value == 'object' && value !== null) {
                value = JSON.stringify(value);
            } else if (value === undefined || value == null) {
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
    let tableRows = document.querySelectorAll('#frameworkBuildContainer tbody tr');
    tableRows.forEach(function(row, index) {
        row.addEventListener('click', function() {
            tableRows.forEach(function(r) {
                r.classList.remove('selected');
            });
            row.classList.add('selected');
            let selectedRecord = records[index];
            onFrameworkSelected(selectedRecord);
        });
    });
}


function displaySelectedFrameworkFromBuildListSelection(selectedFramework) {
    let fields = selectedFramework.fields;
    let name = fields['Name'] || 'N/A';
    let uatStage = fields['UAT_Stage'] || 'N/A';
    let prodNumber = fields['Production Framework Number'] || 'N/A';
    let stageNumber = fields['Stage Framework Number'] || 'N/A';

    let detailsDiv = document.getElementById('selectedFrameworkDetails');
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
    let modal = document.getElementById('frameworkTablesModal');
    modal.style.display = 'none';
}

function displayUpdateFrameworkTablesModal(tables) {
    document.getElementById('closeFrameworkTablesModal').addEventListener('click',function() {
        closeFrameworkTablesModalX();
    });

    let modal = document.getElementById('frameworkTablesModal');
    let container = document.getElementById('frameworkTablesContainer');

    let content = '<h4>Select a Framework Table and View</h4><ul>';

    tables.forEach(function (table, index) {
        content += `<li class="table-item" data-index="${index}">
            <span class="table-name">${table.name}</span>
            <span class="table-id" style="display: none;">${table.id}</span>
            <ul class="views-list" id="views-${index}" style="display: none;">`;

        table.views.forEach(function(view) {
            content += `<li class="view-item" data-table-id="${table.id}" data-table-name="${table.name}" data-view-name="${view.name}">
                ${view.name}
            </li>`;
        });
        content += `</ul></li>`;
    });

    content += '</ul>';

    container.innerHTML = content;
    modal.style.display = 'block';

    let tableItems = container.querySelectorAll('.table-item');
    tableItems.forEach(function(item) {
        let index = item.getAttribute('data-index');
        let viewsList = document.getElementById(`views-${index}`);

        item.querySelector('.table-name').addEventListener('click', function() {
            if (viewsList.style.display === 'none') {
                viewsList.style.display = 'block';
            } else {
                viewsList.style.display = 'none';
            }
        });
    });

    let viewItems = container.querySelectorAll('.view-item');
    viewItems.forEach(function(item) {
        item.addEventListener('click', function() {
            let tableName = item.getAttribute('data-table-name');
            let viewName = item.getAttribute('data-view-name');
            let tableID = item.getAttribute('data-table-id');
            handleUpdateTableViewSelector(tableName, viewName, tableID);
        });
    });
}


function handleUpdateTableViewSelector(tableName, viewName, tableID) {
    window.selectedFrameworkTable = tableName;
    window.selectedFrameworkView = viewName;
    window.selectedFrameworkTableID = tableID;

    let modal = document.getElementById('frameworkTablesModal');
    modal.style.display = 'none';

    updateSelectFrameworkModalWithTableView();
}


function updateSelectFrameworkModalWithTableView() {
    let detailsDiv = document.getElementById('selectedFrameworkDetails');

    let tableViewInfo = `
    <br><strong>Selected Framework Table and View</strong>
    <br><strong>Table Name:</strong> ${window.selectedFrameworkTable}
    <br><strong>View Name:</strong> ${window.selectedFrameworkView}
    <span style="display: none;">${window.selectedFrameworkTableID}</span>
    <br><button id="okButton">OK</button>
    `;

    let existingInfo = detailsDiv.querySelector('.table-view-info');

    if (existingInfo) {
        existingInfo.innerHTML = tableViewInfo;
    } else {
        let div = document.createElement('div');
        div.classList.add('table-view-info');
        div.innerHTML = tableViewInfo;
        detailsDiv.appendChild(div);
    }

    document.getElementById('okButton').addEventListener('click', function() {
        updateSelectedFrameworkLookup();
    });
}

function updateSelectedFrameworkLookup() {
    let data = {
        missingFrameworkName: window.selectedMissingFramework,
        cename: window.selectedFrameworkDetails.name,
        uatStage: window.selectedFrameworkDetails.uatStage,
        prodNumber: window.selectedFrameworkDetails.prodNumber,
        stageNumber: window.selectedFrameworkDetails.stageNumber,
        tableID: window.selectedFrameworkTableID,
        tableName: window.selectedFrameworkTable,
        tableView: window.selectedFrameworkView,
    };
    window.go.main.App.UpdateFrameworkLookup(data)
        .then(function (result) {
            alert('Framework Lookup updated successfully.');
            if (window.selectedFrameworkDetails.name !== "N/A" && window.selectedFrameworkTable !== "" && window.selectedFrameworkView !== "") {
                fetchFrameworkFromAirtable(data)
            }
            // // Close the Missing Frameworks modal
            // let modal = document.getElementById('recordsModal');
            // modal.style.display = 'none';
        })
        .catch(function(err) {
            console.error('Error updating Framework Lookup', err);
            alert('Failed to update Framework Lookup.');
        });
}

function fetchFrameworkFromAirtable(data) {
    // Display the modal and update the framework name
    let modal = document.getElementById('updateFrameworkModal');
    let frameworkNameElement = document.getElementById('frameworkName');
    frameworkNameElement.innerHTML = `Working on updating ${data.tableName}`;
    modal.style.display = 'block'; // Show the modal

    window.go.main.App.GetFrameworkRecords(data)
        .then(function (result) {
            alert('Framework updated successfully.');
            let modal2 = document.getElementById('updateFrameworkModal')
            modal2.style.display = 'none';
            let modal = document.getElementById('recordsModal');
            modal.style.display = 'none';
        })
        .catch(function(err) {
            console.error('Error updating Framework table', err);
            alert('Failed to update Framework table.');
        })

}

document.getElementById('updateAllFrameworksButton').addEventListener('click',function() {
    var modal = document.getElementById('recordsModal');
    let recordsContainer = document.getElementById('recordsContainer');
    window.go.main.App.UpdateAllFrameworks()
        .then(function(message) {
            alert(message);
            recordsContainer.innerHTML = "";
            modal.style.display = 'none';
        })
        .catch(function(err) {
            console.error('Error fetching records:', err);
            alert('Failed to retrieve records.');
            recordsContainer.innerHTML = "";
            modal.style.display = 'none';
        });
});