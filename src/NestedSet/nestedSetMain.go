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
package main

import (
    "syscall"
    "os"
    "os/signal"
    "fmt"
    "NestedSet/logger"
    "NestedSet/sys"
    "NestedSet/restAPI"
    "NestedSet/dataStore/dataSetImpl"
)
////////////////////////////////////////////////////////////////////////////////
const (
    APP_DIR = "/tmp/nestedSet"
    LOGLEVEL_TYPE = logger.Trace
    LOGFILE = APP_DIR + "/nestedSetLog.log"
    SERVER_IP = "127.0.0.1"
    SERVER_PORT = "8080"
    DB_PATH = APP_DIR + "/nestedSet.db"
)
///////////////////////////////////////////////////////////////////////////////

func startLoggerService() {
    createDirectory(APP_DIR, os.FileMode(0755))
    logger := new(logger.Logging)
    logger.LogInitSingleton(LOGLEVEL_TYPE, LOGFILE)
    logger.Trace("Logging service is started..")
}

func setupDataService() error {
    var err error
    dataObj := dataSetImpl.GetDataSetObj()
    err = dataObj.CreateDBConnection(DB_PATH)
    if err != nil {
        return err
    }
    err = dataObj.CreateDataStoreTables()
    return err
}

func createDirectory(path string, mode os.FileMode) {
    if _, err := os.Stat(path); os.IsNotExist(err) {
        os.MkdirAll(path, mode)
    }
}

func setupRESTService() error {
    resthandler := new(restAPI.RestAPI)
    err := resthandler.RestAPIMainHandler(SERVER_IP, SERVER_PORT)
    if err != nil {
        return err
    }
    return nil
}

func main() {
    var err error
    startLoggerService()
    syncObj := sys.GetAppSyncObj()
    defer syncObj.JoinAllRoutines()    
    log := logger.GetLoggerInstance()

    err = setupDataService()
    if err != nil {
        log.Error("Failed to start the database, exiting the application")
        panic("Cannot start Database/backend")
    }
    err = setupRESTService()
    if err != nil {
        log.Error("Failed to start REST service")
        panic("Cannot start REST service")
    }

    // Exit the main thread on Ctrl C
    fmt.Println("\n\n\n *** Press Ctrl+C to Exit *** \n\n\n")
    exitsignal := make(chan os.Signal, 1)
    signal.Notify(exitsignal, syscall.SIGINT, syscall.SIGTERM)
    syncObj.AddRoutineInWaitGroup()
    go func() {
        // Blocking the routine for the exit signal.
        <- exitsignal
        syncObj.ExitRoutineInWaitGroup()
        //Send exit signal to all the goroutines.
        syncObj.DestoryAllRoutines()
        log.Trace("Sending the exit signal to the application")
    }()
    log.Trace("Exiting the application now")
}
