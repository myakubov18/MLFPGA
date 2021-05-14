package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/csv"
	"flag"
	"fmt"
	"image"
	"image/png"
	"io"
	"math/rand"
	"os"
	"strconv"
	"time"
	//"golang.org/x/exp/rand"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	//"gonum.org/v1/plot/vg/draw"
)

func main() {
	// 784 inputs - 28 x 28 pixels, each pixel is an input
	// 100 hidden nodes - an arbitrary number
	// 10 outputs - digits 0 to 9
	// 1/256 is the learning rate
	net := CreateNetwork(784, 100, 10, 0x02000000, 1);

	mnist := flag.String("mnist", "", "Either train or predict to evaluate neural network");
	file := flag.String("file", "", "File name of 28 x 28 PNG file to evaluate");
	flag.Parse();

	// train or mass predict to determine the effectiveness of the trained network
	switch *mnist {
	case "train":
		mnistTrain(&net);
		save(net);
	case "predict":
		load(&net);
		mnistPredict(&net);
	default:
		// don't do anything
	}

	
	// predict individual digit images
	if *file != "" {
		// print the image out nicely on the terminal
		printImage(getImage(*file));
		// load the neural network from file
		load(&net);
		// predict which number it is
		fmt.Println("prediction:", predictFromImage(net, *file));
	}
}

func mnistTrain(net *Network) {
	rand.Seed(time.Now().UTC().UnixNano());
	t1 := time.Now();
	minWeights := []int64 {};
	maxWeights := []int64 {};
	numWeightsDisplay := net.hiddens;
	firstFewWeights := make([][]int64, numWeightsDisplay);
	//fmt.Println("\n\nHidden Weights: ", net.hiddenWeights.At(0,500), "\n\n");
	sample := 0;
	for epochs := 0; epochs < 1; epochs++ {
		testFile, _ := os.Open("mnist_dataset/mnist_train_1000.csv");
		r := csv.NewReader(bufio.NewReader(testFile));
		for {
			record, err := r.Read();
			if err == io.EOF {
				break;
			}

			inputs := make([]int64, net.inputs);
			for i := range inputs {
				//inputs[i], _ = strconv.Atoi(record[i]);
				//BIT SHIFTED HERE
				inputs[i], _ = strconv.ParseInt(record[i],10,64);
				inputs[i] = inputs[i] << 32;
				//inputs[i] = (x / 255.0 * 9.99) + 0.01;
				//inputs[i] = x + 1
			}
			//fmt.Println("inputs: ", inputs);

			targets := make([]int64, net.outputs);
			for i := range targets {
				//targets[i] = 0.01;
				targets[i] = 1 << 32;

			}
			x, _ := strconv.Atoi(record[0]);
			targets[x] = 255 << 32;
			//fmt.Println("Hidden Weights: ", net.hiddenWeights.At(0,500), "\n");
			//fmt.Println("--------------------------------------");
			//fmt.Println(targets);
			fmt.Println("Data Sample: ", sample)
			sample++;
			net.Train(inputs, targets);
			maxWeights = append(maxWeights, net.hiddenWeights.max);
			minWeights = append(minWeights, net.hiddenWeights.min);
			for i:=0; i<numWeightsDisplay; i++{
				firstFewWeights[i] = append(firstFewWeights[i], net.outputWeights.data[5][i])
			}

			//fmt.Println("Hidden Weights: ", net.hiddenWeights.At(0,500), "\n");
		}
		//fmt.Println("Epoch ", epochs, "\n\n", net.hiddenWeights);
		testFile.Close();
	}
	elapsed := time.Since(t1);

	//fmt.Println("Hidden Weights: ", net.outputWeights, "\n");
	//Sort by MIN/MAX range, plot the ones that go to 0.


	fmt.Printf("\nTime taken to train: %s\n", elapsed);
	//fmt.Println("maxWeights: ", maxWeights);
	//fmt.Println("minWeights: ", minWeights);

	//plotWeights("Weights Graph", "Iterations", "Magnitude", [maxWeights], "Max Weights", "maxWeightsPlot.png");
	//plotWeights("Weights Graph", "Iterations", "Magnitude", [minWeights], "Min Weights", "minWeightsPlot.png");
	maxOut := [][]int64 {maxWeights};
	minOut := [][]int64 {minWeights};
	plotWeights("Weights Graph", "Iterations", "Magnitude", firstFewWeights, "weight", "weightGraphs/");
	plotWeights("Weights Graph", "Iterations", "Magnitude", maxOut, "maxWeight", "weightGraphs/");
	plotWeights("Weights Graph", "Iterations", "Magnitude", minOut, "minWeight", "weightGraphs/");
}

func plotWeights(title string, xlabel string, ylabel string, data [][]int64, dataName string, filename string){
	

	for dataLine:=0; dataLine<len(data); dataLine++{
		newPlot := plot.New();
		newPlot.Title.Text = title;
		newPlot.X.Label.Text = xlabel;
		newPlot.Y.Label.Text = ylabel;
		//newPlot.Y.Scale = plot.LogScale{};
		dataSet := data[dataLine]
		points := make(plotter.XYs, len(dataSet));
		for i:=0; i<len(dataSet); i++ {
			points[i].X = float64(i);
			points[i].Y = float64(dataSet[i]);
		}
		newName := dataName + strconv.Itoa(dataLine);
		err := plotutil.AddLinePoints(newPlot, newName, points);
		if err != nil{
			panic(err);
		}
		if err := newPlot.Save(4*vg.Inch, 4*vg.Inch, filename+dataName+strconv.Itoa(dataLine)+".png"); err != nil {
			panic(err)
		}
	}
}

func mnistPredict(net *Network) {
	t1 := time.Now();
	checkFile, _ := os.Open("mnist_dataset/mnist_test_1000.csv");
	//checkFile, _ := os.Open("mnist_dataset/mnist_test.csv");
	defer checkFile.Close();
	//fmt.Println("\nhiddenWeights: ", net.outputWeights, "\n");
	score := 0;
	r := csv.NewReader(bufio.NewReader(checkFile));
	for {
		record, err := r.Read();
		if err == io.EOF {
			break;
		}
		inputs := make([]int64, net.inputs);
		for i := range inputs {
			if i == 0 {
				inputs[i] = 1 << 32;
			}
			//inputs[i], _ = strconv.Atoi(record[i]);
			//BIT SHIFTED HERE
			inputs[i], _ = strconv.ParseInt(record[i],10,64);
			inputs[i] = inputs[i] << 32;
			//inputs[i] = (x / 255.0 * 9.99) + 0.01;
		}
		//fmt.Println("inputs: ", inputs);
		outputs := net.Predict(inputs);
		//fmt.Println("outputs: ", outputs)
		best := 0;
		var highest int64 = 0;
		fmt.Println(outputs);
		for i := 0; i < net.outputs; i++ {
			//fmt.Println("%T\n", outputs);
			//fmt.Println(type());
			if outputs.At(i, 0) > highest {
				best = i;
				highest = outputs.At(i, 0);
			}
		}
		target, _ := strconv.Atoi(record[0]);
		fmt.Println("Predicted: ", best, "... Target: ", target);
		if best == target {
			//fmt.Println("Predicted: ", best);
			score++;
		}
	}

	elapsed := time.Since(t1);
	fmt.Printf("Time taken to check: %s\n", elapsed);
	fmt.Println("score:", score);
}

// print out image on iTerm2; equivalent to imgcat on iTerm2
func printImage(img image.Image) {
	var buf bytes.Buffer;
	png.Encode(&buf, img);
	imgBase64Str := base64.StdEncoding.EncodeToString(buf.Bytes());
	fmt.Printf("\x1b]1337;File=inline=1:%s\a\n", imgBase64Str);
}

// get the file as an image
func getImage(filePath string) image.Image {
	imgFile, err := os.Open(filePath);
	defer imgFile.Close();
	if err != nil {
		fmt.Println("Cannot read file:", err);
	}
	img, _, err := image.Decode(imgFile);
	if err != nil {
		fmt.Println("Cannot decode file:", err);
	}
	return img;
}
