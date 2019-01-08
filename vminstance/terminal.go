package vminstance

import "io"

// Terminal is output device
type Terminal interface {
	io.Writer
}
