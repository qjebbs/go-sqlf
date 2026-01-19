package syntax

import (
	"reflect"
	"testing"
)

func TestScanner(t *testing.T) {
	testCases := []struct {
		raw  string
		want []token
	}{
		{
			raw: "$1,$2",
			want: []token{
				{typ: _Ref, lit: "$1", bad: false, kind: _KindRefNumbered, start: 0, end: 2},
				{typ: _Plain, lit: ",", bad: false, kind: 0, start: 2, end: 3},
				{typ: _Ref, lit: "$2", bad: false, kind: _KindRefNumbered, start: 3, end: 5},
				{typ: _EOF, lit: "", bad: false, kind: 0, start: 5, end: 5},
			},
		},
		{
			raw: "a IN (?,?)",
			want: []token{
				{typ: _Plain, lit: "a IN (", bad: false, kind: 0, start: 0, end: 6},
				{typ: _Ref, lit: "?", bad: false, kind: _KindRefPositional, start: 6, end: 7},
				{typ: _Plain, lit: ",", bad: false, kind: 0, start: 7, end: 8},
				{typ: _Ref, lit: "?", bad: false, kind: _KindRefPositional, start: 8, end: 9},
				{typ: _Plain, lit: ")", bad: false, kind: 0, start: 9, end: 10},
				{typ: _EOF, lit: "", bad: false, kind: 0, start: 10, end: 10},
			},
		},
		{
			raw: "'a''b'",
			want: []token{
				{typ: _Plain, lit: "'a''b'", bad: false, kind: 0, start: 0, end: 6},
				{typ: _EOF, lit: "", bad: false, kind: 0, start: 6, end: 6},
			},
		},
		{
			raw: "'a''b",
			want: []token{
				{typ: _Plain, lit: "'a''b", bad: true, kind: 0, start: 0, end: 5},
				{typ: _EOF, lit: "", bad: false, kind: 0, start: 5, end: 5},
			},
		},
		{
			raw: `"a""b"`,
			want: []token{
				{typ: _Literal, lit: `"a""b"`, bad: false, kind: _KindLitString, start: 0, end: 6},
				{typ: _EOF, lit: "", bad: false, kind: 0, start: 6, end: 6},
			},
		},
		{
			raw: "$$1",
			want: []token{
				{typ: _Escape, lit: "$$", bad: false, kind: 0, start: 0, end: 2},
				{typ: _Plain, lit: "1", bad: false, kind: 0, start: 2, end: 3},
				{typ: _EOF, lit: "", bad: false, kind: 0, start: 3, end: 3},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.raw, func(t *testing.T) {
			got := make([]token, 0)
			s := newScanner(tc.raw, false)
			for s.NextToken() {
				// ignore Pos
				s.token.pos = Pos{}
				got = append(got, *s.token)
			}
			if !reflect.DeepEqual(got, tc.want) {
				for _, tk := range got {
					t.Logf("%#v", tk)
				}
				// for i, tk := range got {
				// 	for _, tk := range got {
				// 		t.Logf("%#v", tk)
				// 	}
				// 	var want *token
				// 	if i < len(tc.want) {
				// 		want = &tc.want[i]
				// 	}
				// 	if !reflect.DeepEqual(&tk, want) {
				// 		t.Logf("#%d, want %#v, got %#v", i, want, &tk)
				// 	}
				// }
				t.Fatal("failed")
			}
		})
	}
}
