<a href="https://goreportcard.com/report/github.com/tintinnabulate/registration"><img src="https://goreportcard.com/badge/github.com/tintinnabulate/registration" /></a>

# registration

## TODO

* Run `errcheck` on codebase and make sure it has no output! If it has output, handle those errors! `go get -u github.com/kisielk/errcheck`
* PCI compliance: https://en.wikipedia.org/wiki/Payment_Card_Industry_Data_Security_Standard

## Go checklist

Make sure that

* `goimports -d` is silent on all files
* `go vet` is silent on all files
* `errcheck` is silent on all files (`go get -u github.com/kisielk/errcheck`)

You could make these git pre-commit hooks...

* This is pretty great: https://github.com/golang/go/wiki/CodeReviewComments
