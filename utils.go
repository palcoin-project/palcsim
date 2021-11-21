// Copyright (c) 2014-2016 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	btcutil "github.com/palcoin-project/palcutil"
)

// Row represents a row in the CSV file
// and holds the key and value ints
type Row struct {
	utxoCount int
	txCount   int
}

// readCSV reads the given filename and
// returns a slice of rows
func readCSV(r io.Reader) (map[int32]*Row, error) {
	m := make(map[int32]*Row)
	reader := csv.NewReader(r)
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		b, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, err
		}
		u, err := strconv.Atoi(row[1])
		if err != nil {
			return nil, err
		}
		t, err := strconv.Atoi(row[2])
		if err != nil {
			return nil, err
		}
		m[int32(b)] = &Row{u, t}
	}
	return m, nil
}

func getLogFile(prefix string) (*os.File, error) {
	return os.Create(filepath.Join(AppDataDir, fmt.Sprintf("%s.log", prefix)))
}

// genCertPair generates a key/cert pair to the paths provided.
func genCertPair(certFile, keyFile string) error {
	org := "btcsim autogenerated cert"
	validUntil := time.Now().Add(10 * 365 * 24 * time.Hour)
	cert, key, err := btcutil.NewTLSCertPair(org, validUntil, nil)
	if err != nil {
		return err
	}

	// Write cert and key files.
	if err = ioutil.WriteFile(certFile, cert, 0666); err != nil {
		return err
	}
	if err = ioutil.WriteFile(keyFile, key, 0600); err != nil {
		os.Remove(certFile)
		return err
	}

	return nil
}

// filesExists reports whether the named file or directory exists.
func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}