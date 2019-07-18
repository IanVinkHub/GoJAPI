# db Folder

DB stands for Database, these files are mostly used for connecting and querying the database. In this version its is very limited and only mysql is supported.

## db/connection.go

The connection file contains all functions that are related to connecting with the database. This is also where you config the database for now.

## db/querybuilder.go

The querybuilder file contains all functions that are related to building querys using functions. If you need to change or add SQL syntaxes this is where you would find it.