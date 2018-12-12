# NestedSet
Implementation of nestedset model on RDBMS using go. It offers REST APIs to
operate on nested set data. please read in the following links to know more
about what is nestedset and why we needed in RDBMS.

[NestedSet wiki](https://en.wikipedia.org/wiki/Nested_set_model)

[Nestedset Implementation in RDBMS](http://mikehillyer.com/articles/managing-hierarchical-data-in-mysql/)

# Prerequisite
Install sqlite3, as the application uses sqlite3 as a default backend.

```
        apt install sqlite3
```

The program has its own library dependencies to interact with different modules.
We use 'dep' package manager to manage the library dependacy. Install the 'dep'
with following command,

```
        apt install go-dep
```

Under the source tree(not under the application root folder), please run the
following command once to download all the relevant dependanices. You may get
compilation errors, if all the libraries are not present in the system.

```
    dep ensure
```

# Run the application
compile the application as below in the top directory.

```
    make
```

After the compilation, simply run the following command to start the application

```
    ./bin/NestedSet
```

# Supported REST APIs

#### Get all the records in the system.

* Request(GET)

```
http://localhost:8080/data
```

* Response

```

[
    {
        "uid": "00112233-4455-6677-8899-aabbccddeeff",
        "puid": "",
        "name": "root",
        "desc": "The default root node in the hierarchy",
        "lftId": 1,
        "rgtId": 6
    },
    {
        "uid": "bc5ca89d-696a-45f1-914d-e9d7d78b2067",
        "puid": "00112233-4455-6677-8899-aabbccddeeff",
        "name": "B",
        "desc": "Sugesh is the record",
        "lftId": 2,
        "rgtId": 5
    },
    {
        "uid": "8fad71a0-bae3-49fb-a587-37abfd414554",
        "puid": "bc5ca89d-696a-45f1-914d-e9d7d78b2067",
        "name": "B",
        "desc": "Sugesh is the record",
        "lftId": 3,
        "rgtId": 4
    }
]    
```

#### Get a record with ID

* Request (GET)

```
http://localhost:8080/data/id/bc5ca89d-696a-45f1-914d-e9d7d78b2067
```

* Response

```
{
    "uid": "bc5ca89d-696a-45f1-914d-e9d7d78b2067",
    "puid": "00112233-4455-6677-8899-aabbccddeeff",
    "name": "B",
    "desc": "Sugesh is the record",
    "lftId": 2,
    "rgtId": 5
}
```

#### Get records with Name

* Request(GET)

```
http://localhost:8080/data/name/B
```

* Response

```
[
    {
        "uid": "bc5ca89d-696a-45f1-914d-e9d7d78b2067",
        "puid": "00112233-4455-6677-8899-aabbccddeeff",
        "name": "B",
        "desc": "Sugesh is the record",
        "lftId": 2,
        "rgtId": 5
    },
    {
        "uid": "8fad71a0-bae3-49fb-a587-37abfd414554",
        "puid": "bc5ca89d-696a-45f1-914d-e9d7d78b2067",
        "name": "B",
        "desc": "Sugesh is the record",
        "lftId": 3,
        "rgtId": 4
    }
]
```

#### Add a new record to the system

* Request(POST)
 
 ```
 http://localhost:8080/data
 
 {  "name":"M",
    "desc":"Sugesh is the record",
    "puid":"bc5ca89d-696a-45f1-914d-e9d7d78b2067"
}
 ```
 
 * Response
 
 ```
    201 Created 
    {Error code otherwise}
```

#### Delete a record from a system

* Request(DELETE)

```
http://localhost:8080/data/id/ee0dce03-3be5-47ed-9807-d291ed3cfbb7
```

* Response

```
    200 STATUS OK
    { Error code incase operation failed}
```
