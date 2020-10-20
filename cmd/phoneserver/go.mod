module github.com/sujit-baniya/phonenumbers/cmd/phoneserver

go 1.15

replace github.com/sujit-baniya/phonenumbers => ../../

require (
	github.com/aws/aws-lambda-go v1.19.1
	github.com/sujit-baniya/phonenumbers v1.0.60
	github.com/urfave/cli v1.21.0 // indirect
)
