module github.com/PowerDNS/lmdb-go/v2/tests/coexist

go 1.21

require (
	github.com/PowerDNS/lmdb-go v1.9.3
	github.com/PowerDNS/lmdb-go/v2 v2.0.0
)

replace github.com/PowerDNS/lmdb-go/v2 => ../../
