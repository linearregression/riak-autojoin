#!/bin/sh
VERSION=1.4

BIN_FILE=riak-autojoin

# Verify the file
echo "Verifying the sha256 checksum of the downloaded file:"
shasum -c ${BIN_FILE}.sha256sum || (echo "Warning, file does not match, exitting!" && exit -1)

fpm \
    -m "Tradeshift Operations <operations@tradeshift.com>" \
    --description "Auto join riak nodes to a cluster using consul" \
    --license "Mozilla Public License, version 2.0" \
    --vendor "Tradeshift" \
    --url "https://github.com/Tradeshift/riak-autojoin" \
    -s dir \
    -t deb \
    -n "riak-autojoin" \
    --prefix /opt/tradeshift/riak_auto_join/bin/ \
    -v ${VERSION} \
    riak-autojoin

shasum riak-autojoin_${VERSION}_amd64.deb > riak-autojoin_${VERSION}_amd64.deb.sha256sum
