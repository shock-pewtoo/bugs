package main

import(
    "math"
)

type Network struct {
    inputs []float64
    inputsize uint32
    outputs []float64
    outputsize uint32
}

func NewNetwork(inputsize uint32, hiddensize uint32, outputsize uint32) *Network {
    inputs := make([]float64, inputsize + hiddensize)
    outputs := make([]float64, outputsize)
    n := &Network{
        inputs: inputs,
        inputsize: inputsize,
        outputs: outputs,
        outputsize: outputsize,
    }

    return n
}

func (n Network) absaddr(atype, addr uint32) uint32 {
    var size uint32
    if atype == 0 {
        size = n.inputsize
    } else {
        size = n.outputsize
    }

    a := addr % uint32(size)
    return a
}

func (n Network) get(atype, addr uint32) float64 {
    abs := n.absaddr(atype, addr)
    var v float64
    if atype == 0 {
        v = n.inputs[abs]
    } else {
        v = n.outputs[abs]
    }
    return math.Tanh(v)
}

func (n Network) add(atype, addr uint32, value float64) {
    abs := n.absaddr(atype, addr)
    if atype == 0 {
        n.inputs[abs] = n.inputs[abs] + value
    } else {
        n.outputs[abs] = n.outputs[abs] + value
    }
}

func (n Network) set(atype, addr uint32, value float64) {
    abs := n.absaddr(atype, addr)
    if atype == 0 {
        n.inputs[abs] = value
    } else {
        n.outputs[abs] = value
    }
}

func (n Network) clear() {
    var i uint32
    for i = 0; i < n.inputsize; i++ {
        n.inputs[i] = 0
    }

    for i = 0; i < n.outputsize; i++ {
        n.outputs[i] = 0
    }
}
