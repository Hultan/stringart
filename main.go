package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	size      = 400
	margin    = 30
	maxLines  = 10000
	nailCount = 100
)

var (
	nails     = make([]rl.Vector2, nailCount)
	lineIndex = make([]int, 0, maxLines)
	col       []color.RGBA
)

func main() {
	rl.InitWindow(size, size, "String art")

	img := rl.LoadImage("einstein.jpg")
	col = rl.LoadImageColors(img)
	tex := rl.LoadTextureFromImage(img)

	for i := 0; i < nailCount; i++ {
		angle := rl.Pi * 2 / float64(nailCount) * float64(i)
		r := float64(size/2 - margin)
		x := size/2 + float32(r*math.Cos(angle))
		y := size/2 + float32(r*math.Sin(angle))
		nails[i] = rl.NewVector2(x, y)
	}

	startIndex := rand.Intn(nailCount)
	lineIndex = append(lineIndex, startIndex)

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		//rl.DrawTexture(tex, 0, 0, rl.White)

		for i := 0; i < nailCount; i++ {
			rl.DrawCircle(int32(nails[i].X), int32(nails[i].Y), 2, rl.Red)
		}

		for i := 1; i < len(lineIndex); i++ {
			nail1 := nails[lineIndex[i-1]]
			nail2 := nails[lineIndex[i]]
			rl.DrawLineEx(nail1, nail2, 0.5, rl.Fade(rl.Black, 0.25))
		}

		if len(lineIndex) < maxLines {
			current := lineIndex[len(lineIndex)-1]
			next, ok := findNextIndex(current)

			if ok {
				lineIndex = append(lineIndex, next)
				updateImage(current, next)
			} else {
				fmt.Println("No valid index found, stopping...")
			}

		} else {
			fmt.Println("Max lines reached")
		}

		rl.EndDrawing()
	}

	rl.UnloadTexture(tex)
	rl.CloseWindow()
}

func updateImage(current, next int) {
	nail1 := nails[current]
	nail2 := nails[next]
	steps := 100
	bright := uint8(10)

	for i := 0; i < steps; i++ {
		l := float32(i) / float32(steps)
		x := lerp(nail1.X, nail2.X, l)
		y := lerp(nail1.Y, nail2.Y, l)
		pixelIndex := int(x) + int(y)*size

		if pixelIndex >= 0 && pixelIndex < len(col) {
			if col[pixelIndex].R < 255 {
				col[pixelIndex].R += bright
				col[pixelIndex].G += bright
				col[pixelIndex].B += bright
			}
		}
	}

}

func findNextIndex(current int) (int, bool) {
	next := -1
	highestContrast := -1

	for i := 0; i < len(nails); i++ {
		if i != current {
			contrast := evaluateContrast(current, i)
			if contrast > highestContrast {
				highestContrast = contrast
				next = i
			}
		}
	}

	if next == -1 {
		next = rand.Intn(nailCount)
		fmt.Println("Finding random next index")
	}

	return next, true
}

func evaluateContrast(current, next int) int {
	total := 0
	nail1 := nails[current]
	nail2 := nails[next]

	steps := 100

	for i := 0; i < steps; i++ {
		x := lerp(nail1.X, nail2.X, float32(i)/float32(steps))
		y := lerp(nail1.Y, nail2.Y, float32(i)/float32(steps))

		if valid(x, y) {
			pixelIndex := int(x) + int(y)*size
			brightness := int(col[pixelIndex].R)
			total += 255 - brightness
		}
	}

	return total / steps
}

func valid(x float32, y float32) bool {
	return x >= 0 && x < size && y >= 0 && y < size
}

func lerp(a, b, t float32) float32 {
	return a + (b-a)*t
}
