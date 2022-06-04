package jobs

import (
	"encoding/json"

	"github.com/philcantcode/localmapper/database"
	"github.com/philcantcode/localmapper/utils"
)

func SELECT_JobSpec_All() []JobSpec {
	utils.Log("SELECT_JobSpec_All from JobSpec Db (sqlite)", false)
	stmt, err := database.Con.Prepare("SELECT `id`, `job` FROM `JobSpecs`")
	utils.ErrorLog("Couldn't select all from JobSpec", err, true)

	rows, err := stmt.Query()
	utils.ErrorLog("Couldn't recieve rows from SELECT_JobSpec_All", err, true)
	defer rows.Close()

	jobs := []JobSpec{}

	for rows.Next() {
		job := JobSpec{}

		id := -1
		jobString := ""

		rows.Scan(&id, &jobString)

		json.Unmarshal([]byte(jobString), &job)
		job.JobID = id

		jobs = append(jobs, job)
	}

	return jobs
}
