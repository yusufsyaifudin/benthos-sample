#!/bin/bash

seedKafka () {
    if [ "$#" -ne 1 ]; then
      echo "kcat: illegal number of parameters"
      echo "Need {message}"
      exit 1
    fi

    MESSAGE=$1

    # send into specific partition because if push data into different partition,
    # consumers will get ordered message foreach partition, but not all partition.
    # this causes different calculation will be made, because message from each partition may come in parallel
    echo $MESSAGE | kcat -P -t calcservice -b localhost:9092 -p 0
}


seedCurl () {
    if [ "$#" -ne 1 ]; then
      echo "curl: illegal number of parameters"
      echo "Need {message}"
      exit 1
    fi

    MESSAGE=$1
    msg=$(printf '{"mathString":"%s"}' "$MESSAGE")
    curl -L -X POST 'localhost:8080/' -H 'Content-Type: application/json' --data-raw "$msg"
}

if [ "$#" -ne 1 ]; then
  echo "Illegal number of parameters"
  echo "Need {file_input}"
  exit 1
fi

FILE_INPUT=$1
while IFS= read -r line
do
  seedKafka $line
  echo "$line"
done < "$FILE_INPUT"