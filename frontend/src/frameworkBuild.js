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

document.getElementById('getFrameworkBuildListButton').addEventListener('click',function() {
    window.go.main.App.GetFrameworkLookup()
        .then(function(records) {
            // console.log('Records:', records)
            displayRecords(records);
        })
        .catch(function(err) {
            console.error('Error fetching records:', err);
            alert('Failed to retrieve records.');
        });
});

// Function to display records in the modal
function displayRecords(records) {
    var modal = document.getElementById('recordsModal');
    var recordsContainer = document.getElementById('recordsContainer');


    var content = '<h2>Frameworks Build List</h2>';
    content += '<div><button class="lcars-button" id="updateFrameworkLookupTableButton">Update Framework Lookup Table</button></div>';
    content += '<div id="selectedRecordLabel"></div>';
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
            // Handle different data types
            if (Array.isArray(value)) {
                value = value.join(', ');
            } else if (typeof value == 'object' && value !== null) {
                value = JSON.stringify(value);
            } else if (value == undefined || value === null) {
                value = '';
            }
            content += '<td>' + value + '</td>';
        });
        content += '</tr>';
    });
    content += '</tbody></table></div>';
    recordsContainer.innerHTML = content;
    modal.style.display = 'block';

    document.getElementById('updateFrameworkLookupTableButton').addEventListener('click', function(){
        updateFrameworkLookupTable(records);
    });

    addRowEventListeners(records);
}

// Close modal when the 'x' is clicked
document.getElementById('closeModal').onclick = function() {
    var modal = document.getElementById('recordsModal');
    modal.style.display = 'none';
};

// Close modal when clicking outside of the modal content
window.onclick = function(event) {
    var modal = document.getElementById('recordsModal');
    if (event.target == modal) {
        modal.style.display = 'none';
    }
};

function addRowEventListeners(records) {
    var tableRows = document.querySelectorAll('#recordsContainer tbody tr');
    tableRows.forEach(function(row, index) {
        row.addEventListener('click', function() {
            tableRows.forEach(function(r) {
                r.classList.remove('selected');
            });
            row.classList.add('selected');
            var selectedRecord = records[index];
            displaySelectedRecordDetails(selectedRecord);
        });
    });
}

function displaySelectedRecordDetails(record) {
    var fields = record.fields;
    var name = fields['Name'] || 'N/A';
    var uatStage = fields['UAT_Stage'] || 'N/A';
    var prodNumber = fields['Production Framework Number'] || 'N/A';
    var stageNumber = fields['Stage Framework Number'] || 'N/A';

    var label = document.getElementById('selectedRecordLabel');
    label.innerHTML = `
        <strong>Selected Framework</strong>
        <br><strong>Name:</strong> ${name}
        <br><strong>UAT Stage:</strong> ${uatStage}
        <br><strong>Production Number:</strong> ${prodNumber}
        <br><strong>Staging Number:</strong> ${stageNumber}                
    `;
}

function updateFrameworkLookupTable(records) {
    var data = records.map(function(record) {
        return record.fields;
    });
    window.go.main.App.UpdateBuildFrameworkLookupTable(data)
        .then(function (response) {
            alert('Framework Lookup table updated successfully.');
            let modal = document.getElementById('recordsModal');
            modal.style.display = 'none';
        })
        .catch(function(err) {
            console.error('Error updating Framework Lookup table:', err);
            alert('Failed to update Framework Lookup table.');
        })
}

// Close modal when the 'x' is clicked
document.getElementById('closeFrameworkBuildModal').onclick = function() {
    var modal = document.getElementById('frameworkBuildModal');
    modal.style.display = 'none';
};

