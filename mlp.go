package main

import (
	"fmt"
	"image"
	"image/png"
	"math"
	"os"
	//"github.com/gonum/stat"

	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat/distuv"
	//"gonum.org/v1/gonum/stat"
)

// Network is a neural network with 3 layers
type Network struct {
	inputs        int;
	hiddens       int;
	outputs       int;
	hiddenWeights *mat.Dense;
	outputWeights *mat.Dense;
	learningRate  float64;
}

// CreateNetwork creates a neural network with random weights
func CreateNetwork(input, hidden, output int, rate float64) (net Network) {
	net = Network{
		inputs:       input,
		hiddens:      hidden,
		outputs:      output,
		learningRate: rate,
	};
	net.hiddenWeights = mat.NewDense(net.hiddens, net.inputs, randomArray((net.inputs)*(net.hiddens), float64(net.inputs)));
	//fmt.Println("Initial Hidden Weights: ", net.hiddenWeights);
	//fmt.Println("\n\n");
	net.outputWeights = mat.NewDense(net.outputs, net.hiddens, randomArray((net.hiddens)*(net.outputs), float64(net.hiddens)));
	//fmt.Println("Initial Output Weights: ", net.outputWeights);
	//fmt.Println("\n\n");
	return;
}

// Train the neural network
func (net *Network) Train(inputData []float64, targetData []float64) {
	// feedforward
	var inputs *mat.Dense;
	var hiddenInputs mat.Matrix;
	var hiddenOutputs mat.Matrix;
	var finalInputs mat.Matrix;
	var finalOutputs mat.Matrix;
	var outputErrors mat.Matrix;
	var hiddenErrors mat.Matrix;

	inputs = mat.NewDense(len(inputData), 1, inputData);

	hiddenInputs = scale(0.1, dot(net.hiddenWeights, inputs));
	hiddenOutputs = apply(sigmoid, hiddenInputs);

	finalInputs = scale(1, dot(net.outputWeights, hiddenOutputs));
	finalOutputs = apply(sigmoid, finalInputs);

	// find errors
	//targets := mat.NewDense(len(targetData), 1, targetData);
	outputErrors = subtract(mat.NewDense(len(targetData), 1, targetData), finalOutputs);
	hiddenErrors = dot(net.outputWeights.T(), outputErrors);

	// backpropagate
	net.outputWeights = add(net.outputWeights,
		scale(net.learningRate,
			dot(multiply(outputErrors, sigmoidPrime(finalOutputs)),
				hiddenOutputs.T()))).(*mat.Dense);

	net.hiddenWeights = add(net.hiddenWeights,
		scale(net.learningRate,
			dot(multiply(hiddenErrors, sigmoidPrime(hiddenOutputs)),
				inputs.T()))).(*mat.Dense);
}

// Predict uses the neural network to predict the value given input data
func (net Network) Predict(inputData []float64) mat.Matrix {
	// feedforward
	var inputs *mat.Dense;
	var hiddenInputs mat.Matrix;
	var hiddenOutputs mat.Matrix;
	var finalInputs mat.Matrix;
	var finalOutputs mat.Matrix;

	inputs = mat.NewDense(len(inputData), 1, inputData);
	//biasedInputs := addBiasNodeTo(inputs, 1);
	//fmt.Println("Inputs: ", inputs);
	hiddenInputs = scale(0.1, dot(net.hiddenWeights, inputs));
	//fmt.Println("hiddenInputs: ", hiddenInputs);
	hiddenOutputs = apply(sigmoid, hiddenInputs);
	//fmt.Println("hiddenOutputs: ", hiddenOutputs);
	finalInputs = scale(1, dot(net.outputWeights, hiddenOutputs));
	//fmt.Println("finalInputs: ", finalInputs);
	finalOutputs = apply(sigmoid, finalInputs);
	//fmt.Println("finalOutputs: ", finalOutputs);
	return finalOutputs;
}

//REPLACED SIGMOID WITH RELU THIS IS ACTUALLY RELU WE WERE JUST LAZY


func sigmoid(r, c int, z float64) float64{
    return math.Max(z, 0); //simple ReLU activation function
}

func relu2(r, c int, z float64) float64{
    if z > 0 {
        return 1;
    }else{
        return 0;
    }
}

func sigmoidPrime(m mat.Matrix) mat.Matrix{
    //x := apply(relu2, m);
    return apply(relu2, m);
}

//THIS IS THE ACTUAL SIGMOID BELOW

/*
func sigmoid(r, c int, z float64) float64 {
	return 1.0 / (1 + math.Exp(-1*z))
}

func sigmoidPrime(m mat.Matrix) mat.Matrix {
	rows, _ := m.Dims()
	o := make([]float64, rows)
	for i := range o {
		o[i] = 1
	}
    //make an r x 1 matrix of 1's for the purpose of subtracting
	ones := mat.NewDense(rows, 1, o)
	return multiply(m, subtract(ones, m)) // m * (1 - m)
}*/

//
// Helper functions to allow easier use of Gonum
//

func batchNorm(m mat.Matrix) mat.Matrix{
	bias := 0.01;
	avg, stDev := getStats(m);
	n := addScalar(-1*avg, m);
	n = scale(1/(math.Sqrt(stDev*stDev + bias)), n);
	return n;
}

func getStats(m mat.Matrix) (avg, stDev float64){
	r, c := m.Dims();
	avg = mat.Sum(m)/float64(r*c);
	data := make([]float64, r*c);
	for i:= range data{
		data[i] = avg;
	}
	n := mat.NewDense(r, c, data);
	//n = subtract(m, n);
	stDev = mat.Sum(subtract(m, n))/float64(r*c);
	return avg, stDev;
}

func dot(m, n mat.Matrix) mat.Matrix {
	r, _ := m.Dims();
	_, c := n.Dims();
	o := mat.NewDense(r, c, nil);
	o.Product(m, n);
	return o;
}

/*func negate (r, c int, v float64) float64{
	return -1*v
}*/

func apply(fn func(i, j int, v float64) float64, m mat.Matrix) mat.Matrix {
	r, c := m.Dims();
	o := mat.NewDense(r, c, nil);
	o.Apply(fn, m);
	return o;
}

func scale(s float64, m mat.Matrix) mat.Matrix {
	r, c := m.Dims();
	o := mat.NewDense(r, c, nil);
	o.Scale(s, m);
	return o;
}

func multiply(m, n mat.Matrix) mat.Matrix {
	r, c := m.Dims();
	o := mat.NewDense(r, c, nil);
	o.MulElem(m, n);
	return o;
}

func add(m, n mat.Matrix) mat.Matrix {
	r, c := m.Dims();
	o := mat.NewDense(r, c, nil);
	o.Add(m, n);
	return o;
}

func addScalar(i float64, m mat.Matrix) mat.Matrix {
	r, c := m.Dims();
	a := make([]float64, r*c);
	for x := 0; x < r*c; x++ {
		a[x] = i;
	}
	n := mat.NewDense(r, c, a);
	return add(m, n);
}

func subtract(m, n mat.Matrix) mat.Matrix {
	r, c := m.Dims();
	o := mat.NewDense(r, c, nil);
	o.Sub(m, n);
	return o;
}

// randomly generate a float64 array
func randomArray(size int, v float64) (data []float64) {
	dist := distuv.Uniform{
		Min: 0 / math.Sqrt(v),
		Max: 1 / math.Sqrt(v),
	};

	data = make([]float64, size);
	for i := 0; i < size; i++ {
		// data[i] = rand.NormFloat64() * math.Pow(v, -0.5)
		data[i] = dist.Rand();
	}
	return;
}

func addBiasNodeTo(m mat.Matrix, b float64) mat.Matrix {
	r, _ := m.Dims();
	a := mat.NewDense(r+1, 1, nil);

	a.Set(0, 0, b);
	for i := 0; i < r; i++ {
		a.Set(i+1, 0, m.At(i, 0));
	}
	return a;
}

// pretty print a Gonum matrix
func matrixPrint(X mat.Matrix) {
	fa := mat.Formatted(X, mat.Prefix(""), mat.Squeeze());
	fmt.Printf("%v\n", fa);
}

func save(net Network) {
	h, err := os.Create("data/hweights.model");
	defer h.Close();
	if err == nil {
		net.hiddenWeights.MarshalBinaryTo(h);
	}
	o, err := os.Create("data/oweights.model");
	defer o.Close();
	if err == nil {
		net.outputWeights.MarshalBinaryTo(o);
	}
}

// load a neural network from file
func load(net *Network) {
	h, err := os.Open("data/hweights.model");
	defer h.Close();
	if err == nil {
		net.hiddenWeights.Reset();
		net.hiddenWeights.UnmarshalBinaryFrom(h);
	}
	o, err := os.Open("data/oweights.model");
	defer o.Close();
	if err == nil {
		net.outputWeights.Reset();
		net.outputWeights.UnmarshalBinaryFrom(o);
	}
	return;
}

// predict a number from an image
// image should be 28 x 28 PNG file
func predictFromImage(net Network, path string) int {
	input := dataFromImage(path);
	output := net.Predict(input);
	matrixPrint(output);
	best := 0;
	highest := 0.0;
	for i := 0; i < net.outputs; i++ {
		if output.At(i, 0) > highest {
			best = i;
			highest = output.At(i, 0);
		}
	}
	fmt.Println("Predicted: ", best);
	return best;
}

// get the pixel data from an image
func dataFromImage(filePath string) (pixels []float64) {
	// read the file
	imgFile, err := os.Open(filePath);
	defer imgFile.Close();
	if err != nil {
		fmt.Println("Cannot read file:", err);
	}
	img, err := png.Decode(imgFile);
	if err != nil {
		fmt.Println("Cannot decode file:", err);
	}

	// create a grayscale image
	bounds := img.Bounds();
	gray := image.NewGray(bounds);

	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			var rgba = img.At(x, y);
			gray.Set(x, y, rgba);
		}
	}
	// make a pixel array
	pixels = make([]float64, len(gray.Pix));
	// populate the pixel array subtract Pix from 255 because that's how
	// the MNIST database was trained (in reverse)
	for i := 0; i < len(gray.Pix); i++ {
		pixels[i] = (float64(255-gray.Pix[i]) / 255.0 * 0.999) + 0.001;
	}
	return;
}
