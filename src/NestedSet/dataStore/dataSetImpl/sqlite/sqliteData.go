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
    "github.com/jmoiron/sqlx"
    "NestedSet/dataStore"
    "NestedSet/logger"
    "NestedSet/appErrors"
    "NestedSet/sys"
)

const (
    SQL_DATA_TABLE_NAME = "dataSet"
    DATA_UID = "Uid"
    PARENT_UID = "Puid"
    DATA_NAME = "Name"
    DATA_DESC = "Desc"
    DATA_LFTID = "LftId"
    DATA_RGTID = "RgtId"
)

var (
    dataSchema = fmt.Sprintf(
                 `CREATE TABLE IF NOT EXISTS %s (%s TEXT PRIMARY KEY,
                 %s TEXT,
                 %s TEXT NOT NULL,
                 %s TEXT,
                 %s INTEGER DEFAULT %d,
                 %s INTEGER DEFAULT %d)`,
                 SQL_DATA_TABLE_NAME,
                 DATA_UID,
                 PARENT_UID,
                 DATA_NAME,
                 DATA_DESC,
                 DATA_LFTID, dataStore.DEFAULT_LFTID,
                 DATA_RGTID, dataStore.DEFAULT_RGTID)
    //Create a entry without any nested set parameters
    dataCreate = fmt.Sprintf(`INSERT INTO %s
                                (%s, %s, %s, %s)
                                VALUES (?, ?, ?, ?)`,
                                SQL_DATA_TABLE_NAME,
                                DATA_UID,
                                PARENT_UID,
                                DATA_NAME,
                                DATA_DESC)
    dataDeleteOnId = fmt.Sprintf(`DELETE FROM %s WHERE %s=(?)`,
                                 SQL_DATA_TABLE_NAME, DATA_UID)
    // Update the nested set parameters for the tree hierarchy.
    dataUpdateTree = fmt.Sprintf(`UPDATE %s SET %s=(?),%s=(?)
                                  WHERE %s=(?)`,
                                SQL_DATA_TABLE_NAME,
                                DATA_LFTID,
                                DATA_RGTID,
                                DATA_UID)
    dataGetAllRec = fmt.Sprintf(`SELECT * FROM %s`,
                                    SQL_DATA_TABLE_NAME)
    dataGetwthUid = fmt.Sprintf(`SELECT * FROM %s WHERE %s=(?)`,
                                SQL_DATA_TABLE_NAME, DATA_UID)
    //Get an entry with name
    dataGetwthName = fmt.Sprintf(`SELECT * FROM %s WHERE %s=(?)`,
                               SQL_DATA_TABLE_NAME,
                               DATA_NAME)
    //Get an entry with name under specific parent.
    dataGetwthNamePID = fmt.Sprintf(`SELECT * FROM %s WHERE %s=(?) AND
                               %s=(?)`,
                               SQL_DATA_TABLE_NAME,
                               DATA_NAME,
                               PARENT_UID)
    //Get all the records greater than the specific lftID
    dataGetAllRecsGrtrId = fmt.Sprintf(`SELECT * FROM %s WHERE %s>(?)`,
                                        SQL_DATA_TABLE_NAME,
                                        DATA_LFTID)
    dataGetAllChildrens = fmt.Sprintf(`SELECT * FROM %s WHERE %s>(?)
                                       AND %s<(?)`, SQL_DATA_TABLE_NAME,
                                       DATA_LFTID, DATA_LFTID)

)

// Anonymous pointer to Data struct. It is decided to use pointer instead of
// Data value as it reduce the overhead of copying in later stages.
// Creation of sqlData must need to allocate Data as well to avoid runtime
// panic. However end to end sqlData can keep the same Data entry struct
// for whole operations without copying the values each time.
// It is assumed copying pointer faster than copying the entire structure.
type sqlData struct {
    *dataStore.Data
}

func(dataObj *sqlData)IsdataRoot() bool {
    if dataObj.Uid == dataStore.ROOT_UID {
        return true
    }
    return false
}
func(dataObj *sqlData)CreateTable(conn *sqlx.DB) error {
    var err error
    log := logger.GetLoggerInstance()
    _, err = conn.Exec(dataSchema)
    if err != nil {
        log.Error("Failed to create data table %s", err)
        return err
    }
    log.Trace("Table %s created successfully", SQL_DATA_TABLE_NAME)
    return nil
}

func (dataObj *sqlData)GetAllChildrens(conn *sqlx.DB)([]dataStore.Data,
                                               error) {
    var err error
    log := logger.GetLoggerInstance()
    rows := []dataStore.Data{}
    err = conn.Select(&rows, dataGetAllChildrens, dataObj.LftId, dataObj.RgtId)
    if err != nil {
        log.Error("Failed to retereive all the childrens for %d - %d",
                    dataObj.LftId, dataObj.RgtId)
        return nil, err
    }
    return rows, nil
 }

func (dataObj *sqlData)GetAllRecords(conn *sqlx.DB)([]dataStore.Data,
                                               error) {
    var err error
    log := logger.GetLoggerInstance()
    rows := []dataStore.Data{}
    err = conn.Select(&rows, dataGetAllRec)
    if err != nil {
        log.Error("Failed to retereive all data  from table err : %s",
                    err)
        return nil, err
    }
    return rows, nil
 }

func (dataObj *sqlData)GetAllRecsGreaterThanId(conn *sqlx.DB)([]dataStore.Data,
                                               error) {
    var err error
    log := logger.GetLoggerInstance()
    rows := []dataStore.Data{}
    err = conn.Select(&rows, dataGetAllRecsGrtrId, dataObj.LftId)
    if err != nil {
        log.Error("Failed to retereive the records greater than lftId %d",
                    dataObj.LftId)
        return nil, err
    }
    return rows, nil
 }

//Delete a record from a table using its ID
func (dataObj *sqlData)deleteDataOnId(conn *sqlx.DB)(error) {
    var err error
    log := logger.GetLoggerInstance()
    if len(dataObj.Uid) == 0 {
        log.Error("Cannot retrieve a record with empty Uid" +
                  "Delete is failed")
        return appErrors.INVALID_INPUT
    }
    if dataObj.IsdataRoot() == true {
        //System doesnt allow to delete root node, as its owned by application.
        log.Error("Cannot delete root node, as its not owned by user")
        return appErrors.INVALID_INPUT
    }
    _, err = conn.Exec(dataDeleteOnId, dataObj.Uid)
    if err != nil {
        log.Error("Failed to delete record with uid %s", dataObj.Uid)
        return err
    }
    return nil
}
//Retrieve a record with record UID
func (dataObj *sqlData)GetdataById(conn *sqlx.DB)(*dataStore.Data, error) {
    var err error
    log := logger.GetLoggerInstance()
    rows := []dataStore.Data{}
    if len(dataObj.Uid) == 0 {
        log.Error("Cannot retrieve a record with empty Uid")
        return nil,appErrors.INVALID_INPUT
    }
    err = conn.Select(&rows, dataGetwthUid, dataObj.Uid)
    if err != nil {
        log.Error("Failed to get record with uid %s", dataObj.Uid)
        return nil, err
    }
    if len(rows) != 1 {
        log.Error("%d records only present in the system", len(rows))
        return nil,appErrors.DATA_NOT_UNIQUE_ERROR
    }
    return &rows[0], nil
}

//Retreive a record with specific name and parent ID.
func (dataObj *sqlData)getDataWithName(conn *sqlx.DB) ([]dataStore.Data,
                                                            error) {
    var err error
    log := logger.GetLoggerInstance()
    rows := []dataStore.Data{}
    if len(dataObj.Name) == 0 {
        log.Error("Failed to get the record, name/puid is null")
        return nil, appErrors.INVALID_INPUT
    }
    err = conn.Select(&rows, dataGetwthName, dataObj.Name)
    if err != nil {
        log.Error("Failed to retereive the record with name %s",
                    dataObj.Name)
        return nil, err
    }
    return rows, nil
}

//Retreive a record with specific name and parent ID.
func (dataObj *sqlData)getDataWithNameAndPID(conn *sqlx.DB) ([]dataStore.Data,
                                                            error) {
    var err error
    var puid string
    log := logger.GetLoggerInstance()
    rows := []dataStore.Data{}
    if len(dataObj.Name) == 0 {
        log.Error("Failed to get the record, name/puid is null")
        return nil, appErrors.INVALID_INPUT
    }
    puid = dataObj.Puid
    if len(dataObj.Puid) == 0 {
        //Wanted to insert the record at top level.so set PUID as default
        puid = dataStore.DEFAULT_PUID
        
    }
    err = conn.Select(&rows, dataGetwthNamePID, dataObj.Name,
                      puid)
    if err != nil {
        log.Error("Failed to retereive the record with name %s and pid %s",
                    dataObj.Name, dataObj.Puid)
        return nil, err
    }
    return rows, nil
}

// Update the NS limits with new values.
func(dataObj *sqlData)updateIds(conn *sqlx.DB) error {
    var err error
    log := logger.GetLoggerInstance()
    _, err = conn.Exec(dataUpdateTree, dataObj.LftId,
                        dataObj.RgtId, dataObj.Uid)
    if err != nil {
        log.Error("Failed to update the limits err : %s", err)
        return err
    }
    return nil
}
//Create default ROOT node for the tree. Root node is created when the new
// table is created at first time.
func(dataObj *sqlData)InsertRoot(conn *sqlx.DB) error {
    var err error
    var rootdata *dataStore.Data
    log := logger.GetLoggerInstance()
    dataObj.Uid = dataStore.DEFAULT_PUID
    dataObj.Name = "root"
    dataObj.Desc = "The default root node in the hierarchy"
    dataObj.LftId = 1
    dataObj.RgtId = dataObj.LftId + 1
    //Check if data present before trying to insert them.
    rootdata, err = dataObj.GetdataById(conn)
    if rootdata != nil {
        log.Trace("Root record already present in the system.")
        return nil
    }
    err = dataObj.insertData(conn)
    if err != nil && err != appErrors.DATA_PRESENT_IN_SYSTEM {
        log.Error("Failed to insert the root node, err : %s", err)
        return err
    }
    //Update the nestedset information for the root node.
    err = dataObj.updateIds(conn)
    if err != nil {
        log.Error("Failed to update nestedSet ID for root node err: %s", err)
        return err
    }
    return nil
}

func(dataObj *sqlData)insertData(conn *sqlx.DB) error {
    log := logger.GetLoggerInstance()
    rows, err := dataObj.getDataWithNameAndPID(conn)
    if err != nil {
        log.Error("Failed to get the record from DB, err: %s", err)
        return err
    }
    if len(rows) != 0 {
        log.Info("Cannot insert data, Record already present in the system")
        return appErrors.DATA_PRESENT_IN_SYSTEM
    }
    //TODO :: Update with proper UUID generator.
    var uid string
    uid, err = sys.NewUUIDString()
    if err != nil {
        log.Error("Failed to create UUID, cannot insert a an entry in DB")
        return err
    }
    // Set the PUID to default at first level. and avoid setting for root.
    if dataObj.Uid != dataStore.DEFAULT_PUID {
        if len(dataObj.Puid) == 0 {
            //Wanted to insert the record at top level, use default Puid
            dataObj.Puid = dataStore.DEFAULT_PUID
        }
        dataObj.Uid = uid
    }
    _, err = conn.Exec(dataCreate, dataObj.Uid, dataObj.Puid, dataObj.Name,
                        dataObj.Desc)
    if err != nil {
        log.Error("Failed to insert data %s err %s", dataObj.Name,
                    err);
        return err
    }
    return nil
}

func(dataObj *sqlData)InsertData(conn *sqlx.DB) error {
    err := dataObj.insertData(conn)
    //While adding a new node, we must update the NS values.
    if err == nil {
        nsObj := NewSqliteNestedSet(dataObj, conn)
        err = nsObj.updateNSListLimitsOnAdd()
    }
    return err
}

func(dataObj *sqlData)DeleteData(conn *sqlx.DB) error {
    //Get all the record fields before delete.
    var err error
    log := logger.GetLoggerInstance()
    dataObj.Data, err = dataObj.GetdataById(conn)
    if err != nil {
        log.Error("Failed to get the record %s on delete, Cannot delete",
                   dataObj.Uid)
        return err
    }    
    err = dataObj.deleteDataOnId(conn)
    if err == nil {
        nsObj := NewSqliteNestedSet(dataObj, conn)
        err = nsObj.updateNSListLimitsOnDel()
    }
    return err
}