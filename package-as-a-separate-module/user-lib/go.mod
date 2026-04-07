module github.com/iypetrov/user-lib

go 1.24.5

require (
	github.com/iypetrov/user-lib/apihelpers v0.0.0-unpublished
	github.com/iypetrov/user-lib/apis v0.0.0-unpublished
)

require github.com/goccy/go-yaml v1.19.2 // indirect

replace (
	github.com/iypetrov/user-lib/apihelpers v0.0.0-unpublished => ./apihelpers
	github.com/iypetrov/user-lib/apis v0.0.0-unpublished => ./apis
)
