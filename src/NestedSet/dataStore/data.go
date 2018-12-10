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

package dataStore

import (
)

type Data struct {
    Uid string      `json:"uid" db:"Uid"`
    Puid string     `json:"puid" db:"Puid"`
    Name string     `json:"name" db:"Name"`
    Desc string     `json:"desc" db:"Desc"`
    LftId int64      `json:"lftId" db:"LftId"`//Used for nestedset hierarchy
    RgtId int64      `json:"rgtId" db:"RgtId"`//Used for nestedset hierarchy
}

const (
    DEFAULT_SETID = -1
    DEFAULT_LFTID = DEFAULT_SETID
    DEFAULT_RGTID = DEFAULT_SETID
    ROOT_UID = "00112233-4455-6677-8899-aabbccddeeff"
    DEFAULT_PUID = ROOT_UID
)
