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
    "github.com/gorilla/mux"
    "NestedSet/logger"
)

// Route defines a route
type routeEntry struct {
    Name        string
    Method      string
    Pattern     string
    HandlerFunc http.HandlerFunc
}

type Routes struct {
    entries []routeEntry
    controller *controller
}

func (routeObj *Routes) CreateAllRoutes() {
    log := logger.GetLoggerInstance()
    routeObj.entries = make([]routeEntry, 5)
    routeObj.entries[0] = routeEntry{
                            "getAllRecords",
                            "GET",
                            "/data",
                            routeObj.controller.getAllRecords}
    routeObj.entries[1] = routeEntry{
                            "getRecord",
                            "GET",
                            "/data/id/{record-id}",
                            routeObj.controller.getRecord}
    routeObj.entries[2] = routeEntry{
                            "getRecordByName",
                            "GET",
                            "/data/name/{record-name}",
                            routeObj.controller.getRecordsByName}
    routeObj.entries[3] = routeEntry{
                            "addRecord",
                            "POST",
                            "/data",
                            routeObj.controller.addRecord}
    routeObj.entries[4] = routeEntry{
                            "deleteRecord",
                            "DELETE",
                            "/data/id/{record-id}",
                            routeObj.controller.deleteRecord}
    log.Trace("rest api routes are defined successfully")
}

// NewRouter function configures a new router to the API
func (routeObj *Routes)NewRouter() *mux.Router {
    log := logger.GetLoggerInstance()
    router := mux.NewRouter().StrictSlash(true)
    routeObj.CreateAllRoutes()
    for _, route := range routeObj.entries {
        var handler http.Handler
        handler = route.HandlerFunc
        router.
         Methods(route.Method).
         Path(route.Pattern).
         Name(route.Name).
         Handler(handler)
        log.Trace("Created route for %s", route.Name)
    }
    routeObj.controller = new(controller)
    return router
}