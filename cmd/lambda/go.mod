module T27FundraisingLambda

go 1.24.1

toolchain go1.24.2

require github.com/graphql-go/graphql v0.8.1 // indirect

require (
	github.com/aws/aws-lambda-go v1.48.0
	github.com/cch71/T27FundraisingLambda/frgql v0.0.0
)

require (
	github.com/codingsince1985/geo-golang v1.8.5 // indirect
	github.com/deckarep/golang-set/v2 v2.8.0 // indirect
	github.com/doug-martin/goqu/v9 v9.19.0 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.2 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.7.5 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	golang.org/x/crypto v0.38.0 // indirect
	golang.org/x/sync v0.14.0 // indirect
	golang.org/x/text v0.25.0 // indirect
)

replace github.com/cch71/T27FundraisingLambda/frgql => ../../frgql
