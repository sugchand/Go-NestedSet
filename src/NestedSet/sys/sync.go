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

package sys

import (
    "sync"
)

// Do not create multiple objects of this structure as it used to track the
// threads, so use the global obj everywhere in the code.
type Sync struct {
    // WaitGroup to keep track of threads that are currently running.
    appWaitGroups sync.WaitGroup
}

var appSync = new(Sync)
var once sync.Once

func(syncObj *Sync)InitSyncParams() {
    once.Do(func() {
        //Do the initialization for all the relevant parameters.
    })
}

// Any goroutine invocation must precede with with this function.
// It allows the bookkeeping of currnetly running goroutines in the application.
func (syncObj *Sync)AddRoutineInWaitGroup() {
    syncObj.appWaitGroups.Add(1)
}

// Call when exiting the goroutine after its executing.
// It allows the book-keeping of active gorotuines in the application.
// NEVER INVOKE ExitRoutineInWaitGroup without AddRoutineInWaitGroup
func (syncObj *Sync)ExitRoutineInWaitGroup() {
    syncObj.appWaitGroups.Done()
}

// Function to wait for all the goroutines to complete execution.
// ONLY INVOKED FROM MAIN THREAD AS A LAST STATEMENT.
func (syncObj *Sync)JoinAllRoutines() {
    syncObj.appWaitGroups.Wait()
}

func(syncObj *Sync)DestoryAllRoutines() {
    //Destroy all the routines that are currently running.
}

//Function to get the application level syncObj.
func GetAppSyncObj() *Sync{
    appSync.InitSyncParams()
    return appSync
}