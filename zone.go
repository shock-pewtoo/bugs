package main

import(
    _ "image"
    "image/color"
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Zone struct {
    x1, y1, x2, y2  int32
}

func (z Zone) InZone(x, y int32) bool {
    if x < z.x1 {
        return false
    }

    if x > z.x2 {
        return false
    }

    if y < z.y1 {
        return false
    }

    if y > z.y2 {
        return false
    }

    return true
}

func (z Zone) Dist(x, y int32) (int32, int32) {
    zx := (z.x1 + z.x2) / 2
    zy := (z.y1 + z.y2) / 2
    return zx - x, zy - y
}

func (z *Zone) draw(dst *ebiten.Image) {
    fx1 := float64(z.x1)
    fy1 := float64(z.y1)
    fx2 := float64(z.x2)
    fy2 := float64(z.y2)
    ebitenutil.DrawLine(dst, fx1, fy1, fx2, fy1, color.White)
    ebitenutil.DrawLine(dst, fx1, fy2, fx2, fy2, color.White)
    ebitenutil.DrawLine(dst, fx1, fy1, fx1, fy2, color.White)
    ebitenutil.DrawLine(dst, fx2, fy1, fx2, fy2, color.White)
}
