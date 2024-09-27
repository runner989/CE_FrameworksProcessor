
// Function to load and display the table
window.loadDeleteFrameworkLookupTable = function() {
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
    const container = document.getElementById('deleteFrameworkTableContainer');
    let tableHtml = '<table><thead><tr><th>CE Framework</th><th>Mapped Name</th><th>UAT #</th><th>Staging #</th><th>Prod #</th></tr></thead><tbody>';

    records.forEach((record, index) => {
        tableHtml += `<tr>
            <td><span style="display: none">${record.RowID}</span>  ${record.CEFramework}</td>
            <td>${record.EvidenceLibraryMappedName}</td>
            <td>${record.FrameworkId_UAT}</td>
            <td>${record.FrameworkId_Staging}</td>
            <td>${record.FrameworkId_Prod}</td>
            <td><button class="delete-button" data-index="${index}">Delete</button></td>
        </tr>`;
    });

    tableHtml += '</tbody></table>';
    container.innerHTML = tableHtml;

    // Add event listeners to each edit button
    document.querySelectorAll('.delete-button').forEach(button => {
        button.addEventListener('click', function () {
            const index = this.getAttribute('data-index');
            deleteRecord(records[index]);

            // // Scroll to the form when edit is clicked
            // document.getElementById('deleteFormContainer').scrollIntoView({ behavior: 'smooth' });
        });
    });
}

function deleteRecord(record) {
    let recordDetails = {
        rowID: record.RowID,
        airtableBase: record.AirtableBase,
        airtableTableID: record.AirtableTableID,
        airtableFramework: record.AirtableFramework,
        airtableView: record.AirtableView,
        mappedName: record.EvidenceLibraryMappedName,
        ceFramework: record.CEFramework,
        frameworkId_UAT: record.FrameworkId_UAT,
        frameworkId_Staging: record.FrameworkId_Staging,
        frameworkId_Prod: record.FrameworkId_Prod,
        version: record.Version,
        description: record.Description,
        comments: record.Comments
    }
    window.go.main.App.DeleteSelectedFramework(recordDetails)
        .then(function (result) {
            alert(`${recordDetails.ceFramework} deleted!`);
            loadDeleteFrameworkLookupTable();
        })
        .catch(function (err) {
            console.log(err);
        })
    // alert(`The ${recordDetails.ceFramework} framework will be deleted as soon as the code is finished!`)
}