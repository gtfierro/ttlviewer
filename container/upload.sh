#!/bin/bash
git pull
cp -r ../server .
cd server ; go build ; cd -
docker build --rm -t gtfierro/ttlviewer .
docker push gtfierro/ttlviewer
