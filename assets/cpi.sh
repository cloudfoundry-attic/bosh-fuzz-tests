#!/bin/bash

source /etc/profile.d/chruby.sh
chruby {{ .RubyVersion }}

read -r INPUT

echo $INPUT | {{ .DummyCPIPath }} {{ .BaseDir }}/cpi_config.json
