check_install:
	which swagger || GO111MODULE=off go get -u github.com/go-swagger/go-swagger/cmd/swagger

swagger: check_install
	GO111MODULE=off swagger generate spec -o ./public/docs/swagger.yml --scan-models

generate_client:
	swagger generate client -f ./public/docs/swagger.yml -A product-api

generate_wire:
  cd internal/dic && wire

generate_mock:
  mockgen -source ./app/repository/userRepository.go -destination ./mock/UserRepository.go -package=mock