package ImageTools

import (
	"math"
	"sync"
)

/*
 * Calculates the pixelwise sine of an image
 * The output is normalised in the range 0-1 (inclusive), while preserving dynamic range
 */
func Sin(image [][]float32) [][]float32 {

	imageWidth, imageHeight := Dimensions(image)

	// Create output image
	outputImage := make([][]float32, imageWidth)
	for j := range outputImage {
		outputImage[j] = make([]float32, imageHeight)
	}

	// Create wait group
	var waitGroup sync.WaitGroup
	waitGroup.Add(imageHeight)

	// Iterate over columns
	for j := 0; j < imageHeight; j++ {

		// Process each row on its own goroutine
		go func(j int) {
			defer waitGroup.Done()

			// Iterate over row
			for i := 0; i < imageWidth; i++ {

				// Calculate new pixel value
				outputImage[i][j] = float32(math.Sin(math.Abs(float64(image[i][j]))))
			}
		} (j)
	}

	// Wait for all goroutines to finish
	waitGroup.Wait()

	return Normalise(outputImage)
}

/*
 * Calculates the pixelwise arcsine of an image
 * The output is normalised in the range 0-1 (inclusive), while preserving dynamic range
 */
func Asin(image [][]float32) [][]float32 {

	imageWidth, imageHeight := Dimensions(image)

	// Create output image
	outputImage := make([][]float32, imageWidth)
	for j := range outputImage {
		outputImage[j] = make([]float32, imageHeight)
	}

	// Create wait group
	var waitGroup sync.WaitGroup
	waitGroup.Add(imageHeight)

	// Iterate over columns
	for j := 0; j < imageHeight; j++ {

		// Process each row on its own goroutine
		go func(j int) {
			defer waitGroup.Done()

			// Iterate over row
			for i := 0; i < imageWidth; i++ {

				// Calculate new pixel value
				outputImage[i][j] = float32(math.Asin(math.Abs(float64(image[i][j]))))
			}
		} (j)
	}

	// Wait for all goroutines to finish
	waitGroup.Wait()

	return Normalise(outputImage)
}

/*
 * Calculates the pixelwise cosine of an image
 * The output is normalised in the range 0-1 (inclusive), while preserving dynamic range
 */
func Cos(image [][]float32) [][]float32 {

	imageWidth, imageHeight := Dimensions(image)

	// Create output image
	outputImage := make([][]float32, imageWidth)
	for j := range outputImage {
		outputImage[j] = make([]float32, imageHeight)
	}

	// Create wait group
	var waitGroup sync.WaitGroup
	waitGroup.Add(imageHeight)

	// Iterate over columns
	for j := 0; j < imageHeight; j++ {

		// Process each row on its own goroutine
		go func(j int) {
			defer waitGroup.Done()

			// Iterate over row
			for i := 0; i < imageWidth; i++ {

				// Calculate new pixel value
				outputImage[i][j] = float32(math.Cos(math.Abs(float64(image[i][j]))))
			}
		} (j)
	}

	// Wait for all goroutines to finish
	waitGroup.Wait()

	return Normalise(outputImage)
}

/*
 * Calculates the pixelwise arccosine of an image
 * The output is normalised in the range 0-1 (inclusive), while preserving dynamic range
 */
func Acos(image [][]float32) [][]float32 {

	imageWidth, imageHeight := Dimensions(image)

	// Create output image
	outputImage := make([][]float32, imageWidth)
	for j := range outputImage {
		outputImage[j] = make([]float32, imageHeight)
	}

	// Create wait group
	var waitGroup sync.WaitGroup
	waitGroup.Add(imageHeight)

	// Iterate over columns
	for j := 0; j < imageHeight; j++ {

		// Process each row on its own goroutine
		go func(j int) {
			defer waitGroup.Done()

			// Iterate over row
			for i := 0; i < imageWidth; i++ {

				// Calculate new pixel value
				outputImage[i][j] = float32(math.Acos(math.Abs(float64(image[i][j]))))
			}
		} (j)
	}

	// Wait for all goroutines to finish
	waitGroup.Wait()

	return Normalise(outputImage)
}

/*
 * Calculates the pixelwise tangent of an image
 * The output is normalised in the range 0-1 (inclusive), while preserving dynamic range
 */
func Tan(image [][]float32) [][]float32 {

	imageWidth, imageHeight := Dimensions(image)

	// Create output image
	outputImage := make([][]float32, imageWidth)
	for j := range outputImage {
		outputImage[j] = make([]float32, imageHeight)
	}

	halfPi := math.Pi * 0.5

	// Create wait group
	var waitGroup sync.WaitGroup
	waitGroup.Add(imageHeight)

	// Iterate over columns
	for j := 0; j < imageHeight; j++ {

		// Process each row on its own goroutine
		go func(j int) {
			defer waitGroup.Done()

			// Iterate over row
			for i := 0; i < imageWidth; i++ {

				// Calculate new pixel value (obvs tan(x) is undefined at x=0.5*pi, so there's a special case for that)
				if float64(image[i][j]) == halfPi {
					outputImage[i][j] = 0
				} else {
					outputImage[i][j] = float32(math.Tan(math.Abs(float64(image[i][j]))))
				}
			}
		} (j)
	}

	// Wait for all goroutines to finish
	waitGroup.Wait()

	return Normalise(outputImage)
}

/*
 * Calculates the pixelwise arctangent of an image
 * The output is normalised in the range 0-1 (inclusive), while preserving dynamic range
 */
func Atan(image [][]float32) [][]float32 {

	imageWidth, imageHeight := Dimensions(image)

	// Create output image
	outputImage := make([][]float32, imageWidth)
	for j := range outputImage {
		outputImage[j] = make([]float32, imageHeight)
	}

	// Create wait group
	var waitGroup sync.WaitGroup
	waitGroup.Add(imageHeight)

	// Iterate over columns
	for j := 0; j < imageHeight; j++ {

		// Process each row on its own goroutine
		go func(j int) {
			defer waitGroup.Done()

			// Iterate over row
			for i := 0; i < imageWidth; i++ {

				// Calculate new pixel value
				outputImage[i][j] = float32(math.Atan(math.Abs(float64(image[i][j]))))
			}
		} (j)
	}

	// Wait for all goroutines to finish
	waitGroup.Wait()

	return Normalise(outputImage)
}