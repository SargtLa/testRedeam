# testRedeam

Simple server with PostgreSQL database. 
Server mainly depends on httpgo server app build on fasthttp.
Database handlers uses dbEngine repo as main connector.

There are .ddl files, that are read and executed by application to create database table needed at the start of the server.

Unit tests could be run anywhere when the server is running. If other port for server container is  chosen it should be statet as flag for go test command.

The .env file is a sample. If postgres image already exists on target machine, env variables should equals the real envs from your postgres image or image must be recreated with new credentials. If postgres image is created at mthe first time, any credencials could be used to create new database.

