package main

import(
    "fmt"
    "math/rand"
)

type Gene uint32

// Divide the value field by this to get -4 -> 4
const FloatMagic = 8192.0

/*
0 = in type (1)
1-7 = in address (7)
8 = out type (1)
9-16 = out address (7)
17 = signbit (1)
18-31 = value * 4000 (11)
*/

func RandGene() Gene {
    return Gene(rand.Int31())
}

func typeString(i uint32) string {
    switch i {
    case 0:
        return "source"
        case 1: return "target"
    default:
        return "error"
    }
}

func (g Gene) bits(first, length uint32) uint32 {
    v := uint32(g)
    a := v >> (first - 1)
    return a & ((1 << length) - 1)
}

func (g Gene) sourceType() uint32 {
    return g.bits(0, 1)
}

func (g Gene) sourceAddr() uint32 {
    return g.bits(1, 7)
}

func (g Gene) sourceTypeString() string {
    return typeString(g.sourceType())
}
func (g Gene) destType() uint32 {
    return g.bits(8, 1)
}

func (g Gene) destAddr() uint32 {
    return g.bits(9, 7)
}

func (g Gene) destTypeString() string {
    return typeString(g.destType())
}

func (g Gene) weight() float64 {
    signbit := g.bits(17, 1)
    sign := 1.0
    if signbit == 1 {
        sign = -1.0
    }
    value := g.bits(18, 15)
    return float64(value) / FloatMagic * sign
}

func (g Gene) processGene(n *Network) {
    stype := g.sourceType()
    saddr := g.sourceAddr()

    dtype := g.destType()
    daddr := g.destAddr()

    value := n.get(stype, saddr) * g.weight()
    n.add(dtype, daddr, value)
}

// Return a graphviz representation of this gene
func (g Gene) dot(inputsize, outputsize uint32) string {
    stype := g.sourceType()
    saddr := g.sourceAddr()
    stypestring, saddrstring := addrAsString(stype, saddr, inputsize, outputsize)
    dtype := g.destType()
    daddr := g.destAddr()
    dtypestring, daddrstring := addrAsString(dtype, daddr, inputsize, outputsize)
    weight := g.weight() / FloatMagic
    return fmt.Sprintf("\"%s_%s\" -> \"%s_%s\" [label=\"%f\"];\n", stypestring, saddrstring, dtypestring, daddrstring, weight)
}

func addrAsString(atype, addr uint32, inputsize, outputsize uint32) (string, string) {
    var typestring string
    var addrstring string
    if atype == 0 {
        typestring = "input"
        modaddr := addr % inputsize
        if modaddr < uint32(len(InputStringMap)) {
            addrstring = fmt.Sprintf("%s", InputStringMap[modaddr])
        } else {
            addrstring = fmt.Sprintf("%d", modaddr)
        }
    } else {
        typestring = "output"
        modaddr := addr % outputsize
        if modaddr < uint32(len(OutputStringMap)) {
            addrstring = fmt.Sprintf("%s", OutputStringMap[modaddr])
        } else {
            addrstring = fmt.Sprintf("%d", modaddr)
        }
    }

    return typestring, addrstring
}
