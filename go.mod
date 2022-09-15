module github.com/leosunmo/docker-gen

go 1.19

require (
	github.com/moby/buildkit v0.10.4
	github.com/spf13/cobra v1.5.0
)

require (
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver/v3 v3.1.1 // indirect
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/containerd/typeurl v1.0.2 // indirect
	github.com/docker/docker v20.10.7+incompatible // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/huandu/xstrings v1.3.1 // indirect
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/mitchellh/copystructure v1.0.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.3-0.20211202183452-c5a74bcca799 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/shopspring/decimal v1.2.0 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/stretchr/testify v1.7.2 // indirect
	golang.org/x/crypto v0.0.0-20211202192323-5770296d904e // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gotest.tools/v3 v3.3.0 // indirect
)

require (
	github.com/Masterminds/sprig/v3 v3.2.2
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	gopkg.in/yaml.v3 v3.0.1
)

replace (
	github.com/Sirupsen/logrus => github.com/sirupsen/logrus v1.9.0
	github.com/docker/docker => github.com/docker/docker v20.10.3-0.20220831131523-b5a0d7a188ac+incompatible // 22.06 branch (v22.06-dev)
)
