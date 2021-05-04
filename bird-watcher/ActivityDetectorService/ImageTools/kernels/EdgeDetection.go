package kernels

var SobelX = [][]float32 {
	{-1,  0,  1},
	{-2,  0,  2},
	{-1,  0,  1},
}

var SepSobelXPt1 = [][]float32 {
	{1},
	{2},
	{1},
}

var SepSobelXPt2 = [][]float32 {
	{-1,  0,  1},
}

var SepSobelYPt1 = [][]float32 {
	{1, 2, 1},
}

var SepSobelYPt2 = [][]float32 {
	{-1},
	{0},
	{1},
}
var SobelY = [][]float32 {
	{-1, -2, -1},
	{ 0,  0,  0},
	{ 1,  2,  1},
}

var Laplacian = [][]float32 {
	{-1, -1, -1, -1, -1},
	{-1, -1, -1, -1, -1},
	{-1, -1, 24, -1, -1},
	{-1, -1, -1, -1, -1},
	{-1, -1, -1, -1, -1},
}