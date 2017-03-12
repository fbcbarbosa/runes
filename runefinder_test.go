package main

import "testing"
import "strings"
import "os"
import "reflect"

const lineLetterA = `0041;LATIN CAPITAL LETTER A;Lu;0;L;;;;;N;;;;0061;`
const lineApostrophe = `0027;APOSTROPHE;Po;0;ON;;;;;N;APOSTROPHE-QUOTE;;;;`

const lines3Dto43 = `
003D;EQUALS SIGN;Sm;0;ON;;;;;N;;;;;
003E;GREATER-THAN SIGN;Sm;0;ON;;;;;Y;;;;;
003F;QUESTION MARK;Po;0;ON;;;;;N;;;;;
0040;COMMERCIAL AT;Po;0;ON;;;;;N;;;;;
0041;LATIN CAPITAL LETTER A;Lu;0;L;;;;;N;;;;0061;
0042;LATIN CAPITAL LETTER B;Lu;0;L;;;;;N;;;;0062;
0043;LATIN CAPITAL LETTER C;Lu;0;L;;;;;N;;;;0063;
`

func TestReadLine(t *testing.T) {
	r, name, words := readLine(lineLetterA)
	if r != 'A' {
		t.Errorf("Expected 'A', got %q", r)
	}
	const nameA = "LATIN CAPITAL LETTER A"
	if name != nameA {
		t.Errorf("Expected %q, got %q", nameA, name)
	}
	wordsA := []string{"LATIN", "CAPITAL", "LETTER", "A"}
	if !reflect.DeepEqual(words, wordsA) {
		t.Errorf("\n\tExpected %q, got %q", wordsA, words)
	}
}

func TestReadLineWithUnicode1(t *testing.T) {
	r, name, words := readLine(lineApostrophe)
	if r != '\'' {
		t.Errorf("Expected ''', got %q", r)
	}
	const nameAp = "APOSTROPHE (APOSTROPHE-QUOTE)"
	if name != nameAp {
		t.Errorf("Expected %q, got %q", nameAp, name)
	}
	wordsAp := []string{"APOSTROPHE", "QUOTE"}
	if !reflect.DeepEqual(words, wordsAp) {
		t.Errorf("\n\tExpected %q, got %q", wordsAp, words)
	}
}

func TestContains(t *testing.T) {
	var tests = []struct {
		slice    []string
		query    string
		expected bool
	}{
		{[]string{"A", "B"}, "B", true},
		{[]string{}, "A", false},
		{[]string{"A", "B"}, "Z", false},
	}
	for _, test := range tests {
		got := contains(test.slice, test.query)
		if got != test.expected {
			t.Errorf("contains(%#v, %#v) expected: %v; got: %v",
				test.slice, test.query, test.expected, got)
		}
	}
}

func TestContainsAll(t *testing.T) {
	var tests = []struct {
		slice    []string
		queries  []string
		expected bool
	}{
		{[]string{"A", "B"}, []string{"B"}, true},
		{[]string{}, []string{"A"}, false},
		{[]string{"A"}, []string{}, true},
		{[]string{}, []string{}, true},
		{[]string{"A", "B"}, []string{"Z"}, false},
		{[]string{"A", "B", "C"}, []string{"A", "C"}, true},
		{[]string{"A", "B", "C"}, []string{"A", "Z"}, false},
		{[]string{"A", "B"}, []string{"A", "B", "C"}, false},
	}
	for _, test := range tests {
		got := containsAll(test.slice, test.queries)
		if got != test.expected {
			t.Errorf("containsAll(%#v, %#v)\nexpected: %v, got: %v",
				test.slice, test.queries, test.expected, got)
		}
	}
}

func TestIsSeparator(t *testing.T) {
	var tests = []struct {
		symbol   rune
		expected bool
	}{
		{' ', true},
		{'-', true},
		{'(', true},
		{')', true},
		{'A', false},
		{'_', false},
	}
	for _, test := range tests {
		got := isSeparator(test.symbol)
		if got != test.expected {
			t.Errorf("isSeperator(%q) expected: %v; got: %v",
				test.symbol, test.expected, got)
		}
	}
}

func TestRemoveDuplicates(t *testing.T) {
	var tests = []struct {
		slice    []string
		expected []string
	}{
		{[]string{"A", "B"}, []string{"A", "B"}},
		{[]string{"A", "A"}, []string{"A"}},
		{[]string{""}, []string{""}},
	}
	for _, test := range tests {
		got := removeDuplicates(test.slice)
		if !reflect.DeepEqual(got, test.expected) {
			t.Errorf("removeDuplicates(%#v)\nexpected: %v, got: %v",
				test.slice, test.expected, got)
		}
	}
}

func Example() {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"", "cruzeiro"}
	main()
	// Output:
	// U+20A2	₢	CRUZEIRO SIGN
}

func Example_queryTwoWords() {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"", "cat", "smiling"}
	main()
	// Output:
	// U+1F638	😸	GRINNING CAT FACE WITH SMILING EYES
	// U+1F63A	😺	SMILING CAT FACE WITH OPEN MOUTH
	// U+1F63B	😻	SMILING CAT FACE WITH HEART-SHAPED EYES
}

func Example_withHyphen() {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"", "small", "hyphen"}
	main()
	// Output:
	// U+FE63	﹣	SMALL HYPHEN-MINUS
}

func Example_withUnicode1Name() {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"", "quote"}
	main()
	// Output:
	// U+0027	'	APOSTROPHE (APOSTROPHE-QUOTE)
	// U+2358	⍘	APL FUNCTIONAL SYMBOL QUOTE UNDERBAR
	// U+235E	⍞	APL FUNCTIONAL SYMBOL QUOTE QUAD
}

func ExampleList() {
	text := strings.NewReader(lines3Dto43)
	List(text, "MARK")
	// Output: U+003F	?	QUESTION MARK
}

func ExampleList_twoResults() {
	text := strings.NewReader(lines3Dto43)
	List(text, "SIGN")
	// Output:
	// U+003D	=	EQUALS SIGN
	// U+003E	>	GREATER-THAN SIGN
}

func ExampleList_twoResultsNoOrder() {
	text := strings.NewReader(lines3Dto43)
	List(text, "CAPITAL LATIN")
	// Output:
	// U+0041	A	LATIN CAPITAL LETTER A
	// U+0042	B	LATIN CAPITAL LETTER B
	// U+0043	C	LATIN CAPITAL LETTER C
}

func ExampleList_withUnicode1Name() {
	text := strings.NewReader(lineApostrophe)
	List(text, "QUOTE")
	// Output:
	// U+0027	'	APOSTROPHE (APOSTROPHE-QUOTE)
}
