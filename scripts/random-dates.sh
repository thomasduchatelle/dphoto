#!/bin/bash
# Photos can be found in https://unsplash.com/

if [[ $# -ne 1 || "$1" == "--help" || "$1" == "-h" ]] ; then
  echo "Usage: $0 <dir or file to randomise date-time>"
  exit 0
fi

if ! command -v exiftool > /dev/null ; then
  echo "exiftool is required to use this script. Please install it with 'pacman -S exiftool' or 'brew install exiftool'"
  exit 1
fi

function random_datetime() {
    year=$((2024 + $RANDOM % 1))
    month=$((1 + $RANDOM % 4))
    day=$((1 + $RANDOM % 28))
    hour=$(($RANDOM % 24))
    minute=$(($RANDOM % 60))
    second=$(($RANDOM % 60))

    printf '%d:%02d:%02d %02d:%02d:%02d' $year $month $day $hour $minute $second
}

function set_datetime() {
    echo "Setting $2 -> $1 ..."
    exiftool -AllDates="$2" "$1"
}

if [[ -f "$1" ]] ; then
  datetime=$(random_datetime)
  set_datetime "$1" "$datetime"
fi

if [[ -d "$1" ]] ; then
  for f in "$1"/* ; do
    datetime=$(random_datetime)
    set_datetime "$f" "$datetime"
  done
fi
