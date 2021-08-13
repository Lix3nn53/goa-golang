check_install:
	which swagger || GO111MODULE=off go get -u github.com/go-swagger/go-swagger/cmd/swagger

swagger: check_install
	GO111MODULE=off swagger generate spec -o ./public/docs/swagger.yml --scan-models

generate_client:
	swagger generate client -f ./public/docs/swagger.yml -A product-api

generate_wire:
  cd internal/dic && wire

generate_mock:
  mockgen -source ./app/repository/userRepository/userRepository.go -destination ./mock/UserRepository.go -package=mock
	&& mockgen -source ./app/service/userService/userService.go -destination ./mock/UserService.go -package=mock
	&& mockgen -source ./app/repository/billingRepository/billingRepository.go -destination ./mock/BillingRepository.go -package=mock
	&& mockgen -source ./app/service/billingService/billingService.go -destination ./mock/BillingService.go -package=mock