#!/bin/sh
function execute {
    "$@"
    local status=$?
    if [ $status -ne 0 ]; then
        echo "Error: $1"
    fi
    return $status
}

execute createuser -U postgres -S -D -R meritwiki
execute createdb -U postgres MeritWiki --owner=wikiadmin --tablespace=pg_default