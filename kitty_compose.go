package main

import (
	"bytes"
	"context"
	"image"
	"image/png"
)

type completedFrame struct {
	jobID uint64
	png   []byte
	err   error
}

func encodeWholeFrame(
	ctx context.Context,
	canvas *image.RGBA,
) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	var buffer bytes.Buffer

	encoder := png.Encoder{
		CompressionLevel: png.BestSpeed,
	}

	if err := encoder.Encode(&buffer, canvas); err != nil {
		return nil, err
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
