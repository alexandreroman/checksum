// This command line utility computes checksums for files.
package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/alexandreroman/checksum/logger"
	"gopkg.in/alecthomas/kingpin.v2"
)

// Command line parsing variables.
var (
	app     = kingpin.New("checksum", "A command line utility for computing checksums.")
	verbose = app.Flag("verbose", "Enable verbose mode.").Bool()

	md5Cmd  = app.Command("md5", "Compute MD5 checksum.")
	md5File = md5Cmd.Arg("file", "File to read.").Required().ExistingFileOrDir()

	sha1Cmd  = app.Command("sha1", "Compute SHA-1 checksum.")
	sha1File = sha1Cmd.Arg("file", "File to read.").Required().ExistingFileOrDir()

	sha256Cmd  = app.Command("sha256", "Compute SHA-256 checksum.")
	sha256File = sha256Cmd.Arg("file", "File to read.").Required().ExistingFileOrDir()
)

// Declare generic checksum function signature.
type ChecksumFunc func(string) (string, error)

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
func md5Checksum(path string) (string, error) {
	return genericChecksum(path, "MD5", md5.New())
}

// Compute SHA-1 checksum from a file.
func sha1Checksum(path string) (string, error) {
	return genericChecksum(path, "SHA-1", sha1.New())
}

// Compute SHA-256 checksum from a file.
func sha256Checksum(path string) (string, error) {
	return genericChecksum(path, "SHA-256", sha256.New())
}

// Application entry point.
func main() {
	kingpin.UsageTemplate(kingpin.CompactUsageTemplate).Version("1.0.0").Author("Alexandre Roman")
	kingpin.CommandLine.Help = "A simple application for computing checksums."

	// Base element: can be a file or a directory.
	var root string
	// Checksum function to use (depending on command line argument).
	var algo ChecksumFunc

	// Parse command line arguments.
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case md5Cmd.FullCommand():
		root = *md5File
		algo = md5Checksum
	case sha1Cmd.FullCommand():
		root = *sha1File
		algo = sha1Checksum
	case sha256Cmd.FullCommand():
		root = *sha256File
		algo = sha256Checksum
	default:
		logger.Fatal("Not yet implemented!")
	}

	// Setup logging.
	logger.Verbose = *verbose

	// Initialize a channel storing pending files.
	backlog := make(chan string)

	// This wait group is used to synchronize goroutines:
	// the main goroutine will wait for them.
	var wg sync.WaitGroup

	// Start the file scanner in a goroutine.
	wg.Add(1)
	go func() {
		defer wg.Done()

		logger.Debug("Looking for files in root directory: %s", root)
		// Gotta catch'em all!
		filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				logger.Debug("Found file: %s", path)
				// Add this file to the backlog.
				backlog <- path
			}
			return nil
		})

		// At this point we have found all files to read:
		// we can safely close our backlog.
		close(backlog)
	}()

	// Process backlog (until it is closed by the file scanner, see above).
	for path := range backlog {
		wg.Add(1)

		// Start a new goroutine to execute checksum algorithm.
		go func(path string) {
			defer wg.Done()
			if checksum, err := algo(path); err != nil {
				logger.Debug("Checksum error: %v", err)
			} else {
				logger.Info("%s: %s", path, checksum)
			}
		}(path)
	}

	// Wait for goroutines to complete.
	wg.Wait()
}
