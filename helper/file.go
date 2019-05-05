package helper

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/terraform/configs"
	"github.com/sourcegraph/go-lsp"
	"github.com/zclconf/go-cty/cty"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"
)

func CheckAndGetConfig(parser *configs.Parser, originalFile *os.File, line int, character int) (*configs.File, hcl.Diagnostics, bool) {
	fileText, _ := ioutil.ReadFile(originalFile.Name())
	result := make([]byte, 1)
	pos := FindOffset(string(fileText), line, character)

	tempFile, _ := ioutil.TempFile("/tmp", "check_tf_lsp")
	defer os.Remove(tempFile.Name())

	originalFile.ReadAt(result, int64(pos))

	if string(result) == "." {
		fileText[pos] = ' '

		fileText = []byte(strings.Replace(string(fileText), ". ", "  ", -1))
		fileText = []byte(strings.Replace(string(fileText), ".,", " ,", -1))
		tempFile.Truncate(0)
		tempFile.Seek(0, 0)
		tempFile.Write(fileText)

		resultConfig, diags := parser.LoadConfigFileOverride(tempFile.Name())

		return resultConfig, diags, true
	}

	resultConfig, diags := parser.LoadConfigFileOverride(originalFile.Name())
	return resultConfig, diags, false
}

func FindOffset(fileText string, line, column int) int {
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

func DumpLog(res interface{}) {
	log.Println(spew.Sdump(res))
}

func GetType(t interface{}) reflect.Type {
	return reflect.TypeOf(t)
}

func ParseVariables(vars hcl.Traversal, configVarsType []*configs.Variable, completionItems []lsp.CompletionItem) []lsp.CompletionItem {
	if len(vars) == 0 {
		for _, t := range configVarsType {
			completionItems = append(completionItems, lsp.CompletionItem{
				Label:  t.Name,
				Detail: t.Type.FriendlyName(),
			})
		}
		return completionItems
	}

	for _, v := range configVarsType {
		if vars[0].(hcl.TraverseAttr).Name == v.Name {
			return parseVariables(vars[1:], v.Type, completionItems)
		}
	}
	return completionItems
}

func parseVariables(vars hcl.Traversal, configVarsType cty.Type, completionItems []lsp.CompletionItem) []lsp.CompletionItem {
	if len(vars) == 0 {
		if configVarsType.IsObjectType() {
			for k, t := range configVarsType.AttributeTypes() {
				completionItems = append(completionItems, lsp.CompletionItem{
					Label:  k,
					Detail: t.FriendlyName(),
				})
			}

			return completionItems
		}
	}

	if !configVarsType.IsObjectType() {
		return completionItems
	}

	varAttr := vars[0].(hcl.TraverseAttr)
	//DumpLog(configVarsType.HasAttribute(varAttr.Name))
	if configVarsType.HasAttribute(varAttr.Name) {
		attr := configVarsType.AttributeType(varAttr.Name)
		return parseVariables(vars[1:], attr, completionItems)
	}

	return nil
}
