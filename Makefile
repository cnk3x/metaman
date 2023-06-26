build::
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o bin/metaman

upload:: build
	scp bin/metaman nas:~/apps/metaman/metaman
