#!/bin/bash

source /usr/local/opt/chruby/share/chruby/chruby.sh
source /usr/local/opt/chruby/share/chruby/auto.sh

RUBIES+=(
  /Users/pivotal/.rubies/ruby-2.3.1
)

chruby 2.3.1

read INPUT

echo $INPUT | {{ .DummyCPIPath }} {{ .BaseDir }}/director.yml
