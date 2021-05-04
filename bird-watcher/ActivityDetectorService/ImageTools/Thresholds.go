package ImageTools

import (
	"sync"
)

/*
 * Thresholds an image with a single threshold
 */
func SingleThreshold(image [][]float32, threshold float32) [][]float32 {

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
				currentPixel := image[i][j]
				if currentPixel < threshold {
					outputImage[i][j] = 0
				} else {
					outputImage[i][j] = 1
				}
			}
		} (j)
	}

	return outputImage
}

/*
 * Thresholds an image with 2 thresholds
 * White pixel if it's between the thresholds, otherwise black
 */
func DualThreshold(image [][]float32, thresholdA float32, thresholdB float32) [][]float32 {

	imageWidth, imageHeight := Dimensions(image)

	// Find which threshold is the upper one and which is the lower one
	upperThreshold, lowerThreshold := float32(thresholdA), float32(thresholdB)
	if thresholdA < thresholdB {
		upperThreshold, lowerThreshold = float32(thresholdB), float32(thresholdA)
	}

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

				// Apply thresholds to each pixel
				currentPixel := image[i][j]
				if currentPixel < upperThreshold && currentPixel > lowerThreshold {
					outputImage[i][j] = 1
				} else {
					outputImage[i][j] = 0
				}
			}
		} (j)
	}

	return outputImage
}

func HysteresisThreshold(image [][]float32, thresholdA float32, thresholdB float32) [][]float32 {

	// Find which threshold is the upper one and which is the lower one
	upperThreshold, lowerThreshold := float64(thresholdA), float64(thresholdB)
	if thresholdA < thresholdB {
		upperThreshold, lowerThreshold = float64(thresholdB), float64(thresholdA)
	}

	//Delete this line
	upperThreshold += lowerThreshold

	return image
}