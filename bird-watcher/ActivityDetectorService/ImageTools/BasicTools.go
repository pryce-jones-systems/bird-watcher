package ImageTools

import (
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"sync"
)

/*
 * Loads an image from a file
 */
func LoadImage(path string) ([][]float32, error) {

	// Open file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Decode image
	img, err := jpeg.Decode(file)

	// Convert to 2D slice
	return Image2Slice(img), err
}

/*
 * Saves an image to a file
 */
func SaveImage(path string, image [][]float32) error {

	// Convert to image.Image
	img, err := Slice2Image(image)
	if err != nil {
		return err
	}

	// Open file
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Encode image
	options := jpeg.Options{Quality: 100}
	err = jpeg.Encode(file, img, &options)

	// Return nil, if successful
	return nil
}

/*
 * Converts an RGBA image to Gray16
 */
func RGBA2Gray16(img image.Image) image.Image {
	newImg := image.NewGray16(img.Bounds())
	for j := 0; j < img.Bounds().Max.Y; j++ {
		for i := 0; i < img.Bounds().Max.X; i++ {
			newImg.Set(i, j, color.Gray16Model.Convert(img.At(i, j)))
		}
	}
	return newImg
}

/*
 * Takes an Image and converts it to a 2D slice of 32-bit floating points
 * Also converts the image to grayscale and normialises all pixels in the range 0-1 (inclusive)
 */
func Image2Slice(img image.Image) [][]float32 {

	// Convert image to 16-bit grayscale
	img = RGBA2Gray16(img)

	// Find brightest and dimmest pixels
	max, min := float64(0), float64(65535)
	for j := 0; j < img.Bounds().Max.Y; j++ {
		for i := 0; i < img.Bounds().Max.X; i++ {
			c := color.Gray16Model.Convert(img.At(i, j)).(color.Gray16).Y
			currentPixel := float64(c)
			if currentPixel > max {
				max = currentPixel
			}
			if currentPixel < min {
				min = currentPixel
			}
		}
	}

	// Create empty slice
	normalisedImage := make([][]float32, img.Bounds().Max.X)
	for j := range normalisedImage {
		normalisedImage[j] = make([]float32, img.Bounds().Max.Y)
	}

	// Copy image into slice and normalise values in the range 0-1 (inclusive), while preserving dynamic range
	divisor := float64(max - min)
	if min == max {
		divisor = 1
	}
	for j := 0; j < img.Bounds().Max.Y; j++ {
		for i := 0; i < img.Bounds().Max.X; i++ {
			c := color.Gray16Model.Convert(img.At(i, j)).(color.Gray16).Y
			currentPixel := float64(c)
			normalisedPixel := float32((currentPixel - min) / divisor)
			normalisedImage[i][j] = normalisedPixel
		}
	}

	return normalisedImage
}

/*
 * Takes a 2D slice of floating points representing a grayscale image and returns an Image
 * Note that all pixel values are assumed to be normalised in the range 0-1, so they are all multiplied by the maximum pixel value (65536), so the image can be displayed
 */
func Slice2Image(slice [][]float32) (image.Image, error) {

	width, height := Dimensions(slice)

	img := image.NewGray16(image.Rect(0, 0, width, height))
	for j := 0; j < height; j++ {
		for i := 0; i < width; i++ {
			img.Set(i, j, color.Gray16{uint16(slice[i][j] * 65535)})
		}
	}
	return img, nil
}

func Normalise(image [][]float32) [][]float32 {

	min, max := MinMax(image)
	imageWidth, imageHeight := Dimensions(image)
	quotient := float64(max - min)
	if min == max {
		quotient = 1
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

				// Normalise each pixel in the range 0-1 (inclusive), while preserving dynamic range
				currentPixel := float64(image[i][j])
				normalisedPixel := float32((currentPixel - float64(min)) / quotient)
				image[i][j] = normalisedPixel
			}
		} (j)
	}

	// Wait for all goroutines to finish
	waitGroup.Wait()

	return image
}

func Invert(image [][]float32) [][]float32 {

	imageWidth, imageHeight := Dimensions(image)

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

				// Normalise each pixel in the range 0-1 (inclusive), while preserving dynamic range
				currentPixel := float64(image[i][j])

				// Invert each pixel, making sure we don't do something dumb and end up with NaN or -Inf
				if currentPixel == 0 {
					currentPixel = 1
				} else if currentPixel == 1 {
					currentPixel = 0
				} else {
					currentPixel = 1 / currentPixel
				}
				image[i][j] = float32(currentPixel)
			}
		} (j)
	}

	// Wait for all goroutines to finish
	waitGroup.Wait()

	return image
}