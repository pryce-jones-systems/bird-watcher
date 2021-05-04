package ImageTools

import "math"

/*
 * Calculates the Root Mean Square Error (RMSE), Mean Square Error (MSE), Sum Square Error (SSE) of two images
 * These metrics are combined into one function for computational efficiency
 */
func SquareError(a [][]float32, b[][]float32) (float32, float32, float32) {

	// Make sure both images are exactly the same dimensions
	aw, ah := Dimensions(a)
	bw, bh := Dimensions(b)
	if aw != bw || ah != bh {
		return math.MaxFloat32, math.MaxFloat32, math.MaxFloat32
	}
	imageWidth, imageHeight := aw, ah

	// Calculate SSE
	accumulator := float64(0)
	for j := 0; j < imageHeight; j++ {
		for i := 0; i < imageWidth; i++ {
			pixelA, pixelB := float64(a[i][j]), float64(b[i][j])
			error := pixelA - pixelB
			accumulator += error * error
		}
	}
	sse := accumulator

	// Calculate MSE
	mse := sse / (float64(imageWidth) * float64(imageHeight))

	// Calculate RMSE
	rmse := math.Sqrt(mse)

	return float32(rmse), float32(mse), float32(sse)
}

/*
 * Calculates the Root Mean Absolute Error (RMAE), Mean Absolute Error (MAE), Sum Absolute Error (SAE) of two images
 * These metrics are combined into one function for computational efficiency
 */
func AbsoluteError(a [][]float32, b[][]float32) (float32, float32, float32) {

	// Make sure both images are exactly the same dimensions
	aw, ah := Dimensions(a)
	bw, bh := Dimensions(b)
	if aw != bw || ah != bh {
		return math.MaxFloat32, math.MaxFloat32, math.MaxFloat32
	}
	imageWidth, imageHeight := aw, ah

	// Calculate SAE
	accumulator := float64(0)
	for j := 0; j < imageHeight; j++ {
		for i := 0; i < imageWidth; i++ {
			pixelA, pixelB := float64(a[i][j]), float64(b[i][j])
			error := pixelA - pixelB
			accumulator += math.Abs(error)
		}
	}
	sae := accumulator

	// Calculate MAE
	mae := sae / (float64(imageWidth) * float64(imageHeight))

	// Calculate RMAE
	rmae := math.Sqrt(mae)

	return float32(rmae), float32(mae), float32(sae)
}

/*
 * Calculates the Zero-Normalised CrossCorrelation (ZNCC) of two images
 */
func CrossCorrelation(a [][]float32, b [][]float32) float32 {

	// Make sure both images are exactly the same dimensions
	aw, ah := Dimensions(a)
	bw, bh := Dimensions(b)
	if aw != bw || ah != bh {
		return math.MaxFloat32
	}
	imageWidth, imageHeight := aw, ah

	// Calculate means and standard deviations
	meanA, stdA := meanStd(a)
	meanB, stdB := meanStd(b)

	// Calculate ZNCC
	accumulator := float64(0)
	for j := 0; j < imageHeight; j++ {
		for i := 0; i < imageWidth; i++ {
			pixelA, pixelB := float64(a[i][j]), float64(b[i][j])
			accumulator += (pixelA - meanA) * (pixelB - meanB)
		}
	}
	zncc := (accumulator / (stdA * stdB)) / (float64(imageWidth) * float64(imageHeight))

	return float32(zncc)
}