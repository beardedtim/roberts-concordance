package data

import (
	"encoding/json"
	"log"
	"mckp/roberts-concordance/globals"
	"os"
	"strings"
)

type BibleVerse struct {
	Book    string   `json:"book"`
	Chapter int      `json:"chapter"`
	Verse   int      `json:"verse"`
	Text    string   `json:"text"`
	Notes   []string `json:"notes"`
}

type BibleBook struct {
	Book   string       `json:"book"`
	Verses []BibleVerse `json:"verses"`
}

var jsonFile []BibleBook
var bookIndex map[string]int
var booksInOrder []string

func readAndAssign() {
	// read JSON and assign data
	file, err := os.ReadFile(globals.ArtifactsDir() + "/parsed/bible.json")

	if err != nil {
		log.Println(err)
		return
	}

	json.Unmarshal(file, &jsonFile)

	// read book names and assign data
	bookNames, _ := os.ReadFile(globals.ArtifactsDir() + "/books.txt")
	split := strings.Split(string(bookNames), "\n")
	bi := make(map[string]int)
	ordered := make([]string, len(split))

	for i, name := range split {
		bi[name] = i
		ordered[i] = name
	}

	bookIndex = bi
	booksInOrder = ordered
}

func GetBooks() []string {
	if len(bookIndex) == 0 {
		readAndAssign()
	}

	return booksInOrder
}

func GetText() []BibleBook {
	if len(jsonFile) == 0 {
		readAndAssign()
	}

	log.Println(bookIndex)

	return jsonFile
}
