# paperfilter

`paperfilter` is a command-line tool written in Go that uses the Grok LLM (via OpenAI-compatible API) to determine if an academic paper aligns with specific research interests.

## Features

- **Automated Triage**: Analyzes paper abstracts or full text to check for relevance.
- **Structured Output**: Returns a JSON response containing a boolean match and a detailed justification.
- **Flexible Input**: Supports reading from a file or piping from standard input.
- **Grok Integration**: Specifically designed to work with Grok models.

## Usage

### Prerequisites

- Go 1.21 or later.
- A Grok API key from [x.ai](https://api.x.ai).

### Installation & Setup

#### 1. Go Environment Setup

If you don't have Go installed, follow the official [installation guide](https://go.dev/doc/install). 

To verify your installation:
```bash
go version
```

#### 2. Get the Source
```bash
# Clone the repository (if you haven't already)
git clone https://github.com/xuyangkang/paperfilter.git
cd paperfilter
```

#### 3. Dependencies
This project uses Go modules for dependency management. The primary dependency is:
- [go-openai](https://github.com/sashabaranov/go-openai): A Go client library for OpenAI-compatible APIs.

Install dependencies using:
```bash
go mod download
```

#### 4. Build
```bash
go build -o paperfilter
```

### Running the Tool

You can provide input text through a file:

```bash
./paperfilter --input path/to/paper.txt --apikey YOUR_GROK_API_KEY
```

Or pipe content directly (useful for integration with tools like `arxivfetcher`):

```bash
cat paper.txt | ./paperfilter --apikey YOUR_GROK_API_KEY
```

### Configuration Flags

- `--input`: Path to the input file (defaults to stdin).
- `--apikey`: Your x.ai / Grok API key (**required**).
- `--baseurl`: API base URL (defaults to `https://api.x.ai/v1`).
- `--model`: The Grok model to use (defaults to `grok-4-1-fast-reasoning`).

## Example Output

```json
{
  "match": true,
  "justification": "The paper presents an efficient algorithm for frequent substring mining... aligns closely with the user's interests in string algorithms and data structures."
}
```
