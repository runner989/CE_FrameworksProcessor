document.addEventListener("DOMContentLoaded", function() {
    loadFrameworkLookupTable()
})

// document.getElementById("getFrameworkLookupButton").addEventListener('click', function () {
//     // Step 1: Fetch Framework_Lookup table records from Go
//     loadFrameworkLookupTable();
// });

// Function to load and display the table
function loadFrameworkLookupTable() {
    window.go.main.App.GetFrameworkLookupTable()
        .then(function (records) {
            displayFrameworkTable(records);
        })
        .catch(function (err) {
            console.error('Error fetching Framework Lookup Table:', err);
        });
}

// Step 2: Display the table
function displayFrameworkTable(records) {
    const container = document.getElementById('frameworkTableContainer');
    let tableHtml = '<table><thead><tr><th>CE Framework</th><th>UAT #</th><th>Staging #</th><th>Prod #</th></tr></thead><tbody>';

    records.forEach((record, index) => {
        tableHtml += `<tr>
            <td><span style="display: none">${record.RowID}</span>  ${record.CEFramework}</td>
            <td>${record.FrameworkId_UAT}</td>
            <td>${record.FrameworkId_Staging}</td>
            <td>${record.FrameworkId_Prod}</td>
            <td><button class="edit-button" data-index="${index}">Edit</button></td>
        </tr>`;
    });

    tableHtml += '</tbody></table>';
    container.innerHTML = tableHtml;

    // Add event listeners to each edit button
    document.querySelectorAll('.edit-button').forEach(button => {
        button.addEventListener('click', function () {
            const index = this.getAttribute('data-index');
            editRecord(records[index]);

            // Scroll to the form when edit is clicked
            document.getElementById('editFormContainer').scrollIntoView({ behavior: 'smooth' });
        });
    });
}

// Step 3: Edit record
function editRecord(record) {
    document.getElementById('editFormContainer').style.display = 'block';

    // Populate the form fields with the selected record
    document.getElementById('rowID').value = record.RowID;
    document.getElementById('airtableBase').value = record.AirtableBase;
    document.getElementById('airtableTableID').value = record.AirtableTableID;
    document.getElementById('airtableFramework').value = record.AirtableFramework;
    document.getElementById('airtableView').value = record.AirtableView;
    document.getElementById('evidenceLibraryMappedName').value = record.EvidenceLibraryMappedName;
    document.getElementById('ceFramework').value = record.CEFramework;
    document.getElementById('frameworkId_UAT').value = record.FrameworkId_UAT;
    document.getElementById('frameworkId_Staging').value = record.FrameworkId_Staging;
    document.getElementById('frameworkId_Prod').value = record.FrameworkId_Prod;
    document.getElementById('version').value = record.Version;
    document.getElementById('description').value = record.Description;
    document.getElementById('comments').value = record.Comments;
}

// Step 4: Save changes
document.getElementById('saveButton').addEventListener('click', function () {
    const updatedRecord = {
        rowID: document.getElementById('rowID').value,
        airtableBase: document.getElementById('airtableBase').value,
        airtableTableID: document.getElementById('airtableTableID').value,
        airtableFramework: document.getElementById('airtableFramework').value,
        airtableView: document.getElementById('airtableView').value,
        evidenceLibraryMappedName: document.getElementById('evidenceLibraryMappedName').value,
        ceFramework: document.getElementById('ceFramework').value,
        frameworkId_UAT: parseInt(document.getElementById('frameworkId_UAT').value, 10),
        frameworkId_Staging: parseInt(document.getElementById('frameworkId_Staging').value, 10),
        frameworkId_Prod: parseInt(document.getElementById('frameworkId_Prod').value, 10),
        version: parseInt(document.getElementById('version').value, 10),
        description: document.getElementById('description').value,
        comments: document.getElementById('comments').value,
    };

    window.go.main.App.UpdateFrameworkLookupRecord(updatedRecord)
        .then(function () {
            alert('Record updated successfully!');
            document.getElementById('editFormContainer').style.display = 'none';
            // Optionally, refresh the table or reload the page
            // Scroll back to the top after saving
            window.scrollTo({ top: 0, behavior: 'smooth' });
            // Reload the table to reflect updated data
            loadFrameworkLookupTable();
        })
        .catch(function (err) {
            console.error('Error updating record:', err);
            alert('Failed to update the record.');
        });
});

// Step 5: Cancel editing
document.getElementById('cancelButton').addEventListener('click', function () {
    document.getElementById('editFormContainer').style.display = 'none';
    // Scroll back to the top after saving
    window.scrollTo({ top: 0, behavior: 'smooth' });
});

// Step 6: Navigate back to the main page
document.getElementById('backButton').addEventListener('click', function () {
    window.location.href = 'index.html'; // Assuming 'index.html' is the main page
});

