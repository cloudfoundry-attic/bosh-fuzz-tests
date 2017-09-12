#!/bin/bash

source /etc/profile.d/chruby.sh
chruby ruby-2.4.1

read INPUT

echo $INPUT | {{ .DummyCPIPath }} {{ .BaseDir }}/director.yml
