#!/bin/bash

export HOME=/Users/pivotal

# source /etc/profile.d/chruby.sh
source /usr/local/share/chruby/chruby.sh

# RUBIES+=(/Users/pivotal/.rubies/*)
# export GEM_PATH="~/.gem/ruby/2.4.2/bin"
# export GEM_HOME="~/.gem/ruby/2.4.2/bin"

chruby {{ .RubyVersion }}

gem list > /tmp/foo.txt
gem environment >> /tmp/foo.txt
echo $HOME >> /tmp/foo.txt

read -r INPUT

echo $INPUT | {{ .DummyCPIPath }} {{ .BaseDir }}/cpi_config.json
