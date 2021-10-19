module github.com/pcarion/gonotionapi

go 1.17

require github.com/jomei/notionapi v1.5.2

require github.com/gernest/front v0.0.0-20210301115436-8a0b0a782d0a

require (
	github.com/pkg/errors v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/jomei/notionapi => ../notionapi
