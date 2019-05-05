package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hcl/hclsyntax"
	"github.com/hashicorp/terraform/configs"
	"github.com/zclconf/go-cty/cty"
	"io/ioutil"
)

func main() {
	parser := configs.NewParser(nil)

	//file3, _ := hclsyntax.ParseConfig([]byte(attr), "", hcl.Pos{Line: 1, Column: 1})
	file, _ := parser.LoadConfigFile("test.tf")
	fileText, _ := ioutil.ReadFile("test.tf")
	pos := findOffset(string(fileText), 9, 27)
	spew.Dump(file.ManagedResources[0].DeclRange.Start)
	spew.Dump(file.ManagedResources[0].DeclRange.End)
	config := file.ManagedResources[0].Config.(interface{}).(*hclsyntax.Body)
	attr := config.AttributeAtPos(hcl.Pos{
		Byte: pos,
	})

	if attr.Expr.Variables()[0].RootName() == "var" {
		name := attr.Expr.Variables()[0].SimpleSplit().Rel[0].(hcl.TraverseAttr).Name
		for _, v := range file.Variables {
			if name == v.Name {
			}
		}
	}
}

//func extractObject(c cty.Type) {
//	if c.IsObjectType() {
//
//	}
//}

func findOffset(fileText string, line, column int) int {
	currentCol := 1
	currentLine := 1

	for offset, ch := range fileText {
		if currentLine == line && currentCol == column {
			return offset
		}

		if ch == '\n' {
			currentLine++
			currentCol = 1
		} else {
			currentCol++
		}

	}
	return -1
}
