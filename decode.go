package jpegxl

import (
	"image"
	"io"

	"github.com/jiahuif/go-jpegxl/internal/decode"
)

func init() {
	image.RegisterFormat("jxl", "\xff\x0a", Decode, DecodeConfig)
}

func Decode(r io.Reader) (image.Image, error) {
	img, _, err := decode.FromReader(r, false)
	return img, err
}

func DecodeConfig(r io.Reader) (image.Config, error) {
	_, cfg, err := decode.FromReader(r, true)
	return cfg, err
}
