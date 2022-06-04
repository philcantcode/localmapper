package proposition

import (
	"encoding/json"

	"github.com/philcantcode/localmapper/database"
	"github.com/philcantcode/localmapper/utils"
)

func InsertProposition(proposition Proposition) {
	utils.Log("Inserting Proposition from Proposition DB", false)
	stmt, err := database.Con.Prepare("INSERT INTO `Propositions`" +
		"(`type`, `description`, `proposition`, `correction`, `status`, `user`) VALUES (?, ?, ?, ?, ?, ?);")
	utils.ErrorLog("Couldn't prepare InsertProposition into Proposition", err, true)

	propositionString, err := json.Marshal(proposition.Proposition)
	utils.ErrorLog("Couldn't convert proposition to JSON", err, true)

	correctionString, err := json.Marshal(proposition.Correction)
	utils.ErrorLog("Couldn't convert correction to JSON", err, true)

	_, err = stmt.Exec(proposition.Type, proposition.Description, string(propositionString), string(correctionString), proposition.Status, proposition.User)
	utils.ErrorLog("Error executing InsertProposition on Proposition", err, true)
	stmt.Close()
}

func SelectAllPropositions() []Proposition {
	utils.Log("SelectAllPropositions from Proposotions Db (sqlite)", false)
	stmt, err := database.Con.Prepare("SELECT `id`, `type`, `date`, `description`, `proposition`, `correction`, `status`, `user` FROM `Propositions` ORDER BY `id` DESC")
	utils.ErrorLog("Couldn't select all from Propositions", err, true)

	rows, err := stmt.Query()
	utils.ErrorLog("Couldn't recieve rows from SelectAllPropositions", err, true)
	defer rows.Close()

	props := []Proposition{}

	for rows.Next() {
		prop := Proposition{}

		propString := ""
		correctionString := ""

		rows.Scan(&prop.ID, &prop.Type, &prop.Date, &prop.Description, &propString, &correctionString, &prop.Status, &prop.User)

		json.Unmarshal([]byte(propString), &prop.Proposition)
		json.Unmarshal([]byte(correctionString), &prop.Correction)

		props = append(props, prop)
	}

	return props
}

func SelectPropositionByID(ID int) Proposition {
	utils.Log("SelectPropositionByID from Proposotions Db (sqlite)", false)
	stmt, err := database.Con.Prepare("SELECT `id`, `type`, `date`, `description`, `proposition`, `correction`, `status`, `user` FROM `Propositions` WHERE `id` = ?")
	utils.ErrorLog("Couldn't select all from Propositions", err, true)

	rows, err := stmt.Query(ID)
	utils.ErrorLog("Couldn't recieve rows from SelectPropositionByID", err, true)
	defer rows.Close()

	for rows.Next() {
		prop := Proposition{}

		propString := ""
		correctionString := ""

		rows.Scan(&prop.ID, &prop.Type, &prop.Date, &prop.Description, &propString, &correctionString, &prop.Status, &prop.User)

		json.Unmarshal([]byte(propString), &prop.Proposition)
		json.Unmarshal([]byte(correctionString), &prop.Correction)

		return prop
	}

	utils.ErrorContextLog("Couldn't find proposition by ID", true)
	return Proposition{}
}

/* 0 = Open
   1 = Complete
   2 = Disabled */
func SetPropositionStatusByID(ID int, status int) {
	utils.Log("SetPropositionStatusByID from Proposotions Db (sqlite)", false)
	stmt, err := database.Con.Prepare("UPDATE `Propositions` SET `status` = ? WHERE `id` = ?")
	utils.ErrorLog("Couldn't select all from Propositions", err, true)

	_, err = stmt.Exec(status, ID)
	utils.ErrorLog("Couldn't recieve rows from SetPropositionStatusByID", err, true)
}
