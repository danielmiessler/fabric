import React, { useState, useEffect } from 'react';

const API_BASE_URL = 'http://192.168.4.21:5000';

function App() {
  const [patterns, setPatterns] = useState([]);
  const [selectedPattern, setSelectedPattern] = useState('');
  const [params, setParams] = useState('');
  const [youtubeUrl, setYoutubeUrl] = useState('');
  const [output, setOutput] = useState('');

  useEffect(() => {
    fetch(`${API_BASE_URL}/patterns`)
      .then(response => response.json())
      .then(data => setPatterns(data))
      .catch(error => console.error('Error fetching patterns:', error));
  }, []);

  const handleSubmit = (e) => {
    e.preventDefault();

    let fabricCommand = ['fabric'];
    if (selectedPattern) {
      fabricCommand.push('--pattern', selectedPattern);
    }
    if (youtubeUrl) {
      fabricCommand.push('-y', `"${youtubeUrl}"`);
    }
    if (params) {
      fabricCommand = fabricCommand.concat(params.split(' '));
    }

    const requestData = {
      pattern: selectedPattern,
      params: fabricCommand.slice(1) // Remove 'fabric' from the beginning
    };

    fetch(`${API_BASE_URL}/run-pattern`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(requestData)
    })
      .then(response => response.json())
      .then(data => {
        setOutput(data.output || data.error);
      })
      .catch(error => {
        setOutput(`Error: ${error.message}`);
        console.error('Error running pattern:', error);
      });
  };

  return (
    <div className="App">
      <style>
        {`
          body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            background-color: #f4f4f4;
            margin: 0;
            padding: 0;
          }
          .container {
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
          }
          h1 {
            color: #2c3e50;
            text-align: center;
            margin-bottom: 30px;
          }
          .form-group {
            margin-bottom: 20px;
          }
          label {
            display: block;
            margin-bottom: 5px;
            font-weight: bold;
          }
          select, input[type="text"] {
            width: 100%;
            padding: 10px;
            border: 1px solid #ddd;
            border-radius: 4px;
            box-sizing: border-box;
          }
          button {
            background-color: #3498db;
            color: white;
            padding: 10px 20px;
            border: none;
            border-radius: 4px;
            cursor: pointer;
            font-size: 16px;
          }
          button:hover {
            background-color: #2980b9;
          }
          .output {
            background-color: white;
            border: 1px solid #ddd;
            border-radius: 4px;
            padding: 20px;
            margin-top: 20px;
          }
          pre {
            white-space: pre-wrap;
            word-wrap: break-word;
          }
        `}
      </style>
      <div className="container">
        <h1>Fabric Pattern Runner</h1>
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label htmlFor="pattern">Select Pattern:</label>
            <select
              id="pattern"
              value={selectedPattern}
              onChange={(e) => setSelectedPattern(e.target.value)}
            >
              <option value="">--Select Pattern--</option>
              {patterns.map((pattern) => (
                <option key={pattern} value={pattern}>{pattern}</option>
              ))}
            </select>
          </div>
          <div className="form-group">
            <label htmlFor="youtube">YouTube URL:</label>
            <input
              id="youtube"
              type="text"
              value={youtubeUrl}
              onChange={(e) => setYoutubeUrl(e.target.value)}
              placeholder="Enter YouTube URL (optional)"
            />
          </div>
          <div className="form-group">
            <label htmlFor="params">Additional Parameters:</label>
            <input
              id="params"
              type="text"
              value={params}
              onChange={(e) => setParams(e.target.value)}
              placeholder="Enter additional parameters (separated by space)"
            />
          </div>
          <button type="submit">Run Pattern</button>
        </form>
        <div className="output">
          <h2>Output:</h2>
          <pre>{output}</pre>
        </div>
      </div>
    </div>
  );
}

export default App;
