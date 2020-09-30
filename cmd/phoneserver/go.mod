module github.com/itsursujit/phonenumbers/cmd/phoneserver

go 1.15

replace github.com/itsursujit/phonenumbers => ../../

require (
	github.com/aws/aws-lambda-go v1.19.1
	github.com/itsursujit/phonenumbers v1.0.60
	github.com/urfave/cli v1.21.0 // indirect
)
