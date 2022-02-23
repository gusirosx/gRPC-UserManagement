generate:
	@protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/usermgmt.proto

run_server:
	@echo "---- Running Server ----"
	@go run server/*

run_client:
	@echo "---- Running Client ----"
	@go run client/*

generate_server:
	@protoc --go_out=server/proto --go_opt=paths=source_relative --go-grpc_out=server/proto --go-grpc_opt=paths=source_relative -I proto/ proto/usermgmt.proto
generate_client:
	@protoc --go_out=client/proto --go_opt=paths=source_relative --go-grpc_out=client/proto --go-grpc_opt=paths=source_relative -I proto/ proto/usermgmt.proto