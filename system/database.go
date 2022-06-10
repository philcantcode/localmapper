package system

type Setting struct {
	ID    string
	Key   string
	Value string
}

func SELECT_Settings_All() []Setting {
	stmt, err := Con.Query("SELECT `id`, `key`, `value` FROM `Settings`;")
	Fatal("Couldn't SELECT_ALL from Settings", err)

	var settingsList []Setting

	for stmt.Next() {
		setting := Setting{}

		stmt.Scan(&setting.ID, &setting.Key, &setting.Value)
		settingsList = append(settingsList, setting)
	}

	stmt.Close()
	return settingsList
}

func DELETE_Settings_All() {
	stmt, err := Con.Prepare("DELETE FROM `Settings`;")
	Fatal("Couldn't DELETE_All from Settings", err)

	stmt.Exec()
	stmt.Close()
}

func INSERT_Settings(key string, value string) {
	stmt, err := Con.Prepare("INSERT INTO `Settings` (`key`, `value`) VALUES (?, ?);")
	Fatal("Couldn't INSERT Into Settings", err)

	_, err = stmt.Exec(key, value)
	Fatal("Results error from INSERT Settings", err)

	stmt.Close()
}
