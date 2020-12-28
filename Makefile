
.PHONY:
build-rlpgen-tests:
	go run rlpgen/*.go --path ./rlpgen/tests/types.go --objs Test1,Header,Transaction,Body,Block,Receipt,Log --output ./rlpgen/tests/types_encoding.go
