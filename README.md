
# Atomic Search
---
Atomic Search is an application built to help me quickly view and make use of my atomic notes.

## Motivation 

I have recently become interested in Personal Knowledge Systems, specifically the Zettelkasten system, which involves separating singular ideas into atomic notes for future linking. This system is powerful in its linking features, allowing you to synthesize completely new ideas by discovering novel connections between certain notes. I have become good at creating these notes, but reading them back and absorbing the information has still been difficult. Therefore, I designed and built this application to aid in my note ingestion and usage, reducing my reliance on a Google search or a ChatGPT query.
## Features
- Automatic Syntax Highlighting of code blocks using [Chroma](https://github.com/alecthomas/chroma)
- Indexing system to maximize showing search results
- Flip through the notes as you would any pages with just the press of a key.
- Outputs information in the terminal, allowing for on hand use whenever needed.

## Atomic Search in Use

> [!todo] Insert Pictures

## Installation
To begin using this project here are the steps to install and begin using:

1. **Install [GO](https://go.dev/doc/install)**
The steps vary by operating system. See the appropriate install method at the [official website](https://go.dev/doc/install).

2. **Clone the repository**
```sh
git clone https://github.com/samodon/atomic-search.git
```

3. **Install dependencies**
```sh
cd $HOME/atomic-search
go mod download
```
This will first put you in the right directory then download the required packages from `go.mod`.

4. **Build an Executable**
```sh
go build .
```
This will build the binary that runs the software.
5. **Add the executable to your PATH**(Optional)
```sh
export PATH="$PATH:/usr/local/go/bin
cd <path_to_executable> /usr/local/go/bin
```
Note: The paths might be different depending on your install of GO. Verify the path to you GO installation before adding to path.
## Usage

To use this tool you are only required to run one command:
```sh
search <your search term>
```

If the application isn't in your path environment variable then usage is as follows:
```sh
./search <your search term>
```




