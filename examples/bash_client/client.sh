#!/usr/bin/env bash
# This is a very simple client example to show how to interact with the server.  If you
# run the program without any args (or with an unexpected one), it will just return your
# ID that serves as your unique token for tests and your unique domain for convenience.
# The script requires jq to be installed.

###
# These values are for documentation only. Change them!
# _host="example.com"  # host running BOAST
_host="localhost"
_port="2096"  # the server's API port
_b64secret="872k5eD/lGRbMZ3GqIPB0bUzqRjBlt1lhLH4+/42sKa="  # your secret (44 bytes max.)
# _b64secret could be generated with: `$ openssl rand -base64 32`
###

function usage {
  cat << EOF
Usage: $0 <receiver>
*<receiver> can be set to http, https, dns, or all.
EOF
  exit 1
}

if [ "$1" == "-h" ]; then
  usage;
fi

# GET /events with the "Authorization" header containing a Base64 secret.
# * Authorization header format: "Authorization: Secret <base64 secret>"
# * Base64 secret's maximum size: 44 bytes
_header="Authorization: Secret ${_b64secret}"
_events=$(curl --silent -X GET https://$_host:$_port/events -H "${_header}")

if [ "$1" == "http" ] || [ "$1" == "https" ] || [ "$1" == "dns" ]; then
  echo $_events | jq ".events[] | select(.receiver==\"${1^^}\")"
elif [  "$1" == "all" ]; then
  echo $_events | jq
else
  _id=$(echo $_events | jq -r .id)
  echo "Your id is $_id"
  echo "Your unique domain is: $_id.$_host"
fi
