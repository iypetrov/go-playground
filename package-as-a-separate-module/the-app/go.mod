module github.com/iypetrov/the-app

go 1.24.5

require github.com/iypetrov/user-lib/apis v0.0.0-unpublished

replace (
	github.com/iypetrov/user-lib/apihelpers v0.0.0-unpublished => ../user-lib/apihelpers
	github.com/iypetrov/user-lib/apis v0.0.0-unpublished => ../user-lib/apis
)
