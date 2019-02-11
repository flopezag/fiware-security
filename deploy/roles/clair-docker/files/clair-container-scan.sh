#!/bin/sh

usage() {
    echo "Usage: $0 [-pv] [IMAGE_NAME]"
    echo
    echo "Options:"
    echo " -p : Pull images before running scan"
    echo " -v : verbose output"
    exit 1
}

redirect_stderr() {
    if [ "$VERBOSE" = 1 ]; then
        "$@"
    else
        "$@" 2>/dev/null
    fi
}

redirect_all() {
    if [ "$VERBOSE" = 1 ]; then
        "$@"
    else
        "$@" 2>/dev/null >/dev/null
    fi
}

PULL=0
VERBOSE=0

while getopts ":phv" opt; do
    case $opt in
        p)
            PULL=1
            ;;
        v)
            VERBOSE=1
            ;;
        \?)
            echo "Invalid option: -$OPTARG" >&2
            usage
            ;;
        h)
            usage
            ;;
    esac
done
shift $(($OPTIND -1))

BASEDIR=$(cd $(dirname "$0") && pwd)
cd "$BASEDIR"

if [ ! -f "docker-compose.yml" ]; then
    wget -q https://raw.githubusercontent.com/usr42/clair-container-scan/master/docker-compose.yml
fi

echo "Pulling from "$@"..."
redirect_all docker pull "$@"
echo

id=$(docker images | grep -E "$@" | awk -e '{print $3}')
labels=$(docker inspect --type=image "$@" 2>/dev/null | jq .[].Config.Labels)

if [ "$PULL" -eq 1 ];
then
  echo "Pulling Clair content ..."
  redirect_all docker-compose pull
  echo
fi

echo "Security analysis of "$@" image..."
redirect_stderr docker-compose run --rm scanner "$@" > a.json
ret=$?
echo

echo "Removing docker instances..."
redirect_all docker-compose down
echo

echo "Clean up the docker image..."
redirect_all docker rmi ${id}
echo

exit ${ret}
