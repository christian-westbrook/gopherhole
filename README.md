# gopherhole

![Build & test](https://github.com/christian-westbrook/gopherhole/workflows/test.yaml/badge.svg)

gopherhole is a Go application that converts XML to JSON with user-defined transformations.  
  

```
    Throw your XML into the hole, and the Gophers will toss back JSON!  

        +--------------------+
	| Was that XML raw?  |
	+--------------------+
	  \\
	   \\
	    \\
	         ,_---~~~~~----._         
	  _,,_,*^____      _____``*g*\"*, 
	 / __/ /'     ^.  /      \\ ^@q   f 
	[  @f | @))    |  | @))   l  0 _/  
	 \`/   \~____ / __ \\_____/    \   
	 |           _l__l_           I   
	 }          [______]           I  
	 ]            | | |            |  
	 ]             ~ ~             |  
	 |                            |   
	  |                           |   

    Developed by Christian Westbrook
    https://github.com/christian-westbrook/

    Artwork by belbomemo
    https://gist.github.com/belbomemo
```

# Table of Contents

[Installation](#installation)  
[Configuration](#configuration)  
[Usage](#usage)  
[Limitations & Roadmap](#limitations--roadmap)  

# Installation

### Download
Download the latest version for your system from the [releases page](https://github.com/christian-westbrook/gopherhole/releases)  
Windows 11 and x86 Linux are currently supported
Ensure that you have a configuration file
- You can place a configuration file named `config.json` in the same directory as the gopherhole executable
- You can also pass in the location of a configuration file at runtime, check out the Usage section
- A default config file and example of input XML is bundled with the release

### Build
Clone the repository  
Build an executable for your machine with `go build`
Ensure that you have a configuration file  
- You can place a configuration file named `config.json` in the same directory as the gopherhole executable
- You can also pass in the location of a configuration file at runtime, check out the Usage section
- A default config file and example of input XML is included when you clone the repository  

Specifically if you're on Windows 11, the Makefile included with the repository can be used to build the executable with the command `make`

# Configuration

The transformation from XML to JSON is controlled by a config file that uses valid JSON to specify the output format. This control is made possible by the embedding of find and replace symbols within the config file that specify which fields in the input XML file should be used to populate fields in the output JSON data.  

Find and replace symbols look like `<Patients.Patient.DateOfBirth transform=yearsElapsed>` where `Patients` is the name of a collection of objects, `Patient` is the name of an object within that collection, and `DateOfBirth` is a field within that object. Given this find and replace symbol, the value of the XML tag pair `<DateOfBirth>value</DateOfBirth>` located within a `<Patient></Patient>` XML tag pair that is itself located within a `<Patients></Patients>` XML tag pair would be used to replace this symbol within the output JSON file.  

The notation `transform=yearsElapsed` is a modifier that can be used to specify pre-defined transformations to the given value at runtime. In this example, `transform=yearsElapsed` indicates that the `<DateOfBirth>value</DateOfBirth>` value should be transformed into the count of
years that have elapsed since that date before being added to the output JSON data. Modifiers add another layer of flexibility to gopherhole's transformations and are easy to contribute and make use of.  

If not specified on the command line, the default configuration file should be called `config.json` and be placed alongside the gopherhole executable.

### Simplifying Assumptions  
To enable a flexible and expressive range of object definitions, gopherhole currently makes the simplifying assumption that your XML file is organized as a list of collection keys mapped to lists of object definitions.  

### Examples

**Configuration file**  
```
{
    "Patients": [
        {
            "id": "<Patients.Patient.ID>",
            "name": "<Patients.Patient.FirstName> <Patients.Patient.LastName>",
            "age": "<Patients.Patient.DateOfBirth transform=yearsElapsed>"
        }
    ],

    "Doctors": [
        {
            "id": "<Doctors.Doctor.ID>",
            "first name": "<Doctors.Doctor.FirstName>",
            "last name": "<Doctors.Doctor.LastName>",
            "date of birth": "<Doctors.Doctor.DateOfBirth>"
        }
    ]
}
```

In this example `<Patients.Patient.FirstName>` refers to the value of a `<FirstName>value</FirstName>` XML tag pair located within a `<Patient></Patient>` tag pair that is itself located within a `<Patients></Patients>` tag pair.

Also in this example, `<Patients.Patient.ID>` refers to the value of the attribute `ID` defined in a `<Patient></Patient>` tag pair that is itself located within a `<Patients></Patients>` tag pair.

The find and replace symbol `<Patients.Patient.DateOfBirth transform=yearsElapsed>` contains the modifier `transform=yearsElapsed` that will transform the given date of birth into the number of years that have elapsed since the given date. This will produce different output than the symbol `<Doctors.Doctor.DateOfBirth>` for which there is no transformation.

**Example Output**

```
{
  "Doctors": [
    {
      "date of birth": "1985-07-15",
      "first name": "Ada",
      "id": "12345",
      "last name": "Lovelace"
    },
    {
      "date of birth": "1992-03-22",
      "first name": "Alan",
      "id": "67890",
      "last name": "Turing"
    },
    {
      "date of birth": "1992-03-22",
      "first name": "Stephen",
      "id": "67890",
      "last name": "Hawking"
    }
  ],
  "Patients": [
    {
      "age": "39",
      "id": "12345",
      "name": "John Doe"
    },
    {
      "age": "32",
      "id": "67890",
      "name": "Jane Smith"
    }
  ]
}
```

# Usage

### Example Input XML
```
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

<Doctors>
    <Doctor ID="12345">
        <FirstName>Ada</FirstName>
        <LastName>Lovelace</LastName>
        <DateOfBirth>1985-07-15</DateOfBirth>
    </Doctor>

    <Doctor ID="67890">
        <FirstName>Alan</FirstName>
        <LastName>Turing</LastName>
        <DateOfBirth>1992-03-22</DateOfBirth>
    </Doctor>

    <Doctor ID="67890">
        <FirstName>Stephen</FirstName>
        <LastName>Hawking</LastName>
        <DateOfBirth>1992-03-22</DateOfBirth>
    </Doctor>
</Doctors>
```

### Example Config File
```
{
    "Patients": [
        {
            "id": "<Patients.Patient.ID>",
            "name": "<Patients.Patient.FirstName> <Patients.Patient.LastName>",
            "age": "<Patients.Patient.DateOfBirth transform=yearsElapsed>"
        }
    ],

    "Doctors": [
        {
            "id": "<Doctors.Doctor.ID>",
            "first name": "<Doctors.Doctor.FirstName>",
            "last name": "<Doctors.Doctor.LastName>",
            "date of birth": "<Doctors.Doctor.DateOfBirth>"
        }
    ]
}
```

### Example Execution

```
gopherhole             <- defaults to converting input.xml using config.json
gopherhole myxmlfile.xml                    <- defaults to using config.json
gopherhole myxmlfile.xml myconfigfile.json
```

### Example Output

```
{
  "Doctors": [
    {
      "date of birth": "1985-07-15",
      "first name": "Ada",
      "id": "12345",
      "last name": "Lovelace"
    },
    {
      "date of birth": "1992-03-22",
      "first name": "Alan",
      "id": "67890",
      "last name": "Turing"
    },
    {
      "date of birth": "1992-03-22",
      "first name": "Stephen",
      "id": "67890",
      "last name": "Hawking"
    }
  ],
  "Patients": [
    {
      "age": "39",
      "id": "12345",
      "name": "John Doe"
    },
    {
      "age": "32",
      "id": "67890",
      "name": "Jane Smith"
    }
  ]
}
```

# Limitations & Roadmap

### Limitations
To enable a flexible and expressive range of object definitions, gopherhole currently makes the simplifying assumption that your XML file is organized as a list of collection keys mapped to lists of object definitions, e.g.

```
<Patients>
    <Patient ID="12345">
        <FirstName>John</FirstName>
        <LastName>Doe</LastName>
        <DateOfBirth>1985-07-15</DateOfBirth>
    </Patient>
</Patients>
```

Find and replace symbols can currently only be replaced once each per object.


### Roadmap
- Expanded test coverage
- Adding the ability to export the output JSON to a file
- Adding the ability to pass a config file as a flag option rather than as a command-line argument
- Adding the ability to pass the desired output file path as a flag option
- Adding support for collection key alias' e.g. `<Patients alias=patients>` becoming `patients`
- Adding support for the same find and replace symbol occuring multiple times in one object
- Adding support for numbers in the output JSON, i.e. the age `39` being represented as `39` rather than `"39"`
- Support for Linux systems in the Makefile
- Instructions for contributing new transformations