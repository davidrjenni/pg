// Copyright (c) 2016 David R. Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/davidrjenni/pg/parser"
	"github.com/davidrjenni/pg/printer"
)

func web(args []string) {
	flags := flag.NewFlagSet("", flag.ExitOnError)
	addr := flags.String("http", ":8080", "HTTP service address (e.g., ':8080')")
	flags.Parse(args)

	http.Handle("/", http.FileServer(http.Dir("static")))

	http.HandleFunc("/api/pg/fmt", fmtHandler)
	http.HandleFunc("/api/pg/parse", parseHandler)

	log.Fatal(http.ListenAndServe(*addr, nil))
}

func fmtHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	g, err := parser.Parse(b, "<input>")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := printer.Fprint(w, g); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func parseHandler(w http.ResponseWriter, r *http.Request) {
}
