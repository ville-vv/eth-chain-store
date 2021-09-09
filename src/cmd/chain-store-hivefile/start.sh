#!/usr/bin/env bash

if [ ! -f 'chain-store-hivefiel' ];then
  go build -o chain-store-hivefiel .
fi

./chain-store-hivefiel -max_sync_thread 100 --max_write_thread 100 -rpc_url http://172.16.16.115:8545/ -max_cpu -start_number 3000000 -end_number 3010000

