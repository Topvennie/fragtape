//go:build !linux

package image

import "fmt"

func ToWebp(_ []byte) ([]byte, error) {
	return nil, fmt.Errorf("ToWebp is only supported on linux builds")
}
