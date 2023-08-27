#!/bin/bash

folders=("golang" "python")

for folder in "${folders[@]}"
do
  if [ -d "$folder" ]; then
    cd "$folder"

    case $folder in
      "golang")
        go mod tidy
        ;;
      "python")
        pip install -r requirements.txt
        ;;
    esac

    cd .. 
  fi
done