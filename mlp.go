package main

import (
    "math/rand"
    "os"
    "strconv"
    "io"
    "bufio"
    "encoding/csv"
)

type Network struct {
	inputs        int;
	hiddens       int;
	outputs       int;
	hiddenWeights *Matrix;
	outputWeights *Matrix;
	learningRate  int64;
}

func CreateNetwork(input, hidden, output int, rate int64) (net Network) {
	net = Network{
		inputs:       input,
		hiddens:      hidden,
		outputs:      output,
		learningRate: rate,
	};
	net.hiddenWeights = NewMatrix(net.hiddens, net.inputs, randomArray((net.inputs * net.hiddens), int64(net.inputs)));
	net.outputWeights = NewMatrix(net.outputs, net.hiddens, randomArray((net.hiddens * net.outputs), int64(net.hiddens)));
	return;
}

func (net *Network) Train(inputData []int64, targetData []int64) {
    inputs := NewMatrix(len(inputData), 1, inputData);

    hiddenInputs := net.hiddenWeights.Product(inputs);
    hiddenOutputs := hiddenInputs.Apply(regionPW);

    finalInputs := net.outputWeights.Product(hiddenOutputs);
    finalOutputs := finalInputs.Apply(regionPW);

    outputErrors := NewMatrix(len(targetData),1,targetData).Sub(finalOutputs);
    hiddenErrors := net.outputWeights.T().Product(outputErrors);

    net.outputWeights = net.outputWeights.Add(outputErrors.MulElem(finalOutputs.Apply(regionPW2)).Product(hiddenOutputs.T()).Scale(net.learningRate))
    net.hiddenWeights = net.hiddenWeights.Add(hiddenErrors.MulElem(hiddenOutputs.Apply(regionPW2)).Product(inputs.T()).Scale(net.learningRate))
}

func (net *Network) Predict(inputData []int64) *Matrix {
    inputs := NewMatrix(len(inputData),1,inputData)

    hiddenInputs := net.hiddenWeights.Product(inputs)

    hiddenOutputs := hiddenInputs.Apply(regionPW)

    finalInputs := net.outputWeights.Product(hiddenOutputs)

    finalOutputs := finalInputs.Apply(regionPW)

    return finalOutputs
}

func linearPW(r, c int, z int64) int64{
	if z < -(4 << 56) {
		return 0;
	}else if z > (4 << 56){
		return (1 << 56);
	}else{
		return (0x80000000000000 + Multiply(0x20000000000000,z))
	}
}

func linearPW2(r, c int, z int64) int64 {
	if z > (4 << 56){
		return 0;
	}else if z < -(4 << 56){
		return 0;
	}
	return 0x20000000000000;
}

func regionPW(r, c int, z int64) int64{
    if z < -(8 << 56) {
        return 0
    } else if z > (8 << 56){
        return (1 << 56)
    }
    if z < -(6 << 56) {
        return Multiply(0x005105DDCEA003, z) + 0x02882EE5F5001D
    }
    if z < -(4 << 56) {
        return Multiply(0x01FC5965FC8E58, z) + 0x0C8C241F832A11
    }
    if z < -(2 << 56) {
        return Multiply(0x0CF4AB520DA07E, z) + 0x386D6BCFC772AC
    }
    if z <  (2 << 56) {
        return Multiply(0x30BDF56A29E728, z) + 0x80000000000000
    }
    if z <  (4 << 56) {
        return Multiply(0x0CF4AB520DA07E, z) + 0xC7929430388D53
    }
    if z <  (6 << 56) {
        return Multiply(0x01FC5965FC8E58, z) + 0xF373DBE07CD5EE
    } else {
        return Multiply(0x005105DDCEA003, z) + 0xFD77D111A0AFFE
    }
}

func regionPW2(r, c int, z int64) int64{
    if z < -(8 << 56) {
        return 0
    } else if z > (8 << 56){
        return 0
    }
    if z < -(6 << 56) {
        return 0x005105DDCEA003
    }
    if z < -(4 << 56) {
        return 0x01FC5965FC8E58
    }
    if z < -(2 << 56) {
        return 0x0CF4AB520DA07E
    }
    if z <  (2 << 56) {
        return 0x30BDF56A29E728
    }
    if z <  (4 << 56) {
        return 0x0CF4AB520DA07E
    }
    if z <  (6 << 56) {
        return 0x01FC5965FC8E58
    } else {
        return 0x005105DDCEA003
    }

}

func randomArray(size int, v int64) (data []int64) {
	data = make([]int64, size);
	for i := 0; i < size; i++ {
		data[i] = rand.Int63n(0x1 << 56);
	}
	return;
}

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


func load(net *Network){
	testFile, _ := os.Open("data/hweights.csv");
	r := csv.NewReader(bufio.NewReader(testFile));
	hweights := make([]int64, net.hiddens*net.inputs);
	rowCount := 0;
	for {
			record, err := r.Read();
			if err == io.EOF {
				break;
			}

			for i := range record {
				hweights[rowCount*net.inputs + i], _ = strconv.ParseInt(record[i],10,64);
			}
			rowCount = rowCount+1;
		}
		testFile.Close();
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
