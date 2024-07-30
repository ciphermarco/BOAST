#!/usr/bin/env bash
# This Let's Encrypt post-validation hook is to be used for the DNS-01 challenge with
# BOAST's main domain. It's made to work with the Dockerfile in this directory and may
# need some changes for customised use cases.
#
# This hook will only be run if the certificate is due for renewal, so certbot can be
# run frequently (e.g. as a cron job) without unnecessarily stopping BOAST.
#
# Doc on how to use this with the provided Dockerfile and more:
# https://github.com/ciphermarco/boast/blob/master/docs/deploying.md#deploying-with-docker
#
if [ -z "$RENEWED_LINEAGE"  ]
then
	echo "error: renewed lineage is empty"
	exit -1
fi

_docker_boast_img="boastimg"
_docker_boast_container="boastmain"
_docker_boast_dns_container="boastdns"
_tls_certificate="${RENEWED_LINEAGE}/fullchain.pem"
_tls_privkey="${RENEWED_LINEAGE}/privkey.pem"
_boast_tls="$HOME/boast/tls"  # <- Change this if necessary

# Ignoring errors with `|| true` in case containers are not running or do not exist.

# Make sure everything is stopped.
docker stop ${_docker_boast_container} || true
docker stop ${_docker_boast_dns_container} || true

# Make sure the BOAST container does not exist.
docker container rm ${_docker_boast_container} || true

# Copy TLS files to BOAST's TLS directory.
mkdir -p ${_boast_tls}
cp ${_tls_certificate} ${_tls_privkey} ${_boast_tls}

# Run the BOAST's main container.
docker run -d --name ${_docker_boast_container} -p 53:53/udp -p 80:80 -p 443:443 -p 2096:2096 -p 8080:8080 -p 8443:8443 -v ${_boast_tls}:/go/src/github.com/ciphermarco/BOAST/tls ${_docker_boast_img}
