package main

import (
	"github.com/urfave/cli"
	"os"
	"math/rand"
	"time"
	"math"
	"sort"
)

func main() {
	cliApp := cli.NewApp()
	cliApp.Name = "find alphabet"
	cliApp.Commands = []cli.Command{{
		Name:   "search",
		Usage:  "find alphabet",
		Action: runSearch,
		Flags: []cli.Flag{
			cli.Int64Flag{
				Name:  "sum,s",
				Value: 100,
				Usage: "destination sum",
			},
		},
	}}
	cliApp.Flags = append(cliApp.Flags, []cli.Flag{}...)
	cliApp.Run(os.Args)
}

func runSearch(ctx *cli.Context) {
	var sum = int32(ctx.Int64("sum"))
	println(sum)
	var alphas = []byte{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I'}
	var gridWidth = 5
	var gridHeight = 3
	var scoreTable = [][]int32{{1, 2, 3, 4, 5, 6, 7, 8, 9}, {2, 4, 6, 8, 10, 12, 14, 16, 18}, {4, 8, 12, 16, 20, 24, 28, 32, 36}}
	////println(scoreTable[1][1])
	var lines = [][][]int{
		{{1, 0}, {1, 1}, {1, 2}, {1, 3}, {1, 4}},
		{{0, 0}, {0, 1}, {0, 2}, {0, 3}, {0, 4}},
		{{2, 0}, {2, 1}, {2, 2}, {2, 3}, {2, 4}},
		{{0, 0}, {1, 1}, {2, 2}, {1, 3}, {0, 4}},
		{{2, 0}, {1, 1}, {0, 2}, {1, 3}, {2, 4}},
		{{2, 0}, {2, 1}, {1, 2}, {0, 3}, {0, 4}},
		{{0, 0}, {0, 1}, {1, 2}, {2, 3}, {2, 4}},
		{{1, 0}, {2, 1}, {2, 2}, {2, 3}, {1, 4}},
		{{1, 0}, {0, 1}, {0, 2}, {0, 3}, {1, 4}},
	}
	var populationNum = 50
	var init = initPopulation(populationNum, gridWidth, gridHeight, alphas)
	var fitness, best = calcFitness(init, lines, scoreTable, sum)
	var toCross = selectImages(50, init, fitness)

	for _, x := range init {
		for _, d := range x {
			for _, d1 := range d {
				print(string(d1))
				print(" ")
			}
			print(";")
		}
		print("\n")
	}
	println("best:", best)
	for _, f := range fitness {
		println(f)
	}
	for _, x := range toCross {
		for _, d := range x {
			for _, d1 := range d {
				print(string(d1))
				print(" ")
			}
			print(";")
		}
		print("\n")
	}
}

func selectImages(n int, images [][][]byte, fitness []float64) [][][]byte {
	var totalFitness float64 = 0
	for _, f := range fitness {
		totalFitness += f
	}
	var random = rand.New(rand.NewSource(time.Now().UnixNano()))
	var randoms = make([]float64, n)
	for i := 0; i < n; i++ {
		randoms[i] = random.Float64() * totalFitness
	}
	sort.Float64s(randoms)
	var selected = make([][][]byte, 0, n)
	var accum float64 = 0
	var fitIndex = 0
	for _, r := range randoms {
		for accum+fitness[fitIndex] < r {
			accum += fitness[fitIndex]
			fitIndex += 1
		}
		selected = append(selected, images[fitIndex])
	}
	return selected
}

func calcFitness(alphaImages [][][]byte, lines [][][]int, scoreTable [][]int32, des int32) ([]float64, float64) {
	var fitness = make([]float64, len(alphaImages))
	var bestFitness float64 = 0
	for imgIdx, img := range alphaImages {
		var scores = make([]int32, 0, len(lines))
		for i := 0; i < len(lines); i++ {
			var l = lines[i]
			var alphaLines = make([]byte, 0, len(l))
			for _, indexes := range l {
				var alpha = img[indexes[0]][indexes[1]]
				alphaLines = append(alphaLines, alpha)
			}
			var score int32 = 0
			var count = 0
			var lastAlpha byte = ' '
			for _, alpha := range alphaLines {
				var curScore int32 = 0
				if lastAlpha == ' ' {
					lastAlpha = alpha
					count = 1
				} else {
					if alpha == lastAlpha {
						count += 1
						if count >= 3 {
							curScore = scoreTable[count-3][lastAlpha-'A']
						}
					}
				}
				score += curScore
			}
			scores = append(scores, score)
		}
		var bestScore, _ = calcBestScore(scores, des)
		var f = float64(bestScore) / float64(des)
		fitness[imgIdx] = f
		if f > bestFitness {
			bestFitness = f
		}
	}
	return fitness, bestFitness
}

func calcBestScore(scores []int32, des int32) (int32, []int32) {
	var final = int32(math.Pow(float64(2), float64(len(scores)))) - 1
	var start int32 = 0
	var cur = start
	var bestScore int32 = -1
	var bestLine = make([]int32, len(scores))
	var curLine = make([]int32, len(scores))
	var swap = curLine
	for cur <= final {
		var score int32 = 0
		for idx := 0; idx < len(scores); idx++ {
			var b int32
			if (1 << uint32(idx) & cur) > 0 {
				b = 1
			} else {
				b = 0
			}
			curLine[idx] = b
			score += b * scores[idx]
		}
		if math.Abs(float64(score-des)) < math.Abs(float64(bestScore-des)) {
			swap = curLine
			curLine = bestLine
			bestLine = swap
			bestScore = score
		}
		cur += 1
	}
	return bestScore, bestLine
}

func initPopulation(num, width, height int, alphas []byte) [][][]byte {
	var alphaLen = len(alphas)
	var r = rand.New(rand.NewSource(time.Now().UnixNano()))
	var result [][][]byte
	for i := 0; i < num; i++ {
		var grid [][]byte
		for h := 0; h < height-1; h++ {
			var hs []byte
			for w := 0; w < width; w++ {
				var index = r.Intn(alphaLen)
				hs = append(hs, alphas[index])
			}
			grid = append(grid, hs)
		}
		grid = append(grid, []byte{'A', 'A', 'A', 'A', 'A'})
		result = append(result, grid)
	}
	return result
}