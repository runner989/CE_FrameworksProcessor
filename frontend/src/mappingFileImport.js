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

document.getElementById('importStagingButton').addEventListener('click', function() {
    let modal = document.getElementById('recordsModal');
    // Use Wails' OpenFileDialog to select a file
    window.go.main.App.OpenFileDialog()
        .then(filePath => {
            if (!filePath) {
                alert("No file selected.");
                return;
            }

            // Send the file path to the Go backend
            window.go.main.App.ProcessEvidenceStagingFile(filePath)
                .then(() => {
                    modal.style.display = 'none';
                    alert("File processed successfully.");
                })
                .catch(err => {
                    console.error("Error processing file:", err);
                    alert("Failed to process the file.");
                });
        })
        .catch(err => {
            console.error("Error selecting file:", err);
        });
});


document.getElementById('importProdButton').addEventListener('click', function() {
    let modal = document.getElementById('recordsModal');
    // Use Wails' OpenFileDialog to select a file
    window.go.main.App.OpenFileDialog()
        .then(filePath => {
            if (!filePath) {
                alert("No file selected.");
                return;
            }

            // Send the file path to the Go backend
            window.go.main.App.ProcessEvidenceProdFile(filePath)
                .then(() => {
                    modal.style.display = 'none';
                    alert("File processed successfully.");
                })
                .catch(err => {
                    console.error("Error processing file:", err);
                    alert("Failed to process the file.");
                });
        })
        .catch(err => {
            console.error("Error selecting file:", err);
        });
});