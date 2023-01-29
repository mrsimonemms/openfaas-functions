#!/usr/bin/env bash

set -uo pipefail;

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
FUNCTION_NAME="${1:-}"

if [ -z "${FUNCTION_NAME}" ]; then
  echo "Function name must be set as first argument"
  exit 1
fi

rm -Rf ./build
faas-cli build -f "${SCRIPT_DIR}/../stack.yml" --shrinkwrap
docker build -t "${IMAGE}" "./build/${FUNCTION_NAME}"

if "${PUSH_IMAGE}"; then
  echo "Pushing ${IMAGE}"
  docker push "${IMAGE}"
else
  echo "Not pushing ${IMAGE}"
fi
