package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"time"
	//"golang.org/x/exp/rand"
	//"gonum.org/v1/plot"
	//"gonum.org/v1/plot/plotter"
	//"gonum.org/v1/plot/plotutil"
	//"gonum.org/v1/plot/vg"
	//"gonum.org/v1/plot/vg/draw"
)

func main() {
	// 784 inputs - 28 x 28 pixels, each pixel is an input
	// 100 hidden nodes - an arbitrary number
	// 10 outputs - digits 0 to 9
	// 1/256 is the learning rate
	net := CreateNetwork(784, 100, 10, 0x020000000000);
	
    mnist := flag.String("mnist", "", "Either train or predict to evaluate neural network");
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

}

func mnistTrain(net *Network) {
	rand.Seed(time.Now().UTC().UnixNano());
	t1 := time.Now();
	numWeightsDisplay := net.hiddens;
	firstFewWeights := make([][]int64, numWeightsDisplay);
	//fmt.Println("\n\nHidden Weights: ", net.hiddenWeights.At(0,500), "\n\n");
	sample := 0;
	for epochs := 0; epochs < 20; epochs++ {
		fmt.Printf("Epoch %d\n", epochs)
        testFile, _ := os.Open("mnist_dataset/mnist_train.csv");
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
				inputs[i] = inputs[i] << 40;
				//inputs[i] = (x / 255.0 * 9.99) + 0.01;
				//inputs[i] = x + 1
			}
			//fmt.Println("inputs: ", inputs);

			targets := make([]int64, net.outputs);
			for i := range targets {
				//targets[i] = 0.01;
				targets[i] = 1 << 40;

			}
			x, _ := strconv.Atoi(record[0]);
			targets[x] = 255 << 40;
			//fmt.Println("Hidden Weights: ", net.hiddenWeights.At(0,500), "\n");
			//fmt.Println("--------------------------------------");
			//fmt.Println(targets);
			//fmt.Println("Data Sample: ", sample)
			sample++;
			net.Train(inputs, targets);
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
	//maxOut := [][]int64 {maxWeights};
	//minOut := [][]int64 {minWeights};
	//plotWeights("Weights Graph", "Iterations", "Magnitude", firstFewWeights, "weight", "weightGraphs/");
	//plotWeights("Weights Graph", "Iterations", "Magnitude", maxOut, "maxWeight", "weightGraphs/");
	//plotWeights("Weights Graph", "Iterations", "Magnitude", minOut, "minWeight", "weightGraphs/");
}

func mnistPredict(net *Network) {
	t1 := time.Now();
	checkFile, _ := os.Open("mnist_dataset/mnist_test.csv");
	//checkFile, _ := os.Open("mnist_dataset/mnist_test.csv");
	defer checkFile.Close();
	//fmt.Println("\nhiddenWeights: ", net.outputWeights, "\n");
	score := 0;
	samples := 0;
	r := csv.NewReader(bufio.NewReader(checkFile));
	for {
		record, err := r.Read();
		if err == io.EOF {
			break;
		}
		inputs := make([]int64, net.inputs);
		for i := range inputs {
			if i == 0 {
				inputs[i] = 1 << 40;
			}
			//inputs[i], _ = strconv.Atoi(record[i]);
			//BIT SHIFTED HERE
			inputs[i], _ = strconv.ParseInt(record[i],10,64);
			inputs[i] = inputs[i] << 40;
			//inputs[i] = (x / 255.0 * 9.99) + 0.01;
		}
		//fmt.Println("inputs: ", inputs);
		outputs := net.Predict(inputs);
		//fmt.Println("outputs: ", outputs)
		best := 0;
		var highest int64 = 0;
		//fmt.Println(outputs);
		for i := 0; i < net.outputs; i++ {
			//fmt.Println("%T\n", outputs);
			//fmt.Println(type());
			if outputs.At(i, 0) > highest {
				best = i;
				highest = outputs.At(i, 0);
			}
		}
		target, _ := strconv.Atoi(record[0]);
		//fmt.Println("Predicted: ", best, "... Target: ", target);
		if best == target {
			//fmt.Println("Predicted: ", best);
			score++;
		}
		samples++;
	}

	elapsed := time.Since(t1);
	fmt.Printf("Time taken to check: %s\n", elapsed);
	fmt.Printf("score: %.2f%%\n", float64(score)*100/float64(samples));
}

