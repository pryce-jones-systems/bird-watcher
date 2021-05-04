package ImageTools

import (
	"errors"
	"sync"
)

func Mask(image [][]float32, mask [][]float32) ([][]float32, error) {

	// Check that dimensions match
	if len(image) != len(mask) {
		return nil, errors.New("Width mismatch")
	}
	if len(image[0]) != len(mask[0]) {
		return nil, errors.New("Height mismatch")
	}
	imageWidth, imageHeight := Dimensions(image)

	// Create output image
	outputImage := make([][]float32, imageWidth)
	for j := range outputImage {
		outputImage[j] = make([]float32, imageHeight)
	}

	// Create wait group
	var waitGroup sync.WaitGroup
	if imageWidth > imageHeight {
		waitGroup.Add(imageWidth)
	} else {
		waitGroup.Add(imageHeight)
	}

	// Iterate over columns
	for j := 0; j < imageHeight; j++ {

		// Process each row on its own goroutine
		go func(j int) {
			defer waitGroup.Done()

			// Iterate over row
			for i := 0; i < imageWidth; i++ {

				// Apply threshold to each pixel
				currentImagePixel := image[i][j]
				currentMaskPixel := mask[i][j]
				if currentMaskPixel > 0.5 {
					outputImage[i][j] = currentImagePixel
				} else {
					outputImage[i][j] = 0
				}
			}
		} (j)
	}

	return outputImage, nil
}