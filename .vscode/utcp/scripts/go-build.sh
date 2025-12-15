#!/bin/bash
# Go build wrapper script
output="${1:-.}"
packages="${2:-.}"
exec go build -o "$output" "$packages"
