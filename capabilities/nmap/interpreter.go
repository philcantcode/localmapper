package nmap

import (
	"encoding/xml"

	"github.com/philcantcode/localmapper/utils"
)

// MakeStructured turns the raw output XML into golang structs for processing
func MakeStructured(xmlstr string) NmapRun {
	utils.Log("Converting from XML to an NMAP struct", false)

	var nmapRun NmapRun
	xmlerr := xml.Unmarshal([]byte(xmlstr), &nmapRun)

	if xmlerr != nil {
		utils.ErrorHandle("Error unmarshaling Nmap XML string", xmlerr, false)
	}
	return nmapRun
}
