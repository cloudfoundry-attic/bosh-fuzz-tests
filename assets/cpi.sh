#!/bin/bash

source /etc/profile.d/chruby.sh
chruby ruby-2.3.1

read INPUT

echo $INPUT | {{ .DummyCPIPath }} {{ .BaseDir }}/director.yml
