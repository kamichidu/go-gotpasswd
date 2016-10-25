package main

import (
	"crypto/rand"
	"errors"
	"flag"
	"fmt"
	"math/big"
	"os"
	"strings"
	"unicode"
)

var (
	kinds  = flag.String("k", "alphabet,number,symbol,underscore,space", "Character kinds")
	length = flag.Int("l", 8, "Length of password")
	num    = flag.Int("n", 1, "Number of passwords")
	debug  = flag.Bool("debug", false, "DO NOT USE THIS")
)

type CharacterKind int

const (
	ALPHABET CharacterKind = iota
	NUMBER
	SYMBOL
	UNDERSCORE
	SPACE
)

var (
	dict map[CharacterKind]([]rune)
)

func init() {
	dict = make(map[CharacterKind]([]rune))
	// Auto generate dictionary using ascii printable characters
	for code := 0x20; code <= 0x7e; code++ {
		r := rune(code)
		if !unicode.IsPrint(r) {
			panic("Internal error, cannot construct character dictionary")
		}

		// if r == '_' {
		// 	fmt.Printf("unicode.IsControl('%c') = %v\n", r, unicode.IsControl(r))
		// 	fmt.Printf("unicode.IsDigit('%c') = %v\n", r, unicode.IsDigit(r))
		// 	fmt.Printf("unicode.IsGraphic('%c') = %v\n", r, unicode.IsGraphic(r))
		// 	fmt.Printf("unicode.IsLetter('%c') = %v\n", r, unicode.IsLetter(r))
		// 	fmt.Printf("unicode.IsLower('%c') = %v\n", r, unicode.IsLower(r))
		// 	fmt.Printf("unicode.IsMark('%c') = %v\n", r, unicode.IsMark(r))
		// 	fmt.Printf("unicode.IsNumber('%c') = %v\n", r, unicode.IsNumber(r))
		// 	fmt.Printf("unicode.IsPrint('%c') = %v\n", r, unicode.IsPrint(r))
		// 	fmt.Printf("unicode.IsPunct('%c') = %v\n", r, unicode.IsPunct(r))
		// 	fmt.Printf("unicode.IsSpace('%c') = %v\n", r, unicode.IsSpace(r))
		// 	fmt.Printf("unicode.IsSymbol('%c') = %v\n", r, unicode.IsSymbol(r))
		// 	fmt.Printf("unicode.IsTitle('%c') = %v\n", r, unicode.IsTitle(r))
		// 	fmt.Printf("unicode.IsUpper('%c') = %v\n", r, unicode.IsUpper(r))
		// }

		switch {
		case unicode.IsLetter(r):
			dict[ALPHABET] = append(dict[ALPHABET], r)
		case unicode.IsNumber(r):
			dict[NUMBER] = append(dict[NUMBER], r)
		case unicode.IsSymbol(r):
			dict[SYMBOL] = append(dict[SYMBOL], r)
		case unicode.IsSpace(r):
			dict[SPACE] = append(dict[SPACE], r)
		case r == '_':
			dict[UNDERSCORE] = append(dict[UNDERSCORE], r)
		}
	}
}

type Config struct {
	Kinds  []CharacterKind
	Length int
	Num    int
}

func (self *Config) ParseKinds(s string) ([]CharacterKind, error) {
	kinds := make([]CharacterKind, 0)
	for _, candidate := range strings.Split(s, ",") {
		switch candidate {
		case "alphabet":
			kinds = append(kinds, ALPHABET)
		case "number":
			kinds = append(kinds, NUMBER)
		case "symbol":
			kinds = append(kinds, SYMBOL)
		case "underscore":
			kinds = append(kinds, UNDERSCORE)
		case "space":
			kinds = append(kinds, SPACE)
		default:
			return kinds, errors.New(fmt.Sprintf("Unknown character kind: %s", candidate))
		}
	}
	return kinds, nil
}

func Generate(config *Config) (string, error) {
	charCandidates := make([]rune, 0)
	for _, kindIndex := range config.Kinds {
		charCandidates = append(charCandidates, dict[kindIndex]...)
	}

	if len(charCandidates) == 0 {
		return "", errors.New("Internal error, cannot work with empty candidates")
	}

	chars := make([]rune, config.Length)
	i := 0
	for i < config.Length {
		charIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charCandidates))))
		if err != nil {
			return "", err
		}
		chars[i] = charCandidates[charIndex.Int64()]
		i++
	}
	return string(chars), nil
}

func _main() int {
	flag.Parse()

	if *debug {
		fmt.Fprintf(os.Stderr, "alphabet chars: %v\n", dict[ALPHABET])
		fmt.Fprintf(os.Stderr, "number chars: %v\n", dict[NUMBER])
		fmt.Fprintf(os.Stderr, "symbol chars: %v\n", dict[SYMBOL])
		fmt.Fprintf(os.Stderr, "underscore chars: %v\n", dict[UNDERSCORE])
		fmt.Fprintf(os.Stderr, "space chars: %v\n", dict[SPACE])
	}

	config := &Config{}
	if parsed, err := config.ParseKinds(*kinds); err == nil {
		config.Kinds = parsed
	} else {
		fmt.Fprintln(os.Stderr, err)
		return 128
	}
	if *length > 0 {
		config.Length = *length
	} else {
		fmt.Fprintln(os.Stderr, "Length of password must be positive")
		return 128
	}
	if *num > 0 {
		config.Num = *num
	} else {
		fmt.Fprintln(os.Stderr, "Number of passwords must be positive")
		return 128
	}

	for i := 0; i < config.Num; i++ {
		passwd, err := Generate(config)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		fmt.Println(passwd)
	}

	return 0
}

func main() {
	os.Exit(_main())
}
