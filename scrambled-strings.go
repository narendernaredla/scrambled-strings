package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type Dictionary struct {
	value     string
	charCount map[string]int
}

var dictionaryArr []Dictionary

func main() {
	var msg string
	var err error
	dictionaryFilePtr := flag.String("dictionary", "", "a string")
	inputFilePtr := flag.String("input", "", "a string")
	flag.Parse()
	dictionaryArr, msg, err = validateDictionary(*dictionaryFilePtr)
	if err != nil {
		log.Fatal("Dictonary File Validation Error:", err)
	} else {
		log.Println("Dictonary File Validation:", msg)
		output, inputFileProcessingErr := processInputFile(*inputFilePtr)
		if inputFileProcessingErr != nil {
			log.Fatal("Input File Processing Error:", inputFileProcessingErr)
		} else {
			for _, v := range output {
				fmt.Println(v)
			}
		}
	}
}

func validateDictionary(fileLoc string) (dictionaryArr []Dictionary, msg string, err error) {
	log.Println("Dictionary File Location:", fileLoc)
	file, err := os.Open(fileLoc)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	dicLen := 0
	duplicateVals := make(map[string]int)
	for scanner.Scan() {
		dictionaryVal := fmt.Sprintf(scanner.Text())
		dictionaryLen := len(dictionaryVal)
		dicLen = dicLen + dictionaryLen
		if _, exist := duplicateVals[dictionaryVal]; exist {
			err = errors.New("Dictionary file doesn't match criteria")
			return
		}
		duplicateVals[dictionaryVal] = 1
		if dictionaryLen < 2 || dictionaryLen > 105 || dicLen > 105 {
			err = errors.New("Dictionary file doesn't match criteria")
			return
		}
		var dictionarDet Dictionary
		dictionarDet.value = dictionaryVal
		charCounts := make(map[string]int)
		for _, char := range dictionaryVal {
			dChar := fmt.Sprintf("%c", char)
			if _, exist := charCounts[dChar]; exist {
				charCounts[dChar] += 1
			} else {
				charCounts[dChar] = 1
			}
		}
		dictionarDet.charCount = charCounts
		dictionaryArr = append(dictionaryArr, dictionarDet)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	msg = "Successfully validate dictionary file"
	return
}

func processInputFile(inputFileLoc string) (results []string, err error) {
	log.Println("Input File Location:", inputFileLoc)
	file, err := os.Open(inputFileLoc)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	counter := 1
	finalRes := make([]string, 0)
	for scanner.Scan() {
		inputString := fmt.Sprintf(scanner.Text())
		wordCount := checkScrembledWord(inputString)
		s := fmt.Sprintf("Case #%v: %v", counter, wordCount)
		finalRes = append(finalRes, s)
		counter += 1
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return finalRes, nil
}

func checkScrembledWord(input string) (wCount int) {
	wCount = 0
	for _, v := range dictionaryArr {
		res := strings.Contains(input, v.value)
		if res {
			wCount += 1
		} else {
			result := findScrembledWordCount(input, v)
			if result {
				wCount += 1
			}
		}
	}
	return
}

func findScrembledWordCount(input string, dictionaryDet Dictionary) bool {
	dictionary := dictionaryDet.value
	dicLen := len(dictionary) - 1
	firstChar := fmt.Sprintf("%s", dictionary[0:1])
	lastChar := fmt.Sprintf("%s", dictionary[dicLen:])
	inputLen := len(input)
	for pos, char := range input {
		stChar := fmt.Sprintf("%c", char)
		endCharAt := pos + dicLen
		if endCharAt > inputLen {
			return false
		}
		if stChar != firstChar {
			continue
		}
		endChar := fmt.Sprintf("%c", input[endCharAt])
		if endChar != lastChar {
			continue
		}
		charCounts := make(map[string]int)
		for _, char := range input[pos : endCharAt+1] {
			dChar := fmt.Sprintf("%c", char)
			if _, exist := charCounts[dChar]; exist {
				charCounts[dChar] += 1
			} else {
				charCounts[dChar] = 1
			}
		}
		stringMatched := true
		for char, count := range dictionaryDet.charCount {
			if ct, exist := charCounts[char]; exist {
				if ct != count {
					stringMatched = false
				}
			} else {
				stringMatched = false
			}
		}
		return stringMatched
	}
	return false
}
