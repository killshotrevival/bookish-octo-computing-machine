
UPPER_ENV=$(echo $ENV | tr '[:lower:]' '[:upper:]')

if [[ $UPPER_ENV == "STAGING" ]]
then
    echo "Working on staging env"
elif [[ $UPPER_ENV == "PRODUCTION" ]]
then
    echo "Working with production env"
else
    echo "Invalid environment found, exiting..."
    exit 1
fi

if [[ "$CI_PIPELINE_SOURCE" != "web" ]]
then
    echo "Checking commit message $CI_COMMIT_MESSAGE"

    if [[ $(expr match "$CI_COMMIT_MESSAGE" '\[deploy\].*$') != 0 ]]
    then
    echo "Staging deploy command received"
    FINAL_VERSION="$(echo $(( $RANDOM % 50 + 1 ))).$(echo $(( $RANDOM % 50 + 1 ))).$(echo $(( $RANDOM % 50 + 1 )))-${CI_COMMIT_AUTHOR:0:3}"
    else
    echo "No deploy command received"
    FINAL_VERSION="$VERSION"
    fi
else
    echo "Not checking commit message as the pipeline source is web"
    FINAL_VERSION="$VERSION"
fi

if [[ "$FINAL_VERSION" == "Version Here" ]]
then
    echo "No version provided. Exiting"
    exit 1
fi

echo "Will be building for version $FINAL_VERSION"

echo "FINAL_VERSION=$FINAL_VERSION" >> build.env

echo "If any not intended info found, please stop the scan here..."
echo "Sleeping for 10 seconds"

sleep 10