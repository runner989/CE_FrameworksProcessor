

document.getElementById('getMissingFrameworksButton').addEventListener('click',function() {
    window.go.main.App.GetMissingFramework()
        .then(function(records) {
            // console.log('Records:', records)
            displayMissingRecords(records);
        })
        .catch(function(err) {
            console.error('Error fetching records:', err);
            alert('Failed to retrieve records.');
        });
});