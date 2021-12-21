package main

import(
    _ "image"
    "image/color"
    "fmt"
    "log"
    "math"
    "math/rand"
    "time"
    "os"
    "runtime/pprof"

    "github.com/hajimehoshi/ebiten/v2"
)

const (
    XMax = 1024
    YMax = 1024
    Generations = 100000
    RoundsPerGen = 750
    MaxBugs = 100
)

var (
    bugImage = ebiten.NewImage(3, 3)
    greenZone = ebiten.NewImage(3, 3)
)

type GameColor struct {
    r uint32
    g uint32
    b uint32
    a uint32
}

func (gc GameColor) RGBA() (uint32, uint32, uint32, uint32) {
    return gc.r, gc.b, gc.g, gc.a
}

var Red GameColor = GameColor{128, 0, 0, 1}
var Green GameColor = GameColor{1, 1, 1, 1}

type Game struct {
    bugs []*Bug
    world *World
    generation int
    round int
}

func pow(a, b uint32) uint32 {
    return uint32(math.Pow(float64(a), float64(b)))
}

func main() {
    bugImage.Fill(color.White)
    greenZone.Fill(Green)

    f, err := os.Create("/tmp/bugs.goprofile")
    if err != nil {
        log.Fatal(err)
    }
    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()

    zones := []Zone{
        Zone{100, 100, 200, 200},
        Zone{300, 500, 400, 600},
        Zone{500, 500, 600, 600},
        Zone{700, 500, 800, 600},
    }

    world := NewWorld(XMax, YMax, zones)
    rand.Seed(time.Now().Unix())
    g := &Game{
        world: world,
        bugs: firstgen(world),
        round: 0,
    }

    log.Printf("Starting bugs.")
    ebiten.SetFullscreen(false)
    ebiten.SetWindowSize(XMax, YMax)
    ebiten.SetWindowTitle("Bugs!")

    if err := ebiten.RunGame(g); err != nil {
        log.Fatal(err)
    }
}

func firstgen(world *World) []*Bug {
    bugs := make([]*Bug, MaxBugs)
    for i := 0; i < MaxBugs; i++ {
        bugs[i] = RandBug(GeneCount, world)
    }
    return bugs
}

func (g *Game) roundover() bool {
    if g.round > RoundsPerGen {
        return true
    }
    return false
}

func (g *Game) Update() error {
    if g.roundover() {
        g.nextgen()
        g.round = 0
        g.generation = g.generation + 1
    } else {
        g.round = g.round + 1
        for _, bug := range(g.bugs) {
            bug.step()
        }
    }

    return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
    for _, zone := range(g.world.goals) {
        zone.draw(screen)
    }
    for _, bug := range(g.bugs) {
        bug.draw(screen)
    }
}

func (g *Game) nextgen() {
    log.Printf("Nextgen")
    survivors := make([]*Bug, 0)
    for _, bug := range(g.bugs) {
        if g.world.InGoal(bug.x, bug.y) && !bug.dead && bug.totalsteps > RequiredSteps {
            survivors = append(survivors, bug)
        }
        bug.die()
    }

    children := make([]*Bug, 0)

    if len(survivors) > 0 {
        path := fmt.Sprintf("/tmp/gen_%05d.dot", g.generation)
        survivors[0].dot(path)
    }

    for i := 1; i < len(survivors); i++ {
        p1 := survivors[i-1]
        p2 := survivors[i]
        children = append(children, p1.reproduce(p2))
        children = append(children, p2.reproduce(p1))

        if len(children) >= MaxBugs {
            break
        }
    }

    for len(children) < MaxBugs {
        children = append(children, RandBug(GeneCount, g.world))
    }

    a := len(survivors)
    b := len(g.bugs) - len(survivors)
    pct := float32(len(survivors)) / float32(len(g.bugs)) * 100.0
    log.Printf("gen %d: %d survivors, %d deaths: %f", g.generation, a, b, pct)
    g.bugs = children
    log.Printf("Nextgendone")
}

func DrawRect(dst *ebiten.Image, srcImage *ebiten.Image, x, y, width, height int32) {
    fwidth := float64(width)
    fheight := float64(height)
    fx := float64(x)
    fy := float64(y)

    op := &ebiten.DrawImageOptions{}
    op.GeoM.Scale(fwidth, fheight)
    op.GeoM.Translate(fx, fy)
    dst.DrawImage(srcImage, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
    return XMax, YMax
}
