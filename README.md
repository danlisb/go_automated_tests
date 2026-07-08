# Teste Automático em Go — Contagem de Palavras (`CountWords`)

Atividade individual: construir, com apoio de uma ferramenta de IA, uma função
em Go que conta a frequência de palavras em um texto e validá-la com testes
automáticos usando o pacote padrão `testing`.

O foco da atividade **não** é concorrência, e sim o uso crítico de IA para
compreender a especificação, gerar a função, criar os testes, revisá-los e
verificar se eles realmente cobrem os casos relevantes. O registro do uso da IA
está em [`PROMPT.md`](./PROMPT.md).

## Estrutura do repositório

| Arquivo               | Descrição                                        |
| --------------------- | ------------------------------------------------ |
| `go.mod`              | Definição do módulo Go.                          |
| `wordcount.go`        | Implementação da função `CountWords`.            |
| `wordcount_test.go`   | Testes automáticos (`go test`).                  |
| `README.md`           | Este arquivo.                                    |
| `PROMPT.md`           | Registro do uso da IA durante a atividade.       |

## Como executar os testes

Pré-requisito: [Go](https://go.dev/dl/) instalado (versão 1.21 ou superior).

```bash
# a partir da raiz do repositório
go test          # execução simples
go test -v       # execução detalhada, mostrando cada subteste
```

## A função `CountWords`

```go
func CountWords(text string) map[string]int
```

Recebe um texto e retorna um mapa `palavra -> frequência`, aplicando as regras:

1. **Minúsculas** — todas as palavras são convertidas para minúsculas
   (`strings.ToLower`), de forma que `Casa`, `casa` e `CASA` contam como a mesma
   palavra.
2. **Remoção de pontuação simples** — o texto é separado em palavras com
   `strings.FieldsFunc`, tratando como separador qualquer caractere que **não**
   seja letra nem dígito. Assim, `,`, `.`, `!`, `?`, `;`, etc. são descartados.
3. **Descarte de palavras curtas** — palavras com **menos de 3 caracteres** são
   ignoradas. O tamanho é medido em *runes* (`utf8.RuneCountInString`), e não em
   bytes, para que acentos sejam contados corretamente.
4. **Preservação de acentos** — letras acentuadas (`á`, `é`, `í`, `ó`, `ú`, `ã`,
   `ç`, ...) são preservadas, pois `unicode.IsLetter` as reconhece como letras.
5. **Contagem** — cada palavra restante incrementa sua entrada no mapa.

### Por que *runes* e não *bytes*?

Em Go, `len("é")` retorna `2` (bytes UTF-8), o que poderia atrapalhar a regra de
"menos de 3 caracteres" para textos acentuados. Usar
`utf8.RuneCountInString` garante a contagem correta de caracteres.

## Casos de teste implementados

Os testes são *table-driven*: cada caso compara o mapa produzido pela função com
o mapa esperado usando `reflect.DeepEqual`, de modo que **qualquer** diferença
(chave faltando, chave a mais ou contagem errada) faz o teste falhar.

| Caso                                     | O que valida                                                              |
| ---------------------------------------- | ------------------------------------------------------------------------- |
| `caso minimo do enunciado`               | Caso mínimo obrigatório do enunciado (minúsculas, pontuação, acentos, repetições e descarte de `a`, `é`, `go`, `ia`). |
| `texto vazio`                            | Texto vazio deve produzir um mapa vazio.                                  |
| `apenas palavras curtas`                 | Texto só com palavras de menos de 3 caracteres deve produzir mapa vazio.  |
| `repeticoes com maiusculas e pontuacao`  | Mesma palavra em variações de maiúsculas/minúsculas e pontuação é normalizada e somada. |

O caso mínimo obrigatório usa o texto:

```
Casa, casa! A casa é azul.
Árvore; árvore? verde.
Go go Go. IA é útil, mas IA erra.
```

e verifica o resultado esperado:

```
casa: 3
árvore: 2
azul: 1
verde: 1
útil: 1
mas: 1
erra: 1
```

As palavras `a`, `é`, `go` e `ia` são ignoradas por terem menos de 3 caracteres.

## Resultado obtido ao executar `go test`

```
$ go test -v ./...
=== RUN   TestCountWords
=== RUN   TestCountWords/caso_minimo_do_enunciado
=== RUN   TestCountWords/texto_vazio
=== RUN   TestCountWords/apenas_palavras_curtas
=== RUN   TestCountWords/repeticoes_com_maiusculas_e_pontuacao
--- PASS: TestCountWords (0.00s)
    --- PASS: TestCountWords/caso_minimo_do_enunciado (0.00s)
    --- PASS: TestCountWords/texto_vazio (0.00s)
    --- PASS: TestCountWords/apenas_palavras_curtas (0.00s)
    --- PASS: TestCountWords/repeticoes_com_maiusculas_e_pontuacao (0.00s)
PASS
ok  	github.com/danlisb/go_automated_tests	0.379s
```

Execução simples:

```
$ go test
ok  	github.com/danlisb/go_automated_tests
```

## Limitações conhecidas

- A separação de palavras quebra em qualquer caractere que não seja letra ou
  dígito. Assim, contrações e hifenizações (ex.: `guarda-chuva`) são divididas em
  partes. Isso é aceitável para o escopo desta atividade ("pontuação simples").
- Não há tratamento de *stopwords* além da regra de tamanho mínimo.
