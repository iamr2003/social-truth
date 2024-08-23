package main

import (
	"fmt"
	// "log"
	"net/http"
	"os"
	"os/exec"
	"text/template"
)

type DafnyInput struct {
	Assertion string
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	var htmlfile = "webpage.html"

	//write html file to w
	http.ServeFile(w, r, htmlfile)
}

func verifyAssertion(w http.ResponseWriter, r *http.Request) {
	fmt.Println("called verifyAssertion")

	// get assertion from request
	r.ParseForm()
	assertion := r.Form.Get("assertion")
	if assertion == "" {
		http.Error(w, "Assertion not provided", http.StatusBadRequest)
		return
	}

	fmt.Println("past form parse")

	// verify assertion
	result, output := VerifyAssertion(assertion)
	if result == "yes" {
		fmt.Fprintf(w, "<html><div>%s</div>Assertion is valid</html>", assertion)
	} else {
		fmt.Fprintf(w, "<html><div>%s</div>Assertion is not valid: %s</html>", assertion, output)
	}
}

func main() {
	http.HandleFunc("/", mainPage)
	http.HandleFunc("/verify", verifyAssertion)
	http.ListenAndServe(":8080", nil)
}

// verify assertion and return yes or no
func VerifyAssertion(assertion string) (string, string) {
	var dafnyTmpl = "lemma.dafny.templ"
	tmpl, err := template.New(dafnyTmpl).ParseFiles(dafnyTmpl)
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
