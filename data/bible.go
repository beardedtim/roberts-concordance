package data

import (
	"encoding/json"
	"log"
	"mckp/roberts-concordance/globals"
	"os"
	"strings"

	"github.com/kljensen/snowball"
)

type BibleVerse struct {
	Book    string   `json:"book"`
	Chapter int      `json:"chapter"`
	Verse   int      `json:"verse"`
	Text    string   `json:"text"`
	Notes   []string `json:"notes"`
}

type WordIndex struct {
	Book    string `json:"book"`
	Chapter int    `json:"chapter"`
	Verse   int    `json:"verse"`
}

type BibleChapter struct {
	Book    string       `json:"book"`
	Chapter int          `json:"chapter"`
	Verses  []BibleVerse `json:"verses"`
}

type BibleBook struct {
	Book     string         `json:"book"`
	Verses   []BibleVerse   `json:"verses"`
	Chapters []BibleChapter `json:"chapters"`
}

var jsonFile []BibleBook
var bookIndex map[string]int
var booksInOrder []string
var wordIndex map[string][]WordIndex
var stemIndex map[string][]WordIndex

var BIBLE_JSON_PATH = globals.ArtifactsDir() + "/parsed/bible.json"
var BIBLE_BOOKS_PATH = globals.ArtifactsDir() + "/books.txt"
var INDEX_WORDS_PATH = globals.ArtifactsDir() + "/parsed/index-words.json"
var INDEX_STEM_PATH = globals.ArtifactsDir() + "/parsed/index-stemmed.json"

func readAndAssign() {
	// read JSON and assign data
	file, err := os.ReadFile(BIBLE_JSON_PATH)

	if err != nil {
		log.Println(err)
		return
	}

	json.Unmarshal(file, &jsonFile)

	// read book names and assign data
	bookNames, _ := os.ReadFile(BIBLE_JSON_PATH)
	split := strings.Split(string(bookNames), "\n")
	bi := make(map[string]int)
	ordered := make([]string, len(split))

	for i, name := range split {
		bi[name] = i
		ordered[i] = name
	}

	bookIndex = bi
	booksInOrder = ordered

	// read index and assign data
	indexRaw, _ := os.ReadFile(INDEX_WORDS_PATH)

	json.Unmarshal(indexRaw, &wordIndex)

	stemRaw, _ := os.ReadFile(INDEX_STEM_PATH)

	json.Unmarshal(stemRaw, &stemIndex)
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

func GetBookByName(name string) BibleBook {
	if len(bookIndex) == 0 {
		readAndAssign()
	}

	bookI := bookIndex[name]
	verses := jsonFile[bookI]

	return verses
}

func GetChapterFromBook(book string, chapter int) BibleChapter {
	bibleBook := GetBookByName(book)
	return bibleBook.Chapters[chapter]
}

func GetVerseFromBook(book string, chapter int, start int, end int) []BibleVerse {
	bibleBook := GetBookByName(book)
	chapterList := bibleBook.Chapters[chapter]

	return chapterList.Verses[start:end]
}

func GetVersesByIndexStem(search string) []WordIndex {
	if len(stemIndex) == 0 {
		readAndAssign()
	}

	stemmed, _ := snowball.Stem(search, "english", true)
	stemmedMatch := stemIndex[stemmed]

	return stemmedMatch
}

func GetVersesByIndexExact(search string) []WordIndex {
	if len(stemIndex) == 0 {
		readAndAssign()
	}

	result := wordIndex[search]

	return result
}
