#!/usr/bin/env bash
cd latest
inotifywait \
  -m "." \
  -e create \
  -e modify \
  | while read -r directory events filename
    do 
      echo -e "$filename $events $directory"
      go test
    done
