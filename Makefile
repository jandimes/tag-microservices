proto-service-meta:
	protoc -I services/meta/proto/ meta_service.proto \
	 --go_out=${GOPATH}/src \
	 --go-grpc_out=${GOPATH}/src

proto-repo-meta:
	protoc -I repositories/meta/proto/ meta_repository.proto \
	--go_out=${GOPATH}/src \
	--go-grpc_out=${GOPATH}/src

subscribe-meta-systems:
	grpcurl -plaintext -msg-template -d '@' localhost:5006 repositories.meta.proto.MetaRepository/SubscribeSystems
