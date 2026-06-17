package main

import (
	"fmt"
	"os"
	"runtime"
	"sort"
	"strconv"
	"time"
)

type WordFreq struct {
	Word  string
	Count int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run . <caminho_do_arquivo> [numero_de_workers]")
		os.Exit(1)
	}

	filename := os.Args[1]

	numWorkers := runtime.NumCPU()
	if len(os.Args) >= 3 {
		val, err := strconv.Atoi(os.Args[2])
		if err != nil || val <= 0 {
			fmt.Printf("Aviso: Número de workers inválido '%s'. Usando o padrão (%d).\n", os.Args[2], numWorkers)
		} else {
			numWorkers = val
		}
	}

	// 1. Execução Sequencial
	startSeq := time.Now()
	seqMap, err := RunSequential(filename)
	durationSeq := time.Since(startSeq)

	if err != nil {
		fmt.Printf("Erro na execução sequencial: %v\n", err)
		os.Exit(1)
	}

	// 2. Execução Concorrente
	startConc := time.Now()
	concMap, err := RunConcurrent(filename, numWorkers)
	durationConc := time.Since(startConc)

	if err != nil {
		fmt.Printf("Erro na execução concorrente: %v\n", err)
		os.Exit(1)
	}

	// 3. Comparação de Resultados
	areEqual := compareMaps(seqMap, concMap)
	resString := "não"
	if areEqual {
		resString = "sim"
	}

	// 4. Obter Top 20 palavras mais frequentes
	top20 := getTop20(seqMap)

	// 5. Exibir Saída
	fmt.Printf("Tempo sequencial: %v\n", durationSeq)
	fmt.Printf("Tempo concorrente: %v (Workers: %d)\n", durationConc, numWorkers)
	fmt.Printf("Resultados iguais: %s\n", resString)
	fmt.Println("Top 20 palavras:")
	for i, wf := range top20 {
		fmt.Printf("%2d. %-15s (%d)\n", i+1, wf.Word, wf.Count)
	}
}

// compareMaps verifica se os dois mapas têm as mesmas chaves e valores.
func compareMaps(m1, m2 map[string]int) bool {
	if len(m1) != len(m2) {
		return false
	}
	for k, v1 := range m1 {
		v2, ok := m2[k]
		if !ok || v1 != v2 {
			return false
		}
	}
	return true
}

// getTop20 extrai as 20 palavras mais frequentes de forma ordenada.
func getTop20(freqs map[string]int) []WordFreq {
	list := make([]WordFreq, 0, len(freqs))
	for k, v := range freqs {
		list = append(list, WordFreq{Word: k, Count: v})
	}

	// Ordena decrescente por frequência; se igual, alfabeticamente
	sort.Slice(list, func(i, j int) bool {
		if list[i].Count == list[j].Count {
			return list[i].Word < list[j].Word
		}
		return list[i].Count > list[j].Count
	})

	if len(list) > 20 {
		return list[:20]
	}
	return list
}
