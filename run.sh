#!/bin/bash

go build -o app cmd/web/*.go
./app -dbname=golf_league_app -dbuser=jakedonahue -cache=false -production=false