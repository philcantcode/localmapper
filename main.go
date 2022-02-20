package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/philcantcode/localmapper/adapters/cli"
	"github.com/philcantcode/localmapper/adapters/web"
	"github.com/philcantcode/localmapper/application/database"
	"github.com/philcantcode/localmapper/application/localhost"
	"github.com/philcantcode/localmapper/utils"
)

func main() {
	utils.LoadGlobalConfigs()
	database.InitSqlite()
	database.InitMongo()

	utils.Log("Server hosted at http://localhost:"+utils.Configs["SERVER_PORT"], true)

	switch utils.Configs["MODE"] {
	case "Interactive":
		go interactiveCLI()
	}

	router := mux.NewRouter()

	router.HandleFunc("/capability/run", web.RunCapability)
	router.HandleFunc("/capability/get", web.GetCapabilities)
	router.HandleFunc("/capability/update", web.UpdateCapability)

	fileServer := http.FileServer(http.Dir("application/web/src"))
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
		cli.RunCapability()
	case "help":
		fmt.Println("Available Commands: ip, os, run capability, help")
	}
}
