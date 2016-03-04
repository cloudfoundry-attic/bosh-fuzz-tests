# README

To install:

```
cd ~/workspace
git clone git@github.com:cloudfoundry-incubator/bosh-fuzz-tests.git
cd bosh-fuzz-tests/src/github.com/cloudfoundry-incubator/bosh-fuzz-tests
```

To run ginkgo (to test changes):

```
go get github.com/onsi/ginkgo/ginkgo
go get github.com/onsi/gomega
bin/env ginkgo -r .
```

To run fuzz tests locally with BOSH director from local source (`/Users/pivotal/workspace/bosh` or
as defined in `config.json`):

```
bin/env go run main.go ~/workspace/bosh-fuzz-tests/config.json
```

To re-run fuzz tests with the same input:

* Search for `Seeding with` and copy the seed number
* Run

```
bin/env go run main.go ~/workspace/bosh-fuzz-tests/config.json <SEED_NUMBER>
```

To run fuzz tests 
