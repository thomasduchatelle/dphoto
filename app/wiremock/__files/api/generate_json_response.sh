#!/bin/bash

SEARCH=${1?'search directory is mandatory'}

for f in $(find "$SEARCH" -type f) ; do
  filename="$(basename "$f")"
  id="$(basename "$(dirname "$f")")"

  cat << EOF
{
  "id": "$id",
  "type": "IMAGE",
  "filename": "$filename",
  "time": "$(date +'%Y-%m-%dT12:42:07' -d "2022-02-22 - $((10#$id * 5)) hour")",
  "source": "pixel"
},
EOF
done
