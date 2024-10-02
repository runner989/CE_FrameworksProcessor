document.getElementById('exportStagingEvidenceMapButton').addEventListener('click',function() {
    let modal = document.getElementById('loadingNotification')
    modal.style.display = 'block';
    modal.innerHTML = '<div class="alert alert-info" role="alert"><span class="sr-only">Exporting Staging Evidence Mapping Report</span></div>';
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
    modal.style.display = 'block';
    modal.innerHTML = '<div class="alert alert-info" role="alert"><span class="sr-only">Exporting Prod Evidence Mapping Report</span></div>';
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
