/*
	Parses artifacts/raw.json into a list of verses
	and outputs those verses into artifacts/parsed/verses.json

	Parses artifacts/books.txt into a hash of Books: [Verses]
	and outputs that hash into artifacts/parsed/bible.json

	Parses artifacts/raw.json into an index of where each word
	is used throughout the Bible and outputs into artifacts/parsed/index.json
*/

package main

import (
	"encoding/json"
	"log"
	"mckp/roberts-concordance/globals"

	"os"
	"regexp"
	"strings"

	"github.com/kljensen/snowball"
	"golang.org/x/exp/slices"
)

type Field struct {
	Book    int      `json:"book"`
	Chapter int      `json:"chapter"`
	Verse   int      `json:"verse"`
	Text    string   `json:"text"`
	Notes   []string `json:"notes"`
}

type BibleVerse struct {
	Book    string   `json:"book"`
	Chapter int      `json:"chapter"`
	Verse   int      `json:"verse"`
	Text    string   `json:"text"`
	Notes   []string `json:"notes"`
}

type BibleChapter struct {
	Book    string       `json:"book"`
	Chapter int          `json:"chapter"`
	Verses  []BibleVerse `json:"verses"`
}

type BibleBook struct {
	Book     string         `json:"book"`
	Chapters []BibleChapter `json:"chapters"`
	Verses   []BibleVerse   `json:"verses"`
}

type WordIndex struct {
	Book    string `json:"book"`
	Chapter int    `json:"chapter"`
	Verse   int    `json:"verse"`
}

type Index map[string][]WordIndex

func (f *Field) UnmarshalJSON(p []byte) error {
	var tmp []interface{}

	if err := json.Unmarshal(p, &tmp); err != nil {
		return err
	}

	f.Book = int(tmp[1].(float64)) - 1
	f.Chapter = int(tmp[2].(float64))
	f.Verse = int(tmp[3].(float64))
	f.Text = tmp[4].(string)

	return nil
}

type JSONData struct {
	Resultset struct {
		Row []struct {
			Field Field `json:"field"`
		} `json:"row"`
	} `json:"resultset"`
}

func getBibleJSON() []struct {
	Field Field `json:"field"`
} {
	// Parse JSON into Struct Above
	filePath := globals.ArtifactsDir() + "/raw.json"

	rawJSONContent, err := os.ReadFile(filePath)

	if err != nil {
		log.Panic("Error when opening file: ", err)
	}

	var parsedJSONData JSONData

	err = json.Unmarshal(rawJSONContent, &parsedJSONData)

	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	return parsedJSONData.Resultset.Row
}

func writeVersesToDisk(verses []BibleVerse) {
	file, _ := json.MarshalIndent(verses, "", " ")

	// Write Verses File
	_ = os.WriteFile(globals.ArtifactsDir()+"/parsed/verses.json", file, 0644)
}

func getBooks() []BibleBook {
	bookFilePath := globals.ArtifactsDir() + "/books.txt"

	bookBytes, err := os.ReadFile(bookFilePath)

	if err != nil {
		log.Panic("Cannot read book file path "+bookFilePath, err)
	}

	names := strings.Split(string(bookBytes), "\n")

	books := make([]BibleBook, len(names))

	for i, bookName := range names {
		book := books[i]

		book.Book = bookName
		book.Verses = []BibleVerse{}

		books[i] = book
	}

	return books
}

func groupVersesByChapter(books []BibleBook, verses []BibleVerse) []BibleBook {
	bookIndex := make(map[string]int)

	for _, verse := range verses {
		index, exists := bookIndex[verse.Book]

		if !exists {
			index = slices.IndexFunc(books, func(book BibleBook) bool {
				return book.Book == verse.Book
			})

			bookIndex[verse.Book] = index
		}

		book := books[index]

		// assign all verses to book
		book.Verses = append(book.Verses, verse)

		// if this chapter already exists
		if len(book.Chapters) > verse.Chapter-1 {
			book.Chapters[verse.Chapter-1].Verses = append(book.Chapters[verse.Chapter-1].Verses, verse)
		} else {
			book.Chapters = append(book.Chapters, BibleChapter{
				Book:    verse.Book,
				Chapter: verse.Chapter,
				Verses:  []BibleVerse{verse},
			})
		}

		books[index] = book
	}

	return books
}

func writeBooksToDisk(books []BibleBook) {

	bibleBook, _ := json.MarshalIndent(books, "", " ")

	// Write Verses File
	_ = os.WriteFile(globals.ArtifactsDir()+"/parsed/bible.json", bibleBook, 0644)
}

func fieldtoVerse(field Field, books []BibleBook) BibleVerse {
	bookIndex := field.Book

	book := books[bookIndex]

	return BibleVerse{
		Book:    book.Book,
		Chapter: field.Chapter,
		Verse:   field.Verse,
		Text:    field.Text,
		Notes:   field.Notes,
	}
}

func mapFieldsToVerses(fields []Field, books []BibleBook) []BibleVerse {
	output := make([]BibleVerse, len(fields))

	for i, value := range fields {
		output[i] = fieldtoVerse(value, books)
	}

	return output
}

func createIndex(verses []BibleVerse, books []BibleBook) (Index, Index) {
	wordIndex := make(Index)
	stemIndex := make(Index)

	for _, verse := range verses {
		// make index
		text := verse.Text
		var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9\-]+`)

		for _, word := range strings.Fields(text) {
			index := strings.ToLower(nonAlphanumericRegex.ReplaceAllString(word, ""))

			if strings.Contains(index, "-") {
				// if the word is hyphenated, let's make sure that we
				// index all three options. Ex:
				//
				// wood-offering becomes
				// wood
				// offering
				// wood-offering
				split := strings.Split(index, "-")

				for _, splitWord := range split {
					stemmed, _ := snowball.Stem(splitWord, "english", true)
					stemIndex[stemmed] = append(stemIndex[stemmed], WordIndex{
						Book:    verse.Book,
						Chapter: verse.Chapter,
						Verse:   verse.Verse,
					})

					wordIndex[splitWord] = append(wordIndex[splitWord], WordIndex{
						Book:    verse.Book,
						Chapter: verse.Chapter,
						Verse:   verse.Verse,
					})
				}
			}

			stemmed, _ := snowball.Stem(index, "english", true)

			stemIndex[stemmed] = append(stemIndex[stemmed], WordIndex{
				Book:    verse.Book,
				Chapter: verse.Chapter,
				Verse:   verse.Verse,
			})

			wordIndex[index] = append(wordIndex[index], WordIndex{
				Book:    verse.Book,
				Chapter: verse.Chapter,
				Verse:   verse.Verse,
			})
		}
	}

	return wordIndex, stemIndex
}

func writeIndexToDisk(wordIndex Index, stemIndex Index) {
	windex, _ := json.MarshalIndent(wordIndex, "", " ")
	sindex, _ := json.MarshalIndent(stemIndex, "", " ")

	_ = os.WriteFile(globals.ArtifactsDir()+"/parsed/index-words.json", windex, 0644)
	_ = os.WriteFile(globals.ArtifactsDir()+"/parsed/index-stemmed.json", sindex, 0644)

}

var compiled, _ = regexp.Compile(`\{.*?\}`)

func ParseFieldInformation(field Field) Field {
	matched := compiled.FindStringIndex(field.Text)

	if matched != nil {
		beginIndex := matched[0]
		endIndex := matched[1]

		note := field.Text[beginIndex:endIndex]

		notes := []string{note[1 : len(note)-1]}

		return Field{
			Book:    field.Book,
			Chapter: field.Chapter,
			Verse:   field.Verse,
			Text:    strings.ReplaceAll(field.Text, note, ""),
			Notes:   notes,
		}
	} else {
		return Field{
			Book:    field.Book,
			Chapter: field.Chapter,
			Verse:   field.Verse,
			Text:    field.Text,
			Notes:   []string{},
		}
	}
}

func main() {
	allVerses := getBibleJSON()

	var fields []Field

	for _, field := range allVerses {
		fields = append(fields, ParseFieldInformation(field.Field))
	}
	books := getBooks()

	verses := mapFieldsToVerses(fields, books)

	books = groupVersesByChapter(books, verses)
	wordIndex, stemIndex := createIndex(verses, books)
	writeVersesToDisk(verses)
	writeBooksToDisk(books)
	writeIndexToDisk(wordIndex, stemIndex)
}
