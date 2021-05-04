package ImageTools

import "math"

/*
 * Finds the dimmest and brightest pixels in an image
 */
func MinMax(image [][]float32) (float32, float32) {

	imageWidth, imageHeight := Dimensions(image)

	// Find brightest and dimmest pixels
	max, min := float64(0), float64(65535)
	for j := 0; j < imageHeight; j++ {
		for i := 0; i < imageWidth; i++ {
			currentPixel := float64(image[i][j])
			if currentPixel > max {
				max = currentPixel
			}
			if currentPixel < min {
				min = currentPixel
			}
		}
	}

	return float32(min), float32(max)
}

/*
 * Returns the width and height of an image
 */
func Dimensions(image [][]float32) (int, int) {
	imageWidth, imageHeight := len(image), len(image[0])
	return imageWidth, imageHeight
}

/*
 * Calculates the mean and (population) standard deviation of an image
 * The two metrics are combined because the mean is needed to calculate the std, so it is more efficient to calculate them both together
 */
func MeanStd(image [][]float32) (float32, float32) {
	imageWidth, imageHeight := Dimensions(image)

	// Sum all pixels
	accumulator := float64(0)
	for j := 0; j < imageHeight; j++ {
		for i := 0; i < imageWidth; i++ {
			currentPixel := float64(image[i][j])
			accumulator += currentPixel
		}
	}

	// Divide by the total number of pixels
	mean := accumulator / (float64(imageWidth) * float64(imageHeight))

	// Sum the square difference of each pixel and the mean
	accumulator = float64(0)
	for j := 0; j < imageHeight; j++ {
		for i := 0; i < imageWidth; i++ {
			currentPixel := float64(image[i][j])
			currentPixel -= mean
			currentPixel *= currentPixel
			accumulator += currentPixel
		}
	}

	// Divide by the total number of pixels
	std := math.Sqrt(accumulator / (float64(imageWidth) * float64(imageHeight)))

	return float32(mean), float32(std)
}

/*
 * Calculates the mean and (population) standard deviation of an image
 * The two metrics are combined because the mean is needed to calculate the std, so it is more efficient to calculate them both together
 * This version doesn't round the result to fit into a float32
 */
func meanStd(image [][]float32) (float64, float64) {
	imageWidth, imageHeight := Dimensions(image)

	// Sum all pixels
	accumulator := float64(0)
	for j := 0; j < imageHeight; j++ {
		for i := 0; i < imageWidth; i++ {
			currentPixel := float64(image[i][j])
			accumulator += currentPixel
		}
	}

	// Divide by the total number of pixels
	mean := accumulator / (float64(imageWidth) * float64(imageHeight))

	// Sum the square difference of each pixel and the mean
	accumulator = float64(0)
	for j := 0; j < imageHeight; j++ {
		for i := 0; i < imageWidth; i++ {
			currentPixel := float64(image[i][j])
			currentPixel -= mean
			currentPixel *= currentPixel
			accumulator += currentPixel
		}
	}

	// Divide by the total number of pixels
	std := math.Sqrt(accumulator / (float64(imageWidth) * float64(imageHeight)))

	return mean, std
}