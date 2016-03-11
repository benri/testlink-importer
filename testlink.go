package main

import (
	"encoding/csv"
	"encoding/xml"
	"flag"
	"fmt"
	// "math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const AppVersion = "0.0.1"

func printError(e error) {
	fmt.Println("Error:", e)
}

type Testcases struct {
	XMLName       xml.Name    `xml:"testcases"`
	TestcaseList  []Testcase  `xml:"testcase"`
}

type Testcase struct {
	XMLName 			xml.Name  		`xml:"testcase"`
	Name    			string    		`xml:"name,attr"`
	Summary 			string    		`xml:"summary"`
	Steps					Steps					`xml:"steps"`
	CustomFields	CustomFields	`xml:"custom_fields"`
}

type Steps struct {
	XMLName   xml.Name  `xml:"steps"`
	StepList  []Step   	`xml:"step"`
}

type Step struct {
	XMLName     xml.Name  `xml:"step"`
	StepNumber  int       `xml:"step_number"`
	Actions     string    `xml:"actions"`
	Results     string    `xml:"expectedresults"`
}

type CustomFields struct {
	XMLName           xml.Name       `xml:"custom_fields"`
	CustomFieldList   []CustomField  `xml:"custom_field"`
}

type CustomField struct {
	XMLName xml.Name  `xml:"custom_field"`
	Name    string    `xml:"name"`
	Value   string    `xml:"value"`
}

func ExportAsTestcases(records [][]string) ([]byte, error) {
	
	xml_data := &Testcases{}
	
	for linenum, record := range records {
		
		// skip first line or if testcase has no name provided
		if linenum == 0 || record[0] == "" {
			continue
		}
		
		testcase := Testcase{Name: record[0], Summary: record[1]}
		
		num_columns := len(record)
		steps := Steps{}
		cfs := CustomFields{}
		
		for i, j, k := 3, 4, 1; i < num_columns; i, j, k = i+2, j+2, k+1 {
			if (i < num_columns && record[i] != "") || (j < num_columns && record[j] != "") {

				step := Step{StepNumber: k, Actions: record[i], Results: record[j]}
				steps.StepList = append(steps.StepList, step)
			}
		}
		
		cf := CustomField{Name: "Comments", Value: record[2]}
		cfs.CustomFieldList = append(cfs.CustomFieldList, cf)
		
		testcase.Steps = steps
		testcase.CustomFields = cfs
		
		xml_data.TestcaseList = append(xml_data.TestcaseList, testcase)
	}
	return xml.MarshalIndent(xml_data, "  ", "    ")
}

func leftPadZero(s string, totalLength int) string {
	var padCount int
	padCount = totalLength - len(s)
	if padCount >= 0 {
		return strings.Repeat("0", padCount) + s
	} else {
		return s
	}
}

func ExportAsRequirements(records [][]string, docPrefix string) (string, error) {
	xml_file := "<?xml version=\"1.0\" encoding=\"UTF-8\"?><requirements>"
	
	docPrefixHyphen := ""
	if docPrefix != "" {
		docPrefixHyphen = docPrefix + "-"
	}
	
	for linenum, record := range records {
		
		// skip first line
		if linenum == 0 {
			continue
		}
		
		xml_file += "<requirement>"
		xml_file += "<docid>" + docPrefixHyphen + leftPadZero(strconv.Itoa(linenum), 4) + "</docid>"
		xml_file += "<title>" + record[0] + "</title>"

		xml_file += "<description>" + record[1] + "</description>"
		
		// if tsNum > -1 {
		// 	xml_file += "<expected_coverage>ts-" + strconv.Itoa(tsNum) + "</expected_coverage>"
		// 	tsNum++
		// }
		
		// num_columns := len(record)
		// xml_file += "<steps>"
		// for i, j := 3, 4; i < num_columns; i, j = i+2, j+2 {
		// 	if (i < num_columns && record[i] != "") || (j < num_columns && record[j] != "") {
		// 		xml_file += "<step>"
		// 		xml_file += "<step_number>" + strconv.Itoa(i+1) + "</step_number>"
		// 		xml_file += "<actions>" + record[i] + "</actions>" 
		// 		xml_file += "<expectedresults>" + record[j] + "</expectedresults>"
		// 		xml_file += "</step>"
		// 	}
		// }
		// xml_file += "</steps>"

		// xml_file += "<custom_fields><custom_field>"
		// xml_file += "<name>Comments</name>"
		// xml_file += "<value>" + record[2] + "</value>"
		// xml_file += "</custom_field></custom_fields>"

		xml_file += "</requirement>"
	}
	xml_file += "</requirements>"
	return xml_file, nil
}

func main() {
	
	requirementsFlagPtr := flag.Bool("r", false, "specify to import as requirements")
	filenamePtr := flag.String("f", "", "filename")
	// docPrefixPtr := flag.String("prefix", "", "doc id prefix (for requirements)")
	// testcaseStartNumPtr := flag.Int("ts", -1, "the first testcase ts-[id] from which this set of requirements will cover")
	versionFlagPtr := flag.Bool("v", false, "print current version")
	flag.Parse()
	
	var filename string = *filenamePtr
	requirementsFlag := *requirementsFlagPtr
	// docPrefix := *docPrefixPtr
	// testcaseStartNum := *testcaseStartNumPtr
	
	if *versionFlagPtr {
		fmt.Println(AppVersion)
		os.Exit(0)
	}
	
	if filename != "" {
		fmt.Print("Parsing ", filename, "... Importing as ")
		if (requirementsFlag) {
			fmt.Print("requirements")
		} else {
			fmt.Print("testcases")
		}
		fmt.Println()
	} else {
		fmt.Println("no filename!")
		return
	}
	
	f, err := os.Open(filename)
	if err != nil {
		// err is printable
		// elements passed are separated by space automatically
		printError(err)
		return
	}
	// automatically call Close() at the end of current method
	defer f.Close()
	
	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		printError(err)
		return
	}
	var extension = filepath.Ext(filename)
	var filenameNoExt = filepath.Base( filename[0:len(filename)-len(extension)] )
	var filenameNoExtBase = strings.Replace(filenameNoExt, "Testcases", "", 1)
	var outfilename string
	var outfilexml []byte
	if (requirementsFlag) {
		outfilename = "Requirements"+filenameNoExtBase+".xml"
		// outfilexml, err = ExportAsRequirements(records, docPrefix)
	} else {
		outfilename = filenameNoExt+".xml"
		outfilexml, err = ExportAsTestcases(records)
	}
	if err != nil {
		printError(err)
	} else {
		outfile, err := os.Create(outfilename)
		if err != nil {
			printError(err)
		} else {
			outfile.Write([]byte(xml.Header))
			outfile.Write(outfilexml)
			outfile.Sync()
			outfile.Close()
			fmt.Println("Success!")
		}
	}
	
}