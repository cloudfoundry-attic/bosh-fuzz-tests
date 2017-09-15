#!/bin/bash

source /etc/profile.d/chruby.sh
chruby {{ .RubyVersion }}

read INPUT

echo $INPUT | {{ .DummyCPIPath }} {{ .BaseDir }}/director.yml
