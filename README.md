# score-implementation-sample

This score-implementation-sample is a template repo for creating a new Score implementation following the conventions laid out in [Score Compose](https://github.com/score-spec/score-compose) and [Score K8s](https://github.com/score-spec/score-k8s).

This sample comes complete with:

1. CLI skeleton including `init` and`generate` subcommands
2. State directory storage
3. TODOs in place of resource provisioning and workload conversion

To adapt this for your target platform, you should:

1. Rename the go module by replacing all instances of `github.com/score-spec/score-implementation-sample` with your own module name.
2. Replace all other instances of `score-implementation-sample` with your own `score-xyz` name including renaming the `cmd/score-implementation-sample` directory.
3. Run the tests with `go test -v ./...`.
4. Change the TODO in [provisioning.go](./internal/provisioners/provisioning.go) to provision resources and set the resource outputs. The existing implementation resolves placeholders in the resource params but does not set any resource outputs.
5. Change the TODO in [convert.go](./internal/convert/convert.go) to convert workloads into the target manifest form. The existing implementation resolves placeholders in the variables and files sections but just returns the workload spec as yaml content in the manifests.

Good luck, and have fun!
