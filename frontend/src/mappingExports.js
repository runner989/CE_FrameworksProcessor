document.getElementById('exportStagingEvidenceMapButton').addEventListener('click',function() {
    window.go.main.App.ExportEvidenceMapReport("Staging")
        .then(function(records) {
            displayMissingRecords(records);
        })
        .catch(function(err) {
            console.error('Error fetching records:', err);
            alert('Failed to retrieve records.');
        });
});
