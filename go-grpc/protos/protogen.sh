#!/bin/bash

set -euo pipefail

# echo error
echo_error() {
    echo "ERROR: $1"
}

# removes trailing slashes for directory paths
remove_trailing_slash() {
    if [[ $# -eq 0 ]]; then
        echo_error "No argument passed to ${FUNCNAME[0]}()" >&2
        return 1
    elif [[ -z $1 ]]; then
        echo_error "Empty argument passed to ${FUNCNAME[0]}()" >&2
        return 1
    fi

    echo "$1" | sed "s/\/$//"
}

# get these values as options
PROTO_SERVICE_NAME=""
PROTO_GEN_GO_MODULE_NAME=""

while getopts "s:m:" opt; do
    case $opt in
    s) PROTO_SERVICE_NAME=$OPTARG ;;
    m) PROTO_GEN_GO_MODULE_NAME=$OPTARG ;;
    *) printf "Invalide option %s\n" "$opt" ;;
    esac
done

if [[ -z $PROTO_SERVICE_NAME ]]; then
    echo_error "provide service name"
    exit 1
else
    PROTO_SERVICE_NAME="$(remove_trailing_slash "$PROTO_SERVICE_NAME")"
fi

if [[ -z $PROTO_GEN_GO_MODULE_NAME ]]; then
    echo_error "provide go module name"
    exit 1
fi

# script directory path
CURR_DIR="$(dirname -- "$(readlink -f "${BASH_SOURCE[0]}")")"
PROTO_SRC_DIR=$CURR_DIR

# relative proto-gen-go path to $CURR_DIR
PROTO_GEN_GO_DIR_NAME="go-proto"
PROTO_GEN_GO_DIR=$(realpath "$CURR_DIR"/../"$PROTO_GEN_GO_DIR_NAME")

# source service proto path and go generated service path
PROTO_SRC_SERVICE_DIR=$(printf "%s/%s" "$PROTO_SRC_DIR" "$PROTO_SERVICE_NAME")
PROTO_GEN_SERVICE_DIR=$(printf "%s/%s" "$PROTO_GEN_GO_DIR" "$PROTO_SERVICE_NAME")

if [[ ! -d $PROTO_SRC_SERVICE_DIR ]]; then
    echo_error "Source service proto directory not exists"
    exit 1
fi

if [[ ! -d $PROTO_GEN_SERVICE_DIR ]]; then
    echo_error "Go protobuf generation service directory not exists."
    exit 1
fi

# all *.proto files in source service directory
PROTO_FILES=$(find "$PROTO_SRC_SERVICE_DIR" -name "*.proto")

# generate protobuf and grpc code
protoc \
    --proto_path="$PROTO_SRC_SERVICE_DIR" \
    --go_out="$PROTO_GEN_SERVICE_DIR" \
    --go_opt=default_api_level=API_OPAQUE \
    --go_opt=module="$PROTO_GEN_GO_MODULE_NAME" \
    --go-grpc_out="$PROTO_GEN_SERVICE_DIR" \
    --go-grpc_opt=module="$PROTO_GEN_GO_MODULE_NAME" \
    "${PROTO_FILES[@]}"
