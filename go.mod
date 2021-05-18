module github.com/uhppoted/uhppoted-app-wild-apricot

go 1.16

require (
	github.com/hyperjumptech/grule-rule-engine v1.8.5
	github.com/sirupsen/logrus v1.8.1
	github.com/uhppoted/uhppote-core v0.6.13-0.20210517175353-3ea261f5ec47
	github.com/uhppoted/uhppoted-api v0.6.13-0.20210517193706-c537c73adc3f
	golang.org/x/sys v0.0.0-20210426230700-d19ff857e887
)

replace github.com/uhppoted/uhppoted-api => ../uhppoted-api
