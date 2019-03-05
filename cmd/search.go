package main

import (
	"github.com/urfave/cli"
	"os"
	"math/rand"
	"time"
	"math"
	"sort"
)

const (
	crossFactor        = 0.3
	mutateFactor       = 0.2
	gridWidth          = 5
	gridHeight         = 3
	iterationTimeMax   = 2000
	iterationTimeMin   = 50
	iterationFactor    = 100
	processTimeSeconds = 2
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
	var populationNum = 500
	var bestFit float64 = 0
	var bestImage [][]byte
	var bestLine []int32

	var cur = initPopulation(populationNum, gridWidth, gridHeight, alphas)
	var iterationNum = 0
	var convergenceNum = 0
	var stop = false
	var fitness []float64
	var genBest float64
	var genBestImage [][]byte
	var genLine []int32
	var bestScore int32
	fitness, genBest, genBestImage, genLine, bestScore = calcImagesFitness(cur, lines, scoreTable, sum)
	var startTime = time.Now()
	for !stop {
		if iterationNum > iterationTimeMin && convergenceNum > iterationFactor {
			cur = initPopulation(populationNum, gridWidth, gridHeight, alphas)
			convergenceNum = 0
		}
		var toCross = selectImages(50, cur, fitness)
		var crossed = crossMatch(toCross, lines, scoreTable, sum)
		fitness, genBest, genBestImage, genLine, bestScore = calcImagesFitness(cur, lines, scoreTable, sum)
		if genBest > bestFit {
			bestFit = genBest
			bestImage = genBestImage
			bestLine = genLine
			convergenceNum = 0
			println("best:", bestFit)
		} else {
			convergenceNum += 1
		}
		mutate(crossed, alphas)
		cur = crossed
		fitness, genBest, genBestImage, genLine, bestScore = calcImagesFitness(cur, lines, scoreTable, sum)
		if genBest > bestFit {
			bestFit = genBest
			bestImage = genBestImage
			bestLine = genLine
			convergenceNum = 0
			println("best:", bestFit, " score:", bestScore)
		} else {
			convergenceNum += 1
		}
		iterationNum += 1
		if iterationNum > iterationTimeMax {
			stop = true
		}
		if genBest == 1 {
			stop = true
		}
		var n = time.Now()
		if n.Sub(startTime).Seconds() > processTimeSeconds {
			stop = true
		}
	}
	println("best")
	for _, d := range bestImage {
		for _, d1 := range d {
			print(string(d1))
			print(" ")
		}
		print(";")
	}
	print("\n")
	println("best score:", bestFit)
	println("best line")
	for _, l := range bestLine {
		print(l)
		print(" ")
	}
}

func crossMatch(images [][][]byte, lines [][][]int, scoreTable [][]int32, des int32) [][][]byte {
	var result = make([][][]byte, 0, len(images))
	var n = len(images)
	var random = rand.New(rand.NewSource(time.Now().UnixNano()))
	random.Shuffle(n, func(i, j int) {
		var s = images[i]
		images[i] = images[j]
		images[j] = s
	})
	var tuples = n / 2
	for i := 0; i < tuples; i++ {
		var left = images[i*2]
		var right = images[i*2+1]
		var leftChild = make([][]byte, len(left))
		for idx := range leftChild {
			var line = make([]byte, len(left[idx]))
			copy(line, left[idx])
			leftChild[idx] = line
		}
		var rightChild = make([][]byte, len(right))
		for idx := range rightChild {
			var line = make([]byte, len(right[idx]))
			copy(line, right[idx])
			rightChild[idx] = line
		}
		for h := 0; h < gridHeight; h++ {
			for w := 0; w < gridWidth; w++ {
				var r = random.Float64()
				if r > crossFactor {
					leftChild[h][w] = right[h][w]
					rightChild[h][w] = left[h][w]
				}
			}
		}
		var leftFitness, _, _ = calcImageFitness(left, lines, scoreTable, des)
		var leftChildFitness, _, _ = calcImageFitness(leftChild, lines, scoreTable, des)
		var rightFitness, _, _ = calcImageFitness(right, lines, scoreTable, des)
		var rightChildFitness, _, _ = calcImageFitness(rightChild, lines, scoreTable, des)
		if leftFitness > leftChildFitness {
			result = append(result, left)
		} else {
			result = append(result, leftChild)
		}
		if rightFitness > rightChildFitness {
			result = append(result, right)
		} else {
			result = append(result, rightChild)
		}
	}
	return result
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

func calcImageFitness(image [][]byte, lines [][][]int, scoreTable [][]int32, des int32) (float64, []int32, int32) {
	var scores = make([]int32, 0, len(lines))
	for i := 0; i < len(lines); i++ {
		var l = lines[i]
		var alphaLines = make([]byte, 0, len(l))
		for _, indexes := range l {
			var alpha = image[indexes[0]][indexes[1]]
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
				} else {
					lastAlpha = alpha
					count = 1
				}
			}
			score = curScore
		}
		scores = append(scores, score)
	}
	var bestScore, bestLine = calcBestScore(scores, des)
	if bestScore == des {
		println(bestScore)
	}
	var f float64
	if bestScore < des {
		f = float64(bestScore) / float64(des)
	} else {
		f = float64(des) / float64(bestScore)
	}
	return f, bestLine, bestScore
}

func mutate(images [][][]byte, alphas []byte) {
	for _, img := range images {
		var mr = rand.Float64()
		if mr > mutateFactor {
			var h = rand.Int31n(gridHeight)
			var w = rand.Int31n(gridWidth)
			var mutateIndex = rand.Int31n(int32(len(alphas)))
			img[h][w] = alphas[mutateIndex]
		}
	}
}

func calcImagesFitness(alphaImages [][][]byte, lines [][][]int, scoreTable [][]int32, des int32) ([]float64, float64, [][]byte, []int32, int32) {
	var fitness = make([]float64, len(alphaImages))
	var bestFitness float64 = 0
	var bestLine []int32
	var bestImage = make([][]byte, gridHeight)
	var bestScore int32
	for idx := range bestImage {
		bestImage[idx] = make([]byte, gridWidth)
	}
	for imgIdx, img := range alphaImages {
		var f, line, score = calcImageFitness(img, lines, scoreTable, des)
		fitness[imgIdx] = f
		if f > bestFitness {
			bestFitness = f
			for idx, l := range img {
				copy(bestImage[idx], l)
			}
			bestLine = line
			bestScore = score
		}
	}
	return fitness, bestFitness, bestImage, bestLine, bestScore
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
		for h := 0; h < height; h++ {
			var hs []byte
			for w := 0; w < width; w++ {
				var index = r.Intn(alphaLen)
				hs = append(hs, alphas[index])
			}
			grid = append(grid, hs)
		}
		//grid = append(grid, []byte{'A', 'A', 'A', 'A', 'A'})
		result = append(result, grid)
	}
	return result
}
