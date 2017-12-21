#!/bin/bash

source /etc/profile.d/chruby.sh
chruby ruby-2.3.6

read INPUT

echo $INPUT | {{ .DummyCPIPath }} {{ .BaseDir }}/director.yml
