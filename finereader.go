// +build windows

package main

import (
	"os/exec"
)

func ocr(finereader string, in string, out string) error {
	cmd := exec.Command(finereader, in, `/out`, out, `/quit`)
	return cmd.Run()
}
