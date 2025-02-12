package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// Package level constants
const FindAndReplaceExpression = "<[a-zA-Z.]+>" // Regex for use in replacing the find and replace symbols

func main() {

	intro()

	// -------------------------------------------------------------------------
	// TODO: Replace these examples with file input
	// -------------------------------------------------------------------------
	rawXMLInput, err := os.ReadFile("input.xml")

	if err != nil {
		fmt.Println("Error opening the input XML file:", err)
	}

	// Example of a config file
	rawConfigInput, err := os.ReadFile("config.json")

	if err != nil {
		fmt.Println("Error opening the config file:", err)
	}
	// -------------------------------------------------------------------------

	// -------------------------------------------------------------------------
	// Config extraction
	// -------------------------------------------------------------------------
	// Create a map to store the input JSON configuration
	configMap := make(map[string]interface{})

	// Unmarshal the configuration data
	err = json.Unmarshal(rawConfigInput, &configMap)

	if err != nil {
		fmt.Println("Invalid configuration JSON:", err)
	}

	// Hard-coded transformations
	// TODO: Use this to generalize

	// TODO: If you find all of the inner maps, you could automate replacement symbols
	// beyond Patients.Patient tags
	// TODO: Change from Patient to patients
	innerMap := configMap["Patients"].([]interface{})[0].(map[string]interface{})
	findAndReplaceMap := generateFindAndReplaceMap(innerMap)

	// If we encounter the top level, i.e. Patients
	// Then we need to create a list of objects
	// How will we store active lists?
	// By mapping a string representing the list to the list
	//
	// Example: Patient -> List of maps
	parentKeyMap := make(map[string][]map[string]interface{})

	// If we encounter an object within an existing list, i.e. Patients.Patient,
	// we need to create the object and add it to the list
	// How will we store objects?
	// As anonymous dictionaries within parent lists

	// If we encounter a tracked field within an existing object, we need to
	// correctly find and replace it
	// How will we find the right field?
	// We'll index into the right map of lists and then into the right field
	// using the findAndReplaceMap
	// We can assume that the most recently added object is the right object

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

			// We need to track how deep into the hierarchy we are at this point
			switch len(xmlKeySlice) {
			case 1: // If we encounter a new parent key, initialize it with an empty map
				parentKeyMap[xmlKeySlice[0]] = []map[string]interface{}{}
			case 2: // If we encounter a new object within a parent, add an empty object to the parent list

				// Confirm that the parent exists
				_, ok := parentKeyMap[xmlKeySlice[0]]

				if !ok {
					fmt.Println("Came across an object for which there was no parent key")
					continue
				}

				// Generate a map to contain the new object
				xmlKey := strings.Join(xmlKeySlice, ".")
				outputObjectMap := generateOutputObjectMap(configMap, xmlKey)

				// Add the new map to the list of maps
				outputObjects := parentKeyMap[xmlKeySlice[0]]
				outputObjects = append(outputObjects, outputObjectMap)
				parentKeyMap[xmlKeySlice[0]] = outputObjects

				// If we come across a tracked attribute
				for _, a := range t.Attr {
					xmlKey := strings.Join(xmlKeySlice, ".") + "." + a.Name.Local

					// If the given attribute is tracked in our input JSON configuration
					outputJSONKey, ok := findAndReplaceMap[xmlKey]

					if ok {
						// We need to find and replace the patient key with
						// this XML token's value in the given output patient field
						outputField := outputObjects[len(outputObjects)-1][outputJSONKey].(string)
						outputField = strings.Replace(outputField, "<"+xmlKey+">", a.Value, 1)
						outputObjects[len(outputObjects)-1][outputJSONKey] = outputField
					}
				}

			case 3:
			default:
				// TODO: You could just look at the last three tiers
				// and replicate the upper tiers
				fmt.Println("Unhandled XML hierarchy")
			}

		case xml.CharData:

			// If we encounter whitespace, ignore it
			if isWhitespace(string(t)) {
				break
			}

			// If we come across one of the configured patient keys
			xmlKey := strings.Join(xmlKeySlice, ".")
			outputJSONKey, ok := findAndReplaceMap[xmlKey]

			outputObjects := parentKeyMap[xmlKeySlice[0]]

			if ok {
				// We need to find and replace the patient key with
				// this XML token's value in the given output patient field
				outputField := outputObjects[len(outputObjects)-1][outputJSONKey].(string)
				outputField = strings.Replace(outputField, "<"+xmlKey+">", string(t), 1)
				outputObjects[len(outputObjects)-1][outputJSONKey] = outputField
			}

		case xml.EndElement:
			xmlKeySlice = xmlKeySlice[:len(xmlKeySlice)-1] // Pop the closed element
		case xml.Comment:
		case xml.Directive:
		default:
			fmt.Println("Unhandled token encountered")
		}
	}

	// For each list of objects
	// parentKeyMap := make(map[string][]map[string]interface{})
	outputJSON := make(map[string][]map[string]interface{})
	for k, v := range parentKeyMap {
		outputJSON[k] = v
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

func generateOutputObjectMap(configMap map[string]interface{}, xmlKey string) map[string]interface{} {

	// Create a map to represent the new output object
	outputObjectMap := make(map[string]interface{})

	// Split the input xmlKey into its tokens
	// 0: Parent slice key, i.e. Patients
	// 1: Object key, i.e. Patients.Patient
	tokens := strings.Split(xmlKey, ".")

	// Get the definition of this object type
	parentKey := tokens[0]
	objectMapConfig := configMap[parentKey].([]interface{})[0]

	for k, v := range objectMapConfig.(map[string]interface{}) {
		outputObjectMap[k] = v
	}

	return outputObjectMap
}

// Generate a map of replacement tokens and where to find them
func generateFindAndReplaceMap(m map[string]interface{}) map[string]string {
	findAndReplaceRegex := regexp.MustCompile(FindAndReplaceExpression)

	findAndReplaceMap := map[string]string{}

	// For each field in the input map
	for k, v := range m {

		switch v.(type) {
		case string:
		case float64:
			// By default a straight assignment to v shadows the outer variable
			// This allows the value string(v) to persist beyond this case block
			f, ok := v.(float64)

			if ok {
				str := strconv.FormatFloat(f, 'f', -1, 64) // Convert JSON float to a string
				vPtr := &v                                 // Get a reference to v
				*vPtr = str                                // Assign the new string to v in a way that will persist beyond this block
			}

		default:
			fmt.Println("Unhandled configuration value type encountered")
		}

		// Search for any replacement symbols
		matches := findAndReplaceRegex.FindAllString(v.(string), -1)

		// Assign any discovered replacement symbols to the find and replace map
		if matches != nil {
			for _, match := range matches {
				findAndReplaceMap[match[1:len(match)-1]] = k
			}
		}

	}

	return findAndReplaceMap
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
