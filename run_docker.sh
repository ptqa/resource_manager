#!/bin/bash
export GIN_MODE=release
PORT=${PORT:-3000}
LIMIT=${LIMIT:-10}
WORKERS=${WORKERS:-16}

cat << EOF > /go/src/github.com/ptqa/resource_manager/config.json
{
	"Port":$PORT,
	"Limit":$LIMIT,
	"Workers":$WORKERS
}
EOF
exec /go/bin/resource_manager
