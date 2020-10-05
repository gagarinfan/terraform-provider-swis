#NAME=terraform-provider-ipam
NAME=terraform-provider-swis

prepare:
	go mod init ${NAME} && go mod vendor
build:
	go build -o ${NAME}
install: build
	cp -pr ${NAME} ~/.terraform.d/plugins/swis/0.1/darwin_amd64/${NAME}