# CE Frameworks Processor (CEFP)

### This is a special use application I created to automate the workflow.

The CE Frameworks Processor is a desktop application for the Coalfire Compliance Essentials frameworks and mappings processing from Airtable to CE. 

This is written in Go, using WAILS for the frontend (HTML and Vanilla Javascript) and SQLite for the database.  

An API key is required to run the application. There is a "secret" folder that contains the compile time API Key. If a new API key is required, it can be loaded through a .env file. 

