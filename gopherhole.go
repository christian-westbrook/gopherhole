package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"unicode"
)

func main() {

	// Introduction
	fmt.Println()

	// Example of input XML data
	xmlData := []byte(`
		<?xml version="1.0" encoding="UTF-8"?>
		<Patients>
			<Patient ID="12345">
				<FirstName>John</FirstName>
				<LastName>Doe</LastName>
				<DateOfBirth>1985-07-15</DateOfBirth>
			</Patient>

			<Patient ID="67890">
				<FirstName>Jane</FirstName>
				<LastName>Smith</LastName>
				<DateOfBirth>1992-03-22</DateOfBirth>
			</Patient>

		</Patients>
	`)

	// Create an XML decoder
	xmlDataReader := bytes.NewReader(xmlData)
	decoder := xml.NewDecoder(xmlDataReader)

	// Iterate over tokens in the decoder
	for {

		// Unpack the next token
		token, error := decoder.Token()

		if error != nil {
			break
		}

		// Switch on the token's asserted type
		switch t := token.(type) {
		case xml.ProcInst:
			fmt.Println("Processing instruction")
			fmt.Println("Target:", t.Target)
			fmt.Println("Instruction:", string(t.Inst))
		case xml.StartElement:
			fmt.Println("Start:", t.Name.Local)
		case xml.CharData:

			// If we encounter whitespace, ignore it
			if isWhitespace(string(t)) {
				break
			}

			fmt.Println("Payload:", string(t))
		case xml.EndElement:
			fmt.Println("End:", t.Name.Local)
		case xml.Comment:
			fmt.Println("Ignoring comment:", t)
		case xml.Directive:
			fmt.Println("Ignoring directive:", t)
		default:
			fmt.Println("Unhandled token encountered")
		}
	}
}

// Determine whether a given string is alphanumeric
func isAlphaNumeric(s string) bool {
	// For each rune in the input string
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}

	return true
}

// Determine whether a given string is whitespace
func isWhitespace(s string) bool {

	// For each rune in the input string
	for _, r := range s {
		if !unicode.IsSpace(r) {
			return false
		}
	}

	return true
}
