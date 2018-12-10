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
    "NestedSet/appErrors"
    "NestedSet/dataStore"
    "NestedSet/logger"
    "github.com/jmoiron/sqlx"
)

//Structure to perform nested set data update for record add/delete.
type sqliteNSOP struct {
    dataObj *sqlData
    conn *sqlx.DB
}

// Update the right end of parent record to accomodate new child.
// Do it recursively until we update all the parents upto the root.
func (nsOp *sqliteNSOP)updateParentRecord(parentObj *sqlData,
                                          updateVal int64) error {
    var err error
    log := logger.GetLoggerInstance()
    parentObj.RgtId = parentObj.RgtId + updateVal
    err = parentObj.updateIds(nsOp.conn)
    if err != nil {
        log.Error("Failed to update the parent %s, err : %s",
                    parentObj.Uid, err)
        return err
    }
    if len(parentObj.Puid) == 0 {
        //No more parent, return now.
        return nil
    }
    ppObj := new(sqlData)
    ppObj.Data = new(dataStore.Data)
    ppObj.Uid = parentObj.Puid
    ppObj.Data, err = ppObj.GetdataById(nsOp.conn)
    if err != nil {
        log.Error("Failed to retrieve parent record %s, err: %s", ppObj.Uid,
                    err)
        return err
    }
    return nsOp.updateParentRecord(ppObj, updateVal)
}

func(nsOp *sqliteNSOP)updateAllRecords(recs []dataStore.Data,
                                updateVal int64) error {
    var err error
    log := logger.GetLoggerInstance()
    var recObj *sqlData
    recObj = new(sqlData)
    err = nil
    for _,row := range recs {
        row.LftId = row.LftId + updateVal
        row.RgtId = row.RgtId + updateVal
        recObj.Data = &row // Assign row to update.
        err = recObj.updateIds(nsOp.conn)
        if err != nil {
            log.Error("Failed to update the record %s, err %s",
                        recObj.Uid, err)
            //Dont return now, lets continue update to other records.
        }
    }
    return err
}

//Update every record after the specific record is added to the database.
//Its an expensive operation as it need to go through every record to update.
func(nsOp *sqliteNSOP)UpdateTreeOnAdd() error {
    log := logger.GetLoggerInstance()
    //Get all the records to be updated on a add.
    recs, err := nsOp.dataObj.GetAllRecsGreaterThanId(nsOp.conn)
    if err != nil {
        log.Error("Failed to get all the recods with NS-id > %d, err: %s",
                    nsOp.dataObj.LftId, err)
        return err
    }
    if len(recs) == 0 {
        log.Trace("No records to update with NS-id > %d", nsOp.dataObj.LftId)
        return nil
    }
    //Update the LftId, rgtId of all the nodes on add with +2.
    nsOp.updateAllRecords(recs, 2)
    return nil
}

//Set the left and right child with its new limit values.
func(nsOp *sqliteNSOP)updateNSListLimitsOnAdd() error{
    log := logger.GetLoggerInstance()
    var parentObj *sqlData
    var err error
    err = nil
    if len(nsOp.dataObj.Puid) == 0 {
        //Something is wrong. Cannot find the parent for new node.
        log.Error("Invalid parent for record to update the NS list limits")
        return appErrors.INVALID_INPUT
    }
    parentObj = new(sqlData)
    parentObj.Data = new(dataStore.Data)
    parentObj.Uid = nsOp.dataObj.Puid
    parentObj.Data, err = parentObj.GetdataById(nsOp.conn)
    if err != nil {
        log.Error("Failed to get the parent record err : %s", err)
        return err
    }
    // Adding a new node with left and right id.
    // Update the record with limits.
    nsOp.dataObj.LftId = parentObj.RgtId
    nsOp.dataObj.RgtId = nsOp.dataObj.LftId + 1
    err = nsOp.dataObj.updateIds(nsOp.conn)
    if err == nil {
        err = nsOp.updateParentRecord(parentObj, 2)
    }
    //Update all other records in the table.
    err = nsOp.UpdateTreeOnAdd()
    return err
}

func(nsOp *sqliteNSOP)deleteAllChildrens(recs []dataStore.Data) error {
    var err error
    log := logger.GetLoggerInstance()
    var recObj *sqlData
    recObj = new(sqlData)
    err = nil
    for _,row := range recs {
        recObj.Data = &row // Assign row to update.
        err = recObj.deleteDataOnId(nsOp.conn)
        if err != nil {
            log.Error("Failed to delete the record %s, err %s",
                        recObj.Uid, err)
            //Dont return now, lets continue delete other records.
        }
    }
    return err
}

//Function to update the nested set values on deleting an node from the db.
// On deleting a node, we must need to delete all its children and update other
// records in the system.
func(nsOp *sqliteNSOP)updateNSListLimitsOnDel() error {
    var diff int64
    var err error
    var rows []dataStore.Data
    log := logger.GetLoggerInstance()
    diff = nsOp.dataObj.RgtId - nsOp.dataObj.LftId
    if diff > 1 {
        // Deleting a node with childrens, Need to delete all its children first
        rows, err = nsOp.dataObj.GetAllChildrens(nsOp.conn)
        if err != nil {
            log.Error("Failed to get childrens of %s record err : %s",
                       nsOp.dataObj.Uid, err)
            return err
        }
        if len(rows) == 0 {
            log.Error("Corrupted tree in the system, as " +
                      "parent node %s dont have a children", nsOp.dataObj.Uid)
            return appErrors.INVALID_STATE
        }
        //Delete all the childrens now.
        nsOp.deleteAllChildrens(rows)
    }
    //Update the parent records for the delete.
    var parentObj *sqlData
    parentObj = new(sqlData)
    parentObj.Data = new(dataStore.Data)
    parentObj.Uid = nsOp.dataObj.Puid
    parentObj.Data, err = parentObj.GetdataById(nsOp.conn)
    //Continue the error as its need to update other records.
    if err != nil {
        log.Error("Failed to get the parent record %s for deleting record %s", 
                   parentObj.Uid, nsOp.dataObj.Uid)
    }
    nsOp.updateParentRecord(parentObj, -(diff+1))
    //Get all records that need to update for the delete.
    rows, err = nsOp.dataObj.GetAllRecsGreaterThanId(nsOp.conn)
    if err != nil {
        log.Error("Failed to find records to update in delete op of %s err: %s",
                   nsOp.dataObj.Uid, err)
        return err
    }
    if len(rows) != 0 {
        // Lets update all the relevant records for the delete.
        err = nsOp.updateAllRecords(rows, -(diff + 1))
         if err != nil {
             log.Error("Failed to update records in the sytem on delete of %s" +
                       " System can be in invalid state err :%s",
                       nsOp.dataObj.Uid, err)
         }
    }
    return nil
}

//use this function to initialize the sqliteNSOP. Do not use just new to create
// the objs.
func NewSqliteNestedSet(dataObj *sqlData, conn *sqlx.DB) *sqliteNSOP {
    NSObj := new(sqliteNSOP)
    NSObj.dataObj = dataObj
    NSObj.conn = conn
    return NSObj
}
