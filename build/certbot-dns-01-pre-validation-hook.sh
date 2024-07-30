#!/usr/bin/env bash
# This Let's Encrypt pre-validation hook is to be used for the DNS-01 challenge with
# BOAST's main domain. It's made to work with the Dockerfile in this directory and may
# need some changes for customised use cases.
#
# This hook will only be run if the certificate is due for renewal, so certbot can be
# run frequently (e.g. as a cron job) without unnecessarily stopping BOAST.
#
# Doc on how to use this with the provided Dockerfile (and more):
# https://github.com/ciphermarco/boast/blob/master/docs/deploying.md#deploying-with-docker
#
if [ -z "$CERTBOT_VALIDATION"  ]
then
	echo "error: validation is empty"
	exit -1
fi

_docker_boast_img="boastimg"
_docker_boast_container="boastmain"
_docker_boast_dns_container="boastdns"
_docker_boast_bin="/go/src/github.com/ciphermarco/BOAST/boast"

# Ignoring errors with `|| true` in case containers are not running or do not exist.

# Make sure everything is stopped
docker stop ${_docker_boast_container} || true
docker stop ${_docker_boast_dns_container} || true

# Make sure the BOAST's DNS temporary container does not exist.
docker container rm ${_docker_boast_dns_container} || true

# Run the DNS receiver with the challenge's TXT record.
docker run -d --name ${_docker_boast_dns_container} -p 53:53/udp ${_docker_boast_img} ${_docker_boast_bin} -dns_only -dns_txt ${CERTBOT_VALIDATION}
