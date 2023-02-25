//nolint:lll
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/moemoe89/file-storage/cmd/cli/client"
)

// usageText is a message to describe how to use the CLI.
const usageText = `file-store is a command line interface (CLI) program designed for file operations such as uploading, listing, and deleting files. The program allows users to upload a file either from a local path or a URL, list the uploaded files on the server, and delete them..

In addition to file operations, the program offers an optional --bucket flag to select the target bucket. If the flag is not defined, the program uses the default bucket.

This program offers a convenient way to manage files by providing simple commands to upload, list, and delete files. For example:

Example:
	fs-store upload-file logo.jpeg
	fs-store delete-file logo.jpeg
	fs-store list-files

`

const (
	// uploadFileArg is an argument for upload file
	uploadFileArg = "upload-file"
	// deleteFileArg is an argument for list file
	deleteFileArg = "delete-file"
	// listFileArg is an argument for delete file
	listFileArg = "list-files"
)

var flagMap = map[string]struct{}{
	"bucket":       {},
	"source":       {},
	"filename":     {},
	"content_type": {},
	"max_size":     {},
}

var (
	// bucket is a flag to do file operation to the specific bucket. If the bucket not specified, the CLI will use `default` bucket.
	bucket = flag.String("bucket", "default", "Specify the file operation to the specific bucket, `default` bucket will use if the value not specified.")
	// source is a flag to select source file when do upload. If the source not specified, the CLI will use from `file` source.
	source = flag.String("source", "file", "Specify the source of upload file. Can be from file or URL. By default it will be upload from file.")
	// filename is a flag to naming the file when do upload. If not specified, automatically get from file path or URL.
	filename = flag.String("filename", "", "Specify the file name of upload file. If not specified, automatically get from file path or URL.")
	// filename is a flag to naming the file when do upload. If not specified, automatically get from file path or URL.
	contentType = flag.String("content_type", "", "Specify the file name of upload file. If not specified, automatically get from file path or URL.")
	// max_size is a flag to validates the max size of file when do upload. If not specified, the value is 0 which mean don't need to validate.
	maxSize = flag.Int64("max_size", 0, "Specify the validation for max size file of upload file. If not specified, validation will be ignored")
)

func main() {
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() < 1 {
		usage()
		log.Fatal("Expected minimum one argument")
	}

	args := flag.Args()

	if args[0] != listFileArg && flag.NArg() < 2 {
		usage()
		log.Fatal("Expected minimum two argument")
	}

	switch args[0] {
	case uploadFileArg:
		var contentTypes []string

		if *contentType != "" {
			contentTypes = strings.Split(*contentType, ",")
		}

		if strings.ToLower(*source) == "file" {
			client.UploadFromFile(args[1], *filename, contentTypes, *maxSize)
		} else if strings.ToLower(*source) == "url" {
			client.UploadFromURL(args[1], *filename, contentTypes, *maxSize)
		}
	case deleteFileArg:
		client.DeleteFile(*bucket, args[1])
	case listFileArg:
		client.ListFile(*bucket)
	default:
		log.Fatalf("Invalid argument! Argument should be one of these: %s, %s, %s", uploadFileArg, deleteFileArg, listFileArg)
	}

	os.Exit(0)
}

func usage() {
	_, _ = io.WriteString(os.Stderr, usageText)

	flag.VisitAll(func(f *flag.Flag) {
		if _, ok := flagMap[f.Name]; ok {
			_, _ = fmt.Fprintf(os.Stderr, "  -%s %s\n", f.Name, f.DefValue)
			_, _ = fmt.Fprintf(os.Stderr, "\t%s\n", f.Usage)
		}
	})
}
