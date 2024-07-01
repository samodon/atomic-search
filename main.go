package main

import "C"
import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"samodon/search/indexing"
	"samodon/search/pkg"
	"samodon/search/searching"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/alecthomas/chroma/v2/quick"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// DisplayNote processes and displays the note with syntax highlighting for code blocks using tview.
func DisplayNote(html string, langauge string, fileName string) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Fatal(err)
	}

	// Initialize tview components
	app := tview.NewApplication()
	textView := tview.NewTextView().SetTextAlign(tview.AlignLeft).SetDynamicColors(true)
	textView.SetBackgroundColor(tcell.ColorDefault)
	textView.SetDynamicColors(true)
	textView.SetBorder(true)
	textView.SetScrollable(true)
	textView.SetBorderColor(tcell.NewHexColor(0x658594))
	textView.SetTitle(fileName)
	textView.SetTitleColor(tcell.NewHexColor(0xE6C384))
	textView.SetTextColor(tcell.NewHexColor(0xDCD7BA))
	flex := tview.NewFlex().SetDirection(tview.FlexRow).AddItem(textView, 0, 1, true)
	textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' {
			app.Stop()
		}
		return event
	})
	// Iterate over all elements in the document
	doc.Find("*").Each(func(index int, item *goquery.Selection) {
		// Handle text nodes
		if item.Is("p, span, div") {
			text := strings.TrimSpace(item.Text())
			if text != "" {
				fmt.Fprintf(textView, "%s\n", text)
			}
		}

		// Handle code blocks
		if item.Is("pre code") {
			// Get the raw HTML content of the code block
			codeHTML, err := item.Html()
			if err != nil {
				log.Println(err)
				return
			}
			// TODO: Automatically choose the langauge for synax highlighting
			writer := tview.ANSIWriter(textView)
			quick.Highlight(writer, fmt.Sprint(codeHTML), langauge, "terminal256", "vim")
		}
	})

	// Set root and run tview application
	if err := app.SetRoot(flex, true).Run(); err != nil {
		log.Fatal(err)
	}
}

func SaveFilePath(filePath, configFile string) error {
	return os.WriteFile(configFile, []byte(filePath), 0644)
}

// LoadFilePath loads the file path from a text file.
func LoadFilePath(configFile string) string {
	data, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Println("Make sure you have added the correct vault path with --path <vaultpath>")
		fmt.Println(err)
	}
	return strings.TrimSpace(string(data))
}

func main() {
	usr, err := user.Current()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if len(os.Args) < 2 {
		fmt.Println("Usage: search <term>")
		return
	}
	termarr := os.Args[1:]
	indexpath := filepath.Join(usr.HomeDir, "/index/wordindex.json")
	fileInfo, _ := os.Stat(indexpath)

	timedifference := time.Now().Sub(fileInfo.ModTime())
	if os.Args[1] == "--path" {
		err := SaveFilePath(os.Args[2], filepath.Join(usr.HomeDir, "/index/config.txt"))
		if err != nil {
			fmt.Println(err)
		}

		return
	}
	vaultPath := LoadFilePath(filepath.Join(usr.HomeDir, "/index/config.txt"))
	if timedifference.Hours() > 24 {
		indexing.CreateIndex(vaultPath)
	}

	if os.Args[1] == "--index" {
		indexing.CreateIndex(vaultPath)
		return
	}

	searchTerm := strings.Join(termarr, " ")
	searchTerm = pkg.RemoveWords(searchTerm)

	sortedResults := searching.GetSearchRanking(searchTerm, filepath.Join(usr.HomeDir, "/index/wordindex.json"), filepath.Join(usr.HomeDir, "/index/tagindex.json"))
	if len(sortedResults) > 1 {
		fileName := filepath.Base(sortedResults[0].NoteLocation)

		content, tags, _ := (pkg.ParseMdData(sortedResults[0].NoteLocation))
		language := fmt.Sprint(tags["Language"])

		language = strings.Trim(language, "[]")
		DisplayNote(content, language, fileName)
	} else {
		fmt.Print("No notes found ")
	}
}
