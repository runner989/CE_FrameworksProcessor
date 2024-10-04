document.getElementById('exportStagingEvidenceMapButton').addEventListener('click',function() {
    let modal = document.getElementById('loadingNotification')
    modal.innerHTML = '<div class="alert alert-info" role="alert" style="padding: 20px; font-size: 18px;">Exporting Staging Evidence Mapping Report...</div>';
    modal.style.display = 'block';
    window.go.main.App.ExportEvidenceMapReport("Staging")
        .then(function() {
            modal.style.display = 'none';
            alert('Finished exporting staging evidence mapping report')
        })
        .catch(function(err) {
            modal.style.display = 'none';
            console.error('Error exporting staging evidence mapping report:', err);
            alert('Failed to export staging evidence mapping report.');
        });
});

document.getElementById('exportProdEvidenceMapButton').addEventListener('click',function() {
    let modal = document.getElementById('loadingNotification')
    modal.innerHTML = '<div class="alert alert-info" role="alert" style="padding: 20px; font-size: 18px;">Exporting Prod Evidence Mapping Report...</div>';
    modal.style.display = 'block';
    window.go.main.App.ExportEvidenceMapReport("Prod")
        .then(function() {
            modal.style.display = 'none';
            alert('Finished exporting prod evidence mapping report')
        })
        .catch(function(err) {
            modal.style.display = 'none';
            console.error('Error exporting prod evidence mapping report:', err);
            alert('Failed to export staging evidence mapping report.');
        });
});

document.getElementById('exportUATEvidenceMapReport').addEventListener('click',function() {
    let modal = document.getElementById('loadingNotification')
    // Update modal content first, then show it
    modal.innerHTML = '<div class="alert alert-info" role="alert" style="padding: 20px; font-size: 18px;">Exporting UAT Evidence Mapping Report...</div>';
    modal.style.display = 'block';

    window.go.main.App.ExportEvidenceMapReport("UAT")
        .then(function() {
            modal.style.display = 'none';
            alert('Finished exporting UAT Evidence Mapping Report');
        })
        .catch(function(err) {
            modal.style.display = 'none';
            console.error('Error exporting UAT Evidence Mapping Report:', err);
            alert('Failed to export UAT evidence mapping report.');
        });
});

document.getElementById('reviewStagingDeleted').addEventListener('click',function() {
    window.go.main.App.GetDeletionsList("Staging")
        .then(function(results) {
            displayDeletedList(results,"Staging")
        })
        .catch(function(err) {
            console.error('Error fetching Staging Deleted list');
        })
})
document.getElementById('reviewProdDeleted').addEventListener('click',function() {
    window.go.main.App.GetDeletionsList("Prod")
        .then(function(results) {
            displayDeletedList(results, "Prod")
        })
        .catch(function(err) {
            console.error('Error fetching Prod Deleted list');
        })
})

// Example function to display deletions in a table format
function displayDeletedList(deletions, table) {
    let modal = document.getElementById('deletionsModal');  // Assuming you have a modal for deletions
    let deletionsContainer = document.getElementById('deletionsContainer'); // Table container for deletions

    let content = '<div><h2>Deletions List for: ' + table + '</h2></div>';
    content += '<table><thead><tr>';
    content += '<th>EvidenceID</th><th>Framework</th><th>FrameworkID</th><th>Requirement</th><th>Delete</th>';
    // content += '<th>Description</th><th>Guidance</th><th>RequirementType</th><th>Delete</th>';
    content += '</tr></thead><tbody>';

    deletions.forEach(function(deletion) {
        content += `<tr>
      <td>${deletion.EvidenceID}</td>
      <td>${deletion.Framework}</td>
      <td>${deletion.FrameworkID}</td>
      <td>${deletion.Requirement}</td>
      <td>${deletion.Delete}</td>
    </tr>`;
    });

    content += '</tbody></table>';
    deletionsContainer.innerHTML = content;
    modal.style.display = 'block';  // Display modal with deletions
}

document.getElementById('closeDeletionsModal').addEventListener('click', function(){
    let modal = document.getElementById('deletionsModal');
    modal.style.display = 'none';
})