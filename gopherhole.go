package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"
	"unicode"
)

func main() {

	intro()

	// -------------------------------------------------------------------------
	// TODO: Replace these examples with file input
	// -------------------------------------------------------------------------
	// Example of input XML data
	rawXMLInput := []byte(`
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

	// Example of a config file
	rawConfigInput := []byte(`
		{
			"patients": [
				{
					"id": "<Patients.Patient.ID>",
					"name": "<Patients.Patient.FirstName> <Patients.Patient.LastName>",
					"age": 39
				}
			]
		}
	`)
	// -------------------------------------------------------------------------

	// -------------------------------------------------------------------------
	// Config extraction
	// -------------------------------------------------------------------------
	// Create a map to store the input JSON configuration
	configMap := make(map[string]interface{})

	// Unmarshal the configuration data
	err := json.Unmarshal(rawConfigInput, &configMap)

	if err != nil {
		fmt.Println("Invalid configuration JSON:", err)
	}

	// Hard-coded transformations
	// TODO: Use this to generalize
	// The trigger when we know we found a new patient
	patientCreationTrigger := "Patients.Patient"
	// Each time we encounter an XML key that is a patient creation trigger, add a new patient to this slice
	patientsList := make([]map[string]interface{}, 0)
	// This tells us where in a patient map we will need to find and replace a configuration key
	xmlKeyToJSONKeyMap := map[string]string{"Patients.Patient.FirstName": "name", "Patients.Patient.LastName": "name", "Patients.Patient.ID": "id"}

	// -------------------------------------------------------------------------

	// -------------------------------------------------------------------------
	// XML conversion
	// -------------------------------------------------------------------------
	// Create a slice to track the current XML token key
	// while iterating over XML tokens based on the token's
	// location in the hierarchy of tags
	//
	// Example key/value pair: Patients.Patient.FirstName = Jane
	//
	// Store each piece of the current key as an element of
	// a slice of strings
	xmlKeySlice := []string{}

	// Create an XML decoder
	xmlReader := bytes.NewReader(rawXMLInput)
	decoder := xml.NewDecoder(xmlReader)

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
		case xml.StartElement:

			// Push the new element name to the key slice
			xmlKeySlice = append(xmlKeySlice, t.Name.Local)

			// If we come across the opening tag for a Patients.Patient,
			// we need to create a new entry in the Patients list.
			xmlKey := strings.Join(xmlKeySlice, ".")

			if xmlKey == patientCreationTrigger {
				newPatientMap := generateOutputPatientMap(configMap)
				patientsList = append(patientsList, newPatientMap)
			}

			// If we come across a tracked attribute
			for _, a := range t.Attr {
				xmlKey := strings.Join(xmlKeySlice, ".") + "." + a.Name.Local

				// If the given attribute is tracked in our input JSON configuration
				outputJSONKey, ok := xmlKeyToJSONKeyMap[xmlKey]

				if ok {
					// We need to find and replace the patient key with
					// this XML token's value in the given output patient field
					outputPatientField := patientsList[len(patientsList)-1][outputJSONKey].(string)
					outputPatientField = strings.Replace(outputPatientField, "<"+xmlKey+">", a.Value, 1)
					patientsList[len(patientsList)-1][outputJSONKey] = outputPatientField
				}
			}

		case xml.CharData:

			// If we encounter whitespace, ignore it
			if isWhitespace(string(t)) {
				break
			}

			// If we come across one of the configured patient keys
			xmlKey := strings.Join(xmlKeySlice, ".")
			outputJSONKey, ok := xmlKeyToJSONKeyMap[xmlKey]

			if ok {
				// We need to find and replace the patient key with
				// this XML token's value in the given output patient field
				outputPatientField := patientsList[len(patientsList)-1][outputJSONKey].(string)
				outputPatientField = strings.Replace(outputPatientField, "<"+xmlKey+">", string(t), 1)
				patientsList[len(patientsList)-1][outputJSONKey] = outputPatientField
			}

		case xml.EndElement:
			xmlKeySlice = xmlKeySlice[:len(xmlKeySlice)-1] // Pop the closed element
		case xml.Comment:
		case xml.Directive:
		default:
			fmt.Println("Unhandled token encountered")
		}
	}

	outputJSON := configMap
	outputJSON["patients"] = make([]interface{}, 0)

	for _, p := range patientsList {
		outputJSON["patients"] = append(outputJSON["patients"].([]interface{}), p)
	}

	// Marshal the output to JSON
	jsonData, err := json.MarshalIndent(outputJSON, "", "  ")

	if err != nil {
		fmt.Println("Error marshaling the output JSON:", err)
	}

	fmt.Println("Output JSON")
	fmt.Println(string(jsonData))
}

// Introduce the application
func intro() {
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
}

// -----------------------------------------------------------------------------
// Generation
// -----------------------------------------------------------------------------
// Create a new map for storing an output patient
func generateOutputPatientMap(configMap map[string]interface{}) map[string]interface{} {
	outputPatientMap := make(map[string]interface{})

	for key, value := range configMap["patients"].([]interface{})[0].(map[string]interface{}) {
		outputPatientMap[key] = value
	}

	return outputPatientMap
}

// -----------------------------------------------------------------------------
// Transformation
// -----------------------------------------------------------------------------

// -----------------------------------------------------------------------------
// Utility
// -----------------------------------------------------------------------------

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
