# Proto destination folder
SRV_PROTO = server/proto
CLI_PROTO = client/proto
# Proto source
PROTO_FLAG = -I proto/ proto/usermgmt.proto
#GO Flags
GO_OPT = --go_opt=paths=source_relative
GRPC_OPT = --go-grpc_opt=paths=source_relative

generate_server:
	@protoc --go_out=$(SRV_PROTO) $(GO_OPT) --go-grpc_out=$(SRV_PROTO) $(GRPC_OPT) $(PROTO_FLAG)
	@echo '--- Generating Server ---'
generate_client:
	@protoc --go_out=$(CLI_PROTO) $(GO_OPT) --go-grpc_out=$(CLI_PROTO) $(GRPC_OPT) $(PROTO_FLAG)
	@echo '--- Generating Client ---'

generate: presentation generate_server generate_client
	@echo '--- The proto files were successfully generated. ---'

run_server:
	@echo "---- Running Server ----"
	@go run server/*

run_client:
	@echo "---- Running Client ----"
	@go run client/*

presentation:
	@echo '                                                            '
	@echo '              ███████████   ███████████    █████████        '
	@echo '             ░░███░░░░░███ ░░███░░░░░███  ███░░░░░███       '
	@echo '      ███████ ░███    ░███  ░███    ░███ ███     ░░░        '
	@echo '     ███░░███ ░██████████   ░██████████ ░███                '
	@echo '    ░███ ░███ ░███░░░░░███  ░███░░░░░░  ░███                '
	@echo '    ░███ ░███ ░███    ░███  ░███        ░░███     ███       '
	@echo '    ░░███████ █████   █████ █████        ░░█████████        '
	@echo '     ░░░░░███░░░░░   ░░░░░ ░░░░░          ░░░░░░░░░         '
	@echo '     ███ ░███                                               '
	@echo '    ░░██████                                                '
	@echo '    ░░░░░░     Protobuf                                     '
	@echo '                                                            '
