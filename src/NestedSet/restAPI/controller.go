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

package restAPI

import (
    "net/http"
    "encoding/json"
    "io"
    "io/ioutil"
    "github.com/gorilla/mux"
    "NestedSet/logger"
    "NestedSet/dataStore"
    "NestedSet/dataStore/dataSetImpl"
    "NestedSet/appErrors"
)

type controller struct { }

func (ctrl *controller) getAllRecords(w http.ResponseWriter, r *http.Request) {
    log := logger.GetLoggerInstance()
    dbObj := dataSetImpl.GetDataSetObj()
    if dbObj == nil {
        log.Error("Empty datastore handle")
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("500-Server Error "))
    }
    rows, err := dbObj.GetAllRecords()
    if err != nil {
        log.Trace("Failed to get the records from DB")
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("500-Server Error "+ err.Error()))
        return
    }
    data, _ := json.Marshal(rows)
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.WriteHeader(http.StatusOK)
    w.Write(data)
}

func (ctrl *controller) getRecord(w http.ResponseWriter, r *http.Request) {
    var err error
    vars := mux.Vars(r)
    log := logger.GetLoggerInstance()
    Uid := vars["record-id"]
    if len(Uid) == 0 {
        log.Error("Empty record id , cannot find it")
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    var dataObj *dataStore.Data
    dbObj := dataSetImpl.GetDataSetObj()
    dataObj,err = dbObj.GetRecord(Uid)
    if err != nil || dataObj == nil {
        log.Error(`Failed to retrieive the record object %s` +
                    `err : %s`, Uid, err)
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    data, _ := json.Marshal(dataObj)
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.WriteHeader(http.StatusOK)
    w.Write(data)
    log.Trace("Getting a single record in the system")
}

func (ctrl *controller) getRecordsByName(w http.ResponseWriter,
                                         r *http.Request) {
    var err error
    vars := mux.Vars(r)
    log := logger.GetLoggerInstance()
    name := vars["record-name"]
    if len(name) == 0 {
        log.Error("Empty record name , cannot find it")
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    var rows []dataStore.Data
    dbObj := dataSetImpl.GetDataSetObj()
    rows,err = dbObj.GetRecordByName(name)
    if err != nil || len(rows) == 0{
        log.Error(`Failed to retrieive the record object: %s
                    err : %s`, name, err)
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    data, _ := json.Marshal(rows)
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.WriteHeader(http.StatusOK)
    w.Write(data)
}

func (ctrl *controller) addRecord(w http.ResponseWriter, r *http.Request) {
    log := logger.GetLoggerInstance()
    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
    if err != nil {
        log.Error("Failed to read request,")
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    if err := r.Body.Close(); err != nil {
        log.Error("Failed to close the request.")
    }
    dataObj := new(dataStore.Data)
    if err := json.Unmarshal(body, &dataObj); err != nil {
        w.WriteHeader(422)
        log.Error("Failed to Unmarshal the camera input err:%s", err)
        if err := json.NewEncoder(w).Encode(err); err != nil {
            log.Error("Failed to encode marshaling err : %s", err)
            w.WriteHeader(http.StatusInternalServerError)
            return
        }
    }
    dbObj := dataSetImpl.GetDataSetObj()
    err = dbObj.CreateRecord(dataObj)
    if err != nil {
        log.Error("REST API failed to create data entry in table err :%s", err)
        if err == appErrors.DATA_PRESENT_IN_SYSTEM {
            w.WriteHeader(400) //Bad Request.
            return
        }
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusCreated)
    log.Trace("Added a record successfully")
}

func (ctrl *controller) deleteRecord(w http.ResponseWriter, r *http.Request) {
    var err error
    vars := mux.Vars(r)
    log := logger.GetLoggerInstance()
    Uid := vars["record-id"]
    if len(Uid) == 0 {
        log.Error("Empty record id , cannot find it")
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    var dataObj *dataStore.Data
    dbObj := dataSetImpl.GetDataSetObj()
    dataObj,err = dbObj.GetRecord(Uid)
    if err != nil || dataObj == nil {
        log.Error(`Failed to retrieive the record object, cannot delete %s
                    err : %s`, Uid, err)
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    err = dbObj.DeleteRecord(Uid)
    if err != nil {
        log.Error("Failed to delete the data record err : %s", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
}