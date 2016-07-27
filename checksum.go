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
