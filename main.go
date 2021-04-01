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
)

func main() {
	// 784 inputs - 28 x 28 pixels, each pixel is an input
	// 100 hidden nodes - an arbitrary number
	// 10 outputs - digits 0 to 9
	// 0.1 is the learning rate
	net := CreateNetwork(784, 100, 10, 1000, 1000000);

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
	//fmt.Println("\n\nHidden Weights: ", net.hiddenWeights.At(0,500), "\n\n");
	for epochs := 0; epochs < 1; epochs++ {
		testFile, _ := os.Open("mnist_dataset/mnist_train_1000.csv");
		r := csv.NewReader(bufio.NewReader(testFile));
		for {
			record, err := r.Read();
			if err == io.EOF {
				break;
			}

			inputs := make([]int, net.inputs);
			for i := range inputs {
				inputs[i], _ = strconv.Atoi(record[i]);
				//inputs[i] = (x / 255.0 * 9.99) + 0.01;
				//inputs[i] = x + 1
			}
			//fmt.Println("inputs: ", inputs);

			targets := make([]int, net.outputs);
			for i := range targets {
				//targets[i] = 0.01;
				targets[i] = 1;

			}
			x, _ := strconv.Atoi(record[0]);
			targets[x] = 255;
			//fmt.Println("Hidden Weights: ", net.hiddenWeights.At(0,500), "\n");
			net.Train(inputs, targets);
			//fmt.Println("Hidden Weights: ", net.hiddenWeights.At(0,500), "\n");
		}
		//fmt.Println("Epoch ", epochs, "\n\n", net.hiddenWeights);
		testFile.Close();
	}
	elapsed := time.Since(t1);
	fmt.Printf("\nTime taken to train: %s\n", elapsed);
	//fmt.Println("Weights: ", net.hiddenWeights)
	//fmt.Println("\noutputWeights: ", net.outputWeights)
}

func mnistPredict(net *Network) {
	t1 := time.Now();
	checkFile, _ := os.Open("mnist_dataset/mnist_test.csv");
	//checkFile, _ := os.Open("mnist_dataset/mnist_test.csv");
	defer checkFile.Close();

	score := 0;
	r := csv.NewReader(bufio.NewReader(checkFile));
	for {
		record, err := r.Read();
		if err == io.EOF {
			break;
		}
		inputs := make([]int, net.inputs);
		for i := range inputs {
			if i == 0 {
				inputs[i] = 1;
			}
			inputs[i], _ = strconv.Atoi(record[i]);
			//inputs[i] = (x / 255.0 * 9.99) + 0.01;
		}
		//fmt.Println("inputs: ", inputs);
		outputs := net.Predict(inputs);
		//fmt.Println("outputs: ", outputs)
		best := 0;
		highest := 0;
		for i := 0; i < net.outputs; i++ {
			//fmt.Println("%T\n", outputs);
			//fmt.Println(type());
			if outputs.At(i, 0) > highest {
				best = i;
				highest = outputs.At(i, 0);
			}
		}
		target, _ := strconv.Atoi(record[0]);
		fmt.Println("Predicted: ", best, 	"... Target: ", target);
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
