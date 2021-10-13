package jpegxl

import (
	"bytes"
	"encoding/base64"
	"image/png"
	"testing"
)

func TestDecodeSmallImage(t *testing.T) {
	for _, tc := range []struct {
		name            string
		referenceBase64 string
		imageBase64     string
	}{
		{
			name: "16bit",
			referenceBase64: "iVBORw0KGgoAAAANSUhEUgAAAGQAAAAyEAYAAAD6paL9AAAAlklEQVR42u3YsQ2DMBAF0DPKAoyQ" +
				"ksoCKasxRFaLFOSKkhGyAc4G0CCL4r0FbJ38fT5HAAAAAAAAAAAAAAAAAAAAAAAAjaS7bOT5fS2/" +
				"tdZW623TZ+yHlBwBjnRKAAICAgJXeyjBRTfNHDny+Qy1v6NEMfvoICAg4InVhG9XdBAQEBAQEBAQ" +
				"EEBAAAAAAAAAAAAAAOB2/ljJEGk1KBNVAAAAAElFTkSuQmCC",
			imageBase64: "/wqIcfwIfoAECBAQAOAASzhpmMqD9ysoCID7/28Mjo/xY/wLyACw+5CMwY5hOqwPuGkEknXpQND1" +
				"FID+M5eCdr4JTtFZNw4=",
		},
		{
			name: "8bit",
			referenceBase64: "iVBORw0KGgoAAAANSUhEUgAAAGQAAAAyCAYAAACqNX6+AAAAaElEQVR42u3WuxGAIBBAQWBsgBYs" +
				"gv5ji7AFOxBjh8whOdnNIeDxSwkAAAAAAAAAAAguz5hkb1f/Mu48apbgrVgCQRAkji3IrhneqHvS" +
				"++eEIMhy315XlhPiykIQQRAEQQAAAAAAfugBhxgIHWoqA24AAAAASUVORK5CYII=",
			imageBase64: "/wqIcbASCBAQALgASygl1uFWUBAA9fx8jLvxLyADACD/qAIQYQDAzP5UJue/7r8ywQ56GCtE1KPz" +
				"AA==",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			referenceImageBytes, _ := base64.StdEncoding.DecodeString(tc.referenceBase64)
			ref, _ := png.Decode(bytes.NewReader(referenceImageBytes))
			imageBytes, _ := base64.StdEncoding.DecodeString(tc.imageBase64)
			image, _ := Decode(bytes.NewReader(imageBytes))
			for i := 0; i < ref.Bounds().Dx(); i++ {
				for j := 0; j < ref.Bounds().Dy(); j++ {
					rR, rG, rB, rA := ref.At(i, j).RGBA()
					iR, iG, iB, iA := image.At(i, j).RGBA()
					if rR != iR || rG != iG || rB != iB || rA != iA {
						t.Errorf("mismatched pixel at %d, %d: expected %v, but got %v", i, j, ref.At(i, j), image.At(i, j))
					}
				}
			}
		})
	}

}
