#!/bin/bash

gofmt -w *.go && clear && go build && nodemon -x ./backend
