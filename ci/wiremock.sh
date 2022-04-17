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
Usage: $0 [--version [version]] [wiremock args]

Download from maven, and run wiremock with specified parameters.

Options:
  --help                 print this help and exit ; do not download wiremock, but display wiremock help if binary is present.
  --debug                enable debug logging
  --version <version>    use wiremock version
  --port <port>          start wiremock on specific port (default $PORT)
  --project <path>       use wiremock definitions (default $PROJECT)
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
VERSION="2.33.1"
PORT="8080"
PROJECT="$(dirname "$DIR")/app/wiremock"
TRACK=0
ARGS=()

# Parse options
while [[ $# -ge 1 ]]
do
	arg="$1"
	case $arg in
	--help|-h)
	  print_help
	  JAR="$DIR/../dist/wiremock-jre8-standalone-${VERSION}.jar"
	  if [[ -e "$JAR" ]] ; then
	    echo ""
	    echo "----------"
	    echo ""
	    java -jar "$JAR" --help
	  fi
	  exit 0
	  ;;
	--debug) DEBUG=1 ;;

	--version)
    VERSION=$2
    shift
		;;

  --port)
    PORT="$2"
    shift
    ;;

  --project)
    PROJECT="$2"
    shift
    ;;

  --track)
    TRACK=1
    ;;

	*)
		ARGS+=("$arg")
		;;
	esac

	shift
done

# Download and run wiremock
JAR="$(dirname "$DIR")/dist/wiremock-jre8-standalone-${VERSION}.jar"
SOURCE="https://repo1.maven.org/maven2/com/github/tomakehurst/wiremock-jre8-standalone/${VERSION}/wiremock-jre8-standalone-${VERSION}.jar"

if [[ ! -f "$JAR" ]] ; then
  info "Download wiremock from $CYAN$SOURCE$WHITE..."
  mkdir -p "$(dirname "$JAR")"
  curl "$SOURCE" -o "$JAR"
fi

cmd=(java -jar "$JAR" --port "$PORT" --root-dir "$PROJECT" "${ARGS[@]}")
if [[ "$TRACK" -eq 0 ]] ; then
  debug "running: ${cmd[*]}"
  "${cmd[@]}"
else
#  trap "trap - SIGTERM && kill -- -$$" SIGINT SIGTERM EXIT
  debug "running: ${cmd[*]}"
  "${cmd[@]}" &
  pid=$!
  debug "PID $pid"
  trap "kill $pid" SIGINT SIGTERM EXIT
  fswatch -o "$PROJECT" | \
      while read; do
        kill "$pid" || echoe "Failed to stop wiremock" RED 1
        "${cmd[@]}" &
        pid=$!
        debug "PID $pid"
        trap "kill $pid" SIGINT SIGTERM EXIT
      done
fi

