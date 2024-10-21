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
    // Get all frameworks from Framework_Lookup table - if Framework Build list has been imported, it will be every framework from that
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
    content += '<p>This list is the Framework Build list if it was imported.</p>';
    content += '<strong>NOTE: </strong>Not all frameworks in this list exist in Airtable. If you do not see it in the next selection, it is not yet ready.</p>';
    content += '<div id="selectedRecordLabel"></div>';
    // search input box
    content += '<input type="text" id="frameworkSearch" placeholder="Search Framework..." style="margin-bottom: 10px; width: 100%;">';
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

    // event listener for search input
    document.getElementById('frameworkSearch').addEventListener('input', function(e) {
        let searchTerm = e.target.value.toLowerCase();
        let firstMatchIndex = -1;

        tableRows.forEach(function(row, index) {
            let framework = frameworks[index].toLowerCase();
            if (framework.startsWith(searchTerm)) {
                if (firstMatchIndex === -1) {
                    firstMatchIndex = index;
                }
                row.style.display = '';  // Show matching row
            } else {
                row.style.display = 'none'; // Hide non-matching row
            }
        });

        // Scroll to the first matching row
        if (firstMatchIndex !== -1) {
            let firstMatchRow = tableRows[firstMatchIndex];
            firstMatchRow.scrollIntoView({ behavior: 'smooth' , block: "center"});
        }
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
    document.getElementById('loadingNotification').style.display = 'block';
    window.go.main.App.GetFrameworkDetails(selectedFramework)
        .then(function(frameworkDetails) {
            if(!frameworkDetails.EvidenceLibraryMappedName) {
                console.log("Missing EvidenceLibraryMappedName")
                window.go.main.App.GetMappedFrameworks()
                    .then(function(records) {
                        document.getElementById('loadingNotification').style.display = 'none';
                        displayMappedFrameworkRecords(records);
                    })
                    .catch(function(err) {
                        document.getElementById('loadingNotification').style.display = 'none';
                        console.error('Error fetching records:', err);
                        alert('Failed to retrieve records.');
                    });
            } else if (!frameworkDetails.AirtableFramework || !frameworkDetails.AirtableView) {
                console.log("Missing AirtableFramework and/or AirtableView")
                let name = frameworkDetails.CEName || 'N/A';
                let mappedName = frameworkDetails.EvidenceLibraryMappedName || 'N/A';
                let uatStage = frameworkDetails.FrameworkId_UAT || 'N/A';
                let prodNumber = frameworkDetails.FrameworkId_Prod || 'N/A';
                let stageNumber = frameworkDetails.FrameworkId_Staging || 'N/A';
                window.selectedFrameworkDetails = {
                    name: name,
                    ceName: name,
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
                    ceName: frameworkDetails.CEName,
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
            document.getElementById('loadingNotification').style.display = 'none';
            console.error('Error fetching framework details:', err);
            alert('Failed to retrieve framework details.');
        });
}

function displayMappedFrameworkRecords(records) {
    let modal = document.getElementById('recordsModal');
    let recordsContainer = document.getElementById('recordsContainer');
    let selectedFramework = window.selectedFramework
    let content = '<h3>Frameworks From Mapping Table</h3>';
    content += '<p>Framework is missing the Framework name from the Mapping table.</p>'
    content += `<div id="selectedFrameworkLabel"><br><strong>Selected Framework From Framework_Lookup:</strong> ${selectedFramework}\n</div>`;
    content += '<div id="selectedRecordLabel"></div>';
    content += '<input type="text" id="frameworkSearch" placeholder="Search for Framework Name..." style="margin-bottom: 10px; width: 100%;"/>';
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

    // addFrameworkRowEventListeners(records);
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
    // event listener for search input
    document.getElementById('frameworkSearch').addEventListener('input', function(e) {
        let searchTerm = e.target.value.toLowerCase();
        let firstMatchIndex = -1;

        tableRows.forEach(function(row, index) {
            let framework = records[index].toLowerCase();
            if (framework.startsWith(searchTerm)) {
                if (firstMatchIndex === -1) {
                    firstMatchIndex = index;
                }
                row.style.display = '';  // Show matching row
            } else {
                row.style.display = 'none'; // Hide non-matching row
            }
        });

        // Scroll to the first matching row
        if (firstMatchIndex !== -1) {
            let firstMatchRow = tableRows[firstMatchIndex];
            firstMatchRow.scrollIntoView({ behavior: 'smooth' , block: "center"});
        }
    });
}

function displaySelectedFrameworkInfo(record) {
    window.selectedMissingFramework = record;
    let label = document.getElementById('selectedRecordLabel');
    label.innerHTML = `
        <strong>Selected Framework</strong>
        <br><strong>Name:</strong> ${record}
        <div id="selectedFrameworkDetails"></div>
    `;
    openFrameworkBuildListSelectionModal();
}

function openFrameworkBuildListSelectionModal() {
    document.getElementById('loadingNotification').style.display = 'block';
    window.go.main.App.GetFrameworkLookup()
        .then(function(records) {
            displayFrameworkBuildListSelection(records, function(selectedFramework) {
                document.getElementById('loadingNotification').style.display = 'none';
                closeFrameworkBuildListModal();
                displaySelectedFrameworkFromBuildListSelection(selectedFramework);
            });
        })
        .catch(function(err) {
            document.getElementById('loadingNotification').style.display = 'none';
            console.error('Error fetching framework build list:', err);
            alert('Failed to retrieve framework build list.');
        });
}

function closeFrameworkBuildListModal() {
    let modal = document.getElementById('frameworkBuildModal');
    // let container = document.getElementById('frameworkBuildContainer');
    modal.style.display = 'none';
    // container.innerHTML = '';
}

function displayFrameworkBuildListSelection(records, onFrameworkSelected) {
    document.getElementById('loadingNotification').style.display = 'none';
    let modal = document.getElementById('frameworkBuildModal');
    let recordsContainer = document.getElementById('frameworkBuildContainer');
    let selectedFramework = window.selectedFramework

    let content = '<h3>Frameworks Build List</h3>';
    content += `<div id="selectedFrameworkLabel"><br><strong>Selected Framework from Mapping:</strong> ${selectedFramework}\n</div>`;
    content += '<input type="text" id="frameworkSearch2" placeholder="Search for Framework Name..." style="margin-bottom: 10px; width: 100%;"/>';
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

    // addFrameworkBuildSelectionRowEventListeners(records, onFrameworkSelected);
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

    // event listener for search input
    document.getElementById('frameworkSearch2').addEventListener('input', function(e) {
        let searchTerm = e.target.value.toLowerCase();
        let firstMatchIndex = -1;
        tableRows.forEach(function(row, index) {
            let frameworkName = records[index].fields["Name"].toLowerCase();
            if (frameworkName.startsWith(searchTerm)) {
                if (firstMatchIndex === -1) {
                    firstMatchIndex = index;
                }
                row.style.display = '';  // Show matching row
            } else {
                row.style.display = 'none'; // Hide non-matching row
            }
        });

        // Scroll to the first matching row
        if (firstMatchIndex !== -1) {
            let firstMatchRow = tableRows[firstMatchIndex];
            firstMatchRow.scrollIntoView({ behavior: 'smooth' , block: "center"});
        }
    });
}

function displaySelectedFrameworkFromBuildListSelection(selectedFramework) {
    let fields = selectedFramework.fields;
    let name = fields['Name'] || 'N/A';
    let uatStage = fields['UAT_Stage'] || 'N/A';
    let prodNumber = fields['Production Framework Number'] || 'N/A';
    let stageNumber = fields['Stage Framework Number'] || 'N/A';

    let searchBox = document.getElementById('frameworkSearch');
    let searchBox2 = document.getElementById('frameworkSearch2');
    searchBox.style.display = 'none';
    searchBox2.style.display = 'none';
    let detailsDiv = document.getElementById('selectedFrameworkDetails');
    detailsDiv.innerHTML = `
        <strong>Selected Framework from Build List</strong>
        <br><strong>Name:</strong> ${name}
        <br><strong>UAT Stage:</strong> ${uatStage}
        <br><strong>Staging Number:</strong> ${stageNumber}
        <br><strong>Production Framework Number:</strong> ${prodNumber}
        <br><button class="lcars-button" id="selectTableViewButton">Select Framework Table and View</button>
    `;

    window.selectedFrameworkDetails = {
        name: name,
        uatStage: uatStage,
        stageNumber: stageNumber,
        prodNumber: prodNumber
    }

    document.getElementById('selectTableViewButton').addEventListener('click', function() {
        document.getElementById('loadingNotification').style.display = 'block';
        fetchAirtableTablesandViewsUpdate();
    });
}

function fetchAirtableTablesandViewsUpdate() {
    Promise.all([
        window.go.main.App.GetAirtableBaseTables(),
        window.go.main.App.GetAvailableAirtableBases()
    ])
    // window.go.main.App.GetAirtableBaseTables()
        .then(function([tablesResponse, basesResponse]) {
            document.getElementById('loadingNotification').style.display = 'none';
            displayUpdateFrameworkTablesModal(tablesResponse.tables, basesResponse);
        })
        .catch(function(err) {
            document.getElementById('loadingNotification').style.display = 'none';
            console.error('Error fetching Airtable tables and views.', err);
            alert('Failed to fetch Airtable tables and views.');
        })
}

function closeFrameworkTablesModalX() {
    let modal = document.getElementById('frameworkTablesModal');
    let container = document.getElementById('frameworkTablesContainer');
    modal.style.display = 'none';
    container.innerHTML = '';
}

function displayUpdateFrameworkTablesModal(tables, bases) {
    document.getElementById('closeFrameworkTablesModal').addEventListener('click', function () {
        closeFrameworkTablesModalX();
    });

    let modal = document.getElementById('frameworkTablesModal');
    let container = document.getElementById('frameworkTablesContainer');

    let content = '<h4>Select a Framework Table and View</h4>';
    content += '<div id="selectedFramework"><p>Looking for Framework: ' + window.selectedMissingFramework + '</p></div>';

    content += '<label for="baseSelect">Select Base: </label> ';
    content += '<select id="baseSelect">';
    bases.forEach(function (base) {
        content += `<option value="${base.id}">${base.name}</option>`;
    });
    content += '</select><hr>';

    // search input for filtering the table
    content += '<input type="text" id="tableSearch" placeholder="Search Tables..." style="margin-bottom: 10px; width: 100%;">';

    content += '<div id="tableContainer">';
    content += '<ul id="tablesList"></ul></div>';

    container.innerHTML = content
    modal.style.display = 'block';

    let initialBase = bases.find(base => base.id === 'app5fTueYfRM65SzX') || bases[0];
    let initialBaseID = initialBase.id;

    document.getElementById('baseSelect').value = initialBaseID;

    fetchTablesForBase(initialBaseID)

    document.getElementById('baseSelect').addEventListener('change', function () {
        let selectedBaseID = this.value;
        document.getElementById('loadingNotification').style.display = 'block';
        fetchTablesForBase(selectedBaseID);
    });
}

function fetchTablesForBase(baseID) {
    window.go.main.App.GetAirtableTables(baseID)
        .then(function (response) {
            document.getElementById('loadingNotification').style.display = 'none';
            let tablesList = document.getElementById('tablesList');
            let tableSearch = document.getElementById('tableSearch');
            let content = '';

            response.tables.forEach(function (table, index) {
                content += `<li class="table-item" data-index="${index}">
                <span class="table-name">${table.name}</span>
                <span class="table-id" style="display: none;">${table.id}</span>
                <ul class="views-list" id="views-${index}" style="display: none;">`;

                table.views.forEach(function (view) {
                    content += `<li class="view-item" data-base-id="${baseID}" data-table-id="${table.id}" data-table-name="${table.name}" data-view-name="${view.name}">
                    ${view.name}</li>`;
                });
                content += `</ul></li>`;
            });

            tablesList.innerHTML = content;

            let tableItems = document.querySelectorAll('.table-item');
            tableItems.forEach(function (item) {
                let index = item.getAttribute('data-index');
                let viewsList = document.getElementById(`views-${index}`);

                item.querySelector('.table-name').addEventListener('click', function () {
                    if (viewsList.style.display === 'none') {
                        viewsList.style.display = 'block';
                    } else {
                        viewsList.style.display = 'none';
                    }
                });
            });

            let viewItems = document.querySelectorAll('.view-item');
            viewItems.forEach(function (item) {
                item.addEventListener('click', function () {
                    let tableName = item.getAttribute('data-table-name');
                    let viewName = item.getAttribute('data-view-name');
                    let tableID = item.getAttribute('data-table-id');
                    handleUpdateTableViewSelector(baseID, tableName, viewName, tableID);
                });
            });

            // Filter tables based on search input
            tableSearch.addEventListener('input', function(e) {
                let searchTerm = e.target.value.toLowerCase();

                tableItems.forEach(function(item) {
                    let tableName = item.querySelector('.table-name').innerText.toLowerCase();

                    if (tableName.startsWith(searchTerm)) {
                        item.style.display = ''; // Show matching table
                    } else {
                        item.style.display = 'none'; // Hide non-matching table
                    }
                });
            });
        })
        .catch(function (err) {
            console.error('Error fetching Airtable tables:', err);
            alert('Failed to fetch tables for the selected base.');
        });
}


function handleUpdateTableViewSelector(baseID, tableName, viewName, tableID) {
    window.selectedFrameworkTable = tableName;
    window.selectedFrameworkView = viewName;
    window.selectedFrameworkTableID = tableID;
    window.selectedFrameworkBase = baseID;

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
    <br><button class="lcars-button" id="okButton">OK</button>
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
        mappedName: window.selectedMissingFramework,
        ceName: window.selectedFrameworkDetails.name,
        uatStage: window.selectedFrameworkDetails.uatStage,
        prodNumber: window.selectedFrameworkDetails.prodNumber,
        stageNumber: window.selectedFrameworkDetails.stageNumber,
        tableID: window.selectedFrameworkTableID,
        tableName: window.selectedFrameworkTable,
        tableView: window.selectedFrameworkView,
        baseID: window.selectedFrameworkBase,
    };
    window.go.main.App.UpdateFrameworkLookup(data)
        .then(function (result) {
            alert('Framework Lookup updated successfully.');
            // Close the Missing Frameworks modal
            let modal = document.getElementById('recordsModal');
            let recordsContainer = document.getElementById('recordsContainer');
            recordsContainer.innerHTML = "";
            modal.style.display = 'none';
            if (window.selectedFrameworkDetails.name !== "N/A" && window.selectedFrameworkTable !== "" && window.selectedFrameworkView !== "") {
                console.log("running fetchFrameworkFromAirtable")
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
            document.getElementById('loadingNotification').style.display = 'none';
            alert('Framework updated successfully.');
            let modal2 = document.getElementById('updateFrameworkModal')
            modal2.style.display = 'none';
            let modal = document.getElementById('recordsModal');
            let recordsContainer = document.getElementById('recordsContainer');
            recordsContainer.innerHTML = '';
            modal.style.display = 'none';
        })
        .catch(function(err) {
            document.getElementById('loadingNotification').style.display = 'none';
            console.error('Error updating Framework table', err);
            alert('Failed to update Framework table.');
        })

}

document.getElementById('updateAllFrameworksButton').addEventListener('click',function() {
    let modal = document.getElementById('recordsModal');
    let recordsContainer = document.getElementById('recordsContainer');
    recordsContainer.innerHTML = "";
    window.go.main.App.UpdateAllFrameworks()
        .then(function(message) {
            document.getElementById('loadingNotification').style.display = 'none';
            alert(message);
            recordsContainer.innerHTML = "";
            modal.style.display = 'none';
        })
        .catch(function(err) {
            document.getElementById('loadingNotification').style.display = 'none';
            console.error('Error fetching records:', err);
            alert('Failed to retrieve records.');
            recordsContainer.innerHTML = "";
            modal.style.display = 'none';
        });
});