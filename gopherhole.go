package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"
	"unicode"
)

func main() {

	// Introduction
	fmt.Println()
	fmt.Println("Welcome to the Gopher Hole v0.1.0!")
	fmt.Println()
	fmt.Println("Throw your XML into the hole, and the")
	fmt.Println("Gophers will toss back JSON!")
	fmt.Println()

	fmt.Println("         ,_---~~~~~----._         ")
	fmt.Println("  _,,_,*^____      _____``*g*\"*, ")
	fmt.Println(" / __/ /'     ^.  /      \\ ^@q   f ")
	fmt.Println("[  @f | @))    |  | @))   l  0 _/  ")
	fmt.Println(" \\`/   \\~____ / __ \\_____/    \\   ")
	fmt.Println(" |           _l__l_           I   ")
	fmt.Println(" }          [______]           I  ")
	fmt.Println(" ]            | | |            |  ")
	fmt.Println(" ]             ~ ~             |  ")
	fmt.Println(" |                            |   ")
	fmt.Println("  |                           |   ")
	fmt.Println()
	fmt.Println("Developed by Christian Westbrook ")
	fmt.Println("https://github.com/christian-westbrook/")
	fmt.Println()

	fmt.Println("Artwork by belbomemo")
	fmt.Println("https://gist.github.com/belbomemo")
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

	// Create a slice to track the current XML token key
	// based on the token's location in the hierarchy
	//
	// Example key/value pair: Patients.Patient.FirstName = Jane
	//
	// Store each piece of the current key as an element of
	// a slice of strings
	xmlKeySlice := []string{}

	// Create a mapping from an XML key to its value
	xmlKeyMap := map[string]string{}

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
			xmlKeySlice = append(xmlKeySlice, t.Name.Local) // Push the new element
		case xml.CharData:

			// If we encounter whitespace, ignore it
			if isWhitespace(string(t)) {
				break
			}

			fmt.Println("Payload:", string(t))

			// Add this value to the map
			xmlKey := strings.Join(xmlKeySlice, ".")
			xmlKeyMap[xmlKey] = string(t)
		case xml.EndElement:
			fmt.Println("End:", t.Name.Local)
			xmlKeySlice = xmlKeySlice[:len(xmlKeySlice)-1] // Pop the closed element
		case xml.Comment:
			fmt.Println("Ignoring comment:", t)
		case xml.Directive:
			fmt.Println("Ignoring directive:", t)
		default:
			fmt.Println("Unhandled token encountered")
		}
	fmt.Println(configMap["patients"].([]interface{})[0].(map[string]interface{})["age"])
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
