# sensu-rri-check

```
# set this:
#
export VERSION=v0.0.3


# just paste this
#
export REPONAME=sensu-rri
CGO_ENABLED=0 go build -o bin/sensu-rri cmd/main.go
tar czf /tmp/${REPONAME}_${VERSION}_linux_amd64.tar.gz bin/
export SUM=$(sha512sum /tmp/${REPONAME}_${VERSION}_linux_amd64.tar.gz | cut -d ' ' -f 1)
export FILE=$(basename ${REPONAME}_${VERSION}_linux_amd64.tar.gz)
echo $SUM $FILE > /tmp/${REPONAME}_${VERSION}_sha512_checksums.txt
rm bin/sensu-rri

```