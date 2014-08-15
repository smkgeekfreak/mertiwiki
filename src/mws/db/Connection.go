// Connection.go
package db

import (
	"bitbucket.org/liamstask/goose/lib/goose"
	"database/sql"
	"errors"
	"log"
	"os"
	//"path"
)

/*
 GetConnection will return a connection from the specified environment
*/
func GetConnection() (*sql.DB, error) {
	return loadConnection()
}

func loadConnection() (*sql.DB, error) {
	log.Printf("Loading connection")
	wikiEnv := os.Getenv("WIKIENV")
	if len(wikiEnv) <= 0 {
		log.Fatalf("Could not load connection, wiki env not set")
		return nil, errors.New("Wiki Environment not set")
	}
	log.Printf("WIKIENV = %s", wikiEnv)
	wikiConfigPath := os.Getenv("WIKICONFIGROOT")
	if len(wikiConfigPath) <= 0 {
		log.Fatalf("Could not load connection, wiki config path not set")
		return nil, errors.New("Wiki Config Path not set")
	}
	log.Printf("WIKICONFIGROOT = %s", wikiConfigPath)
	config, err := getDBConfig(wikiEnv, wikiConfigPath)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	//log.Printf("Getting Connection to Database: %s , %s", config.Driver.Name, config.Driver.OpenStr)
	db, err := sql.Open(config.Driver.Name, config.Driver.OpenStr)
	if err != nil {
		//panic(err)
		return nil, err
	}
	db.SetMaxIdleConns(0)
	db.SetMaxOpenConns(0)
	//log.Printf(" Connected to Successfully to Database: %s , %s", config.Driver.Name, config.Driver.OpenStr)
	return db, nil
}

func getDBConfig(env string, configPath string) (*goose.DBConf, error) {
	//log.Printf("Goosing config for  %#v", env)
	//configPath := path.Join("../../../../deploy/db/")
	schema := "public"
	dbconf, err := goose.NewDBConf(configPath, env, schema)
	if err != nil {
		log.Printf("Problem with config %#v", err)
		return nil, err
	}
	//log.Printf("Conf %#v:", dbconf)
	return dbconf, nil
}
