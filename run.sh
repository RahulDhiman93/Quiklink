#!/bin/bash

go install github.com/gobuffalo/pop/v6/soda@latest
soda migrate
go build -o Quiklink_BE cmd/*.go
./Quiklink_BE -dbname=Quiklink_BE -dbuser=postgres