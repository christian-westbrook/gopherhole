package main

import "fmt"

func main() {
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

	fmt.Println(xmlData)
}
