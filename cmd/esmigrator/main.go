package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"time"

	goelasticmigrator "github.com/xoticdsign/go-elasticmigrator"
)

var (
	ErrEnvsMustBeSpecified = fmt.Errorf("you have to specify the envs")
)

// getEnv() retrieves the necessary environment variables required for migration.
// These environment variables must be set before running the application. The function ensures
// all variables are present and returns an error if any are missing.
//
// Environment Variables (envs) explained:
//   - DIRECTION: This environment variable specifies the direction of the migration.
//     It can either be "up" or "down" and determines whether the migration should be applied
//     or reverted. If not set, the program will throw an error.
//   - MIGRATIONS: This variable holds the path or identifier for the migration scripts or
//     configurations. It is critical for the program to know where to look for migration files.
//   - MIGRATE_DOWN_WITH_INDEX: A flag (true/false) indicating if the index should be removed
//     when rolling back a migration (down migration). If not provided, an error will occur.
//   - ADDRESS: The address (URL) of the Elasticsearch server. This is used to connect to the
//     server for performing migration operations. Without this, the connection cannot be established.
//   - USERNAME: The username for Elasticsearch authentication. It is required to authenticate
//     and access Elasticsearch for migration operations.
//   - PASSWORD: The password corresponding to the USERNAME, used to authenticate access to Elasticsearch.
//
// If any of these environment variables are missing, the function will return an error (ErrEnvIsNotSet),
// indicating which variable was not set properly.
func getEnv() (map[string]string, error) {
	direction := os.Getenv("DIRECTION")
	migrations := os.Getenv("MIGRATIONS")
	migrationDownBool := os.Getenv("MIGRATE_DOWN_WITH_INDEX")
	esAddress := os.Getenv("ES_ADDRESS")
	esUsername := os.Getenv("ES_USERNAME")
	esPassword := os.Getenv("ES_PASSWORD")

	envs := make(map[string]string)

	envs["direction"] = direction
	envs["migrations"] = migrations
	envs["migrate_down_with_index"] = migrationDownBool
	envs["es_address"] = esAddress
	envs["es_username"] = esUsername
	envs["es_password"] = esPassword

	return envs, nil
}

// main() is the entry point of the Migrator.
func main() {
	envs, err := getEnv()
	if err != nil {
		panic(err)
	}

	m := goelasticmigrator.New(goelasticmigrator.MigratorConfig{
		Client: http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
				TLSHandshakeTimeout: time.Second * 10,
			},
			Timeout: time.Second * 10,
		},
		ElasticSearch: goelasticmigrator.ElasticSearch{
			Address: envs["address"],
			Index: goelasticmigrator.Index{
				Name: "products",
				Definition: map[string]interface{}{
					"mappings": map[string]interface{}{
						"properties": map[string]interface{}{
							"category": map[string]interface{}{
								"type": "keyword",
							},
							"created_at": map[string]interface{}{
								"type": "date",
							},
							"description": map[string]interface{}{
								"type": "text",
							},
							"id": map[string]interface{}{
								"type": "long",
							},
							"name": map[string]interface{}{
								"type": "text",
							},
							"price": map[string]interface{}{
								"type": "float",
							},
							"stock": map[string]interface{}{
								"type": "integer",
							},
						},
					},
				},
			},
			Credentials: goelasticmigrator.Credentials{
				Username: envs["username"],
				Password: envs["password"],
			},
		},
	})

	switch envs["direction"] {
	case "up":
		err := m.MigrateUp(envs["migrations"])
		if err != nil {
			panic(err)
		}

	case "down":
		if envs["migrate_down_with_index"] == "true" {
			err := m.MigrateDown(true)
			if err != nil {
				panic(err)
			}
		}
		if envs["migrate_down_with_index"] == "false" {
			err := m.MigrateDown(false)
			if err != nil {
				panic(err)
			}
		}
	}
}
