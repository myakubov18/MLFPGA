package main

import (
	"fmt"
	"image"
	"image/png"
	"math/rand"
	"encoding/csv"
	"os"
	"strconv"
	"io"
	"bufio"
	//"github.com/gonum/stat"

	//"gonum.org/v1/gonum/mat"
	//"gonum.org/v1/gonum/stat/distuv"
	//"gonum.org/v1/gonum/stat"
)

// Network is a neural network with 3 layers
type Network struct {
	inputs        int;
	hiddens       int;
	outputs       int;
	hiddenWeights *Matrix;
	outputWeights *Matrix;
	learningRate  int64;
	scalingFactor int64;
}

// CreateNetwork creates a neural network with random weights
func CreateNetwork(input, hidden, output int, rate int64, scale int64) (net Network) {
	net = Network{
		inputs:       input,
		hiddens:      hidden,
		outputs:      output,
		learningRate: rate,
		scalingFactor: scale,
	};
	net.hiddenWeights = NewMatrix(net.hiddens, net.inputs, randomArray((net.inputs * net.hiddens), int64(net.inputs)));
	//fmt.Println("Initial Hidden Weights: ", net.hiddenWeights);
	//fmt.Println("\n\n");
	net.outputWeights = NewMatrix(net.outputs, net.hiddens, randomArray((net.hiddens * net.outputs), int64(net.hiddens)));
	//fmt.Println("Initial Output Weights: ", net.outputWeights);
	//fmt.Println("\n\n");
	return;
}

// Train the neural network
func (net *Network) Train(inputData []int64, targetData []int64) {
	// feedforward
	var inputs *Matrix;
	var hiddenInputs *Matrix;
	var hiddenOutputs *Matrix;
	var finalInputs *Matrix;
	var finalOutputs *Matrix;
	var outputErrors *Matrix;
	var hiddenErrors *Matrix;

	inputs = NewMatrix(len(inputData), 1, inputData);
	//fmt.Println("\nInputs: ", inputs);

	//fmt.Println("\nhiddenWeights: ", net.hiddenWeights, "\n\n\n");
	hiddenInputs = dot(net.hiddenWeights, inputs);
	//fmt.Println("\nHidden Inputs: ", hiddenInputs);

	hiddenOutputs = apply(sigmoid, hiddenInputs);
	//fmt.Println("\nHidden Outputs: ", hiddenOutputs);

	finalInputs = dot(net.outputWeights, hiddenOutputs);
	//fmt.Println("\nfinalInputs: ", finalInputs);

	finalOutputs = apply(sigmoid, finalInputs);
	fmt.Println("\nfinalOutputs: ", finalOutputs);
	// find errors
	//targets := *Matrix(len(targetData), 1, targetData);
	outputErrors = subtract(NewMatrix(len(targetData), 1, targetData), finalOutputs);
	hiddenErrors = dot(net.outputWeights.T(), outputErrors);

	//fmt.Println("\ntargetData: ", targetData);
	/*best := 0;
	var highest int64 = 0;
	for i := 0; i < finalOutputs; i++ {
			//fmt.Println("%T\n", outputs);
			//fmt.Println(type());
		if outputs.At(i, 0) > highest {
			best = i;
			highest = outputs.At(i, 0);
		}
	}
	target, _ := strconv.Atoi(record[0]);
	fmt.Println("Predicted: ", best, 	"... Target: ", target);*/

	/*for i:=0; i< net.outputWeights.row; i++{
		for j:= 0; j<net.outputWeights.col; j++{
			if (net.outputWeights.At(i,j)<=0){
				fmt.Println(net.outputWeights.At(i,j));
			}
		}
	}*/
	//fmt.Println("Output Weights: ", net.outputWeights , "\n");

	// backpropagate
	net.outputWeights = add(net.outputWeights,
		scale(net.learningRate,
			dot(multiply(outputErrors, sigmoidPrime(finalOutputs)),
				hiddenOutputs.T())));

	//outputErrors*w1*x

	net.hiddenWeights = add(net.hiddenWeights,
		scale(net.learningRate,
			dot(multiply(hiddenErrors, sigmoidPrime(hiddenOutputs)),
				inputs.T())));
}

// Predict uses the neural network to predict the value given input data
func (net Network) Predict(inputData []int64) *Matrix {
	// feedforward
	var inputs *Matrix;
	var hiddenInputs *Matrix;
	var hiddenOutputs *Matrix;
	var finalInputs *Matrix;
	var finalOutputs *Matrix;

	inputs = NewMatrix(len(inputData), 1, inputData);
	//biasedInputs := addBiasNodeTo(inputs, 1);
	//fmt.Println("Inputs: ", inputs);
	hiddenInputs = dot(net.hiddenWeights, inputs);
	//fmt.Println("hiddenInputs: ", hiddenInputs);
	hiddenOutputs = apply(sigmoid, hiddenInputs);
	//fmt.Println("hiddenOutputs: ", hiddenOutputs);
	finalInputs = dot(net.outputWeights, hiddenOutputs);
	//fmt.Println("finalInputs: ", finalInputs);
	finalOutputs = apply(sigmoid, finalInputs);
	//fmt.Println("finalOutputs: ", finalOutputs);
	return finalOutputs;
}

//REPLACED SIGMOID WITH RELU THIS IS ACTUALLY RELU WE WERE JUST LAZY


/*func sigmoid(r, c int, z int64) int64{
    if z > 0 {
        return z;
    }else{
        return 0;
    } //simple ReLU activation function
}

func relu2(r, c int, z int64) int64{
    if z > 0 {
        return 1;
    }else{
        return 0;
    }
}

func sigmoidPrime(m *Matrix) *Matrix{
    //x := apply(relu2, m);
    return apply(relu2, m);
}*/

//THIS IS THE ACTUAL SIGMOID BELOW

//Fixed point Sigmoid
func sigmoid(r, c int, z int64) int64{
	var num int64 = 1 << 32;
	var denom int64 = (1 << 32) + exp(MultiplyFixed(-1 << 32, z));
	return DivideFixed(num, denom);
}
//Fixed point Sigmoid'
func sigmoidPrime(m *Matrix) *Matrix {
	rows, _ := m.Dims();
	o := make([]int64, rows);
	for i := range o {
		o[i] = 1;
	}
	ones := NewMatrix(rows,1,o);
	return multiply(m,subtract(ones,m)) //m * (1-m) 
}

//Regular Sigmoid
/*func sigmoid(r, c int, z int64) int64 {
	return 1 / (1 + math.Exp(-1*z)
}

//Regular Sigmoid'
func sigmoidPrime(m *Matrix) *Matrix {
	rows, _ := m.Dims()
	o := make([]int64, rows)
	for i := range o {
		o[i] = 1
	}
    //make an r x 1 *matrix of 1's for the purpose of subtracting
	ones := mat.NewDense(rows, 1, o)
	return multiply(m, subtract(ones, m)) // m * (1 - m)
}*/


//
// Helper functions to allow easier use of Gonum
//

/*func batchNorm(m *Matrix) *Matrix{
	bias := 0.01;
	avg, stDev := getStats(m);
	n := addScalar(-1*avg, m);
	n = scale(1/(math.Sqrt(stDev*stDev + bias)), n);
	return n;
}

func getStats(m *Matrix) (avg, stDev float64){
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
}*/

func dot(m, n *Matrix) *Matrix {
	o := Product(m, n);
	return o;
}

/*func negate (r, c int, v float64) float64{
	return -1*v
}*/

func apply(fn func(i, j int, v int64) int64, m *Matrix) *Matrix {
	r, c := m.Dims();
	o := NewMatrix(r, c, nil);
	o.Apply(fn, m);
	return o;
}

func scale(s int64, m *Matrix) *Matrix {
	//r, c := m.Dims();
	//o := NewMatrix(r, c, nil);
	m.Scale(s);
	return m;
}

func multiply(m, n *Matrix) *Matrix {
	r, c := m.Dims();
	o := NewMatrix(r, c, nil);
	o.MulElem(m, n);
	return o;
}

func add(m, n *Matrix) *Matrix {
	r, c := m.Dims();
	o := NewMatrix(r, c, nil);
	o.Add(m, n);
	return o;
}

func addScalar(i int64, m *Matrix) *Matrix {
	r, c := m.Dims();
	a := make([]int64, r*c);
	for x := 0; x < r*c; x++ {
		a[x] = i;
	}
	n := NewMatrix(r, c, a);
	return add(m, n);
}

func subtract(m, n *Matrix) *Matrix {
	r, c := m.Dims();
	o := NewMatrix(r, c, nil);
	o.Sub(m, n);
	return o;
}


// randomly generate a float64 array
func randomArray(size int, v int64) (data []int64) {
	/*dist := distuv.Uniform{
		Min: 0,
		Max: 1000000,
	};*/

	data = make([]int64, size);
	for i := 0; i < size; i++ {
		// data[i] = rand.NormFloat64() * math.Pow(v, -0.5)
		data[i] = rand.Int63n(0x100000000);
	}
	return;
}

func addBiasNodeTo(m *Matrix, b int64) *Matrix {
	r, _ := m.Dims();
	a := NewMatrix(r+1, 1, nil);

	a.Set(0, 0, b);
	for i := 0; i < r; i++ {
		a.Set(i+1, 0, m.At(i, 0));
	}
	return a;
}

// pretty print a Gonum *matrix
func matrixPrint(X *Matrix) {
	//fa := mat.Formatted(X, mat.Prefix(""), mat.Squeeze());
	//fmt.Printf("%v\n", fa);
}

/*func save(net Network) {
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
}*/

func save(net Network){
	hidden, err := os.Create("data/hweights.csv");
	defer hidden.Close();
	if err == nil {
		writer :=  csv.NewWriter(hidden)
		defer writer.Flush()
		for x := range net.hiddenWeights.data{
			val_string := make([]string, len(net.hiddenWeights.data[x]))
			for i, value := range net.hiddenWeights.data[x]{
				val_string[i] = strconv.FormatInt(value,10)
			}
			writer.Write(val_string)
		}
	}

	out, err := os.Create("data/oweights.csv");
	defer out.Close();
	if err == nil {
		writer :=  csv.NewWriter(out)
		defer writer.Flush()
		for x := range net.outputWeights.data{
			val_string := make([]string, len(net.outputWeights.data[x]))
			for i, value := range net.outputWeights.data[x]{
				val_string[i] = strconv.FormatInt(value,10)
			}
			writer.Write(val_string)	
		}
	}
}

// load a neural network from file
/*func load(net *Network) {
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
}*/

func load(net *Network){
	testFile, _ := os.Open("data/hweights.csv");
	r := csv.NewReader(bufio.NewReader(testFile));
	hweights := make([]int64, net.hiddens*net.inputs);
	rowCount := 0;
	for {
			record, err := r.Read();
			//fmt.Println(record);
			if err == io.EOF {
				break;
			}

			for i := range record {
				//hweights[i] = make([]int, net.inputs)
				//fmt.Println(i);
				hweights[rowCount*net.inputs + i], _ = strconv.ParseInt(record[i],10,64);
			}
			rowCount = rowCount+1;
		}
		testFile.Close();
	//fmt.println(hweights);
	//fmt.Println(len(hweights));
	net.hiddenWeights = NewMatrix(net.hiddens, net.inputs, hweights);

	testFile, _ = os.Open("data/oweights.csv");
	r = csv.NewReader(bufio.NewReader(testFile));
	oweights := make([]int64, net.outputs*net.hiddens);
	rowCount = 0;
	for {
			record, err := r.Read();
			if err == io.EOF {
				break;
			}

			for i := range record {
				oweights[rowCount*net.hiddens + i], _ = strconv.ParseInt(record[i],10,64);
			}
			rowCount = rowCount+1;
		}
		testFile.Close();

	net.outputWeights = NewMatrix(net.outputs, net.hiddens, oweights);
}

// predict a number from an image
// image should be 28 x 28 PNG file
func predictFromImage(net Network, path string) int {
	input := dataFromImage(path);
	output := net.Predict(input);
	matrixPrint(output);
	var best int = 0;
	var highest int64= 0;
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
func dataFromImage(filePath string) (pixels []int64) {
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
	pixels = make([]int64, len(gray.Pix));
	// populate the pixel array subtract Pix from 255 because that's how
	// the MNIST database was trained (in reverse)
	for i := 0; i < len(gray.Pix); i++ {
		pixels[i] = (int64(255-gray.Pix[i]));
	}
	return;
}
