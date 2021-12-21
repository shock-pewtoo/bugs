package main

import(
    "math/rand"
    "log"
    "os"
    "github.com/hajimehoshi/ebiten/v2"
)

const RequiredSteps = 100

const (
    TInput uint32 = iota
    TOutput
)

const (
    DistX uint32 = iota
    DistY
    TotalSteps
    // BlockedLeft
    // BlockedRight
    // BlockedUp
    // BlockedDown
)
const InputCount = 3
var InputStringMap = []string{
    "DistX",
    "DistY",
    "TotalSteps",
    // "BlockedLeft"
    // "BlockedRight"
    // "BlockedUp"
    // "BlockedDown"
}

const GeneCount = 8
const HiddenCount = 6

const (
    GoLeft uint32 = iota
    GoRight
    GoUp
    GoDown
)
var OutputStringMap = []string{
    "GoLeft",
    "GoRight",
    "GoUp",
    "GoDown",
}

const OutputCount = 4

func bool2f(v bool) float64 {
    if v {
        return 1.0
    }
    return 0.0
}

type Bug struct {
    dead bool
    genes []Gene
    network *Network
    world *World
    x int32
    y int32
    totalsteps uint32
}

func RandBug(genecount uint32, world *World) *Bug {
    genes := make([]Gene, 0)
    var i uint32
    for i = 0; i < genecount; i++ {
        genes = append(genes, RandGene())
    }
    return NewBug(genes, world)
}

func NewBug(genes []Gene, world *World) *Bug {
    x, y := world.blockrand()
    bug := Bug{
        x: x,
        y: y,
        world: world,
        genes: genes,
    }

    world.block(bug.x, bug.y)

    bug.network = NewNetwork(InputCount, HiddenCount, OutputCount)
    return &bug
}

func (b *Bug) step() {
    b.network.clear()
    b.totalsteps = b.totalsteps + 1
    xd, yd := b.world.NearestZoneXY(b.x, b.y)
    b.network.set(TInput, DistX, float64(xd))
    b.network.set(TInput, DistY, float64(yd))
    b.network.set(TInput, TotalSteps, float64(b.totalsteps / RequiredSteps))
    /*
    b.network.set(TInput, BlockedLeft, b.world.blockedf(b.x - 1, b.y))
    b.network.set(TInput, BlockedRight, b.world.blockedf(b.x + 1, b.y))
    b.network.set(TInput, BlockedUp, b.world.blockedf(b.x, b.y - 1))
    b.network.set(TInput, BlockedDown, b.world.blockedf(b.x, b.y + 1))
    */

    for _, gene := range(b.genes) {
        gene.processGene(b.network)
    }

    var xmove, ymove int32
    xmove = 0
    ymove = 0
    if b.network.get(TOutput, GoLeft) > 0.5 {
        xmove = xmove - 1
    }

    if b.network.get(TOutput, GoRight) > 0.5 {
        xmove = xmove + 1
    }

    if b.network.get(TOutput, GoUp) > 0.5 {
        ymove = ymove - 1
    }

    if b.network.get(TOutput, GoDown) > 0.5 {
        ymove = ymove + 1
    }

    b.move(xmove, ymove)
}

func (b *Bug) move(i, j int32) {
    newx := b.x + i
    newy := b.y + j

    if newx < 0 {
        return
    }

    if newx > XMax {
        return
    }

    if newy < 0 {
        return
    }

    if newy > YMax {
        return
    }

    if b.world.blocked(newx, newy) {
        return
    }

    b.world.unblock(b.x, b.y)
    b.x = newx
    b.y = newy
    b.world.block(b.x, b.y)
}

func (b Bug) mutate(prob float64) {
    if rand.Float64() < prob {
        gene := rand.Int31n(int32(len(b.genes)))
        bit := rand.Int31n(8)
        log.Printf("mutating gene %d, bit %d", gene, bit)
        b.genes[gene] = b.genes[gene] ^ 1 << bit
    }
}

func (b1 *Bug) reproduce(b2 *Bug) *Bug {
    genes := make([]Gene, len(b1.genes))
    for i := 0; i < len(b1.genes); i++ {
        p := rand.Int31n(2)
        if p == 0 {
            genes[i] = b1.genes[i]
        } else {
            genes[i] = b2.genes[i]
        }
    }

    b := NewBug(genes, b1.world)
    b.mutate(0.01)
    return b
}

func (b *Bug) die() {
    b.world.unblock(b.x, b.y)
    b.dead = true
}

func (b *Bug) draw(dst *ebiten.Image) {
    // log.Printf("Drawing bug at (%d, %d)", b.x, b.y)
    DrawRect(dst, bugImage, b.x, b.y, 1, 1)
}

// Convert a bug's genes to graphviz representation
func (b *Bug) dot(file string ) {
    f, err := os.Create(file)
    if err != nil {
        log.Fatal(err)
    }

    f.Write([]byte("digraph {\n"))
    for _, g := range(b.genes) {
        f.Write([]byte(g.dot(b.network.inputsize, b.network.outputsize)))
    }
    f.Write([]byte("}\n"))
}
