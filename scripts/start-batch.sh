#!/usr/bin/env bash

go get -v

if [ ${IS_DEBUG} == "true" ]
then
   cd cmd && dlv debug --headless --listen=:2345 --api-version=2
else
   fresh -c configs/fresh-batch.conf
fi