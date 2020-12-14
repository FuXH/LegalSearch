module LegalSearch

go 1.13

require (
	github.com/elastic/go-elasticsearch/v7 v7.10.0 // indirect
	github.com/gin-gonic/gin v1.6.3
	github.com/olivere/elastic/v7 v7.0.22
	github.com/stretchr/testify v1.5.1
	gopkg.in/yaml.v2 v2.2.8
)

replace github.com/olivere/elastic/v7 => gopkg.in/olivere/elastic.v7 v7.0.22
