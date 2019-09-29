package langserver

//func TextDocumentHover(ctx context.Context, vs lsp.TextDocumentPositionParams) (lsp.Hover, error) {
//
//	parser := configs.NewParser(nil)
//	file, _, column, _, _ := helper.CheckAndGetConfig(parser, tempFile, vs.Position.Line+1, vs.Position.Character)
//	fileText, _ := ioutil.ReadFile(tempFile.Name())
//	pos := helper.FindOffset(string(fileText), vs.Position.Line+1, column)
//	posHCL := hcl.Pos{
//		Byte: pos,
//	}
//	config, _, _ := tfstructs.GetConfig(file, posHCL)
//	if config == nil {
//		return lsp.Hover{
//			Contents: []lsp.MarkedString{},
//		}, nil
//	}
//	attr := config.AttributeAtPos(posHCL)
//	if attr != nil && attr.Expr != nil {
//		scope := lang.Scope{}
//
//		s, w := scope.EvalExpr(attr.Expr, cty.DynamicPseudoType)
//
//		val := ""
//
//		if w != nil {
//			return lsp.Hover{
//				Contents: []lsp.MarkedString{},
//			}, nil
//		}
//
//		if s.CanIterateElements() {
//		} else {
//			val = s.AsString()
//		}
//
//		return lsp.Hover{
//			Contents: []lsp.MarkedString{
//				lsp.MarkedString{
//					Language: "Terraform",
//					Value:    val,
//				},
//			},
//		}, nil
//	}
//	return lsp.Hover{
//		Contents: []lsp.MarkedString{},
//	}, nil
//}
