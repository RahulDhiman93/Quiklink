#!/bin/bash

service postgresql start &&

soda migrate &&

go build -o Quiklink_BE cmd/*.go && ./Quiklink_BE -dbname=Quiklink_BE -dbuser=postgres