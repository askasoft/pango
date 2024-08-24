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
	err := ExtractStringFromPdfFile(ctx, bw, name)
	return bw.String(), err
}

func ExtractTextFromPdfReader(ctx context.Context, r io.Reader) (string, error) {
	bw := &bytes.Buffer{}
	err := ExtractStringFromPdfReader(ctx, bw, r)
	return bw.String(), err
}

func buildPdfToTextArgs(input string) []string {
	// See "man pdftotext" for more options.
	args := []string{
		"-layout",  // Maintain (as best as possible) the original physical layout of the text
		"-nopgbrk", // Don't insert page breaks (form feed characters) between pages
		input,      // The input file (-: stdin)
		"-",        // The output file (stdout)
	}
	return args
}

func ExtractStringFromPdfFile(ctx context.Context, w io.Writer, name string) error {
	args := buildPdfToTextArgs(name)
	cmd := exec.CommandContext(ctx, "pdftotext", args...)
	cmd.Stdout = w
	return cmd.Run()
}

func ExtractStringFromPdfReader(ctx context.Context, w io.Writer, r io.Reader) error {
	args := buildPdfToTextArgs("-")
	cmd := exec.CommandContext(ctx, "pdftotext", args...)
	cmd.Stdin = r
	cmd.Stdout = w
	return cmd.Run()
}
