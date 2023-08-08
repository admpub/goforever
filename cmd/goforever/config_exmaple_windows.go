//go:build windows

package main

func init() {
	exampleConfig += `[process.Options]
HideWindow = false
`
}
