package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/philcantcode/localmapper/api"
	"github.com/philcantcode/localmapper/console"
	"github.com/philcantcode/localmapper/database"
	"github.com/philcantcode/localmapper/discovery"
	"github.com/philcantcode/localmapper/installers"
	"github.com/philcantcode/localmapper/utils"
)

func main() {
	utils.LoadGlobalConfigs()
	installers.Check3rdPartyPrerequisites()
	database.Initialise()

	utils.Log("Server hosted at http://localhost:"+utils.Configs["SERVER_PORT"], true)

	switch utils.Configs["MODE"] {
	case "Interactive":
		go interactiveCLI()
	}

	router := mux.NewRouter()

	router.HandleFunc("/get/nmap/pingsweep", api.NmapPingScan)
	router.HandleFunc("/get/nmap/osdetection", api.NmapOSDetectionScan)

	fileServer := http.FileServer(http.Dir("web/src"))
	router.PathPrefix("/").Handler(http.StripPrefix("/", fileServer))

	http.ListenAndServe(":"+utils.Configs["SERVER_PORT"], router)
}

func interactiveCLI() {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Text() != "q" {
		fmt.Print("[>] ")
		scanner.Scan()
		RunCMD(scanner.Text())
	}
}

func RunCMD(cmd string) {
	switch cmd {
	case "ip":
		utils.PrettyPrint(discovery.IpInfo())
	case "os":
		utils.PrettyPrint(discovery.OSInfo())
	case "register command":
		console.RegisterCmdCapability()
	}
}
