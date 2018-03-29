/*
 * Checksum utility
 * Copyright (c) 2018 Alexandre Roman <alexandre.roman@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

// Utility functions for computing checksums from files.
package checksum

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
	"os"

	"github.com/alexandreroman/checksum/logger"
)

func genericChecksum(path string, algoName string, algoHash hash.Hash) (string, error) {
	logger.Debug("Computing %s checksum for file: %s", algoName, path)
	if file, err := os.Open(path); err != nil {
		return "", fmt.Errorf("error while opening file: %v", err)
	} else {
		defer file.Close()
		if _, err := io.Copy(algoHash, file); err != nil {
			return "", fmt.Errorf("error while reading file: %v", err)
		} else {
			// Convert hash result to hexadecimal string.
			r := fmt.Sprintf("%x", algoHash.Sum(nil))
			logger.Debug("%s for %s: %s", algoName, path, r)
			return r, nil
		}
	}
}

// Compute MD5 checksum from a file.
func MD5(path string) (string, error) {
	return genericChecksum(path, "MD5", md5.New())
}

// Compute SHA-1 checksum from a file.
func SHA1(path string) (string, error) {
	return genericChecksum(path, "SHA-1", sha1.New())
}

// Compute SHA-256 checksum from a file.
func SHA256(path string) (string, error) {
	return genericChecksum(path, "SHA-256", sha256.New())
}
