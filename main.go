package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/philcantcode/localmapper/capabilities/localhost"
	"github.com/philcantcode/localmapper/database"
	"github.com/philcantcode/localmapper/execute"
	"github.com/philcantcode/localmapper/installers"
	"github.com/philcantcode/localmapper/utils"
)

func main() {
	utils.LoadGlobalConfigs()
	installers.Check3rdPartyPrerequisites()
	database.Initialise()
	database.MongoConnect()

	utils.Log("Server hosted at http://localhost:"+utils.Configs["SERVER_PORT"], true)

	switch utils.Configs["MODE"] {
	case "Interactive":
		go interactiveCLI()
	}

	router := mux.NewRouter()

	router.HandleFunc("/capability/run", execute.CapabilityRunAPI)
	router.HandleFunc("/capability/list", execute.CapabilityListAPI)
	router.HandleFunc("/capability/displayfields/update", execute.CapabilityDisplayFieldsUpdateAPI)

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
		utils.PrettyPrint(localhost.IpInfo())
	case "os":
		utils.PrettyPrint(localhost.OSInfo())
	case "run capability":
		execute.RunCapability()
	case "help":
		fmt.Println("Available Commands: ip, os, run capability, help")
	}
}
