var orderedFields = [
    'name',
    'id',
];

document.getElementById('airtableBasesButton').addEventListener('click', function() {
    window.go.main.App.GetAvailableAirtableBases()
        .then(function(records) {
            // console.log('Records:', records)
            displayBasesRecords(records);
        })
        .catch(function(err) {
            console.error('Error fetching records:', err);
            alert('Failed to retrieve records.');
        });
});

// Function to display records in the modal
function displayBasesRecords(records) {
    var modal = document.getElementById('recordsModal');
    var recordsContainer = document.getElementById('recordsContainer');

    var content = '<h2>Available Airtable Bases for this API Key</h2>';
    content += '<div><button id="updateAirtableBaseTableButton">Update Airtable Base Table</button></div>';
    content += '<div id="selectedRecordLabel"></div>';
    content += '<div id="tableContainer"><table><thead><tr>';

    orderedFields.forEach(function(field){
        content += '<th>' + field + '</th>';
    });
    content += '</tr></thead><tbody>';

    records.forEach(function(record, index) {
        content += '<tr data-index="' + index + '">';
        orderedFields.forEach(function(field) {
            var value = record[field];
            if (value === undefined || value === null) {
                value = '';
            }
            content += '<td>' + value + '</td>';
        });
        content += '</tr>';
    });
    content += '</tbody></table></div>';
    recordsContainer.innerHTML = content;
    modal.style.display = 'block';

    document.getElementById('updateAirtableBaseTableButton').addEventListener('click', function(){
        updateAirtableBaseTable(records);
    });
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

function updateAirtableBaseTable(records) {
    window.go.main.App.UpdateAirtableBasesTable(records)
        .then(function (result) {
            alert('Framework Lookup updated successfully.');

            // Close the Missing Frameworks modal
            let modal = document.getElementById('recordsModal');
            modal.style.display = 'none';
        })
        .catch(function(err) {
            console.error('Error updating Framework Lookup', err);
            alert('Failed to update Framework Lookup.');
        });
}