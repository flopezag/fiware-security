#!/bin/sh

usage() {
    echo "Usage: $0 [-pv] [IMAGE_NAME]"
    echo
    echo "Options:"
    echo " -p : Pull images before running scan"
    echo " -v : verbose output"
    echo
    echo "[IMAGE_NAME] : (Optional) Docker image file to be analysed."
    echo "               If it is not provided the Docker images are "
    echo "               obtained from the enablers.json file."
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

security_analysis() {
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
    filename=$(echo "$@" | awk -F '/' '{print $2 ".json"}')

    redirect_stderr docker-compose run --rm scanner "$@" > ${filename}
    ret=$?
    echo

    echo "Removing docker instances..."
    redirect_all docker-compose down
    echo

    echo "Clean up the docker image..."
    redirect_all docker rmi ${id}
    echo
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
    wget -q https://raw.githubusercontent.com/flopezag/fiware-clair/develop/docker/docker-compose.yml
fi

if [ ! -f "enablers.json" ]; then
    wget -q https://raw.githubusercontent.com/flopezag/fiware-clair/develop/docker/enablers.json
fi

if [ -n "$1" ]; then
    security_analysis "$1"
else
    for ge in `more enablers.json | jq .enablers[].image | sed 's/"//g'`
    do
      security_analysis ${ge}
      echo
      echo
    done
fi

exit ${ret}
