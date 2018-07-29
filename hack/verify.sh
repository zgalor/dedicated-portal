#!/bin/bash

# Usage:
# verify.sh <source directory list>
#
# Example:
#    verify.sh cmd pkg
SOURCE=$@

# Check go file format, return 0 if found offending file.
# list all offending files into stderr.
function check_fmt {
  local bad_files=$(gofmt -s -l $1/)
  if [[ -n "${bad_files}" ]]; then
      echo "gofmt needs to be run on the listed files" >&2
      echo "${bad_files}" >&2
      echo "Try running 'gofmt -w -d [path]'" >&2

      return 0
  fi

  return 1
}

# Check go file linting, return 0 if found offending file.
# list all offending files into stderr.
function check_lint {
  local bad_files=$(golint -set_exit_status ${source}/...)
  if [[ -n "${bad_files}" ]]; then
      echo "golint found some problems" >&2
      echo "${bad_files}" >&2
      echo "Try fix linting" >&2

      return 0
  fi

  return 1
}

# Run checks on all source directories.
# Exit on problems.
for source in ${SOURCE}; do
  if check_fmt ${source}; then
    exit 1
  fi

  if check_lint ${source}; then
    exit 1
  fi
done

exit 0
