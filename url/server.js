const http = require('http');
const fs = require('fs');
const path = require('path');

const server = http.createServer((req, res) => {
    // Define the file to be served
    const filePath = path.join(__dirname, 'assets', 'index.html');

    // Read the file and send its content
    fs.readFile(filePath, (err, content) => {
        if (err) {
            res.writeHead(500, { 'Content-Type': 'text/plain' });
            res.end('Server Error');
            return;
        }

        // Serve the HTML file
        res.writeHead(200, { 'Content-Type': 'text/html' });
        res.end(content, 'utf-8');
    });
});

// Define the port to listen on
const PORT = 8080;

server.listen(PORT, () => {
    console.log(`Server running at http://localhost:${PORT}/`);
});
