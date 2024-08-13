// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package template

import (
	"bytes"
	"html/template"
	"strings"

	"github.com/cilium/cilium-cli/connectivity/tests"
)

// Render executes temp template with data and returns the result
func Render(temp string, data any) (string, error) {
	fns := template.FuncMap{
		"trimSuffix": func(in, suffix string) string { return strings.TrimSuffix(in, suffix) },
		"isIPv6":     func(in string) bool { return tests.FormattedAsIPv6(in) },
	}

	tm, err := template.New("template").Funcs(fns).Parse(temp)
	if err != nil {
		return "", err
	}

	buf := bytes.NewBuffer(nil)
	if err := tm.Execute(buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
