// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package static

import (
	"bytes"
	"io/ioutil"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"unicode"
)

func TestStaticIsUpToDate(t *testing.T) {
	if runtime.GOOS == "android" {
		t.Skip("files not available on android")
	}
	oldBuf, err := ioutil.ReadFile("static.go")
	if err != nil {
		t.Errorf("error while reading static.go: %v\n", err)
	}

	newBuf, err := Generate()
	if err != nil {
		t.Errorf("error while generating static.go: %v\n", err)
	}

	// Generate stamps time.Now().Year() into the copyright header (the first
	// line), so the committed static.go goes "stale" every new year. Compare
	// the content after the copyright line only.
	if bytes.Compare(dropFirstLine(oldBuf), dropFirstLine(newBuf)) != 0 {
		t.Error(`static.go is stale.  Run:
  $ go generate golang.org/x/tools/godoc/static
  $ git diff
to see the differences.`)

	}
}

// TestJQueryHTMLManipulation checks that the jquery.js served by godoc does
// not mutate untrusted HTML passed to .html()/.append() and friends:
// <option> markup is wrapped in <select multiple> only where the browser
// cannot parse it bare (CVE-2020-11023), and self-closing tags are no longer
// expanded with "<$1></$2>" (CVE-2020-11022).
func TestJQueryHTMLManipulation(t *testing.T) {
	js := Files["jquery.js"]
	for _, bad := range []string{
		`option:[1,"<select multiple='multiple'>","</select>"]`,
		`"<$1></$2>"`,
	} {
		if strings.Contains(js, bad) {
			t.Errorf("jquery.js contains vulnerable code %q", bad)
		}
	}
	if !strings.Contains(js, `.support.option||(`) {
		t.Error("jquery.js does not gate the <option> wrapper on support.option")
	}
}

// dropFirstLine returns b without its first line (the copyright header).
func dropFirstLine(b []byte) []byte {
	if i := bytes.IndexByte(b, '\n'); i >= 0 {
		return b[i+1:]
	}
	return b
}

// TestAppendQuote ensures that AppendQuote produces a valid literal.
func TestAppendQuote(t *testing.T) {
	var in, out bytes.Buffer
	for r := rune(0); r < unicode.MaxRune; r++ {
		in.WriteRune(r)
	}
	appendQuote(&out, in.Bytes())
	in2, err := strconv.Unquote(out.String())
	if err != nil {
		t.Fatalf("AppendQuote produced invalid string literal: %v", err)
	}
	if got, want := in2, in.String(); got != want {
		t.Fatal("AppendQuote modified string") // no point printing got/want: huge
	}
}
