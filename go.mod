module codeberg.org/vula/highctidh

go 1.19

require codeberg.org/vula/highctidh/ctidh2048 v0.0.0 // XXX FIX ME LATER

require codeberg.org/vula/highctidh/ctidh1024 v0.0.0 // XXX FIX ME LATER

require codeberg.org/vula/highctidh/ctidh512 v0.0.0 // XXX FIX ME LATER

require codeberg.org/vula/highctidh/ctidh511 v0.0.0 // XXX FIX ME LATER


require (
	github.com/stretchr/testify v1.8.4
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/mattn/go-pointer v0.0.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace codeberg.org/vula/highctidh/ctidh2048 => ./src/ctidh2048 // XXX REMOVE ME LATER

replace codeberg.org/vula/highctidh/ctidh1024 => ./src/ctidh1024 // XXX REMOVE ME LATER

replace codeberg.org/vula/highctidh/ctidh512 => ./src/ctidh512 // XXX REMOVE ME LATER

replace codeberg.org/vula/highctidh/ctidh511 => ./src/ctidh511 // XXX REMOVE ME LATER
