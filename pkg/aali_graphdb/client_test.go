// Copyright (C) 2025 ANSYS, Inc. and/or its affiliates.
// SPDX-License-Identifier: MIT
//
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package aali_graphdb

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"testing"
	"time"

	"cloud.google.com/go/civil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// StdoutLogConsumer is a LogConsumer that prints the log to stdout
type StdoutLogConsumer struct{}

// Accept prints the log to stdout
func (lc *StdoutLogConsumer) Accept(l testcontainers.Log) {
	fmt.Print(string(l.Content))
}

var imageName string

func init() {
	flag.StringVar(&imageName, "imagename", "ghcr.io/ansys/aali-graphdb:edge", "Name of the aali-graphdb image to run the tests against")
}

// if you are using Colima see https://golang.testcontainers.org/system_requirements/using_colima/
func getTestClient(t *testing.T) *Client {
	ctx := context.Background()

	fmt.Printf("Running test against image: %q\n", imageName)

	req := testcontainers.ContainerRequest{
		Image:        imageName,
		ExposedPorts: []string{"8080/tcp"},
		WaitingFor:   wait.ForHTTP("/health").WithPort("8080/tcp").WithStartupTimeout(30 * time.Second),
		LogConsumerCfg: &testcontainers.LogConsumerConfig{
			Opts:      []testcontainers.LogProductionOption{testcontainers.WithLogProductionTimeout(10 * time.Second)},
			Consumers: []testcontainers.LogConsumer{&StdoutLogConsumer{}},
		},
		Env: map[string]string{"RUST_LOG": "debug"},
	}
	aaliDbCont, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req, Started: true,
	})
	defer testcontainers.CleanupContainer(t, aaliDbCont)
	require.NoError(t, err)

	port, err := aaliDbCont.MappedPort(ctx, "8080/tcp")
	require.NoError(t, err)
	host, err := aaliDbCont.Host(ctx)
	require.NoError(t, err)

	tr := &http.Transport{
		DisableKeepAlives: true,
	}
	httpClient := http.Client{
		Transport: tr,
	}

	address := fmt.Sprintf("http://%s:%s", host, port.Port())
	client, err := NewClient(address, &httpClient)
	require.NoError(t, err)
	return client
}

func TestGetHealth(t *testing.T) {
	client := getTestClient(t)
	dbs, err := client.GetHealth()
	require.NoError(t, err)
	assert.Equal(t, true, dbs)
}

func TestGetVersion(t *testing.T) {
	client := getTestClient(t)
	version, err := client.GetVersion()
	require.NoError(t, err)
	assert.NotEmpty(t, version.Version)
	assert.NotEmpty(t, version.KuzuVersion)
	assert.NotEmpty(t, version.KuzuStorageVersion)
}

func TestGetDatabases(t *testing.T) {
	client := getTestClient(t)
	dbs, err := client.GetDatabases()
	require.NoError(t, err)
	assert.Equal(t, []string{}, dbs)
}

func TestCreateDatabase(t *testing.T) {
	client := getTestClient(t)

	// make sure no dbs
	dbs, err := client.GetDatabases()
	require.NoError(t, err)
	assert.Equal(t, []string{}, dbs)

	// now insert 1
	const NEWDBNAME = "test-create-db"
	err = client.CreateDatabase(NEWDBNAME)
	require.NoError(t, err)

	// now check the dbs again
	dbs, err = client.GetDatabases()
	require.NoError(t, err)
	assert.Equal(t, []string{NEWDBNAME}, dbs)
}

func TestDeleteDatabase(t *testing.T) {
	client := getTestClient(t)

	// insert db
	const NEWDBNAME = "test-delete-db"
	err := client.CreateDatabase(NEWDBNAME)
	require.NoError(t, err)

	// check its there
	dbs, err := client.GetDatabases()
	require.NoError(t, err)
	assert.Equal(t, []string{NEWDBNAME}, dbs)

	// delete it
	err = client.DeleteDatabase(NEWDBNAME)
	require.NoError(t, err)

	// now check the dbs again
	dbs, err = client.GetDatabases()
	require.NoError(t, err)
	assert.Equal(t, []string{}, dbs)
}

func TestReadWriteData(t *testing.T) {
	client := getTestClient(t)
	const DBNAME = "my-db"

	// create db
	err := client.CreateDatabase(DBNAME)
	require.NoError(t, err)

	// write some data in there
	queries := []string{
		// setup schema
		"CREATE NODE TABLE User(name STRING, age INT64, PRIMARY KEY (name))",
		"CREATE NODE TABLE City(name STRING, population INT64, PRIMARY KEY (name))",
		"CREATE REL TABLE Follows(FROM User TO User, since INT64)",
		"CREATE REL TABLE LivesIn(FROM User TO City)",

		// add a few users
		"CREATE (:User {name: 'Adam', age: 30});",
		"CREATE (:User {name: 'Karissa', age: 40});",
		"CREATE (:User {name: 'Zhang', age: 50});",
		"CREATE (:User {name: 'Noura', age: 25});",

		// create a few cities
		"CREATE (:City {name: 'Waterloo', population: 150000});",
		"CREATE (:City {name: 'Kitchener', population: 200000});",
		"CREATE (:City {name: 'Guelph', population: 75000});",

		// add a few follows relationships
		"MATCH (u1:User), (u2:User) WHERE u1.name = 'Adam' AND u2.name = 'Karissa' CREATE (u1)-[:Follows {since: 2020}]->(u2);",
		"MATCH (u1:User), (u2:User) WHERE u1.name = 'Adam' AND u2.name = 'Zhang' CREATE (u1)-[:Follows {since: 2020}]->(u2);",
		"MATCH (u1:User), (u2:User) WHERE u1.name = 'Karissa' AND u2.name = 'Zhang' CREATE (u1)-[:Follows {since: 2021}]->(u2);",
		"MATCH (u1:User), (u2:User) WHERE u1.name = 'Zhang' AND u2.name = 'Noura' CREATE (u1)-[:Follows {since: 2022}]->(u2);",

		// add a few lives-in relationships
		"MATCH (u:User), (c:City) WHERE u.name = 'Adam' AND c.name = 'Waterloo' CREATE (u)-[:LivesIn {}]->(c);",
		"MATCH (u:User), (c:City) WHERE u.name = 'Karissa' AND c.name = 'Waterloo' CREATE (u)-[:LivesIn {}]->(c);",
		"MATCH (u:User), (c:City) WHERE u.name = 'Zhang' AND c.name = 'Kitchener' CREATE (u)-[:LivesIn {}]->(c);",
		"MATCH (u:User), (c:City) WHERE u.name = 'Noura' AND c.name = 'Guelph' CREATE (u)-[:LivesIn {}]->(c);",
	}
	for _, query := range queries {
		_, err := client.CypherQueryWrite(DBNAME, query, nil)
		require.NoError(t, err)
	}

	// read it back
	res, err := client.CypherQueryRead(DBNAME, "MATCH (a:User)-[e:Follows]->(b:User) RETURN a.name, e.since, b.name ORDER BY a.name ASC, b.name ASC", nil)
	require.NoError(t, err)
	expected := []map[string]any{
		{"a.name": "Adam", "b.name": "Karissa", "e.since": float64(2020)},
		{"a.name": "Adam", "b.name": "Zhang", "e.since": float64(2020)},
		{"a.name": "Karissa", "b.name": "Zhang", "e.since": float64(2021)},
		{"a.name": "Zhang", "b.name": "Noura", "e.since": float64(2022)},
	}
	assert.Equal(t, expected, res)
}

func TestReadWriteGeneric(t *testing.T) {
	client := getTestClient(t)
	const DBNAME = "my-db-generics"

	// create db
	err := client.CreateDatabase(DBNAME)
	require.NoError(t, err)

	// write some data in there
	type Result struct {
		Result string `json:"result"`
	}
	queries := []string{
		// setup schema
		"CREATE NODE TABLE User(name STRING, age INT64, PRIMARY KEY (name))",
		"CREATE NODE TABLE City(name STRING, population INT64, PRIMARY KEY (name))",
		"CREATE REL TABLE Follows(FROM User TO User, since INT64)",
		"CREATE REL TABLE LivesIn(FROM User TO City)",

		// add a few users
		"CREATE (:User {name: 'Adam', age: 30});",
		"CREATE (:User {name: 'Karissa', age: 40});",
		"CREATE (:User {name: 'Zhang', age: 50});",
		"CREATE (:User {name: 'Noura', age: 25});",

		// create a few cities
		"CREATE (:City {name: 'Waterloo', population: 150000});",
		"CREATE (:City {name: 'Kitchener', population: 200000});",
		"CREATE (:City {name: 'Guelph', population: 75000});",

		// add a few follows relationships
		"MATCH (u1:User), (u2:User) WHERE u1.name = 'Adam' AND u2.name = 'Karissa' CREATE (u1)-[:Follows {since: 2020}]->(u2);",
		"MATCH (u1:User), (u2:User) WHERE u1.name = 'Adam' AND u2.name = 'Zhang' CREATE (u1)-[:Follows {since: 2020}]->(u2);",
		"MATCH (u1:User), (u2:User) WHERE u1.name = 'Karissa' AND u2.name = 'Zhang' CREATE (u1)-[:Follows {since: 2021}]->(u2);",
		"MATCH (u1:User), (u2:User) WHERE u1.name = 'Zhang' AND u2.name = 'Noura' CREATE (u1)-[:Follows {since: 2022}]->(u2);",

		// add a few lives-in relationships
		"MATCH (u:User), (c:City) WHERE u.name = 'Adam' AND c.name = 'Waterloo' CREATE (u)-[:LivesIn {}]->(c);",
		"MATCH (u:User), (c:City) WHERE u.name = 'Karissa' AND c.name = 'Waterloo' CREATE (u)-[:LivesIn {}]->(c);",
		"MATCH (u:User), (c:City) WHERE u.name = 'Zhang' AND c.name = 'Kitchener' CREATE (u)-[:LivesIn {}]->(c);",
		"MATCH (u:User), (c:City) WHERE u.name = 'Noura' AND c.name = 'Guelph' CREATE (u)-[:LivesIn {}]->(c);",
	}
	for _, query := range queries {
		_, err := CypherQueryWriteGeneric[Result](client, DBNAME, query, nil)
		require.NoError(t, err)
	}

	type Person struct {
		A     string `json:"a.name"`
		B     string `json:"b.name"`
		Since int64  `json:"e.since"`
	}

	// read it back
	res, err := CypherQueryReadGeneric[Person](client, DBNAME, "MATCH (a:User)-[e:Follows]->(b:User) RETURN a.name, e.since, b.name ORDER BY a.name ASC, b.name ASC", nil)
	require.NoError(t, err)
	expected := []Person{
		{"Adam", "Karissa", 2020},
		{"Adam", "Zhang", 2020},
		{"Karissa", "Zhang", 2021},
		{"Zhang", "Noura", 2022},
	}
	assert.Equal(t, expected, res)
}

func TestReadWriteDataWithParameters(t *testing.T) {
	client := getTestClient(t)
	const DBNAME = "my-db-with-params"

	// create db
	err := client.CreateDatabase(DBNAME)
	require.NoError(t, err)

	// create schema
	queries := []string{
		// setup schema
		"CREATE NODE TABLE User(name STRING, age INT64, PRIMARY KEY (name))",
		"CREATE REL TABLE Follows(FROM User TO User, since INT64)",

		// add a few follows relationships
		"MATCH (u1:User), (u2:User) WHERE u1.name = 'Adam' AND u2.name = 'Karissa' CREATE (u1)-[:Follows {since: 2020}]->(u2);",
		"MATCH (u1:User), (u2:User) WHERE u1.name = 'Adam' AND u2.name = 'Zhang' CREATE (u1)-[:Follows {since: 2020}]->(u2);",
		"MATCH (u1:User), (u2:User) WHERE u1.name = 'Karissa' AND u2.name = 'Zhang' CREATE (u1)-[:Follows {since: 2021}]->(u2);",
		"MATCH (u1:User), (u2:User) WHERE u1.name = 'Zhang' AND u2.name = 'Noura' CREATE (u1)-[:Follows {since: 2022}]->(u2);",
	}
	for _, query := range queries {
		_, err := client.CypherQueryWrite(DBNAME, query, nil)
		require.NoError(t, err)
	}

	// insert user data w/ parameters
	userQuery := "CREATE (:User {name: $name, age: $age});"
	users := []ParameterMap{
		{"name": StringValue("Adam"), "age": Int64Value(30)},
		{"name": StringValue("Karissa"), "age": Int64Value(40)},
		{"name": StringValue("Zhang"), "age": Int64Value(50)},
		{"name": StringValue("Noura"), "age": Int64Value(25)},
	}
	for _, user := range users {
		_, err := client.CypherQueryWrite(DBNAME, userQuery, user)
		require.NoError(t, err)
	}

	// insert relationships data w/ parameters
	followsQuery := "MATCH (u1:User{name:$from}), (u2:User{name:$to}) CREATE (u1)-[:Follows {since: $since}]->(u2);"
	follows := []ParameterMap{
		{"from": StringValue("Adam"), "to": StringValue("Karissa"), "since": Int64Value(2020)},
		{"from": StringValue("Adam"), "to": StringValue("Zhang"), "since": Int64Value(2020)},
		{"from": StringValue("Karissa"), "to": StringValue("Zhang"), "since": Int64Value(2021)},
		{"from": StringValue("Zhang"), "to": StringValue("Noura"), "since": Int64Value(2022)},
	}
	for _, follow := range follows {
		_, err := client.CypherQueryWrite(DBNAME, followsQuery, follow)
		require.NoError(t, err)
	}

	// now read it back w/ parameters
	query := "MATCH (a:User)-[e:Follows]->(b:User) WHERE a.name = $from OR e.since > $after RETURN a.name, e.since, b.name ORDER BY a.name ASC, b.name ASC"
	params := ParameterMap{
		"from": StringValue("Adam"), "after": Int64Value(2021),
	}
	res, err := client.CypherQueryRead(DBNAME, query, params)
	require.NoError(t, err)
	expected := []map[string]any{
		{"a.name": "Adam", "b.name": "Karissa", "e.since": float64(2020)},
		{"a.name": "Adam", "b.name": "Zhang", "e.since": float64(2020)},
		{"a.name": "Zhang", "b.name": "Noura", "e.since": float64(2022)},
	}
	assert.Equal(t, expected, res)
}

type ParamsStruct struct {
	name    string
	boolean bool
	date    civil.Date
	age     uint64
}

func (ps ParamsStruct) AsParameters() (map[string]Value, error) {
	return map[string]Value{
		"name":    StringValue(ps.name),
		"boolean": BoolValue(ps.boolean),
		"date":    DateValue(ps.date),
		"age":     UInt64Value(ps.age),
	}, nil
}

func TestParametersStruct(t *testing.T) {
	client := getTestClient(t)
	const DBNAME = "my-db-from-paramstruct"

	// create db
	err := client.CreateDatabase(DBNAME)
	require.NoError(t, err)

	// create schema
	_, err = client.CypherQueryWrite(DBNAME, "CREATE NODE TABLE User(name STRING, boolean BOOL, date DATE, age INT64, PRIMARY KEY (name))", nil)
	require.NoError(t, err)

	// insert user data w/ parameters
	userQuery := "CREATE (:User {name: $name, age: $age, boolean: $boolean, date: $date});"
	users := []ParamsStruct{
		{"Adam", true, civil.Date{Year: 2024, Month: time.August, Day: 22}, 30},
		{"Karissa", true, civil.Date{Year: 2022, Month: time.January, Day: 7}, 40},
		{"Zhang", false, civil.Date{Year: 2025, Month: time.July, Day: 3}, 50},
		{"Noura", true, civil.Date{Year: 2023, Month: time.October, Day: 15}, 25},
	}
	for _, user := range users {
		_, err := client.CypherQueryWrite(DBNAME, userQuery, user)
		require.NoError(t, err)
	}

	// now read it back w/ parameters
	query := "MATCH (u:User) RETURN u.* ORDER BY u.name ASC"
	res, err := client.CypherQueryRead(DBNAME, query, nil)
	require.NoError(t, err)
	expected := []map[string]any{
		{"u.name": "Adam", "u.boolean": true, "u.date": "2024-08-22", "u.age": float64(30)},
		{"u.name": "Karissa", "u.boolean": true, "u.date": "2022-01-07", "u.age": float64(40)},
		{"u.name": "Noura", "u.boolean": true, "u.date": "2023-10-15", "u.age": float64(25)},
		{"u.name": "Zhang", "u.boolean": false, "u.date": "2025-07-03", "u.age": float64(50)},
	}
	assert.Equal(t, expected, res)
}

func TestErrorsReturned(t *testing.T) {
	client := getTestClient(t)
	const DBNAME = "test-errors"

	err := client.CreateDatabase(DBNAME)
	require.NoError(t, err)

	query := "not a real cypher query"
	expected := `unexpected status code: 500 "Query execution failed: Parser exception: extraneous input 'not' expecting {ALTER, ATTACH, BEGIN, CALL, CHECKPOINT, COMMENT, COMMIT, COPY, CREATE, DELETE, DETACH, DROP, EXPLAIN, EXPORT, IMPORT, INSTALL, LOAD, MATCH, MERGE, OPTIONAL, PROFILE, RETURN, ROLLBACK, SET, UNWIND, USE, WITH, SP} (line: 1, offset: 0)\n\"not a real cypher query\"\n ^^^"`

	t.Run("Read", func(t *testing.T) {
		_, err = client.CypherQueryRead(DBNAME, query, nil)
		require.Error(t, err)
		assert.Equal(t, expected, fmt.Sprint(err))
	})
	t.Run("Write", func(t *testing.T) {
		_, err = client.CypherQueryWrite(DBNAME, query, nil)
		require.Error(t, err)
		assert.Equal(t, expected, fmt.Sprint(err))
	})
}

func TestParameterMapJson(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	pMap := ParameterMap{"myStr": StringValue("hi"), "myBool": BoolValue(true), "myInt": Int16Value(23)}
	jsonParamMap, err := json.Marshal(pMap)
	require.NoError(err)

	expected := map[string]any{
		"myStr":  map[string]any{"String": "hi"},
		"myBool": map[string]any{"Bool": true},
		"myInt":  map[string]any{"Int16": float64(23)},
	}
	var unmarshalledMap map[string]any
	require.NoError(json.Unmarshal(jsonParamMap, &unmarshalledMap))
	assert.Equal(expected, unmarshalledMap)

	var unmarshalledParamMap ParameterMap
	require.NoError(json.Unmarshal(jsonParamMap, &unmarshalledParamMap))
	assert.Equal(pMap, unmarshalledParamMap)
}
