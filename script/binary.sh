#!/bin/bash

# template for Bourne shell
# which is default shell for CentOS
# the purpose of this template is for scripts to be used in CentOS container

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${0}")" && pwd -P)"

. "$SCRIPT_DIR"/_functions.sh

# define constants

# script usage
function usage() {
  cat <<EOF | xargs -0 printf '%b'

$(colorize 'NAME' 0 1)
  $0

$(colorize 'USAGE' 0 1)
  $0 [OPTIONS]

$(colorize 'ARGUMENTS' 0 1)

$(colorize 'OPTIONS' 0 1)
  -h, --help    show usage
  --dry-run     run without any changes

EOF
  exit 0
}

# SIGNALS
# HUP 1
# INT 2 CTRL-C
# QUIT 3 CTRL-\
# TERM 15
trap cleanup HUP INT QUIT TERM

# cleanup steps
function cleanup() {
  exit 0
}

# define params to parse
function parse_params() {
  # default values of options
  # built-in variables
  local options_count=0
  local args_count=0
  local args=()
  dry_run=0

  # define options

  while :; do
    case "${1-}" in
    -h | --help)
      usage
      ;;

    --dry-run)
      # shellcheck disable=SC2034
      dry_run=1
      ;;

    # parse options

    -?*)
      usage
      ;;

    *)
      args+=("${1-}")

      if [[ -z "${1-}" ]]; then
        break
      fi
      ;;
    esac

    if [[ -n "${1-}" ]]; then
      shift
    fi
  done

  # define arguments

  # parse args
  for arg in "${args[@]}"; do
    if [[ -n "$arg" ]]; then
      (( args_count += 1 ))

      # define positional arguments here
    fi
  done

  if (( options_count == 0 && args_count == 0 )); then
    usage
  fi

  # validate params
}

# define more functions
function build() {
  go build \
    -o dist/snowflake \
    ./cmd

  return 0
}

# run application
function run() {
  parse_params "$@"

  if [[ -d "dist" ]]; then
    rm -r dist
  fi

  build

  quit 'Done :)'
}

run "$@"
