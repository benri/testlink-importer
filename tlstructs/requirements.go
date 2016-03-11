package tlstructs

import (
  "encoding/xml"
)

type Requirements struct {
	XMLName           xml.Name      `xml:"requirements"`
	RequirementList  []Requirement  `xml:"requirement"`
}

type Requirement struct {
	XMLName 			    xml.Name      `xml:"requirement"`
	Docid 			      string   	    `xml:"docid"`
	Title					    string		    `xml:"title"`
	Description	      string	      `xml:"description"`
  ExpectedCoverage  string        `xml:"expected_coverage"`
  CustomFields      CustomFields  `xml:"custom_fields"`
}