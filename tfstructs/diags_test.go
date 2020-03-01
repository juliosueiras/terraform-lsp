package tfstructs

import (
	"testing"

	v2 "github.com/hashicorp/hcl/v2"
	oldHCL2 "github.com/hashicorp/hcl2/hcl"
)

func TestRangeOfV2Range(t *testing.T) {
	r := v2.Range{
		Start: v2.Pos{
			Line:   1,
			Column: 2,
		},
		End: v2.Pos{
			Line:   3,
			Column: 4,
		},
	}

	res := rangeOf(r)

	if res.Start.Line != r.Start.Line-1 {
		t.Errorf("V2 Range Start Line does not match expected. Got %d, Expected %d", res.Start.Line, r.Start.Line)
	}

	if res.Start.Character != r.Start.Column-1 {
		t.Errorf("LSP Range Start Character does not match expected. Got %d, Expected %d", res.Start.Character, r.Start.Column)
	}

	if res.End.Line != r.End.Line-1 {
		t.Errorf("LSP Range End Line does not match expected. Got %d, Expected %d", res.End.Line, r.End.Line)
	}

	if res.End.Character != r.End.Column-1 {
		t.Errorf("LSP Range End Character does not match expected. Got %d, Expected %d", res.End.Character, r.End.Column)
	}
}

func TestCastToV2Range(t *testing.T) {
	r := oldHCL2.Range{
		Start: oldHCL2.Pos{
			Line:   5,
			Column: 6,
		},
		End: oldHCL2.Pos{
			Line:   7,
			Column: 8,
		},
	}

	res := castToV2Range(r)

	if res.Start.Line != r.Start.Line {
		t.Errorf("V2 Range Start Line does not match expected. Got %d, Expected %d", res.Start.Line, r.Start.Line)
	}

	if res.Start.Column != r.Start.Column {
		t.Errorf("V2 Range Start Column does not match expected. Got %d, Expected %d", res.Start.Column, r.Start.Column)
	}

	if res.End.Line != r.End.Line {
		t.Errorf("V2 Range End Line does not match expected. Got %d, Expected %d", res.End.Line, r.End.Line)
	}

	if res.End.Column != r.End.Column {
		t.Errorf("V2 Range End Column does not match expected. Got %d, Expected %d", res.End.Column, r.End.Column)
	}
}
