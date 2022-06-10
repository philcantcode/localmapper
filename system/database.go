package system

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

func INSERT_LogEntry(log LogEntry) {
	if !MONGO_INITIALISED {
		return
	}

	log.ID = primitive.NewObjectID()
	_, err := System_Logs_DB.InsertOne(context.Background(), log)

	if err != nil {
		fmt.Println("Fatal Error INSERT_LogEntry")
	}
}

func SELECT_LogEntry(filter bson.M, projection bson.M) []LogEntry {
	cursor, err := System_Logs_DB.Find(context.Background(), filter, options.Find().SetProjection(projection))

	if err != nil {
		fmt.Println("Couldn't find SELECT_LogEntry")
	}

	defer cursor.Close(context.Background())

	results := []LogEntry{}

	for cursor.Next(context.Background()) {
		var log LogEntry

		err = cursor.Decode(&log)

		if err != nil {
			fmt.Println("Couldn't decode SELECT_LogEntry")
		}

		results = append(results, log)
	}

	return results
}
