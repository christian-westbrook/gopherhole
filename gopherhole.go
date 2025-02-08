package main

import (
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

	// Map unmarshaled XML tags to their values
	var tagValueMap map[string]interface{}

	err := xml.Unmarshal(xmlData, &tagValueMap)

	if err != nil {
		fmt.Println("Error unmarshaling XML:", err)
	}
}
