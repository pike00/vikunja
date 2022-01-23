#!/bin/sh

# Set the uid and gid of the vikunja run user
usermod --non-unique --uid ${PUID} vikunja
groupmod --non-unique --gid ${PGID} vikunja

nginx -g 'daemon off;' & 

su vikunja -c '/app/vikunja/vikunja'