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
	ErrEnvIsNotSet = fmt.Errorf("following envs must be set: DIRECTION, MIGRATIONS, MIGRATE_DOWN_WITH_INDEX (optional), ADDRESS, USERNAME, PASSWORD")
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
func getEnv() (string, string, string, string, string, string, error) {
	direction := os.Getenv("DIRECTION")
	if direction == "" {
		return "", "", "", "", "", "", ErrEnvIsNotSet
	}
	defer os.Unsetenv("DIRECTION")

	migrations := os.Getenv("MIGRATIONS")
	if direction == "" {
		return "", "", "", "", "", "", ErrEnvIsNotSet
	}
	defer os.Unsetenv("MIGRATIONS")

	migrationDownBool := os.Getenv("MIGRATE_DOWN_WITH_INDEX")
	defer os.Unsetenv("MIGRATE_DOWN_WITH_INDEX")

	addr := os.Getenv("ADDRESS")
	if direction == "" {
		return "", "", "", "", "", "", ErrEnvIsNotSet
	}
	defer os.Unsetenv("ADDRESS")

	username := os.Getenv("USERNAME")
	if direction == "" {
		return "", "", "", "", "", "", ErrEnvIsNotSet
	}
	defer os.Unsetenv("USERNAME")

	password := os.Getenv("PASSWORD")
	if direction == "" {
		return "", "", "", "", "", "", ErrEnvIsNotSet
	}
	defer os.Unsetenv("PASSWORD")

	return direction, migrations, migrationDownBool, addr, username, password, nil
}

// main() is the entry point of the Migrator.
func main() {
	direction, migrations, migrationDownBool, addr, username, password, err := getEnv()
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
			Address: addr,
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
				Username: username,
				Password: password,
			},
		},
	})

	switch direction {
	case "up":
		err := m.MigrateUp(migrations)
		if err != nil {
			panic(err)
		}

	case "down":
		if migrationDownBool == "true" {
			err := m.MigrateDown(true)
			if err != nil {
				panic(err)
			}
		}
		if migrationDownBool == "false" {
			err := m.MigrateDown(false)
			if err != nil {
				panic(err)
			}
		}
	}
}
