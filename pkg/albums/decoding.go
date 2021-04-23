package albums

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

type DecoderFunc func(io.Reader) (image.Image, error)

type Decoder struct {
	decoders map[string]DecoderFunc
}

func NewDecoder() *Decoder {
	decoders := defaultDecoders()
	return &Decoder{decoders}
}

func (d *Decoder) Decode(buff io.Reader, format string) (image.Image, error) {
	b := bytes.NewBuffer([]byte{})
	io.Copy(b, buff)
	return d.decode(b, format)
}

func (d *Decoder) decode(buff io.Reader, format string) (image.Image, error) {
	f, err := d.Lookup(format)
	if err != nil {
		log.Println(err)
		return d.tryDecode(buff)
	}
	return f(buff)
}

func (d *Decoder) tryDecode(buff io.Reader) (image.Image, error) {
	img, codec, err := image.Decode(buff)
	log.Printf("used codec \"%s\" to decode image", codec)
	return img, err
}

func (d *Decoder) Format(s string) string {
	return strings.TrimPrefix(filepath.Ext(s), ".")
}

func (d *Decoder) Lookup(format string) (DecoderFunc, error) {
	dec, ok := d.decoders[format]
	if !ok {
		return nil, errors.Errorf("decoder func not found for format %s", format)
	}
	return dec, nil
}

func defaultDecoders() map[string]DecoderFunc {
	d := make(map[string]DecoderFunc)
	d["png"] = png.Decode
	d["jpg"] = jpeg.Decode
	d["jpeg"] = jpeg.Decode
	return d
}
