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
	utils.FatalErrorHandle("Couldn't convert CapabilityDisplayFieldsUpdateAPI displayFields to JSON", err)

	if r.FormValue("displayFields") == "" {
		utils.ErrorHandleLog("No displayFields specified for CapabilityDisplayFieldsUpdateAPI", true)
		os.Exit(0)
	}

	if id == "" {
		utils.ErrorHandleLog("No id specified for CapabilityDisplayFieldsUpdateAPI", true)
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
		utils.ErrorHandleLog("No id specified for CapabilityListAPI", true)
		os.Exit(0)
	}

	capabilities := GetAllCapabilities()

	if id == "all" || id == "-1" {
		json.NewEncoder(w).Encode(capabilities)
		return
	}

	idInt, err := strconv.Atoi(id)
	utils.FatalErrorHandle("Couldn't convert ID in CapabilityListAPI", err)

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
		utils.ErrorHandleLog("No params specified for CapabilityRunAPI", true)
		os.Exit(0)
	}

	if id == "" {
		utils.ErrorHandleLog("No id specified for CapabilityRunAPI", true)
		os.Exit(0)
	}

	capabilities := GetAllCapabilities()
	var capID int

	capID, err := strconv.Atoi(id)
	utils.FatalErrorHandle("Couldn't convert ID", err)

	for _, k := range capabilities {
		if capID == k.ID {
			result, success := Run(k.Interpreter, k.Params[0], swapCapabilityParamsWithAPIValues(k.Params, params)...)

			utils.Log(fmt.Sprintf("Command status: %v\n", success), true)
			utils.PrettyPrint(result)
			json.NewEncoder(w).Encode(result)
		}
	}

	utils.Log("RunCapability done", true)
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
