# NestedSet
Implementation of nestedset model on RDBMS using go. It offers REST APIs to
operate on nested set data. The know more about what is nestedset and why we
need, please follow the following link

[NestedSet wiki](https://en.wikipedia.org/wiki/Nested_set_model)

[Nestedset Implementation in RDBMS](http://mikehillyer.com/articles/managing-hierarchical-data-in-mysql/)

#Prerequisite
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

#Run the application
To run the application, you must compile the application as below in the top
directory.

```
    make
```

After the compilation, simply run the following command to start the application

```
    ./bin/NestedSet
```