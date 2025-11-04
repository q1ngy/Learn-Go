.PHONY: gen

gen:
	protoc \
	--proto_path=pb \
	--go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	$(shell find pb -iname "*.proto")

