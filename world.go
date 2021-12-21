package main

import(
    "lukechampine.com/frand"
)

func Abs32(v int32) int32 {
    if v < 0 {
        return -v
    }
    return v
}

type World struct {
    cells [XMax * YMax]bool
    goals []Zone
    width int32
    height int32
}

func NewWorld(width, height int32, goals []Zone) *World {
    return &World{
        width: width,
        height: height,
        goals: goals,
    }
}

func (w World) block(x, y int32) {
    w.cells[x + y * w.width] = true
}

func (w World) unblock(x, y int32) {
    w.cells[x + y * w.width] = false
}

func (w World) blocked(x, y int32) bool {
    if x < 0 || x >= w.width || y < 0 || y >= w.height {
        return true
    }
    return w.cells[x + y * w.width]
}

func (w World) blockedf(x, y int32) float64 {
    if x < 0 || x >= w.width {
        return 1.0
    }

    if y < 0 || y >= w.height {
        return 1.0
    }

    if w.cells[x + y * w.width] {
        return 1.0
    } else {
        return 0.0
    }
}

func randint32(v int32) int32 {
    return int32(frand.Intn(int(v)))
}

func (w *World) blockrand() (int32, int32) {
    x := randint32(w.width)
    y := randint32(w.height)
    for w.blocked(x, y) {
        x = randint32(w.width)
        y = randint32(w.height)
    }

    w.block(x, y)
    return x, y
}

func (w *World) NearestZoneXY(x, y int32) (int32, int32) {
    minxd, minyd := w.goals[0].Dist(x, y)
    mindist := Abs32(minxd) + Abs32(minyd)
    for _, z := range(w.goals) {
        xd, yd := z.Dist(x, y)
        dist := Abs32(xd) + Abs32(yd)
        if dist < mindist {
            mindist = dist
            minxd = xd
            minyd = yd
        }
    }

    return minxd, minyd
}

func (w *World) InGoal(x, y int32) bool {
    for _, z := range(w.goals) {
        if z.InZone(x, y) {
            return true
        }
    }

    return false
}
