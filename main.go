package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/appscode/mergo"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/printer"
	jsonParser "github.com/hashicorp/hcl/json/parser"
)

// VERSION is what is returned by the `-v` flag
var Version = "development"

func main() {
	version := flag.Bool("version", false, "Prints current app version")
	// reverse := flag.Bool("reverse", false, "Input HCL, output JSON")
	srcFile := flag.String("src-file", "", "Name of the source file")
	dstFile := flag.String("dst-file", "", "Name of the destination file, in merge operation non-empty dst file attributes will be overridden by non-empty src file attribute values.")
	writeFile := flag.String("write-file", "", "Name of the file to write resultant merged file")
	flag.Parse()

	if *version {
		fmt.Println(Version)
		return
	}

	var err error
	//if *reverse {
	//	err = toJSON()
	//} else {
	//	err = toHCL()
	//}

	err = mergeHCLFile(*dstFile, *srcFile, *writeFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// mergeHCLFile will merge dstFile and srcFile and write
// the resultant merged file will be written in writeFile.
// In merge operation non-empty dstFile attributes will be
// overridden by non-empty srcFile attribute values.
func mergeHCLFile(dstFile, srcFile, writeFile string) error {
	dst, err := ioutil.ReadFile(dstFile)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %s", dstFile, err)
	}

	var d map[string]interface{}
	err = hcl.Unmarshal(dst, &d)
	if err != nil {
		return fmt.Errorf("failed to unmarshal: %s", err)
	}

	src, err := ioutil.ReadFile(srcFile)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %s", srcFile, err)
	}

	var s map[string]interface{}
	err = hcl.Unmarshal(src, &s)
	if err != nil {
		return fmt.Errorf("failed to unmarshal: %s", err)
	}

	err = mergo.Map(&d, s, mergo.WithOverride)
	if err != nil {
		return fmt.Errorf("failed to merge file: %s", err)
	}

	json, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return fmt.Errorf("unable to marshal json: %s", err)
	}

	fmt.Println(string(json))

	ast, err := jsonParser.Parse(json)
	if err!=nil {
		return err
	}

	f, err := os.Create(writeFile)
	if err != nil {
		return err
	}

	err = printer.Fprint(f, ast)
	if err != nil {
		return fmt.Errorf("failed to write file: %s", err)
	}

	return nil
}

func toJSON() error {
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return fmt.Errorf("unable to read from stdin: %s", err)
	}

	var v map[string]interface{}
	err = hcl.Unmarshal(input, &v)
	if err != nil {
		return fmt.Errorf("unable to parse HCL: %s", err)
	}

	json, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("unable to marshal json: %s", err)
	}

	fmt.Println(string(json))

	return nil
}

func toHCL() error {
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return fmt.Errorf("unable to read from stdin: %s", err)
	}

	ast, err := jsonParser.Parse([]byte(input))
	if err != nil {
		return fmt.Errorf("unable to parse JSON: %s", err)
	}

	err = printer.Fprint(os.Stdout, ast)
	if err != nil {
		return fmt.Errorf("unable to print HCL: %s", err)
	}

	return nil
}

func testAST(dstFile, srcFile, writeFile string) error {
	dst, err := ioutil.ReadFile(dstFile)
	if err!=nil {
		return err
	}

	dast,err := hcl.ParseBytes(dst)
	if err!=nil {
		return err
	}

	src, err := ioutil.ReadFile(dstFile)
	if err!=nil {
		return err
	}

	sast,err := hcl.ParseBytes(src)
	if err!=nil {
		return err
	}

	err = mergo.Map(dast, *sast, mergo.WithOverride)
	if err!=nil {
		return err
	}

	err = printer.Fprint(os.Stdout, dast)
	if err != nil {
		return fmt.Errorf("unable to print HCL: %s", err)
	}
	return err
}

func check(file string) error {
	dst, err := ioutil.ReadFile("original.hcl")
	if err != nil {
		return fmt.Errorf("failed to read file %s: %s",file, err)
	}

	ast, err := hcl.ParseBytes(dst)
	err = printer.Fprint(os.Stdout, ast)
	if err != nil {
		return fmt.Errorf("unable to print HCL: %s", err)
	}
	return err
}
