package jpegxl

import (
	"image"
	"io"

	"github.com/indeedplusplus/go-jpegxl/internal/decoder"
)

func init() {
	image.RegisterFormat("jxl", "\xff\x0a", Decode, DecodeConfig)
}

func Decode(r io.Reader) (image.Image, error) {
	img, _, err := decoder.DecodeReader(r, false)
	return img, err
}

func DecodeConfig(r io.Reader) (image.Config, error) {
	_, cfg, err := decoder.DecodeReader(r, true)
	return cfg, err
}
