#!/bin/bash -x

function cleanup {
  rm -rf $1 $2
}

function main {
  local bosh_src_path
  local seed
  local ssl_assets_dir

  bosh_src_path=$1
  seed=$2

  cp "${PWD}/assets/ssl/*" /tmp/
  tmp_config_file=$(mktemp /tmp/fuzz-config-.XXXXXX)
  trap "cleanup ${tmp_config_file} /tmp/ssl" EXIT

  sed  -e "s#BOSH_SRC_PATH#${bosh_src_path}#g" -e "s#PWD#${PWD}#g" \
    ./ci/concourse-config.json > "${tmp_config_file}"

  DB=postgresql go run src/github.com/cloudfoundry-incubator/bosh-fuzz-tests/main.go \
    "${tmp_config_file}" \
    "${seed:-}"
}

main "$@"
