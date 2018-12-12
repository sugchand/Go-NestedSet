// Copyright 2018 Sugesh Chandran
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package sqlite

import (
    "fmt"
    "sync"
    "path/filepath"
    "github.com/jmoiron/sqlx"
    _ "github.com/mattn/go-sqlite3"
    "NestedSet/logger"
    "NestedSet/dataStore"
)

var dbOnce sync.Once
var sqlObj *SqliteDataStore

type SqliteDataStore struct {
    dblogger *logger.Logging
    DBConn *sqlx.DB
}


//Create a sql connection and store in the datastore object.
// Return '0' on success and errorcode otherwise.
// It is advised to make single handle in entire application as every handle
// uses a connection pool to manage multiple DB requests.
func (sqlds *SqliteDataStore)CreateDBConnection(
                                dbPath string) error{
    dbDriver := "sqlite3"
    dbFile, err := filepath.Abs(dbPath)
    if err != nil {
        sqlds.dblogger.Error("Failed to open DB file, %s", err.Error())
        return err
    }
    var dbHandle *sqlx.DB
    dbHandle, err = sqlx.Open(dbDriver, dbFile)
    if err != nil {
        sqlds.dblogger.Error("Failed to connect DB %s", err.Error())
        return err
    }
    sqlds.DBConn = dbHandle
    // Serialize the DB access by limiting open connections to 1.
    // This will ensure there are no issues when concurrent threads are
    // accessing the DB file.
    //sqlds.DBConn.SetMaxOpenConns(1)
    sqlds.dblogger.Trace("Created sqlite3 DB connection to %s", dbFile)
    //Start the single thread to update the nested set data now. Its safe to do
    // it here as DB connection is created only once in the entire application.
    go updateNestedSetLimitsInDB()
    return nil
}

//Create all the sqlite tables for the application.
func (sqlds *SqliteDataStore)CreateDataStoreTables() error {
    var err error
    if sqlds.DBConn == nil {
        return fmt.Errorf("Null DB connection, cannot create tables")
    }
    dataObj := new(sqlData) 
    dataObj.Data = new(dataStore.Data) //Must allocate internal pointer too.
    err = dataObj.CreateTable(sqlds.DBConn)
    if err != nil {
        return err
    }
    //Create the root node if not exisits.
    return dataObj.InsertRoot(sqlds.DBConn)
}

func (sqlds *SqliteDataStore)CreateRecord(rec *dataStore.Data) error {
    sqlDataObj := new(sqlData)
    sqlDataObj.Data = rec
    return sqlDataObj.InsertData(sqlds.DBConn)
}

func (sqlds *SqliteDataStore)DeleteRecord(recid string) error {
    sqlDataObj := new(sqlData)
    sqlDataObj.Data = new(dataStore.Data)
    sqlDataObj.Uid = recid
    sqlDataObj.DeleteData(sqlds.DBConn)
    return nil
}

func (sqlds *SqliteDataStore)GetRecord(recid string) (*dataStore.Data, error) {
    sqlDataObj := new(sqlData)
    sqlDataObj.Data = new(dataStore.Data)
    sqlDataObj.Uid = recid
    row, err := sqlDataObj.GetdataById(sqlds.DBConn)
    return row, err
}

func (sqlds *SqliteDataStore)GetRecordByName(name string)([]dataStore.Data, error) {
    sqlDataObj := new(sqlData)
    sqlDataObj.Data = new(dataStore.Data)
    sqlDataObj.Name = name
    row, err := sqlDataObj.getDataWithName(sqlds.DBConn)
    return row, err
}

 func (sqlds *SqliteDataStore)GetAllRecords()([]dataStore.Data, error) {
    sqlDataObj := new(sqlData)
    sqlDataObj.Data = new(dataStore.Data)
    rows, err := sqlDataObj.GetAllRecords(sqlds.DBConn)
     return rows, err
}
 
 // Only one SQL datastore object can be present in the system as connection
//pool can be handled inside the database connection itself
func GetsqliteDataStoreObj() *SqliteDataStore {
    //Initialize the global variable.
    dbOnce.Do(func() {
        sqlObj = new(SqliteDataStore)
        sqlObj.dblogger = logger.GetLoggerInstance()
    })
    return sqlObj
}
