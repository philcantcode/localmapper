package database

import (
	"encoding/json"

	"github.com/philcantcode/localmapper/core"
	"github.com/philcantcode/localmapper/utils"
)

func InsertCMDBItem(cmdb core.CMDBItem) {
	utils.Log("InsertCMDBItem from CMDB DB", false)
	stmt, err := connection.Prepare("INSERT INTO `CMDB`" +
		"(`osiLayer`, `dateSeen`, `description`, `statusTags`, `userTags`, `infoTags`) VALUES (?, ?, ?, ?, ?, ?);")
	utils.ErrorLog("Couldn't prepare InsertCMDBItem into CMDB", err, true)

	statusTags, err := json.Marshal(cmdb.StatusTags)
	utils.ErrorFatal("Couldn't unmarshall cmdb.StatusTags", err)

	userTags, err := json.Marshal(cmdb.UserTags)
	utils.ErrorFatal("Couldn't unmarshall cmdb.UserTags", err)

	infoTags, err := json.Marshal(cmdb.InfoTags)
	utils.ErrorFatal("Couldn't unmarshall cmdb.InfoTags", err)

	dateSeen, err := json.Marshal(cmdb.DateSeen)
	utils.ErrorFatal("Couldn't unmarshall cmdb.DateSeen", err)

	_, err = stmt.Exec(cmdb.OSILayer, string(dateSeen), cmdb.Description, string(statusTags), string(userTags), string(infoTags))
	utils.ErrorLog("Error executing InsertCMDBItem on CMDB", err, true)
	stmt.Close()
}

// SelectCMDBItemByLayer returns all items at an OSI layer
func SelectCMDBItemByLayer(osiLayer int) []core.CMDBItem {
	utils.Log("SelectCMDBItemByLayer from CMDB DB (sqlite)", false)
	stmt, err := connection.Prepare("SELECT `id`, `osiLayer`, `dateSeen`, `description`, `statusTags`, `userTags`, `infoTags` FROM `CMDB` WHERE `osiLayer` = ?")
	utils.ErrorLog("Couldn't select CMDB items by iosLayer from CMDB", err, true)

	rows, err := stmt.Query(osiLayer)
	utils.ErrorLog("Couldn't recieve rows from SelectCMDBItemByLayer", err, true)
	defer rows.Close()

	cmdbs := []core.CMDBItem{}

	for rows.Next() {
		cmdb := core.CMDBItem{}

		rows.Scan(&cmdb.ID, &cmdb.OSILayer, &cmdb.DateSeen, &cmdb.Description, &cmdb.StatusTags, &cmdb.UserTags, &cmdb.InfoTags)

		cmdbs = append(cmdbs, cmdb)
	}

	return cmdbs
}

func SelectAllCMDB() []core.CMDBItem {
	utils.Log("SelectAllCMDB from CMDB DB (sqlite)", false)
	stmt, err := connection.Prepare("SELECT `id`, `osiLayer`, `dateSeen`, `description`, `statusTags`, `userTags`, `infoTags` FROM `CMDB`")
	utils.ErrorLog("Couldn't select SelectAllCMDB from CMDB", err, true)

	rows, err := stmt.Query()
	utils.ErrorLog("Couldn't recieve rows from SelectAllCMDB", err, true)
	defer rows.Close()

	cmdbs := []core.CMDBItem{}

	for rows.Next() {
		cmdb := core.CMDBItem{}
		var statusTags string
		var userTags string
		var infoTags string
		dateSeen := ""

		rows.Scan(&cmdb.ID, &cmdb.OSILayer, &dateSeen, &cmdb.Description, &statusTags, &userTags, &infoTags)

		json.Unmarshal([]byte(dateSeen), &cmdb.DateSeen)
		json.Unmarshal([]byte(statusTags), &cmdb.StatusTags)
		json.Unmarshal([]byte(userTags), &cmdb.UserTags)
		json.Unmarshal([]byte(infoTags), &cmdb.InfoTags)

		cmdbs = append(cmdbs, cmdb)
	}

	return cmdbs
}
