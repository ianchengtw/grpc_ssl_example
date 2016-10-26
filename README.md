# grpc_ssl_example

Generate OpenSSL Keys (private and public)

	openssl req -x509 -newkey rsa:2048 -keyout key.pem -out cert.pem -days XXX -nodes

Convert proto to Golang file

	protoc -I helloworld/ helloworld/helloworld.proto --go_out=plugins=grpc:helloworld

