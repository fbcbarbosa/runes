package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/user"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const ucdURL = "http://www.unicode.org/Public/UNIDATA/UnicodeData.txt"

func main() {
	ucd, err := openUCD(obtainUCDPath())
	if err != nil {
		log.Fatal(err.Error())
	}
	defer ucd.Close()
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

func obtainUCDPath() string {
	ucdPath := os.Getenv("UCD_PATH")
	if ucdPath == "" {
		user, err := user.Current()
		check(err)
		ucdPath = user.HomeDir + "/UnicodeData.txt"
	}
	return ucdPath
}

func openUCD(path string) (*os.File, error) {
	ucd, err := os.Open(path)
	if os.IsNotExist(err) {
		fmt.Printf("%s not found\nDownloading %s...\n", path, ucdURL)
		done := make(chan bool)
		go downloadUCD(ucdURL, path, done)
		progress(done)
		ucd, err = os.Open(path)
	}
	return ucd, err
}

func downloadUCD(url, path string, done chan<- bool) {
	resp, err := http.Get(url)
	check(err)
	defer resp.Body.Close()
	file, err := os.Create(path)
	check(err)
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	check(err)
	done <- true
}

func progress(done <-chan bool) {
	for {
		select {
		case <-done:
			fmt.Println()
			return
		default:
			fmt.Print(".")
			time.Sleep(200 * time.Millisecond)
		}
	}
}

func restore(name, value string, exists bool) {
	if exists {
		os.Setenv(name, value)
	} else {
		os.Unsetenv(name)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
