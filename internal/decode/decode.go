package decode

import (
	"fmt"
	"image"
	"image/color"
	"io"
	"unsafe"

	"github.com/indeedplusplus/go-jpegxl/internal/decode/glue"
)

var ErrDecodeError = fmt.Errorf("JXL_DEC_ERROR")

const BufferSize = 4096

type BasicInfo struct {
	Width, Height int
}

func FromReader(r io.Reader, headerOnly bool) (image.Image, image.Config, error) {
	decoder := glue.JxlDecoderCreate(glue.SwigcptrJxlMemoryManager(0))
	defer glue.JxlDecoderDestroy(decoder)

	requested := glue.JXL_DEC_BASIC_INFO
	if !headerOnly {
		requested = requested | glue.JXL_DEC_FULL_IMAGE
	}
	if glue.JxlDecoderSubscribeEvents(decoder, requested) != glue.JxlDecoderStatus(glue.JXL_DEC_SUCCESS) {
		return nil, image.Config{}, ErrDecodeError
	}

	pixelFormat := glue.NewJxlPixelFormat()
	defer glue.DeleteJxlPixelFormat(pixelFormat)
	pixelFormat.SetNum_channels(4) // RGBA always
	pixelFormat.SetEndianness(glue.JxlEndianness(glue.JXL_NATIVE_ENDIAN))
	pixelFormat.SetData_type(glue.JxlDataType(glue.JXL_TYPE_UINT8))
	pixelFormat.SetAlign(0)

	var config image.Config
	buf := make([]byte, BufferSize)
	var imageBuffer []uint8
	var bitsPerSample uint
loop:
	for {
		switch glue.JxlDecoderProcessInput(decoder) {
		case glue.JxlDecoderStatus(glue.JXL_DEC_ERROR):
			return nil, image.Config{}, ErrDecodeError
		case glue.JxlDecoderStatus(glue.JXL_DEC_NEED_MORE_INPUT):
			remaining := int(glue.JxlDecoderReleaseInput(decoder))
			for i := 0; i < remaining; i++ {
				buf[i] = buf[len(buf)-remaining+i]
			}
			n, err := r.Read(buf[remaining:])
			if err != nil {
				return nil, image.Config{}, err
			}
			if glue.JxlDecoderSetInput(decoder, &buf[0], int64(n+remaining)) != glue.JxlDecoderStatus(glue.JXL_DEC_SUCCESS) {
				return nil, image.Config{}, ErrDecodeError
			}
		case glue.JxlDecoderStatus(glue.JXL_DEC_NEED_IMAGE_OUT_BUFFER):
			var size int64
			if glue.JxlDecoderImageOutBufferSize(decoder, pixelFormat, &size) != glue.JxlDecoderStatus(glue.JXL_DEC_SUCCESS) {
				return nil, image.Config{}, ErrDecodeError
			}
			imageBuffer = make([]uint8, size)
			if glue.JxlDecoderSetImageOutBuffer(decoder, pixelFormat, uintptr(unsafe.Pointer(&imageBuffer[0])), size) != glue.JxlDecoderStatus(glue.JXL_DEC_SUCCESS) {
				return nil, image.Config{}, ErrDecodeError
			}
		case glue.JxlDecoderStatus(glue.JXL_DEC_BASIC_INFO):
			var err error
			config, err = func() (image.Config, error) {
				basicInfo := glue.NewJxlBasicInfo()
				defer glue.DeleteJxlBasicInfo(basicInfo)
				if glue.JxlDecoderGetBasicInfo(decoder, basicInfo) != glue.JxlDecoderStatus(glue.JXL_DEC_SUCCESS) {
					return image.Config{}, ErrDecodeError
				}
				bitsPerSample = basicInfo.GetBits_per_sample()
				if bitsPerSample == 16 {
					pixelFormat.SetData_type(glue.JxlDataType(glue.JXL_TYPE_UINT16))
					pixelFormat.SetEndianness(glue.JxlEndianness(glue.JXL_BIG_ENDIAN))
				}
				return image.Config{
					ColorModel: color.RGBAModel,
					Width:      int(basicInfo.GetXsize()),
					Height:     int(basicInfo.GetYsize()),
				}, nil
			}()
			if err != nil {
				return nil, image.Config{}, err
			}
			if headerOnly {
				return nil, config, err
			}
		case glue.JxlDecoderStatus(glue.JXL_DEC_FULL_IMAGE):
		case glue.JxlDecoderStatus(glue.JXL_DEC_SUCCESS):
			break loop
		}
	}
	var img image.Image
	if bitsPerSample == 16 {
		img = &image.NRGBA64{
			Pix:    imageBuffer,
			Stride: config.Width * 8,
			Rect: image.Rectangle{Min: image.Point{}, Max: image.Point{
				X: config.Width,
				Y: config.Height,
			}},
		}
	} else {
		img = &image.NRGBA{
			Pix:    imageBuffer,
			Stride: config.Width * 4,
			Rect: image.Rectangle{Min: image.Point{}, Max: image.Point{
				X: config.Width,
				Y: config.Height,
			}},
		}
	}
	return img, config, nil
}
