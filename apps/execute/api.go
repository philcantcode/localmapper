package execute

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/philcantcode/localmapper/utils"
)

func CapabilityDisplayFieldsUpdateAPI(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	var displayFields []string

	err := json.Unmarshal([]byte(r.FormValue("displayFields")), &displayFields)
	utils.ErrorFatal("Couldn't convert CapabilityDisplayFieldsUpdateAPI displayFields to JSON", err)

	if r.FormValue("displayFields") == "" {
		utils.ErrorContextLog("No displayFields specified for CapabilityDisplayFieldsUpdateAPI", true)
		os.Exit(0)
	}

	if id == "" {
		utils.ErrorContextLog("No id specified for CapabilityDisplayFieldsUpdateAPI", true)
		os.Exit(0)
	}

	idInt, err := strconv.Atoi(id)
	displayFieldsStr := new(strings.Builder)
	json.NewEncoder(displayFieldsStr).Encode(displayFields)

	UpdateCommandDisplayField(displayFieldsStr.String(), idInt)

	w.WriteHeader(200)
}

func CapabilityListAPI(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	if id == "" {
		utils.ErrorContextLog("No id specified for CapabilityListAPI", true)
		os.Exit(0)
	}

	capabilities := GetAllCapabilities()

	if id == "all" || id == "-1" {
		json.NewEncoder(w).Encode(capabilities)
		return
	}

	idInt, err := strconv.Atoi(id)
	utils.ErrorFatal("Couldn't convert ID in CapabilityListAPI", err)

	for _, capability := range capabilities {
		if capability.ID == idInt {
			json.NewEncoder(w).Encode(capability)
			return
		}
	}
}

func CapabilityRunAPI(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	var params []string
	json.Unmarshal([]byte(r.FormValue("params")), &params)

	if r.FormValue("params") == "" {
		utils.ErrorContextLog("No params specified for CapabilityRunAPI", true)
		os.Exit(0)
	}

	if id == "" {
		utils.ErrorContextLog("No id specified for CapabilityRunAPI", true)
		os.Exit(0)
	}

	capabilities := GetAllCapabilities()
	var capID int

	capID, err := strconv.Atoi(id)
	utils.ErrorFatal("Couldn't convert ID", err)

	for i, k := range capabilities {
		if capID == k.ID {
			capabilities[i].Params = swapCapabilityParamsWithAPIValues(k.Params, params)
			result := Run(capabilities[i])

			utils.PrettyPrint(result)
			json.NewEncoder(w).Encode(result)

			utils.Log(fmt.Sprintf("Capability Complete: [%s] %s", k.Type, k.Name), true)
			return
		}
	}

	utils.ErrorForceFatal("Could not find a patching capability")
	return
}

func swapCapabilityParamsWithAPIValues(params []string, values []string) []string {
	for i, param := range params {
		param = strings.Replace(param, "<", "", -1)
		param = strings.Replace(param, ">", "", -1)

		switch param {
		case "string":
			params[i] = values[i]
		case "string:iprange":
			params[i] = values[i]
		case "string:ip":
			params[i] = values[i]
		}
	}

	return params
}
