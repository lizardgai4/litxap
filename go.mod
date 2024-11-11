module github.com/gissleh/litxap

go 1.22.2

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fwew/fwew-lib/v5 v5.22.2 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/testify v1.9.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

//for testing on a local machine's fwew-lib
replace github.com/fwew/fwew-lib/v5 => ../fwew-lib
