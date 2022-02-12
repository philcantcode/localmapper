package installers

import (
	"github.com/philcantcode/localmapper/execute"
	"github.com/philcantcode/localmapper/utils"
)

func Check3rdPartyPrerequisites() {
	appInstallChecker("Nmap", "NMAP_PATH", "-v")
}

func appInstallChecker(appName string, appPath_Config string, testStrings ...string) {

	var success bool = false

	for !success {
		_, success = execute.Run("default", utils.Configs[appPath_Config], testStrings...)

		if !success {
			utils.Log("NMAP not installed at "+utils.Configs[appPath_Config], false)
			utils.UserPrint("Please install the latest version of " + appName + " then provide the PATH to its executable:")
			utils.Configs[appPath_Config] = utils.UserStringInput()

			continue
		}

		utils.Log("NMAP is installed at "+utils.Configs[appPath_Config], false)
	}
}
