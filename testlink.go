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
  "github.com/benri/testlink/tlstructs"
)

const AppVersion = "0.0.2"

func printError(e error) {
  fmt.Println("Error:", e)
}

func ExportAsTestcases(records [][]string) ([]byte, error) {
  
  xml_data := &tlstructs.Testcases{}
  
  for linenum, record := range records {
    
    // skip first line or if no provided testcase name
    if linenum == 0 || record[0] == "" {
      continue
    }
    
    testcase := tlstructs.Testcase{Name: record[0], Summary: record[1]}
    
    num_columns := len(record)
    steps := tlstructs.Steps{}
    cfs := tlstructs.CustomFields{}
    
    for i, j, k := 3, 4, 1; i < num_columns; i, j, k = i+2, j+2, k+1 {
      if (i < num_columns && record[i] != "") || (j < num_columns && record[j] != "") {

        step := tlstructs.Step{}
        step.StepNumber = k
        step.Actions = record[i]
        step.Results = record[j]
        
        steps.StepList = append(steps.StepList, step)
      }
    }
    
    cf := tlstructs.CustomField{}
    cf.Name = "Comments"
    cf.Value = record[2]
    
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

func ExportAsRequirements(records [][]string, docPrefix string) ([]byte, error) {
  xml_data := &tlstructs.Requirements{}
  
  docPrefixHyphen := ""
  if docPrefix != "" {
    docPrefixHyphen = docPrefix + "-"
  }
  
  for linenum, record := range records {
    
    // skip first line or if no provided requirement name
    if linenum == 0 || record[0] == "" {
      continue
    }
    
    req := tlstructs.Requirement{}
    req.Docid = docPrefixHyphen + leftPadZero(strconv.Itoa(linenum), 4)
    req.Title = record[0]
    req.Description = record[1]
    
    req.CustomFields = tlstructs.CustomFields{}
    cf := tlstructs.CustomField{}
    cf.Name = "Comments"
    cf.Value = record[2]
    req.CustomFields.CustomFieldList = append(req.CustomFields.CustomFieldList, cf)
    
    xml_data.RequirementList = append(xml_data.RequirementList, req)
  }
  return xml.MarshalIndent(xml_data, "  ", "    ")
}

func main() {
  
  requirementsFlagPtr := flag.Bool("r", false, "specify to import as requirements")
  filenamePtr := flag.String("f", "", "filename")
  docPrefixPtr := flag.String("prefix", "", "doc id prefix (for requirements)")
  // testcaseStartNumPtr := flag.Int("ts", -1, "the first testcase ts-[id] from which this set of requirements will cover")
  versionFlagPtr := flag.Bool("v", false, "print current version")
  flag.Parse()
  
  var filename string = *filenamePtr
  requirementsFlag := *requirementsFlagPtr
  docPrefix := *docPrefixPtr
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
    outfilexml, err = ExportAsRequirements(records, docPrefix)
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