#!/bin/bash

if ! command -v soda &> /dev/null; then
    # Package is not installed, install it
    echo "Installing SODA"
    go install github.com/gobuffalo/pop/v6/soda@latest
else
    echo "SODA is already installed."
fi

soda migrate -e "production"
go build -o Quiklink_BE cmd/*.go
./Quiklink_BE -dbhost=quiklink-postgresql.cp4eu20ycadc.us-west-1.rds.amazonaws.com -dbport=5432 -dbname=Quiklink_BE -dbuser=postgres -dbpass=Rahul1234 -dbssl=prefer