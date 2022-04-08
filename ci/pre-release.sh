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
Usage: TODO

Args:
  <foo>         arguments

Options:
  --help        print this help and exit
  --debug       enable debug logging
  --bar=1       with equals
  --baz <baz>   with space
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
REPO_DIR="$(dirname "$DIR")"
VERSION=""

# Parse options
while [[ $# -ge 1 ]]
do
	arg="$1"
	case $arg in
	--help|-h) print_help ; exit 0 ;;
	--debug) DEBUG=1 ;;

	*)
		if [[ -z "$VERSION" ]] ; then
			VERSION=$arg
		else
			echoe "Argument not expected: $arg." RED
			print_help
			exit 1
		fi
		;;
	esac

	shift
done

## Functions

function current_cli_version {
    CLI_VERSION=$(grep 'Version = ' "$REPO_DIR/dphoto/cmd/version.go" | sed 's/Version = "\(.*\)"/\1/')
}

function update_cli_version {
  sed -i "s/Version = \".*\"/Version = \"$VERSION\"/" "$REPO_DIR/dphoto/cmd/version.go"
}

function update_app_version {
  sed -i "s/const appVersion = \".*\"/const appVersion = \"$VERSION\"/" "$REPO_DIR/app/viewer_ui/src/components/app-nav.component/index.tsx"
}

# Updating version ...
current_cli_version
info "Updating version $CYAN$CLI_VERSION$WHITE -> $CYAN$VERSION"

update_cli_version
update_app_version
