// +build windows

package std

func init() {
	Stdin = NewRawReader()
	Stdout = NewANSIWriter(Stdout)
	Stderr = NewANSIWriter(Stderr)
}
