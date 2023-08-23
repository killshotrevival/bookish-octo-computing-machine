GIT_ID="$CI_COMMIT_SHORT_SHA"
echo "Building docker image for $IMAGE_NAME with git ID -> $GIT_ID"

IMAGE_NAME="endgame"
DOCKERFILE="Dockerfile"


echo "Login..."
echo "Version found -> $FINAL_VERSION"
docker login registry.digitalocean.com -u $USER -p $PASSWORD

docker build -t registry.digitalocean.com/getastra/$IMAGE_NAME:$FINAL_VERSION-$GIT_ID -f "$DOCKERFILE" .

echo "Pushing docker image"
docker push registry.digitalocean.com/getastra/$IMAGE_NAME:$FINAL_VERSION-$GIT_ID