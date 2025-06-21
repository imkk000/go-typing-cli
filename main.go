package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"strings"

	"atomicgo.dev/color"
	"atomicgo.dev/cursor"
	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
)

func main() {
	textArea := cursor.NewArea()
	text := resetArea(textArea)

	keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		switch key.Code {
		case keys.CtrlC:
			return true, nil
		case keys.Backspace:
			cursor.Left(1)
			x = max(0, x-1)
			r := rune(text[x])
			s := string(r)
			defaultStyle.Print(s)
			cursor.Left(1)
		case keys.RuneKey, keys.Space:
			k := key.Runes[0]
			r := rune(text[x])
			s := string(r)
			if r == k {
				correctStyle.Print(s)
			} else {
				incorrectStyle.Print(s)
			}
			x = min(len(text), x+1)

			if key.Code == keys.Space && current < x {
				wordCount++
				textArea.StartOfLine()
				fmt.Printf("%-3d", wordCount)
				textArea.Move(x+1, 0)
				current = x
			}
		}
		if x >= len(text) {
			wordCount++
			text = resetArea(textArea)
		}

		return false, nil
	})
}

func randomWords() []string {
	content, err := os.ReadFile("words.json")
	if err != nil {
		log.Fatal("read file:", err)
	}
	var data []string
	if err := json.Unmarshal(content, &data); err != nil {
		log.Fatal("unmarshal:", err)
	}
	var randomWord func() string
	randomWord = func() string {
		i := rand.IntN(len(data))
		if _, ok := dedupedWords[data[i]]; ok {
			return randomWord()
		}
		dedupedWords[data[i]] = struct{}{}
		return data[i]
	}

	words := make([]string, maxWords)
	for i := range maxWords {
		words[i] = randomWord()
	}
	return words
}

func resetArea(a cursor.Area) string {
	words := randomWords()
	text := strings.Join(words, " ")
	a.Clear()
	a.Update(defaultStyle.Sprintf("%-3d %s", wordCount, text))
	a.StartOfLine()
	a.Move(4, 0)
	x = 0
	current = 0
	return text
}

var (
	x              int
	current        int
	wordCount      int
	dedupedWords   = make(map[string]struct{})
	defaultStyle   = color.NewStyle(color.ANSICyan, color.NoColor)
	correctStyle   = color.NewStyle(color.ANSIBlack, color.ANSIYellow)
	incorrectStyle = color.NewStyle(color.ANSICyan, color.ANSIRed)
)

const maxWords = 5
