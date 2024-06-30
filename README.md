
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

https://github.com/samodon/atomic-search/assets/77257036/0a7b7b5f-a276-4c01-af9d-c1777c9d4718

<img width="1470" alt="Screenshot 2024-06-28 at 3 25 58 PM" src="https://github.com/samodon/atomic-search/assets/77257036/7ca1f216-cfe3-413d-841d-efe08a4a97a3">

<img width="1470" alt="Screenshot 2024-06-28 at 4 09 58 PM" src="https://github.com/samodon/atomic-search/assets/77257036/97dbb399-7028-4472-a3c6-bcf57d5f163f">

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

Select the path to your vault of notes:
```sh
search --path <your vault path>
```

To use this tool you are only required to run one command:
```sh
search <your search term>
```

If the application isn't in your path environment variable then usage is as follows:
```sh
./search <your search term>
```

Index updates every 24 hours, to manually force an update run the command:
```sh
search --index
```




