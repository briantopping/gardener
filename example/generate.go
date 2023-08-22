/*
 * Copyright (c) 2023 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *
 *  You may obtain a copy of the License at
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

//go:generate go run generate.go github.com/gardener/gardener/example gardener-local provider-local

package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"text/template"
)

type templateData struct {
	Package     string
	PackagePath string
	Children    []string
}

func walk(path string, dir string) error {
	path = filepath.Join(path, dir)
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	subdirs := []string{}
	for _, file := range files {
		if file.IsDir() {
			subdirs = append(subdirs, filepath.Join(path, file.Name()))
			if err := walk(path, file.Name()); err != nil {
				return err
			}
		}
	}

	fmt.Printf("Processing directory: %s\n", path)
	return generate(path, dir, subdirs)
}

func generate(path string, pkg string, subdirs []string) error {
	subdirFilePath := filepath.Join(cwd, path, "anchor.go")
	subdirFile, err := os.Create(subdirFilePath)
	if err != nil {
		return err
	}
	defer subdirFile.Close()

	w := bufio.NewWriter(subdirFile)

	err = fileTemplate.Execute(w, templateData{
		Package:     packageRegex.ReplaceAllString(pkg, "_"),
		PackagePath: project,
		Children:    subdirs,
	})
	if err != nil {
		return err
	}

	return w.Flush()
}

func main() {
	var err error
	fileTemplate, err = template.New("file").Parse(tmpl)
	check(err)
	cwd, err = os.Getwd()
	check(err)

	project = os.Args[1]
	for _, path := range os.Args[2:] {
		check(walk(".", path))
	}
	check(generate(".", "main", os.Args[2:]))
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var (
	fileTemplate *template.Template
	packageRegex = regexp.MustCompile(`\W`)
	cwd          string
	project      string
	tmpl         = `/*
 * Copyright (c) 2023 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *
 *  You may obtain a copy of the License at
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package {{.Package}}

import (
{{range .Children}}    _ "{{$.PackagePath}}/{{.}}"{{println}}{{end}}
)
`
)
