#!/bin/bash

seedCurl () {
    if [ "$#" -ne 1 ]; then
      echo "curl: illegal number of parameters"
      echo "Need {message}"
      exit 1
    fi

    MESSAGE=$1
    curl -L -X POST 'localhost:8080/' -H 'Content-Type: application/json' --data-raw "$MESSAGE"
}

if [ "$#" -ne 1 ]; then
  echo "Illegal number of parameters"
  echo "Need {file_input}"
  exit 1
fi

FILE_INPUT=$1
while IFS= read -r line
do
  seedCurl "$line"
  echo "$line"
done < "$FILE_INPUT"