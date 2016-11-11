```
cp -r ../server .
docker build --rm -t gtfierro/ttlviewer
docker push gtfierro/ttlviewer
# running
docker run -p 80:80 gtfierro/ttlviewer
```
