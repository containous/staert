#!/bin/sh

curl -i -H "Accept: application/json" -X PUT -d "28"    http://localhost:8500/v1/kv/test/ptrstruct1/s1int
curl -i -H "Accept: application/json" -X PUT -d "28"    http://localhost:8500/v1/kv/test/durationfield
# curl -i -H "Accept: application/json" -X PUT -d "10"                          http://localhost:8500/v1/kv/test
# curl -i -H "Accept: application/json" -X PUT -d "http://172.17.0.3:80"        http://localhost:8500/v1/kv/test
# curl -i -H "Accept: application/json" -X PUT -d "1"                           http://localhost:8500/v1/kv/test
