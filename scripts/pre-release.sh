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
Usage: $0 [<version>]
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
    CLI_VERSION=$(grep 'version = ' "$REPO_DIR/pkg/meta/version.go" | sed 's/.*version = "\(.*\)".*/\1/')
}

function update_cli_version {
  sed -i "s/version = \".*\"/version = \"$VERSION\"/" "$REPO_DIR/pkg/meta/version.go"
}

function update_app_version {
  sed -i "s/const appVersion = \".*\"/const appVersion = \"$VERSION\"/" "$REPO_DIR/web/src/components/AppNav/index.tsx"
}

# Updating version ...
current_cli_version
if [[ -z "$VERSION" ]] ; then
  >&2 info "Current version:"
  echo "$CLI_VERSION"
  exit 0
fi

info "Updating version $CYAN$CLI_VERSION$WHITE -> $CYAN$VERSION"

update_cli_version || echoe "CLI version update failed" RED -1
update_app_version || echoe "APP version update failed" RED -1
