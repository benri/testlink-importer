package tlstructs

import (
  "encoding/xml"
)

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