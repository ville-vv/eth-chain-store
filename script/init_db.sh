#!/usr/bin/env bash

root=$1

env=$2

opt=$3

cd ${root}/src/cmd/db-migrate

if [[ ! ${env} ]];then
    env=test
fi

echo "${opt}"

if [[ "$opt" == "del" ]];then

  DB_ENV=${env}  go run main.go drop
  exit

elif [[ "$opt" == "create" ]]; then
    DB_ENV=${env}  go run main.go create
fi

DB_ENV=${env}  go run main.go migrate
