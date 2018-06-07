#!/bin/bash
git pull
cp -r ../server .
cd server ; make ; cd -
docker build --rm -t gtfierro/ttlviewer .
docker push gtfierro/ttlviewer
