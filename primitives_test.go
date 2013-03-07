package jo

import (
	"testing"
)

var parseIntTests = []struct {
	in   string
	want int64
	err  error
}{
	{`0`, 0, nil},
	{`-0`, 0, nil},
	{`-10.0e2`, -1000, nil},
	{`0.5`, 0, ErrFraction},
	{`100e-3`, 0, ErrFraction},
	{`0e-2`, 0, nil},
	{`-1000`, -1000, nil},
	{`123456789`, 123456789, nil},
	{`-0.9223372036854775808e19`, -9223372036854775808, nil},
	{`0.9223372036854775807e19`, 9223372036854775807, nil},
	{`-0.2147483648e10`, -2147483648, nil},
	{`0.2147483647e10`, 2147483647, nil},
	{`-32768`, -32768, nil},
	{`32767`, 32767, nil},
	{`-128`, -128, nil},
	{`127`, 127, nil},
	{`-000000000000.00000000000000000000000000000000000000000000`, 0, nil},
	{`10000000000000000000000000000000000000000000000000000000e-55`, 1, nil},
	{`-9223372036854775808`, -9223372036854775808, nil},
	{`-9223372036854775809`, 0, ErrOverflow},
	{`9223372036854775807`, 9223372036854775807, nil},
	{`9223372036854775808`, 0, ErrOverflow},
	{`  123  `, 123, nil},
	{` 1 2 3 `, 0, ErrSyntax},
	{``, 0, ErrSyntax},
	{`--10`, 0, ErrSyntax},
	{`10a`, 0, ErrSyntax},
	{`0x20`, 0, ErrSyntax},
}

func TestParseInt(t *testing.T) {
	for _, test := range parseIntTests {
		got, err := ParseInt([]byte(test.in))
		if got != test.want || err != test.err {
			t.Errorf("ParseInt(%q):", test.in)
			t.Errorf("   got %d, %v", got, err)
			t.Errorf("  want %d, %v", test.want, test.err)
		}
	}
}

var parseUintTests = []struct {
	in   string
	want uint64
	err  error
}{
	{`0`, 0, nil},
	{`10`, 10, nil},
	{`101.00`, 101, nil},
	{`0.0101e4`, 101, nil},
	{`0.0101e3`, 0, ErrFraction},
	{`0e-2`, 0, nil},
	{`18446744073709551615`, 18446744073709551615, nil},
	{`1.8446744073709551615e19`, 18446744073709551615, nil},
	{`0.0018446744073709551615e22`, 18446744073709551615, nil},
	{`0.0018446744073709551615e21`, 0, ErrFraction},
	{`0.18446744073709551615e20`, 18446744073709551615, nil},
	{`0.4294967295e10`, 4294967295, nil},
	{`0.65535e5`, 65535, nil},
	{`0.255e3`, 255, nil},
	{`000000000000.000000000000000000000000000000000000000000000`, 0, nil},
	{`10000000000000000000000000000000000000000000000000000000e-55`, 1, nil},
	{` 1 2 3 `, 0, ErrSyntax},
	{``, 0, ErrSyntax},
	{`10a`, 0, ErrSyntax},
	{`0x20`, 0, ErrSyntax},
}

func TestParseUint(t *testing.T) {
	for _, test := range parseUintTests {
		got, err := ParseUint([]byte(test.in))
		if got != test.want || err != test.err {
			t.Errorf("ParseUint(%q):", test.in)
			t.Errorf("   got %d, %v", got, err)
			t.Errorf("  want %d, %v", test.want, test.err)
		}
	}
}
