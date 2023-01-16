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

type BibleBook struct {
	Book   string       `json:"book"`
	Verses []BibleVerse `json:"verses"`
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

func main() {
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

	allVerses := parsedJSONData.Resultset.Row

	reg := `\{.*?\}`
	compiled, _ := regexp.Compile(reg)

	var verses []Field

	for _, field := range allVerses {
		matched := compiled.FindStringIndex(field.Field.Text)

		if matched != nil {
			beginIndex := matched[0]
			endIndex := matched[1]

			note := field.Field.Text[beginIndex:endIndex]

			notes := []string{note[1 : len(note)-1]}

			verses = append(verses, Field{
				Book:    field.Field.Book,
				Chapter: field.Field.Chapter,
				Verse:   field.Field.Verse,
				Text:    strings.ReplaceAll(field.Field.Text, note, ""),
				Notes:   notes,
			})
		} else {
			verses = append(verses, Field{
				Book:    field.Field.Book,
				Chapter: field.Field.Chapter,
				Verse:   field.Field.Verse,
				Text:    field.Field.Text,
				Notes:   []string{},
			})
		}
	}

	file, _ := json.MarshalIndent(verses, "", " ")

	// Write Verses File
	_ = os.WriteFile(globals.ArtifactsDir()+"/parsed/verses.json", file, 0644)

	bookFilePath := globals.ArtifactsDir() + "/books.txt"

	bookBytes, err := os.ReadFile(bookFilePath)

	if err != nil {
		log.Panic("Cannot read book file path "+bookFilePath, err)
	}

	bookNames := strings.Split(string(bookBytes), "\n")

	books := make([]BibleBook, len(bookNames))

	for i, bookName := range bookNames {
		book := books[i]

		book.Book = bookName
		book.Verses = []BibleVerse{}

		books[i] = book
	}

	wordIndex := make(Index)

	for _, field := range verses {
		verse := field

		bookIndex := verse.Book

		book := books[bookIndex]

		// Make book
		verseHash := BibleVerse{
			Book:    book.Book,
			Chapter: verse.Chapter,
			Verse:   verse.Verse,
			Text:    verse.Text,
			Notes:   verse.Notes,
		}
		book.Verses = append(book.Verses, verseHash)

		books[bookIndex] = book

		// make index
		text := verseHash.Text
		var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9\-]+`)

		for _, word := range strings.Fields(text) {
			index := nonAlphanumericRegex.ReplaceAllString(word, "")
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
					wordIndex[splitWord] = append(wordIndex[splitWord], WordIndex{
						Book:    verseHash.Book,
						Chapter: verse.Chapter,
						Verse:   verse.Verse,
					})
				}
			}

			wordIndex[index] = append(wordIndex[word], WordIndex{
				Book:    verseHash.Book,
				Chapter: verse.Chapter,
				Verse:   verse.Verse,
			})
		}
	}

	bibleBook, _ := json.MarshalIndent(books, "", " ")

	// Write Verses File
	_ = os.WriteFile(globals.ArtifactsDir()+"/parsed/bible.json", bibleBook, 0644)
	index, _ := json.MarshalIndent(wordIndex, "", " ")

	// Write Index File
	_ = os.WriteFile(globals.ArtifactsDir()+"/parsed/index.json", index, 0644)
}
