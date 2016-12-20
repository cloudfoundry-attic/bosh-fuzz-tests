#!/usr/bin/env bash

read INPUT

echo $INPUT | /Users/pivotal/workspace/bosh/src/bosh-director/bin/dummy_cpi {{ .DirectorConfigPath }}
