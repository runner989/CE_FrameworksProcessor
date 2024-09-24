
window.runtime.EventsOn("progress", (message) => {
    updateProgressModal(message)
});

function updateProgressModal(message) {
    var modal = document.getElementById('recordsModal');
    var recordsContainer = document.getElementById('recordsContainer');

    var p = document.createElement('p');
    p.textContent = message;

    recordsContainer.appendChild(p);

    while (recordsContainer.childNodes.length > 15) {
        recordsContainer.removeChild(recordsContainer.firstChild);
    }
    recordsContainer.scrollTop = recordsContainer.scrollHeight;

    if (modal.style.display !== 'block') {
        modal.style.display = 'block';
    }
}

document.getElementById('updateEvidenceButton').addEventListener('click',function() {
    var recordsContainer = document.getElementById('recordsContainer');
    recordsContainer.innerHTML = '';
    
    window.go.main.App.ReadAPIEvidenceTable()
        .then(function(message) {
            // console.log('Records:', records)
            displayEvidenceUpdated(message);
        })
        .catch(function(err) {
            console.error('Error fetching records:', err);
            alert('Failed to retrieve records.');
        });
});

function displayEvidenceUpdated(message) {
    var modal = document.getElementById('recordsModal');
    var recordsContainer = document.getElementById('recordsContainer');

    var content = '<h3>Evidence and Mapping Tables updated!</h3>';
    recordsContainer.innerHTML = content;
    modal.style.display = 'block';
}