package pdfx

import (
	"bytes"
	"context"
	"io"
	"os/exec"
)

// sudo apt install poppler-utils

func ExtractTextFromPdfFile(ctx context.Context, name string) (string, error) {
	bw := &bytes.Buffer{}
	err := ExtractStringFromPdfFile(ctx, name, bw)
	return bw.String(), err
}

func ExtractTextFromPdfReader(ctx context.Context, r io.Reader) (string, error) {
	bw := &bytes.Buffer{}
	err := ExtractStringFromPdfReader(ctx, r, bw)
	return bw.String(), err
}

func ExtractStringFromPdfFile(ctx context.Context, name string, w io.Writer) error {
	// See "man pdftotext" for more options.
	args := []string{
		"-layout",  // Maintain (as best as possible) the original physical layout of the text
		"-nopgbrk", // Don't insert page breaks (form feed characters) between pages
		name,       // The input file
		"-",        // The output file (stdout)
	}

	cmd := exec.CommandContext(ctx, "pdftotext", args...)
	cmd.Stdout = w

	return cmd.Run()
}

func ExtractStringFromPdfReader(ctx context.Context, r io.Reader, w io.Writer) error {
	// See "man pdftotext" for more options.
	args := []string{
		"-layout",  // Maintain (as best as possible) the original physical layout of the text
		"-nopgbrk", // Don't insert page breaks (form feed characters) between pages
		"-",        // The input file (stdin)
		"-",        // The output file (stdout)
	}

	cmd := exec.CommandContext(ctx, "pdftotext", args...)
	cmd.Stdin = r
	cmd.Stdout = w

	return cmd.Run()
}
