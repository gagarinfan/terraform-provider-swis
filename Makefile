NAME=terraform-provider-ipam

prepare:
	go mod init terraform-provider-ipam && go mod vendor
build:
	go build -o ${NAME}
install: build
	cp -pr ${NAME} ~/.terraform.d/plugins/ipam/0.1/darwin_amd64/${NAME}