package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func main() {
	ucd, err := os.Open("UnicodeData.txt")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer func() { ucd.Close() }()
	query := strings.Join(os.Args[1:], " ")
	List(ucd, strings.ToUpper(query))
}

// List outputs the code, the rune and the name of the Unicode characters
// that contains the query string
func List(text io.Reader, query string) {
	fields := strings.Fields(query)
	scanner := bufio.NewScanner(text)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		r, name, words := readLine(line)
		if containsAll(words, fields) {
			fmt.Printf("U+%04X\t%[1]c\t%s\n", r, name)
		}
	}
}

func readLine(line string) (rune, string, []string) {
	fields := strings.Split(line, ";")
	name := fields[1]
	uni1Name := fields[10]
	if uni1Name != "" {
		name += fmt.Sprintf(" (%s)", fields[10])
	}
	code, _ := strconv.ParseInt(fields[0], 16, 32)
	words := removeDuplicates(strings.FieldsFunc(name, isSeparator))
	return rune(code), name, words
}

func contains(slice []string, query string) bool {
	for _, v := range slice {
		if v == query {
			return true
		}
	}
	return false
}

func containsAll(slice []string, queries []string) bool {
	for _, q := range queries {
		if !contains(slice, q) {
			return false
		}
	}
	return true
}

func isSeparator(r rune) bool {
	return unicode.IsSpace(r) || r == '-' || r == '(' || r == ')'
}

func removeDuplicates(slice []string) []string {
	var newSlice []string
	for _, v := range slice {
		if !contains(newSlice, v) {
			newSlice = append(newSlice, v)
		}
	}
	return newSlice
}
