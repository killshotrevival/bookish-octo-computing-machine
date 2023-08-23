
GIT_ID="$CI_COMMIT_SHORT_SHA"
IMAGE_NAME="endgame"
COMPLETE_IMAGE="registry.digitalocean.com/getastra/$IMAGE_NAME:$FINAL_VERSION-$GIT_ID"

UPPER_ENV=$(echo $ENV | tr '[:lower:]' '[:upper:]')

if [[ $UPPER_ENV == "STAGING" ]]
then
    KUBE_CONFIG="$STAGING_CONFIG"
elif [[ $UPPER_ENV == "PRODUCTION" ]]
then
    KUBE_CONFIG="$PRODUCTION_CONFIG"
fi

echo "$KUBE_CONFIG" | base64 -d > temp.config

echo "Updating to cluster for service with image -> $COMPLETE_IMAGE"

kubectl --kubeconfig temp.config create configmap endgame-scan-image --from-literal=image-name=$COMPLETE_IMAGE --dry-run=client -o yaml | kubectl --kubeconfig temp.config apply -f -