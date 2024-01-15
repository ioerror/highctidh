module codeberg.org/vula/highctidh

go 1.19

require codeberg.org/vula/highctidh/ctidh2048 v0.0.0 // XXX FIX ME LATER

require codeberg.org/vula/highctidh/ctidh1024 v0.0.0 // XXX FIX ME LATER

require codeberg.org/vula/highctidh/ctidh512 v0.0.0 // XXX FIX ME LATER

require codeberg.org/vula/highctidh/ctidh511 v0.0.0 // XXX FIX ME LATER

require codeberg.org/vula/highctidh/interfaces v0.0.0 // indirect; XXX FIX ME LATER

require codeberg.org/vula/highctidh/schemes v0.0.0 // XXX FIX ME LATER

require (
	codeberg.org/vula/highctidh/unified v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.8.4
)

require (
	codeberg.org/vula/highctidh/types v0.0.0-00010101000000-000000000000 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/mattn/go-pointer v0.0.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace codeberg.org/vula/highctidh/ctidh2048 => ./ctidh2048 // XXX REMOVE ME LATER

replace codeberg.org/vula/highctidh/ctidh1024 => ./ctidh1024 // XXX REMOVE ME LATER

replace codeberg.org/vula/highctidh/ctidh512 => ./ctidh512 // XXX REMOVE ME LATER

replace codeberg.org/vula/highctidh/ctidh511 => ./ctidh511 // XXX REMOVE ME LATER

replace codeberg.org/vula/highctidh/interfaces => ./interfaces // XXX REMOVE ME LATER

replace codeberg.org/vula/highctidh/schemes => ./schemes // XXX REMOVE ME LATER

replace codeberg.org/vula/highctidh/types => ./types // XXX REMOVE ME LATER

replace codeberg.org/vula/highctidh/unified => ./unified // XXX REMOVE ME LATER
