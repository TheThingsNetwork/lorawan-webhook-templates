const fs = require('fs');
const path = require('path');
const yaml = require('yaml');
const { createObjectCsvWriter } = require('csv-writer');

// Define paths to the templates.yml file and the directory containing the template files
const templatesYmlPath = path.join(__dirname, '..', 'templates.yml');
const templatesDir = path.join(__dirname, '..');

// Read and parse the templates.yml file
const templatesYml = fs.readFileSync(templatesYmlPath, 'utf8');
const templateNames = yaml.parse(templatesYml);

// Define the CSV writer
const csvWriter = createObjectCsvWriter({
    path: path.join(__dirname, 'templates.csv'),
    header: [
        { id: 'name', title: 'Name' },
        { id: 'description', title: 'Description' },
        { id: 'documentationUrl', title: 'Documentation URL' },
        { id: 'TTSDocumentationUrl', title: 'TTS Documentation URL' },
        { id: 'format', title: 'Format' },
    ]
});

// Function to read and parse each template file
const readTemplateFile = (templateName) => {
    const templateFilePath = path.join(templatesDir, `${templateName}.yml`);
    const templateYml = fs.readFileSync(templateFilePath, 'utf8');
    const templateData = yaml.parse(templateYml);
    let format = templateData.format;
    if (format === 'json') {
        format = format.toUpperCase();
    }
    return {
        name: templateData.name,
        description: templateData.description,
        documentationUrl: templateData['documentation-url'],
        ttsDocumentationUrl: templateData['tts-documentation-url'],
        format: format,
    };
};

// Collect data from each template file
const templatesData = templateNames.map(templateName => readTemplateFile(templateName));

// Write data to CSV
csvWriter.writeRecords(templatesData)
    .then(() => console.log('Templates data written to templates.csv'))
    .catch(err => console.error('Error writing to CSV file', err));
