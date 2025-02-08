package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
)

func main() {

	// Define an example of input XML data
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

	// Map for storing unmarshaled XML tags and their values
	//var tagValueMap map[string]interface{}

	// Iterate over tokens in the decoder
	for {

		// Unpack the next token
		token, error := decoder.Token()

		if error != nil {
			break
		}

		// Switch on the token's asserted type
		switch t := token.(type) {
		case xml.StartElement:
			fmt.Println("Start of element:", t.Name.Local)
		case xml.CharData:
			fmt.Println("Payload of element:", string(t))
		case xml.EndElement:
			fmt.Println("End of element:", t.Name.Local)
		default:
			fmt.Println("Unhandled token encountered")
		}
	}
}
