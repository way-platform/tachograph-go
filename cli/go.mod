module github.com/way-platform/tachograph-go/cli

go 1.25.0

require (
	github.com/spf13/cobra v1.10.1
	github.com/way-platform/tachograph-go v0.17.0
	google.golang.org/protobuf v1.36.10
)

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.36.10-20250912141014-52f32327d4b0.1 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.9 // indirect
	golang.org/x/text v0.29.0 // indirect
)

replace github.com/way-platform/tachograph-go => ..
