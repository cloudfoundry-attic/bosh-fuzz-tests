# README

To install:

```
cd src/github.com/cloudfoundry-incubator/bosh-fuzz-tests
```

To run ginkgo (to test changes):

```
go get github.com/onsi/ginkgo/ginkgo
go get github.com/onsi/gomega
bin/env ginkgo -r .
```

To run fuzz tests locally with BOSH director from local source or as defined in [`config.json`](config.json):

```
cp ../../../../assets/ssl/* /tmp/
DB=postgresql bin/env go run main.go ../../../../config.json
```

Note, on the local workstation this will leave processes behind.
To clean up those processes you can try:

```
kill $(ps aux | egrep "nats-server|bosh-fuzz-tests|bin/bosh-director|bosh-config-server-executable" | grep -v grep | awk '{print $2}')
```

To re-create failures seen on Concourse:

* Search for `Seeding with` and copy the seed number
* Copy the `parameters` section from `ci/concourse-config.json` to `config.json`
* Run the following command:

```
bin/env go run main.go ../../../../config.json <SEED_NUMBER>
```
