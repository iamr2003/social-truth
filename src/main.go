package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"text/template"
)

type DafnyInput struct {
	Assertion string
}

func main() {
	var dafnyTmpl = "lemma.dafny.templ"
	tmpl, err := template.New(dafnyTmpl).ParseFiles(dafnyTmpl)
	if err != nil {
		panic(err)
	}

	file, err := os.Create("generated.dfy")
	if err != nil {
		log.Fatalf("Failed to create file: %s", err)
	}
	defer file.Close()

	assertion := DafnyInput{Assertion: "5 < 3"}
	err = tmpl.Execute(file, assertion)
	if err != nil {
		panic(err)
	}

	fmt.Println("Initial poc for social truth")
	cmd := exec.Command("dafny", "verify", "generated.dfy")

	output, err := cmd.Output()
	if err != nil {
		log.Printf("dafny fails with %s\n", err)
	}
	fmt.Printf("Output: \n%s \n", output)
}

// verify assertion and return yes or no
func VerifyAssertion(assertion string) (string, string) {
	tmpl, err := template.New("verify").ParseFiles("lemma.dafny.templ")
	if err != nil {
		return "", "Failed to parse template: " + err.Error()
	}

	file, err := os.Create("generated.dfy")
	if err != nil {
		return "", "Failed to create file: " + err.Error()
	}
	defer file.Close()

	dafnyInput := DafnyInput{Assertion: assertion}
	err = tmpl.Execute(file, dafnyInput)
	if err != nil {
		return "", "Failed to execute template: " + err.Error()
	}

	cmd := exec.Command("dafny", "verify", "generated.dfy")
	output, err := cmd.Output()
	if err != nil {
		return "", "Verification failed: " + err.Error()
	}

	return "yes", string(output)
}

// serve webpage on a route that accesses above endpoint
