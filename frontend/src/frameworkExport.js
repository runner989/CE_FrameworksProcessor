
// Fetch distinct frameworks from the backend and populate the dropdown
window.fetchFrameworks = function() {
    window.go.main.App.GetAllFrameworks()
        .then(function (frameworks) {
            console.log(frameworks);
            populateFrameworkDropdown(frameworks);
        })
        .catch(function (err) {
            console.error('Error fetching frameworks:', err);
            alert('Failed to retrieve frameworks.');
        });
}

// Populate the dropdown with framework options
function populateFrameworkDropdown(frameworks) {
    const frameworkSelect = document.getElementById('frameworkSelect');
    frameworkSelect.innerHTML = '<option value="" disabled selected>Select a Framework</option>'; // Reset

    frameworks.forEach(function(framework) {
    const option = document.createElement('option');
    option.value = framework;
    option.textContent = framework;
    frameworkSelect.appendChild(option);
});

// Enable export button when a framework is selected
frameworkSelect.addEventListener('change', function() {
    document.getElementById('exportFrameworkButton').disabled = false;
});
}

// Handle exporting the selected framework
document.getElementById('exportFrameworkButton').addEventListener('click', function() {
    const selectedFramework = document.getElementById('frameworkSelect').value;

    if (selectedFramework) {
        exportFrameworkToExcel(selectedFramework);
    } else {
        alert('Please select a framework to export.');
    }
});

// Call backend function to export the framework to Excel
function exportFrameworkToExcel(framework) {
    console.log('Exporting framework:', framework);
    window.go.main.App.ExportAFramework(framework)
    .then(function() {
        alert('Framework exported successfully!');
    })
    .catch(function(err) {
        console.error('Error exporting framework:', err);
        alert('Failed to export framework.');
    });
}

// // On page load, fetch frameworks
// document.addEventListener('DOMContentLoaded', fetchFrameworks);
