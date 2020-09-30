package main

import (
	"fmt"
	"os"

	"github.com/itsursujit/phonenumbers"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("usage: phoneparser [number] [two letter coutry]")
		os.Exit(1)
	}

	metadata, err := phonenumbers.Parse(os.Args[1], os.Args[2])
	if err != nil {
		fmt.Printf("Error parsing number: %s\n", err)
	}
	fmt.Printf("            E164: %s\n", phonenumbers.Format(metadata, phonenumbers.E164))
	fmt.Printf("National Dialing: %s\n", phonenumbers.Format(metadata, phonenumbers.NATIONAL))
	fmt.Printf("        National: %d\n", *metadata.NationalNumber)
	fmt.Printf("NationalNumber:%d\n", *metadata.NationalNumber)
	fmt.Printf("CountryCode: %d\n", *metadata.CountryCode)
	fmt.Printf("IsPossible: %t\n", phonenumbers.IsPossibleNumber(metadata))
	fmt.Printf("IsValid: %t\n", phonenumbers.IsValidNumber(metadata))
	fmt.Printf("NationalFormatted: %s\n", phonenumbers.Format(metadata, phonenumbers.NATIONAL))
	fmt.Printf("InternationalFormatted: %s\n", phonenumbers.Format(metadata, phonenumbers.INTERNATIONAL))
}
