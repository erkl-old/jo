package jo

import (
	"testing"
)

var scannerTests = []struct {
	in  string
	out []Event
}{
	{
		` { } `,
		[]Event{
			Space,             // ' '
			ObjectStart,       // '{'
			Space,             // ' '
			None,              // '}'
			ObjectEnd | Space, // ' '
			None,              // EOF
		},
	},
	{
		`{ "a": { "b":{} , "c" : 0} }`,
		[]Event{
			ObjectStart,       // '{'
			Space,             // ' '
			KeyStart,          // '"'
			None,              // 'a'
			None,              // '"'
			KeyEnd,            // ':'
			Space,             // ' '
			ObjectStart,       // '{'
			Space,             // ' '
			KeyStart,          // '"'
			None,              // 'b'
			None,              // '"'
			KeyEnd,            // ':'
			ObjectStart,       // '{'
			None,              // '}'
			ObjectEnd | Space, // ' '
			None,              // ','
			Space,             // ' '
			KeyStart,          // '"'
			None,              // 'c'
			None,              // '"'
			KeyEnd | Space,    // ' '
			None,              // ':'
			Space,             // ' '
			NumberStart,       // '0'
			NumberEnd,         // '}'
			ObjectEnd | Space, // ' '
			None,              // '}'
			ObjectEnd,         // EOF
		},
	},
	{
		` [ ] `,
		[]Event{
			Space,            // ' '
			ArrayStart,       // '['
			Space,            // ' '
			None,             // ']'
			ArrayEnd | Space, // ' '
			None,             // EOF
		},
	},
	{
		`[[],[[]]]`,
		[]Event{
			ArrayStart, // '['
			ArrayStart, // '['
			None,       // ']'
			ArrayEnd,   // ','
			ArrayStart, // '['
			ArrayStart, // '['
			None,       // ']'
			ArrayEnd,   // ']'
			ArrayEnd,   // ']'
			ArrayEnd,   // EOF
		},
	},
	{
		`[0,1, 2 ,3 , 4]`,
		[]Event{
			ArrayStart,        // '['
			NumberStart,       // '0'
			NumberEnd,         // ','
			NumberStart,       // '1'
			NumberEnd,         // ','
			Space,             // ' '
			NumberStart,       // '2'
			NumberEnd | Space, // ' '
			None,              // ','
			NumberStart,       // '3'
			NumberEnd | Space, // ' '
			None,              // ','
			Space,             // ' '
			NumberStart,       // '4'
			NumberEnd,         // ']'
			ArrayEnd,          // EOF
		},
	},
	{
		`"foo"`,
		[]Event{
			StringStart, // '"'
			None,        // 'f'
			None,        // 'o'
			None,        // 'o'
			None,        // '"'
			StringEnd,   // EOF
		},
	},
	{
		" \" bar\"\n ",
		[]Event{
			Space,             // ' '
			StringStart,       // '"'
			None,              // ' '
			None,              // 'b'
			None,              // 'a'
			None,              // 'r'
			None,              // '"'
			StringEnd | Space, // '\n'
			Space,             // ' '
			None,              // EOF
		},
	},
	{
		`"\b\f\n\r\t\\\/\""`,
		[]Event{
			StringStart, // '"'
			None,        // '\\'
			None,        // 'b'
			None,        // '\\'
			None,        // 'f'
			None,        // '\\'
			None,        // 'n'
			None,        // '\\'
			None,        // 'r'
			None,        // '\\'
			None,        // 't'
			None,        // '\\'
			None,        // '\\'
			None,        // '\\'
			None,        // '/'
			None,        // '\\'
			None,        // '"'
			None,        // '"'
			StringEnd,   // EOF
		},
	},
	{
		`"\u2603 = ☃"`,
		[]Event{
			StringStart, // '"'
			None,        // '\\'
			None,        // 'u'
			None,        // '2'
			None,        // '6'
			None,        // '0'
			None,        // '3'
			None,        // ' '
			None,        // '='
			None,        // ' '
			None,        // '\xE2'
			None,        // '\x98'
			None,        // '\x83'
			None,        // '"'
			StringEnd,   // EOF
		},
	},
	{
		`0 `,
		[]Event{
			NumberStart,       // '0'
			NumberEnd | Space, // ' '
			None,              // EOF
		},
	},
	{
		`1 `,
		[]Event{
			NumberStart,       // '1'
			NumberEnd | Space, // ' '
			None,              // EOF
		},
	},
	{
		`2.5 `,
		[]Event{
			NumberStart,       // '2'
			None,              // '.'
			None,              // '5'
			NumberEnd | Space, // ' '
			None,              // EOF
		},
	},
	{
		`0.1e+2`,
		[]Event{
			NumberStart, // '0'
			None,        // '.'
			None,        // '1'
			None,        // 'e'
			None,        // '+'
			None,        // '2'
			NumberEnd,   // EOF
		},
	},
	{
		`-0.003`,
		[]Event{
			NumberStart, // '-'
			None,        // '0'
			None,        // '.'
			None,        // '0'
			None,        // '0'
			None,        // '3'
			NumberEnd,   // EOF
		},
	},
	{
		`-1000E10`,
		[]Event{
			NumberStart, // '-'
			None,        // '1'
			None,        // '0'
			None,        // '0'
			None,        // '0'
			None,        // 'E'
			None,        // '1'
			None,        // '0'
			NumberEnd,   // EOF
		},
	},
	{
		`true`,
		[]Event{
			BoolStart, // 't'
			None,      // 'r'
			None,      // 'u'
			None,      // 'e'
			BoolEnd,   // EOF
		},
	},
	{
		`false`,
		[]Event{
			BoolStart, // 'f'
			None,      // 'a'
			None,      // 'l'
			None,      // 's'
			None,      // 'e'
			BoolEnd,   // EOF
		},
	},
	{
		`null`,
		[]Event{
			NullStart, // 'n'
			None,      // 'u'
			None,      // 'l'
			None,      // 'l'
			NullEnd,   // EOF
		},
	},
}

func TestScanner(t *testing.T) {
	for _, test := range scannerTests {
		var s = NewScanner()
		var ev Event

		for i, want := range test.out {
			if i < len(test.in) {
				ev = s.Scan(test.in[i])
			} else {
				ev = s.End()
			}

			if ev != want {
				t.Errorf("Scanner(%#q):", test.in)

				for j, prev := range test.out[:i] {
					if j < len(test.in) {
						t.Errorf("  %4q -> %s", test.in[j], prev)
					} else {
						t.Errorf("   EOF -> %s", prev)
					}
				}

				if i < len(test.in) {
					t.Errorf("  %4q -> %s (want %s)", test.in[i], ev, test.out[i])
				} else {
					t.Errorf("   EOF -> %s (want %s)", ev, test.out[i])
				}

				break
			}
		}
	}
}

func TestScannerErrors(t *testing.T) {
	var s = NewScanner()

	if s.Scan('1') != NumberStart || s.Scan('x') != Error {
		t.Fatalf("setup failed")
	}

	err := s.LastError()
	if err == nil {
		t.Fatalf("Scanner.LastError returned nil after Error event")
	}

	if s.Scan('2') != Error || s.LastError().Error() != err.Error() {
		t.Fatalf("Scanner.Scan did not remember previous error")
	}

	if s.End() != Error || s.LastError().Error() != err.Error() {
		t.Fatalf("Scanner.End did not remember previous error")
	}
}
