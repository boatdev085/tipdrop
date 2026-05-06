#!/bin/sh
set -eu

: "${S3_BUCKET:=tipdrop-local}"
: "${MINIO_ENDPOINT:=http://minio:9000}"
: "${MINIO_ROOT_USER:=tipdrop}"
: "${MINIO_ROOT_PASSWORD:=tipdrop-secret}"

mc alias set local "${MINIO_ENDPOINT}" "${MINIO_ROOT_USER}" "${MINIO_ROOT_PASSWORD}"
mc mb --ignore-existing "local/${S3_BUCKET}"
mc anonymous set none "local/${S3_BUCKET}"
