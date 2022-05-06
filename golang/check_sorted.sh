#!/bin/bash

if [ "$#" -ne 1 ]; then
  echo "Illegal number of parameters"
  echo "Need {file_input}"
  exit 1
fi

FILE_INPUT=$1

JQ_FILE=$(cat $FILE_INPUT | jq '.kafkaOffset | tonumber')
TMP=""
for line in $JQ_FILE
do
  if [ $TMP ]
  then
    # https://stackoverflow.com/a/18668580
    # if [ "$TMP" -gt "$line" ]
    if (( $TMP > $line ))
    then
      echo "kafkaOffset: prev offset " $TMP " larger than current offset " $line
    fi

    TMP=$line
  else
    TMP=$line
  fi

done
