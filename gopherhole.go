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
	"time"
	"unicode"
)

// OUTSTANDING FEATURES
// - Handle the same find and replace key occuring multiple times in the config file
// - Handle time transformations
// - Handle parent key alias'
// - Convert numbers to numbers (rather than strings)

// Package level constants
const FindAndReplaceExpression = "<[a-zA-Z.=\\s]+>" // Regex for use in replacing the config file's find and replace symbols

// -----------------------------------------------------------------------------
// Function : main()
// Input    :
// input.xml   - An XML file to be configurably converted into JSON
// config.json - A configuration file that uses replacement symbols to specify
// an output JSON file format
//
// Output       : Raw JSON data printed to the console
// Side Effects : none
//
// Abstract :
// This function serves as the entry point to Gopher Hole. It reads
// in XML data from input.xml and uses the configuration file
// config.json to convert the input XML data into a specified output
// JSON format.
// -----------------------------------------------------------------------------
func main() {

	// Introduce the application
	intro()

	inputXMLPath := "input.xml" // Default input file name

	// Get command-line arguments
	if len(os.Args) > 1 {
		inputXMLPath = os.Args[1]
	}

	fmt.Println("Processing", inputXMLPath)
	fmt.Println()

	// -------------------------------------------------------------------------
	// OPEN INPUT FILES
	// TODO: Replace these hardcoded file names with command-line arguments
	// -------------------------------------------------------------------------
	rawXMLInput, err := os.ReadFile(inputXMLPath)

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
	// READ CONFIG FILE
	// -------------------------------------------------------------------------
	// Create a map to store the input JSON configuration
	configMap := make(map[string]interface{})

	// Unmarshal configuration data into the map
	err = json.Unmarshal(rawConfigInput, &configMap)

	if err != nil {
		fmt.Println("Invalid configuration JSON:", err)
	}

	// Build a record of find and replace symbols
	findAndReplaceMaps := make(map[string]map[string]string)

	// Build a record of modifiers
	modifierMaps := make(map[string]map[string]map[string]string)

	// Map
	for k, v := range configMap {
		innerMap := v.([]interface{})[0].(map[string]interface{})
		findAndReplaceMap, modifierMap := getReplacementMapAndModifiers(innerMap)

		findAndReplaceMaps[k] = findAndReplaceMap
		modifierMaps[k] = modifierMap
	}

	// Build a record of modifiers

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
	// READ XML
	// -------------------------------------------------------------------------
	// As we iterate over XML tokens, use this key to keep track of where we
	// are in the hierarchy of tags
	//
	// Example key    : Patients.Patient.FirstName
	// Representation : ["Patients", "Patient", "FirstName"]
	xmlKeySlice := []string{}

	// Create an XML decoder
	xmlReader := bytes.NewReader(rawXMLInput)
	decoder := xml.NewDecoder(xmlReader)

	// Iterate over tokens in the XML decoder
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
					parentKey := xmlKeySlice[0]
					xmlKey := strings.Join(xmlKeySlice, ".") + "." + a.Name.Local

					// If the given attribute is related to a find and replace symbol
					outputFieldKey, ok := findAndReplaceMaps[parentKey][xmlKey]

					if ok {
						// We need to find and replace the patient key with
						// this XML token's value in the given output patient field
						outputField := outputObjects[len(outputObjects)-1][outputFieldKey].(string)
						outputField = strings.Replace(outputField, "<"+xmlKey+">", a.Value, 1)
						outputObjects[len(outputObjects)-1][outputFieldKey] = outputField
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
			if IsWhitespace(string(t)) {
				break
			}

			// If we come across one of the configured patient keys
			parentKey := xmlKeySlice[0]
			xmlKey := strings.Join(xmlKeySlice, ".")

			outputJSONKey, ok := findAndReplaceMaps[parentKey][xmlKey]

			outputObjects := parentKeyMap[xmlKeySlice[0]]

			if ok {
				// We need to find and replace the patient key with
				// this XML token's value in the given output patient field
				fieldValue := string(t)

				// First, check if there are any transformations necessary
				transformation, ok := modifierMaps[parentKey][xmlKey]["transform"]

				if ok {
					switch transformation {
					case "yearsElapsed":
						fvPtr := &fieldValue // Get a pointer to fieldValue
						*fvPtr = strconv.Itoa(yearsElapsed(string(t)))
						xkPtr := &xmlKey
						*xkPtr = xmlKey + " transform=" + transformation
					default:
						fmt.Println("Unhandled transformation: ", transformation)
					}
				}

				// Assign the token's value to the output field
				outputField := outputObjects[len(outputObjects)-1][outputJSONKey].(string)
				outputField = strings.Replace(outputField, "<"+xmlKey+">", fieldValue, 1)
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
	// -------------------------------------------------------------------------

	// -------------------------------------------------------------------------
	// CONVERSION TO JSON
	// -------------------------------------------------------------------------
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
	fmt.Println()
	fmt.Println(string(jsonData))
	// -------------------------------------------------------------------------
}

// -----------------------------------------------------------------------------
// TRANSFORMATIONS
// -----------------------------------------------------------------------------
func yearsElapsed(dateOfBirth string) int {

	dob, err := time.Parse("2006-01-02", dateOfBirth)

	if err != nil {
		fmt.Println("Failed to parse input date: ", err)
	}

	now := time.Now()
	age := now.Year() - dob.Year()

	// Handle birthdays
	if now.YearDay() < dob.YearDay() {
		age--
	}

	return age
}

// -----------------------------------------------------------------------------
// DATA STRUCTURES
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

// Generate a map of replacement tokens to the output field where you can find them
func getReplacementMapAndModifiers(m map[string]interface{}) (map[string]string, map[string]map[string]string) {

	findAndReplaceRegex := regexp.MustCompile(FindAndReplaceExpression)

	findAndReplaceMap := map[string]string{}
	modifierMap := make(map[string]map[string]string)

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
			fmt.Println("Unhandled configuration value type encountered:", k, v)
		}

		// Search for any replacement symbols
		matches := findAndReplaceRegex.FindAllString(v.(string), -1)

		// Assign any discovered replacement symbols to the find and replace map
		if matches != nil {
			for _, match := range matches {
				name, modifiers := ParseFindAndReplaceSymbol(match)
				findAndReplaceMap[name] = k
				modifierMap[name] = modifiers
			}
		}

	}

	return findAndReplaceMap, modifierMap
}

// -----------------------------------------------------------------------------

// -----------------------------------------------------------------------------
// UTILITY
// -----------------------------------------------------------------------------
// -----------------------------------------------------------------------------
// Function     : intro()
// Input        : none
// Output       : none
// Side Effects : Prints an introductory message to the console
//
// Abstract :
// This function introduces the application by printing a message to the screen
//
// TODO: Replace hardcoded version string with a constant to be replaced in
// a deployment pipeline
// -----------------------------------------------------------------------------
func intro() {
	// Introduction
	fmt.Println()
	fmt.Println("Welcome to the Gopher Hole v0.1.0!")
	fmt.Println()
	fmt.Println("Throw your XML into the hole, and the")
	fmt.Println("Gophers will toss back JSON!")
	fmt.Println()

	fmt.Println("+--------------------+")
	fmt.Println("| Was that XML raw?  |")
	fmt.Println("+--------------------+")
	fmt.Println("  \\")
	fmt.Println("   \\")
	fmt.Println("    \\")
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

// -----------------------------------------------------------------------------
// Function     : IsWhitespace()
// Input        : s - A string to check for any non-whitespace characters
// Output       :
// A boolean value representing whether the input string consists only of
// whitespace characters, e.g. true indicates a string of nothing but whitespace
//
// Side Effects : none
//
// Abstract :
// This function determines whether a given string consists only of whitespace
// characters, e.g. true indicates a string of nothing but whitespace
// -----------------------------------------------------------------------------
func IsWhitespace(s string) bool {

	// For each rune in the input string
	for _, r := range s {
		if !unicode.IsSpace(r) {
			return false
		}
	}

	return true
}

// -----------------------------------------------------------------------------

// -----------------------------------------------------------------------------
// Function     : ParseFindAndReplaceSymbol
//
// Input :
// s - A string representing a find and replace symbol in the config.json file
// Example: <Patients.Patient.DateOfBirth transform=yearsElapsed>
//
// Output :
// name - A string representing the name of the given find and replace symbol
// modifiers - A map of strings representing modifier names to strings
// representing modifiers
//
// Side Effects : none
//
// Abstract :
// This function takes in a find and replace symbol from the config.json file
// and breaks it into a name and a map of modifiers.
//
// Example: For <Patients.Patient.DateOfBirth transform=yearsElapsed>, name
// is 'Patients.Patient.DateOfBirth'
//
// Example: For <Patients.Patient.DateOfBirth transform=yearsElapsed>, modifiers
// contains the key 'transform' mapped to the value 'yearsElapsed'
// Note: In the case where there is no modifier in a find and replace symbol,
// modifiers is a nil value
// -----------------------------------------------------------------------------
func ParseFindAndReplaceSymbol(s string) (string, map[string]string) {
	tokens := strings.Split(s[1:len(s)-1], " ")
	name := tokens[0]

	// If there aren't any modifiers, return early
	if len(tokens) == 1 {
		return name, nil
	}

	modifiers := map[string]string{}

	for _, t := range tokens[1:] {
		tokens := strings.Split(t, "=")
		modifiers[tokens[0]] = tokens[1]
	}

	return name, modifiers
}

// -----------------------------------------------------------------------------
