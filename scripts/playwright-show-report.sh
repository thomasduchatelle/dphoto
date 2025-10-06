#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
RED='\033[1;31m'
GREEN='\033[1;32m'
BLUE='\033[1;34m'
CYAN='\033[1;36m'
GRAY='\033[0;37m'
YELLOW='\033[1;33m'
WHITE='\033[0;97m'
NC='\033[0m'

function print_help() {
	cat << EOF
Usage: Show the report of a downloaded artifact from a GitHub Actions

Args:
  <report zip>         ZIP file downloaded from GITHUB

Options:
  --help        print this help and exit
  --debug       enable debug logging
EOF
}
# Args: message [color] [exit=-1]
function echoe() {
	MESS=$1
	if [[ $# -ge 2 ]] ; then
		MESS="${!2}$MESS$NC"
	fi
	echo -e $MESS

	if [[ $# -ge 3 ]] ; then
		exit $3
	fi
}
# Args: message
function debug() {
	if [[ $DEBUG -eq 1 ]] ; then
		echo -e "${CYAN}debug:$GRAY $1$NC"
	fi
}
# Args: message
function info() {
	echo -e "${BLUE}info:$WHITE $1$NC"
}

# Defaults
DEBUG=0
ZIP_FILE=""

# Parse options
while [[ $# -ge 1 ]]
do
	arg="$1"
	case $arg in
	--help|-h) print_help ; exit 0 ;;
	--debug) DEBUG=1 ;;

	*)
		if [[ -z "$FOO" ]] ; then
			ZIP_FILE=$arg
		else
			echoe "Argument not expected: $arg." RED
			print_help
			exit -1
		fi
		;;
	esac

	shift
done

# Decompress the document
if [[ -z "$ZIP_FILE" ]] ; then
  echoe "No ZIP file provided." RED
  print_help
  exit -1
fi

if [[ ! -f "$ZIP_FILE" ]] ; then
  echoe "ZIP file not found: $ZIP_FILE" RED
  exit -1
fi

# Create temp dir in /tmp with datetime up to the second, no spaces
DATETIME=$(date +"%Y%m%d%H%M%S")
TMP_DIR="/tmp/playwright-$DATETIME"
mkdir -p "$TMP_DIR"
debug "Created temporary directory: $TMP_DIR"

unzip -q "$ZIP_FILE" -d "$TMP_DIR"
if [[ $? -ne 0 ]] ; then
  echoe "Error unzipping file: $ZIP_FILE" RED
  exit -1
fi

REPORT_DIR="$TMP_DIR"

if [[ -d "$REPORT_DIR/web/playwright-report" ]] ; then
  REPORT_DIR="$REPORT_DIR/web/playwright-report"
fi

info "Opening report from $ZIP_FILE"
cd "$DIR/../web"
debug "npx playwright show-report \"$REPORT_DIR\""
npx playwright show-report "$REPORT_DIR"


rm -rf "$TMP_DIR"
debug "Removed temporary directory: $TMP_DIR"
