module github.com/VertebrateResequencing/wr

require (
	cloud.google.com/go v0.53.0 // indirect
	code.cloudfoundry.org/bytefmt v0.0.0-20200131002437-cf55d5288a48
	github.com/Microsoft/go-winio v0.4.14 // indirect
	github.com/StackExchange/wmi v0.0.0-20190523213315-cbe66965904d // indirect
	github.com/VertebrateResequencing/muxfys/v4 v4.0.2
	github.com/VividCortex/ewma v0.0.0-20170804035156-43880d236f69
	github.com/carbocation/runningvariance v0.0.0-20150817162428-fdcce8a03b6b
	github.com/dgryski/go-farm v0.0.0-20200201041132-a6ae2369ad13
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v0.0.0-20180524003928-df5175e1ee95
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/docker/spdystream v0.0.0-20181023171402-6480d4af844c // indirect
	github.com/elazarl/goproxy v0.0.0-20190421051319-9d40249d3c2f // indirect
	github.com/elazarl/goproxy/ext v0.0.0-20191011121108-aa519ddbe484 // indirect
	github.com/fanatic/go-infoblox v0.0.0-20190709161059-e25f3820238c
	github.com/fatih/color v1.9.0
	github.com/go-ini/ini v1.52.0 // indirect
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/gofrs/uuid v3.2.0+incompatible
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/golang/protobuf v1.3.4 // indirect
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/googleapis/gnostic v0.4.1 // indirect
	github.com/gophercloud/gophercloud v0.8.0
	github.com/gorilla/websocket v1.4.1
	github.com/gotestyourself/gotestyourself v2.2.0+incompatible // indirect
	github.com/grafov/bcast v0.0.0-20190217190352-1447f067e08d
	github.com/hashicorp/go-multierror v1.0.0
	github.com/hashicorp/golang-lru v0.5.4
	github.com/howeyc/gopass v0.0.0-20190910152052-7cb4b85ec19c // indirect
	github.com/imdario/mergo v0.3.8 // indirect
	github.com/inconshreveable/log15 v0.0.0-20200109203555-b30bc20e4fd1
	github.com/jinzhu/configor v1.1.1
	github.com/jpillora/backoff v1.0.0
	github.com/json-iterator/go v1.1.9 // indirect
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0
	github.com/mattn/go-colorable v0.1.6 // indirect
	github.com/minio/minio-go/v6 v6.0.49 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/onsi/ginkgo v1.8.0 // indirect
	github.com/onsi/gomega v1.5.0 // indirect
	github.com/opencontainers/go-digest v0.0.0-20180430190053-c9281466c8b2 // indirect
	github.com/opencontainers/image-spec v0.0.0-20180411145040-e562b0440392 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pkg/sftp v1.11.0
	github.com/ricochet2200/go-disk-usage v0.0.0-20150921141558-f0d1b743428f
	github.com/sasha-s/go-deadlock v0.2.1-0.20190427202633-1595213edefa
	github.com/sb10/l15h v0.0.0-20170510122137-64c488bf8e22
	github.com/sb10/waitgroup v0.0.0-20200305124406-7ed665007efa
	github.com/sevlyar/go-daemon v0.1.5
	github.com/shirou/gopsutil v2.20.2+incompatible
	github.com/smartystreets/goconvey v0.0.0-20190731233626-505e41936337
	github.com/spf13/cobra v0.0.6
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/ugorji/go/codec v1.1.7
	go.etcd.io/bbolt v1.3.3
	golang.org/x/crypto v0.0.0-20200302210943-78000ba7a073
	golang.org/x/net v0.0.0-20200301022130-244492dfa37a // indirect
	golang.org/x/sys v0.0.0-20200302150141-5c8b2ff67527 // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/ini.v1 v1.52.0 // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
	gotest.tools v2.2.0+incompatible // indirect
	k8s.io/api v0.17.3
	k8s.io/apimachinery v0.17.3
	k8s.io/client-go v11.0.0+incompatible
	nanomsg.org/go-mangos v0.0.0-20180815160134-b7ff4263f0d7
)

replace k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20180228050457-302974c03f7e

replace k8s.io/api => k8s.io/api v0.0.0-20180308224125-73d903622b73

replace k8s.io/client-go => k8s.io/client-go v7.0.0+incompatible

// this version of gnostic needed to work with v7 of k8s.io/client-go
replace github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.4.1-0.20200130232022-81b31a2e6e4e

replace github.com/grafov/bcast => github.com/grafov/bcast v0.0.0-20161019100130-e9affb593f6c

replace github.com/sevlyar/go-daemon => github.com/sevlyar/go-daemon v0.1.1-0.20160925164401-01bb5caedcc4

replace sync => github.com/sasha-s/go-deadlock v0.2.1-0.20190427202633-1595213edefa // doesn't do anything?

replace github.com/sasha-s/go-deadlock => github.com/sasha-s/go-deadlock v0.2.1-0.20190427202633-1595213edefa

go 1.14
