//
// Copyright 2018-2019 Solus Project <copyright@getsol.us>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package config

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// ConfigFiles is a list of config files to check for, in order of ascending priority
var ConfigFiles = []string{
	"/usr/share/defaults/etc/aa-lsm-hook.conf",
	"/etc/aa-lsm-hook.conf",
}

// ProfileDirs retriesves a list of directories to search for profiles, listed one per line in each of the ConfigFiles
func ProfileDirs() ([]string, error) {
	dirs := make([]string, 0)

	for _, file := range ConfigFiles { // For each file
		if _, err := os.Stat(file); !os.IsNotExist(err) { // If the file exists
			f, err := os.Open(file) // Open the File
			if err != nil {
				return dirs, err
			}
			r := bufio.NewReader(f)
			// Read each line
			for {
				raw, _, err := r.ReadLine()
				// skip check for empty lines
				if len(raw) > 0 {
					path := string(raw)
					if _, e := os.Stat(path); !os.IsNotExist(e) {
						dirs = append(dirs, path) // Add to list if the directory exists
					}
				}
				// Check for I/O error or last line
				if err != nil {
					if err != io.EOF {
						_ = f.Close()
						return dirs, err
					}
					break
				}
			}
			// Close files
			_ = f.Close()
		}
	}
	if len(dirs) == 0 { // check for no configured profile paths
		return dirs, fmt.Errorf("failed to find any configured profile paths")
	}
	// return all found directories
	return dirs, nil
}
