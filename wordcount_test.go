package wordcount

import (
	"reflect"
	"testing"
)

// TestCountWords cobre o caso mínimo obrigatório e casos adicionais.
//
// Cada caso compara o mapa produzido por CountWords com o mapa esperado
// usando reflect.DeepEqual, de modo que qualquer diferença (chave faltando,
// chave a mais ou contagem errada) faz o teste falhar.
func TestCountWords(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  map[string]int
	}{
		{
			// Caso mínimo obrigatório definido no enunciado.
			name: "caso minimo do enunciado",
			input: "Casa, casa! A casa é azul.\n" +
				"Árvore; árvore? verde.\n" +
				"Go go Go. IA é útil, mas IA erra.",
			want: map[string]int{
				"casa":   3,
				"árvore": 2,
				"azul":   1,
				"verde":  1,
				"útil":   1,
				"mas":    1,
				"erra":   1,
			},
		},
		{
			// Texto vazio deve produzir um mapa vazio.
			name:  "texto vazio",
			input: "",
			want:  map[string]int{},
		},
		{
			// Apenas palavras curtas (< 3 caracteres) ou pontuação:
			// todas devem ser ignoradas, resultando em mapa vazio.
			name:  "apenas palavras curtas",
			input: "A é o de vc, go! ia? um. eu tu.",
			want:  map[string]int{},
		},
		{
			// Mesma palavra repetida com diferentes combinações de
			// maiúsculas/minúsculas e pontuação deve ser normalizada e somada.
			name:  "repeticoes com maiusculas e pontuacao",
			input: "Bola, BOLA! bola. Gato; gato, GATO? rede rede.",
			want: map[string]int{
				"bola": 3,
				"gato": 3,
				"rede": 2,
			},
		},
		{
			// Números são tratados como parte da palavra:
			//   - palavras alfanuméricas (html5, css3, go2) são mantidas;
			//   - um número puro com 3+ dígitos (2024) é contado;
			//   - um número curto (42) é ignorado pela regra de tamanho mínimo;
			//   - acentos continuam preservados (versão).
			name:  "palavras com numeros",
			input: "html5 html5 css3. Versão 2024, 2024! 42 é ok, go2.",
			want: map[string]int{
				"html5":  2,
				"css3":   1,
				"versão": 1,
				"2024":   2,
				"go2":    1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CountWords(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CountWords(%q)\n  obtido:   %v\n  esperado: %v", tt.input, got, tt.want)
			}
		})
	}
}
