// Copyright 2019 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package markup

import (
	"strings"

	"github.com/gohugoio/hugo/markup/org"

	"github.com/gohugoio/hugo/markup/asciidoc"
	"github.com/gohugoio/hugo/markup/blackfriday"
	"github.com/gohugoio/hugo/markup/converter"
	"github.com/gohugoio/hugo/markup/mmark"
	"github.com/gohugoio/hugo/markup/pandoc"
	"github.com/gohugoio/hugo/markup/rst"
)

func NewConverterProvider(cfg converter.ProviderConfig) (ConverterProvider, error) {
	converters := make(map[string]converter.Provider)

	add := func(p converter.NewProvider, aliases ...string) error {
		c, err := p.New(cfg)
		if err != nil {
			return err
		}
		addConverter(converters, c, aliases...)
		return nil
	}

	if err := add(blackfriday.Provider, "md", "markdown", "blackfriday"); err != nil {
		return nil, err
	}
	if err := add(mmark.Provider, "mmark"); err != nil {
		return nil, err
	}
	if err := add(asciidoc.Provider, "asciidoc"); err != nil {
		return nil, err
	}
	if err := add(rst.Provider, "rst"); err != nil {
		return nil, err
	}
	if err := add(pandoc.Provider, "pandoc"); err != nil {
		return nil, err
	}
	if err := add(org.Provider, "org"); err != nil {
		return nil, err
	}

	return &converterRegistry{converters: converters}, nil
}

type ConverterProvider interface {
	Get(name string) converter.Provider
}

type converterRegistry struct {
	// Maps name (md, markdown, blackfriday etc.) to a converter provider.
	// Note that this is also used for aliasing, so the same converter
	// may be registered multiple times.
	// All names are lower case.
	converters map[string]converter.Provider
}

func (r *converterRegistry) Get(name string) converter.Provider {
	return r.converters[strings.ToLower(name)]
}

func addConverter(m map[string]converter.Provider, c converter.Provider, aliases ...string) {
	for _, alias := range aliases {
		m[alias] = c
	}
}
