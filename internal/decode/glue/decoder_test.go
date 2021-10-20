package glue

import "testing"

func TestJxlDecoderVersion(t *testing.T) {
	if version := JxlDecoderVersion(); version < 6000 {
		t.Errorf("unknown JXL decoder version: %v", version)
	}
}

func TestJxlDecoderCreate(t *testing.T) {
	decoder := JxlDecoderCreate(SwigcptrJxlMemoryManager(0))
	JxlDecoderDestroy(decoder)
}

func TestBasicDecoding(t *testing.T) {
	simpleImage := []byte{255, 10, 250, 31, 65, 145, 8, 6, 1, 0, 120, 0, 75, 56, 68, 220, 237, 61, 160, 118, 89, 70,
		41, 244, 211, 223, 2, 49, 241, 0, 0, 200, 160, 144, 66, 112, 255, 3, 36, 9, 7, 72,
	}
	decoder := JxlDecoderCreate(SwigcptrJxlMemoryManager(0))
	defer JxlDecoderDestroy(decoder)
	pixelFormat := NewJxlPixelFormat()
	defer DeleteJxlPixelFormat(pixelFormat)
	pixelFormat.SetNum_channels(4)
	pixelFormat.SetEndianness(JxlEndianness(JXL_NATIVE_ENDIAN))
	pixelFormat.SetData_type(JxlDataType(JXL_TYPE_UINT8))
	pixelFormat.SetAlign(0)

	JxlDecoderSetInput(decoder, &simpleImage[0], int64(len(simpleImage)))
	if JxlDecoderSubscribeEvents(decoder, JXL_DEC_BASIC_INFO) != JxlDecoderStatus(JXL_DEC_SUCCESS) {
		t.Fatalf("fail to subscribe")
	}
loop:
	for {
		switch JxlDecoderProcessInput(decoder) {
		case JxlDecoderStatus(JXL_DEC_BASIC_INFO):
			func() {
				info := NewJxlBasicInfo()
				defer DeleteJxlBasicInfo(info)
				if JxlDecoderGetBasicInfo(decoder, info) != JxlDecoderStatus(JXL_DEC_SUCCESS) {
					t.Fatalf("fail to get basic info")
				}
				if info.GetXsize() != 1024 {
					t.Errorf("expected xsize = 1024 but got %v", info.GetXsize())
				}
				if info.GetYsize() != 1024 {
					t.Errorf("expected ysize = 1024 but got %v", info.GetYsize())
				}
				if info.GetNum_color_channels() != 3 {
					t.Errorf("expected channels = 3 but got %v", info.GetNum_color_channels())
				}
			}()
		case JxlDecoderStatus(JXL_DEC_SUCCESS):
			break loop
		case JxlDecoderStatus(JXL_DEC_ERROR):
			t.Fatalf("decoder error")
		}
	}
}
