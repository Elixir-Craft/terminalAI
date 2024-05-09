# Terminal AI

A command line tool that uses Gemini AI and OpenAI API to generate text and code.


## Installation

```bash
git clone https://github.com/Elixir-Craft/terminalAI.git
cd terminalAI
go mod tidy
```


## Build

```bash
 go build -o terminal-ai
```


## Options

* `-i <input file>` Input file path
* `-o <output file>` Output file path
* `-p <Prompt>` Prompt
* `-c ` Prompt from clipboard
* `--chat` Chat with AI


## Usage

```bash
./terminal-ai "Generate Some Text"
```
```bash
./terminal-ai -p " $(tree -L 2) generate docker-compose file for this project" -o docker-compose.yaml 
```
```bash
./terminal-ai -i input.txt -p "Read the following text and generate a summary" -o output.txt
```


Chat (Currently only supports for Gemini)
```bash
./terminal-ai --chat
```



