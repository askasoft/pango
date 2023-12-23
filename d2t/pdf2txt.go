package d2t

import (
	"context"
	"io"
	"os/exec"
)

// sudo apt install poppler-utils

func ExtractTextFromPdfFile(name string, w io.Writer) error {
	// See "man pdftotext" for more options.
	args := []string{
		"-layout",  // Maintain (as best as possible) the original physical layout of the text
		"-nopgbrk", // Don't insert page breaks (form feed characters) between pages
		name,       // The input file
		"-",        // The output file (stdout)
	}

	cmd := exec.CommandContext(context.Background(), "pdftotext", args...)

	cmd.Stdout = w

	return cmd.Run()
}

func ExtractTextFromPdfReader(r io.Reader, w io.Writer) error {
	// See "man pdftotext" for more options.
	args := []string{
		"-layout",  // Maintain (as best as possible) the original physical layout of the text
		"-nopgbrk", // Don't insert page breaks (form feed characters) between pages
		"-",        // The input file (stdin)
		"-",        // The output file (stdout)
	}

	cmd := exec.CommandContext(context.Background(), "pdftotext", args...)
	cmd.Stdin = r
	cmd.Stdout = w

	return cmd.Run()
}
