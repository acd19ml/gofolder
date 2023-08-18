check_install:
	which swagger || GO111MODULE=OFF go get -u github.com/go-swagger/go-swagger/cmd/swagger	

swagger: 
	swagger generate spec -o ./swagger.yaml --scan-models

	swagger generate client -f ./swagger.yaml -A swagger