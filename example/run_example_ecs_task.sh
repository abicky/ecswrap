#!/bin/bash

set -eo pipefail

REPOSITORY_NAMESPACE=ecswrap-example
TASK_FAMILY=ecswrap-example

dir=$(cd $(dirname $0) && pwd)
cd $dir

$(aws ecr get-login --no-include-email)
for container in logger fluentd; do
  repo=$REPOSITORY_NAMESPACE/$container
  repo_uri=$((aws ecr describe-repositories --repository-names $repo | jq -r '.repositories[] | .repositoryUri') || true)
  if [ -z "$repo_uri" ]; then
    echo "Create ECR repository '$repo'"
    repo_uri=$(aws ecr create-repository --repository-name $repo | jq -r '.repository.repositoryUri')
  fi

  echo "Build logger docker image"
  docker build -t $repo $dir/$container

  echo "Push $repo_uri:latest"
  docker tag $repo:latest $repo_uri:latest
  docker push $repo_uri:latest

  declare ${container}_repo_uri=$repo_uri
done

aws ecs register-task-definition --cli-input-json "$(cat <<JSON
{
  "family": "$TASK_FAMILY",
  "containerDefinitions": [
    {
      "name": "logger",
      "image": "$logger_repo_uri",
      "essential": true,
      "links": ["fluentd:fluentd"]
    },
    {
      "name": "fluentd",
      "image": "$fluentd_repo_uri",
      "essential": true,
      "environment": [{"name": "ECSWRAP_LINKED_CONTAINERS", "value": "logger"}]
    }
  ],
  "cpu": "128",
  "memory": "256"
}
JSON
)"

aws ecs run-task --task-definition $TASK_FAMILY
