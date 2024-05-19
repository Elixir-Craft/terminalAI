[![Go Reference](https://pkg.go.dev/badge/github.com/Elixir-Craft/terminalAI.svg)](https://pkg.go.dev/github.com/Elixir-Craft/terminalAI)

# Terminal AI

A command line tool that uses Gemini AI and OpenAI API to generate text and code.


## Installation


If you using Arch Linux or Arch based distro, you can install the terminalAI using [PKGBUILD](https://github.com/Elixir-Craft/terminalAI-pkg) 


You have to install Go and set up the environment before installing the tool.


```bash
git clone https://github.com/Elixir-Craft/terminalAI.git
cd terminalAI
go mod tidy
```


## Build

```bash
go build -v -o out ./...  
```


## Options

* `-i <input file>` Input file path
* `-o <output file>` Output file path
* `-p <Prompt>` Prompt
* `-c ` Prompt from clipboard
* `chat` Chat with AI
* `config` Configure Services and API Keys


## Set configurations

```bash
./terminalai config init
```



## Usage

```bash
./terminalai "Generate Some Text"
```
```bash
./terminalai -p " $(tree -L 2) generate docker-compose file for this project" -o docker-compose.yaml 
```
```bash
./terminalai -i input.txt -p "Read the following text and generate a summary" -o output.txt
```


Chat (Currently only supports for Gemini)
```bash
./terminalai chat
```



