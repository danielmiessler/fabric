# Using Fabric Qualifiers in the Web Interface

## How Qualifiers Work

The Fabric web interface connects to a local Fabric backend server that processes command-line qualifiers. The backend handles all the same qualifiers that are available in the CLI version.

### Using Qualifiers

You can use any of the CLI qualifiers directly in the chat input field. For example:

1. Language qualifier:
```
hello -g=fr
```

2. Multiple qualifiers:
```
hello -g=fr -t=0.7 -P=0.2
```

### Available Qualifiers

The following qualifiers are supported:

- `-g, --language=` : Specify language code (e.g., -g=fr for French)
- `-t, --temperature=` : Set temperature (default: 0.7)
- `-T, --topp=` : Set top P (default: 0.9)
- `-P, --presencepenalty=` : Set presence penalty (default: 0.0)
- `-F, --frequencypenalty=` : Set frequency penalty (default: 0.0)
- `-r, --raw` : Use model defaults without chat options
- `-m, --model=` : Choose specific model

### Special Cases

Some qualifiers have special frontend handling:

1. YouTube URLs (-y):
   - Just paste the YouTube URL directly
   - The frontend detects it and adds -y automatically
   - The frontend gets the transcript before sending to backend

2. Patterns (-p):
   - Select from the pattern dropdown
   - The frontend adds -p automatically

### How It Works

1. You enter text with qualifiers in the chat input
2. The frontend sends the complete input to the Fabric backend (localhost:8080)
3. The backend parses the input and processes any qualifiers using the same logic as the CLI
4. The backend executes the command with the specified options
5. Results are streamed back to the web interface

### Examples

1. Get response in French:
```
Tell me about AI -g=fr
```

2. Adjust temperature and presence penalty:
```
Tell me about AI -t=0.9 -P=0.2
```

3. Combine with pattern:
(Select pattern from dropdown, then type)
```
Tell me about AI -g=fr -t=0.8
```

The backend will combine the automatically added pattern qualifier (-p) with your manually entered qualifiers (-g=fr -t=0.8).

### Technical Details

The qualifiers are processed by the Fabric backend in this order:

1. Input text is received by the backend
2. The backend parses the input for qualifiers using the same parser as the CLI
3. Qualifiers are converted into appropriate options:
   - Language codes are validated
   - Numeric values are checked for valid ranges
   - Special flags are processed
4. The processed request is sent to the AI model with all the specified options

This means you can use any qualifier that works in the CLI directly in the web interface's chat input.