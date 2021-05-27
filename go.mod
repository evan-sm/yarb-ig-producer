module yarb-ig-producer

go 1.16

require (
	cloud.google.com/go/pubsub v1.10.3
	github.com/go-resty/resty/v2 v2.6.0
	github.com/k0kubun/colorstring v0.0.0-20150214042306-9440f1994b88 // indirect
	github.com/k0kubun/pp v3.0.1+incompatible // indirect
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/tidwall/gjson v1.7.5
	github.com/wmw9/ig v0.0.3
	github.com/wmw9/yarb-struct v0.0.1
	google.golang.org/grpc v1.38.0 // indirect
)

replace github.com/wmw9/go-makaba => /home/wmw/git/go/go-makaba
