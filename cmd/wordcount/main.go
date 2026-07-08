// Command wordcount é uma pequena demonstração executável da função CountWords.
//
// Uso:
//
//	go run ./cmd/wordcount                 # usa o texto de exemplo do enunciado
//	echo "seu texto aqui" | go run ./cmd/wordcount
//	go run ./cmd/wordcount < arquivo.txt
package main

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/danlisb/go_automated_tests"
)

// defaultText é o texto de exemplo do enunciado, usado quando nada é passado
// via stdin.
const defaultText = `Casa, casa! A casa é azul.
Árvore; árvore? verde.
Go go Go. IA é útil, mas IA erra.`

func main() {
	text := defaultText

	// Se houver algo em stdin (pipe ou redirecionamento de arquivo), usa esse
	// conteúdo em vez do texto de exemplo.
	if info, _ := os.Stdin.Stat(); info != nil && info.Mode()&os.ModeCharDevice == 0 {
		if data, err := io.ReadAll(os.Stdin); err == nil && len(data) > 0 {
			text = string(data)
		}
	}

	counts := wordcount.CountWords(text)

	// Ordena as chaves para uma saída estável e legível (mapas em Go não têm
	// ordem garantida).
	words := make([]string, 0, len(counts))
	for w := range counts {
		words = append(words, w)
	}
	sort.Strings(words)

	for _, w := range words {
		fmt.Printf("%s: %d\n", w, counts[w])
	}
}
