check_install:
	which swagger || GO111MODULE=off go get -u github.com/go-swagger/go-swagger/cmd/swagger

swagger: check_install
	GO111MODULE=off swagger generate spec -o ./public/swagger.yml --scan-models

generate_client:
	cd sdk && swagger generate client -f ../public/swagger.yml -A product-api