package lexer

import (
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/llir/llvm/asm/token"
)

func TestParseString(t *testing.T) {
	golden := []struct {
		input string
		want  []token.Token
	}{
		// i=0
		{
			input: ",",
			want: []token.Token{
				{Kind: token.Comma, Val: ",", Pos: 0},
				{Kind: token.EOF, Pos: 1},
			},
		},
		// i=1
		{
			input: "+0.314e+1",
			want: []token.Token{
				{Kind: token.Float, Val: "+0.314e+1", Pos: 0},
				{Kind: token.EOF, Pos: 9},
			},
		},
		// i=2
		{
			input: "$baz@foo%bar$baz!qux@42%37#7",
			want: []token.Token{
				{Kind: token.ComdatVar, Val: "baz", Pos: 0},
				{Kind: token.GlobalVar, Val: "foo", Pos: 4},
				{Kind: token.LocalVar, Val: "bar$baz", Pos: 8},
				{Kind: token.MetadataVar, Val: "qux", Pos: 16},
				{Kind: token.GlobalID, Val: "42", Pos: 20},
				{Kind: token.LocalID, Val: "37", Pos: 23},
				{Kind: token.AttrID, Val: "7", Pos: 26},
				{Kind: token.EOF, Pos: 28},
			},
		},
		// i=3
		{
			input: "...=,*[]{}()<>!",
			want: []token.Token{
				{Kind: token.Ellipsis, Val: "...", Pos: 0},
				{Kind: token.Equal, Val: "=", Pos: 3},
				{Kind: token.Comma, Val: ",", Pos: 4},
				{Kind: token.Star, Val: "*", Pos: 5},
				{Kind: token.Lbrack, Val: "[", Pos: 6},
				{Kind: token.Rbrack, Val: "]", Pos: 7},
				{Kind: token.Lbrace, Val: "{", Pos: 8},
				{Kind: token.Rbrace, Val: "}", Pos: 9},
				{Kind: token.Lparen, Val: "(", Pos: 10},
				{Kind: token.Rparen, Val: ")", Pos: 11},
				{Kind: token.Less, Val: "<", Pos: 12},
				{Kind: token.Greater, Val: ">", Pos: 13},
				{Kind: token.Exclaim, Val: "!", Pos: 14},
				{Kind: token.EOF, Pos: 15},
			},
		},
		// i=4
		{
			input: `"fo\6F":"fo\6F"@"fo\6F"%"fo\6F"$"fo\6F"!fo\6F`,
			want: []token.Token{
				{Kind: token.Label, Val: "foo", Pos: 0},
				{Kind: token.String, Val: "foo", Pos: 8},
				{Kind: token.GlobalVar, Val: "foo", Pos: 15},
				{Kind: token.LocalVar, Val: "foo", Pos: 23},
				{Kind: token.ComdatVar, Val: "foo", Pos: 31},
				{Kind: token.MetadataVar, Val: "foo", Pos: 39},
				{Kind: token.EOF, Pos: 45},
			},
		},
		// i=5
		{
			input: "42.0!42.0foo:;foo",
			want: []token.Token{
				{Kind: token.Float, Val: "42.0", Pos: 0},
				{Kind: token.Exclaim, Val: "!", Pos: 4},
				{Kind: token.Label, Val: "42.0foo", Pos: 5},
				{Kind: token.Comment, Val: "foo", Pos: 13},
				{Kind: token.EOF, Pos: 17},
			},
		},
		// i=6
		{
			input: "i42floatvoidaddxu0x6F",
			want: []token.Token{
				{Kind: token.Type, Val: "i42", Pos: 0},
				{Kind: token.Type, Val: "float", Pos: 3},
				{Kind: token.Type, Val: "void", Pos: 8},
				{Kind: token.KwAdd, Val: "add", Pos: 12},
				{Kind: token.KwX, Val: "x", Pos: 15},
				{Kind: token.Int, Val: "u0x6F", Pos: 16},
				{Kind: token.EOF, Pos: 21},
			},
		},
		// i=7
		{
			input: "i42floatvoidaddxu0x6F:",
			want: []token.Token{
				{Kind: token.Label, Val: "i42floatvoidaddxu0x6F", Pos: 0},
				{Kind: token.EOF, Pos: 22},
			},
		},
		// i=8
		{
			input: "42:-foo:0x1e",
			want: []token.Token{
				{Kind: token.Label, Val: "42", Pos: 0},
				{Kind: token.Label, Val: "-foo", Pos: 3},
				{Kind: token.Float, Val: "0x1e", Pos: 8},
				{Kind: token.EOF, Pos: 12},
			},
		},
		// i=9
		{
			input: "0xK1e 0xL1e 0xM1e 0xH1e",
			want: []token.Token{
				{Kind: token.Float, Val: "0xK1e", Pos: 0},
				{Kind: token.Float, Val: "0xL1e", Pos: 6},
				{Kind: token.Float, Val: "0xM1e", Pos: 12},
				{Kind: token.Float, Val: "0xH1e", Pos: 18},
				{Kind: token.EOF, Pos: 23},
			},
		},
		// i=10
		{
			input: "37-42",
			want: []token.Token{
				{Kind: token.Int, Val: "37", Pos: 0},
				{Kind: token.Int, Val: "-42", Pos: 2},
				{Kind: token.EOF, Pos: 5},
			},
		},
		// i=11
		{
			input: "....foo:",
			want: []token.Token{
				{Kind: token.Label, Val: "....foo", Pos: 0},
				{Kind: token.EOF, Pos: 8},
			},
		},
		// i=12
		{
			input: `!fo\6F!bar`,
			want: []token.Token{
				{Kind: token.MetadataVar, Val: "foo", Pos: 0},
				{Kind: token.MetadataVar, Val: "bar", Pos: 6},
				{Kind: token.EOF, Pos: 10},
			},
		},
		// i=13
		{
			input: ";foo\n;bar\r\n;baz",
			want: []token.Token{
				{Kind: token.Comment, Val: "foo", Pos: 0},
				{Kind: token.Comment, Val: "bar", Pos: 5},
				{Kind: token.Comment, Val: "baz", Pos: 11},
				{Kind: token.EOF, Pos: 15},
			},
		},
		// i=14
		{
			input: `"\foo""\de\ad\be\ef"`,
			want: []token.Token{
				{Kind: token.String, Val: `\foo`, Pos: 0},
				{Kind: token.String, Val: "\xDE\xAD\xBE\xEF", Pos: 6},
				{Kind: token.EOF, Pos: 20},
			},
		},
		// i=15
		{
			input: "$42foo:$foo:",
			want: []token.Token{
				{Kind: token.Label, Val: "$42foo", Pos: 0},
				{Kind: token.Label, Val: "$foo", Pos: 7},
				{Kind: token.EOF, Pos: 12},
			},
		},
		// i=16
		{
			input: `"fooλbar" "baz`,
			want: []token.Token{
				{Kind: token.String, Val: "fooλbar", Pos: 0},
				{Kind: token.Error, Val: "unexpected eof in quoted string", Pos: 15},
				{Kind: token.EOF, Pos: 15},
			},
		},
		// i=17
		{
			input: `+br acall .icmp #void $42 & % @`,
			want: []token.Token{
				{Kind: token.Error, Val: "unexpected '+'", Pos: 0},
				{Kind: token.KwBr, Val: "br", Pos: 1},
				{Kind: token.Error, Val: "unexpected 'a'", Pos: 4},
				{Kind: token.KwCall, Val: "call", Pos: 5},
				{Kind: token.Error, Val: "unexpected '.'", Pos: 10},
				{Kind: token.KwIcmp, Val: "icmp", Pos: 11},
				{Kind: token.Error, Val: "unexpected '#'", Pos: 16},
				{Kind: token.Type, Val: "void", Pos: 17},
				{Kind: token.Error, Val: "unexpected '$'", Pos: 22},
				{Kind: token.Int, Val: "42", Pos: 23},
				{Kind: token.Error, Val: "unexpected '&'", Pos: 26},
				{Kind: token.Error, Val: "unexpected '%'", Pos: 28},
				{Kind: token.Error, Val: "unexpected '@'", Pos: 30},
				{Kind: token.EOF, Pos: 31},
			},
		},
		// i=18
		{
			input: "\uFFFD \"foo\uFFFDbar\" ;foo\uFFFDbar",
			want: []token.Token{
				{Kind: token.Error, Val: "illegal UTF-8 encoding", Pos: 0},
				{Kind: token.Error, Val: "illegal UTF-8 encoding", Pos: 8},
				{Kind: token.String, Val: "foo\uFFFDbar", Pos: 4},
				{Kind: token.Error, Val: "illegal UTF-8 encoding", Pos: 20},
				{Kind: token.Comment, Val: "foo\uFFFDbar", Pos: 16},
				{Kind: token.EOF, Pos: 26},
			},
		},
		// i=19
		{
			input: `@"foo`,
			want: []token.Token{
				{Kind: token.Error, Val: "unexpected eof in quoted string", Pos: 5},
				{Kind: token.EOF, Pos: 5},
			},
		},
		// i=20
		{
			input: `%"foo`,
			want: []token.Token{
				{Kind: token.Error, Val: "unexpected eof in quoted string", Pos: 5},
				{Kind: token.EOF, Pos: 5},
			},
		},
		// i=21
		{
			input: `$"foo`,
			want: []token.Token{
				{Kind: token.Error, Val: "unexpected eof in quoted string", Pos: 5},
				{Kind: token.EOF, Pos: 5},
			},
		},
	}
	for i, g := range golden {
		got := ParseString(g.input)
		if !reflect.DeepEqual(got, g.want) {
			t.Errorf("i=%d: expected %#v, got %#v", i, g.want, got)
			continue
		}
	}
}

func BenchmarkParseString(b *testing.B) {
	buf, err := ioutil.ReadFile("../testdata/for.ll")
	if err != nil {
		b.Fatal(err)
	}
	input := string(buf)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ParseString(input)
	}
}
