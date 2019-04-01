// Copyright © 2018 Moises P. Sena <moisespsena@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"io"
	"os"
	"fmt"
	"bufio"

	"path/filepath"
	"github.com/spf13/cobra"
	"github.com/moisespsena/go-topsort"
	"github.com/moisespsena-go/error-wrap"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "topsort [flags] [file...]",
	Short: "Topological sorting algorithms",
	Long: `Topological sorting algorithms are especially useful for dependency calculation, 
and so this particular implementation is mainly intended for this purpose. 

As a result, the direction of edges and the order of the results may seem reversed 
compared to other implementations of topological sorting.

Home Page: https://github.com/moisespsena/go-topsort

EXAMPLES
--------

$ echo "A-B,B-C,B-D,E-D,F" | topsort
$ echo "A-B:B-C:B-D:E-D:F" | topsort -p :
$ topsort pairs.txt

Ordered input files including STDIN (file name is '-')
$ echo "A-B,B-C,B-D,E-D,F" | topsort pairs1.txt pairs2.txt - pairs3.txt
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		pairSep, err := cmd.Flags().GetString("pair-sep")
		if err != nil {
			return err
		}
		if pairSep == "" {
			return fmt.Errorf("Pair separator is empty.")
		}

		edgeSep, err := cmd.Flags().GetString("edge-sep")
		if err != nil {
			return err
		}
		if edgeSep == "" {
			return fmt.Errorf("Edge separator is empty.")
		}

		topSort, err := cmd.Flags().GetBool("top-sort")
		if err != nil {
			return err
		}

		filesNamesMap := map[string]bool{}
		var filesNames []string

		if len(args) == 0 {
			filesNames = append(filesNames, "-")
		} else {
			for _, f := range args {
				if f != "-" {
					f, err = filepath.Abs(f)
					if err != nil {
						return errwrap.Wrap(err, "Get Abs path")
					}
				}
				if _, ok := filesNamesMap[f]; !ok {
					filesNamesMap[f] = true
					filesNames = append(filesNames, f)
				}
			}
		}

		graph := topsort.NewGraph()

		read := func(f string) (err error) {
			var r io.Reader
			if f == "-" {
				r = os.Stdin
			} else {
				fi, err := os.Open(f)
				if err != nil {
					return err
				}
				defer fi.Close()
				r = fi
			}
			b := bufio.NewReader(r)
			return graph.ParseLines(edgeSep, pairSep, func() (string, error) {
				line, err := b.ReadString('\n')
				if err != nil {
					return "", err
				}
				if line[len(line)-2] == '\r' {
					return line[0 : len(line)-2], nil
				}
				return line[0 : len(line)-1], nil
			})
		}

		for _, f := range filesNames {
			err = read(f)
			if err != nil {
				return errwrap.Wrap(err, "Read from %q failed", f)
			}
		}

		var results []string

		if topSort {
			results, err = graph.TopSort()
			if err != nil {
				return errwrap.Wrap(err, "Top Sort classifier")
			}
		} else {
			results, err = graph.DepthFirst()
			if err != nil {
				return errwrap.Wrap(err, "Depth-First classifier")
			}
		}
		for _, r := range results {
			os.Stdout.WriteString(r + "\n")
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("top-sort", "T", false, "Use Topological node classifier, otherwise, Depth-first classifier.")
	rootCmd.Flags().StringP("pair-sep", "p", ",", "Set the pairs separator")
	rootCmd.Flags().StringP("edge-sep", "e", "-", "Set the edge separator")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
