package nmap

import (
	"encoding/xml"
	"strconv"

	"github.com/philcantcode/localmapper/cmdb"
	"github.com/philcantcode/localmapper/system"
	"github.com/philcantcode/localmapper/utils"
)

// Execute takes a list of parameters to execute against NMAP
func ProcessResults(resultByte []byte) NmapRun {
	system.Log("Converting from []byte to NmapRun struct", false)

	var nmapRun NmapRun
	err := xml.Unmarshal(resultByte, &nmapRun)
	system.Error("Couldn't unmarshal result from Nmap console", err)

	return nmapRun
}

func (result NmapRun) StoreResults() {
	INSERT_Nmap(result)
}

/*
	ConvertToEntry takes in an nmapRun and extracts
	relevant variabels.
*/
func (nmapRun NmapRun) ConvertToEntry() {
	// For each host
	for _, host := range nmapRun.Hosts {
		sysTags := []cmdb.EntryTag{}

		ports := cmdb.EntryTag{
			Label:       "Ports",
			Description: "Open Ports",
			DataType:    system.INTEGER,
			Values:      []string{},
		}

		services := cmdb.EntryTag{
			Label:       "Services",
			Description: "Open Services",
			DataType:    system.STRING,
			Values:      []string{},
		}

		products := cmdb.EntryTag{
			Label:       "Products",
			Description: "Detected products on ports",
			DataType:    system.STRING,
			Values:      []string{},
		}

		for _, port := range host.Ports {
			if port.State.State == "open" {
				if !utils.ArrayContains(strconv.Itoa(port.PortId), ports.Values) {
					ports.Values = append(ports.Values, strconv.Itoa(port.PortId))
				}

				if !utils.ArrayContains(port.Service.Name, services.Values) && port.Service.Name != "" {
					services.Values = append(services.Values, port.Service.Name)
				}

				if !utils.ArrayContains(port.Service.Product, products.Values) && port.Service.Product != "" {
					products.Values = append(products.Values, port.Service.Product)
				}
			}
		}

		vendorTags := cmdb.EntryTag{
			Label:       "MACVendor",
			Description: "Vendor of the MAC code",
			DataType:    system.STRING,
			Values:      []string{},
		}

		osFamily := cmdb.EntryTag{
			Label:       "OS",
			Description: "Operating System",
			DataType:    system.STRING,
			Values:      []string{},
		}

		osGen := cmdb.EntryTag{
			Label:       "OSGen",
			Description: "Operating System Generation/Version",
			DataType:    system.STRING,
			Values:      []string{},
		}

		osAccuracy := cmdb.EntryTag{
			Label:       "OSAccuracy",
			Description: "Confidence of Nmap detection",
			DataType:    system.INTEGER,
			Values:      []string{},
		}

		osVendor := cmdb.EntryTag{
			Label:       "OSVendor",
			Description: "Vendor of detected OS",
			DataType:    system.STRING,
			Values:      []string{},
		}

		osCPE := cmdb.EntryTag{
			Label:       "CPE",
			Description: "http://cpe.mitre.org/",
			DataType:    system.STRING,
			Values:      []string{},
		}

		if len(host.Os.OsMatches) > 0 {
			match := host.Os.OsMatches[0]

			if match.Accuracy != "" {
				osAccuracy.Values = append(osAccuracy.Values, match.Accuracy)
			}

			for _, osClass := range match.OsClasses {
				if osClass.OsFamily != "" && !utils.ArrayContains(osClass.OsFamily, osFamily.Values) {
					osFamily.Values = append(osFamily.Values, osClass.OsFamily)
				}

				if osClass.OsGen != "" && !utils.ArrayContains(osClass.OsGen, osGen.Values) {
					osGen.Values = append(osGen.Values, osClass.OsGen)
				}

				if osClass.Vendor != "" && !utils.ArrayContains(osClass.Vendor, osVendor.Values) {
					osVendor.Values = append(osVendor.Values, osClass.Vendor)
				}

				for _, cpe := range osClass.CPEs {
					if !utils.ArrayContains(string(cpe), osCPE.Values) {
						osCPE.Values = append(osCPE.Values, string(cpe))
					}
				}
			}
		}

		for _, address := range host.Addresses {
			if address.AddrType == "ipv4" {
				sysTags = append(sysTags, cmdb.EntryTag{
					Label:       "IP",
					Description: "IP4 Address",
					DataType:    system.IP,
					Values:      []string{address.Addr},
				})
			}

			if address.AddrType == "mac" {
				sysTags = append(sysTags, cmdb.EntryTag{
					Label:       "MAC",
					Description: "Media Access Control",
					DataType:    system.MAC,
					Values:      []string{address.Addr},
				})
			}

			if address.Vendor != "" {
				vendorTags.Values = append(vendorTags.Values, address.Vendor)
			}
		}

		// Check that they're not empty.
		if len(vendorTags.Values) > 0 {
			sysTags = append(sysTags, vendorTags)
		}

		if len(osFamily.Values) > 0 {
			sysTags = append(sysTags, osFamily)
		}

		if len(osGen.Values) > 0 {
			sysTags = append(sysTags, osGen)
		}

		if len(osVendor.Values) > 0 {
			sysTags = append(sysTags, osVendor)
		}

		if len(osAccuracy.Values) > 0 {
			sysTags = append(sysTags, osAccuracy)
		}

		if len(osCPE.Values) > 0 {
			sysTags = append(sysTags, osCPE)
		}

		if len(ports.Values) > 0 {
			sysTags = append(sysTags, ports)
		}

		if len(services.Values) > 0 {
			sysTags = append(sysTags, services)
		}

		if len(products.Values) > 0 {
			sysTags = append(sysTags, products)
		}

		// Hostnames
		if len(host.Hostnames) > 0 {
			hostNameTag := cmdb.EntryTag{
				Label:       "HostName",
				Description: "Media Access Control",
				DataType:    system.STRING,
				Values:      []string{},
			}

			for _, name := range host.Hostnames {
				hostNameTag.Values = append(hostNameTag.Values, name.Name)
			}

			sysTags = append(sysTags, hostNameTag)
		}

		entry := cmdb.Entry{
			Label:       "Nmap Discovered Device",
			Description: "This device was discovered during an Nmap scan: " + nmapRun.Args,
			OSILayer:    0,
			CMDBType:    cmdb.PENDING,
			DateSeen:    []string{utils.Now()},
			SysTags:     sysTags,
		}

		tag, exists, _ := cmdb.FindSysTag("HostName", entry)

		if exists {
			entry.Label = tag.Values[0]
		}

		cmdb.UpdateOrInsert(entry)
	}
}
