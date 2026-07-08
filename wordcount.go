// Package wordcount fornece uma função para contar a frequência de palavras
// em um texto, seguindo regras simples de normalização.
package wordcount

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// minWordLen é o número mínimo de caracteres (runes) para uma palavra ser contada.
const minWordLen = 3

// CountWords recebe um texto e retorna um mapa com a frequência das palavras.
//
// Regras aplicadas:
//   - converte as palavras para minúsculas;
//   - remove pontuação simples (qualquer caractere que não seja letra ou dígito
//     funciona como separador de palavras);
//   - ignora palavras com menos de 3 caracteres (contados como runes, para
//     tratar corretamente acentos);
//   - preserva acentos (á, é, í, ó, ú, ã, õ, ç, ...).
func CountWords(text string) map[string]int {
	counts := make(map[string]int)

	// FieldsFunc separa o texto em "palavras": qualquer rune que não seja
	// letra nem dígito é tratado como separador, o que remove a pontuação.
	words := strings.FieldsFunc(text, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})

	for _, word := range words {
		word = strings.ToLower(word)

		// Contamos runes (não bytes) para que "árvore", "útil" etc. tenham o
		// tamanho correto e a regra de 3 caracteres funcione com acentos.
		if utf8.RuneCountInString(word) < minWordLen {
			continue
		}

		counts[word]++
	}

	return counts
}
