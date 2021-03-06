/*
   Copyright 2018 Assetnote

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package wordswithext

import (
	"github.com/assetnote/commonspeak2/log"
	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
	"golang.org/x/net/context"
	"io"
	"os"
	"strings"
	"fmt"
)

func query(client *bigquery.Client, ctx context.Context, compiledSql string) (*bigquery.RowIterator, error) {
	query := client.Query(compiledSql)
	return query.Read(ctx)
}


func handleResults(w io.Writer, iter *bigquery.RowIterator, outputFile string, silent bool, verbose bool) error {
	fields := log.Fields{
		"Mode":       "WordsWithExt",
		"Source":     "Github",
	}
	file, err := os.Create(outputFile)
    if err != nil {
    	fields["Filename"] = outputFile
    	fields["Error"] = err.Error()
        log.WithFields(fields).Fatal("Cannot create output file")
    }
    defer file.Close()
    totalRows := 0
	for {
		var row ExtPaths
		err := iter.Next(&row)
		if err == iterator.Done {
			if !silent {
				log.WithFields(fields).Infof("Total rows extracted %d.", totalRows)
			}
			return nil
		}
		if err != nil {
			return err
		}
		// Save to output file
		fmt.Fprintf(file, "%s\n", row.Path)
		// Print to console if verbose mode is on
		if verbose {
			fmt.Fprintf(w, "Path: %s Count: %s\n", row.Path, row.PathCount.String())
		}
		totalRows++
	}
	
}

func convertExtensionsToRegex(extensions string) string {
	extList := strings.Split(extensions, ",")
	convertedList := make([]string, 0)
	for _, elm := range extList {
		convertedElm := fmt.Sprintf("\\.%s", elm)
		convertedList = append(convertedList, convertedElm)
	}
	finalRegex := strings.Join(convertedList, "|")
	return finalRegex
}