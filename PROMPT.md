# Registro do uso de IA — `CountWords`

Este documento registra como a ferramenta de IA foi utilizada e avaliada
criticamente durante a construção da função `CountWords` e de seus testes.

## Ferramenta de IA utilizada

- **Ferramenta:** Claude Code (interface de linha de comando da Anthropic).
- **Modelo:** Claude Opus 4.8.

## Ambiente utilizado

- **Sistema operacional:** macOS (Darwin, Apple Silicon / `arm64`).
- **Linguagem/Toolchain:** Go `1.26.5` (instalado via Homebrew durante a sessão;
  o `go.mod` declara `go 1.21` como versão mínima).
- **Editor/terminal:** VS Code + terminal integrado.
- **Controle de versão:** Git + GitHub.

---

## Roteiro obrigatório (Etapas 1 a 6)

### Etapa 1 — Compreensão do problema

**Prompt utilizado (obrigatório):**

> Quero implementar em Go uma função `CountWords(text string) map[string]int`. A
> função deve contar a frequência das palavras em um texto, converter palavras
> para minúsculas, remover pontuação simples e ignorar palavras com menos de 3
> caracteres. Explique a especificação e os principais cuidados de implementação.

**Resumo da resposta da IA:**

A IA descreveu o fluxo esperado — normalizar (minúsculas), separar em palavras
removendo pontuação, filtrar por tamanho e acumular a contagem em um `map` — e
destacou os cuidados abaixo.

**Cuidados identificados:**

1. **Bytes × runes:** em Go, `len(s)` conta *bytes*, e caracteres acentuados
   (ex.: `é`) ocupam 2 bytes em UTF-8. A regra de "menos de 3 caracteres" deve
   contar *runes* (`utf8.RuneCountInString`), não bytes.
2. **Preservar acentos:** a normalização deve baixar para minúsculas **sem**
   remover acentos, pois o resultado esperado contém `árvore` e `útil`.
   `strings.ToLower` já trata acentos corretamente.
3. **Definir "pontuação simples":** decidir se a remoção seria por lista fixa de
   sinais (`. , ! ? ;`) ou por "tudo que não for letra/dígito é separador".
4. **Comparação de mapas em teste:** mapas não podem ser comparados com `==`;
   é preciso `reflect.DeepEqual` (ou `maps.Equal`).

**Decisões tomadas:**

- Contar **runes** para a regra de tamanho mínimo.
- Separar palavras com `strings.FieldsFunc`, tratando como separador qualquer
  rune que não seja letra nem dígito (`!unicode.IsLetter && !unicode.IsDigit`) —
  mais robusto do que uma lista fixa de pontuação.
- Preservar acentos, apoiando-se em `unicode.IsLetter` e `strings.ToLower`.

---

### Etapa 2 — Implementação da função

**Prompt utilizado (obrigatório):**

> Implemente em Go a função `CountWords(text string) map[string]int`. A função
> deve converter palavras para minúsculas, remover pontuação simples, ignorar
> palavras com menos de 3 caracteres e retornar um `map[string]int` com a
> frequência das palavras.

**Resposta da IA:** gerou a implementação presente em `wordcount.go`, usando
`strings.FieldsFunc`, `strings.ToLower` e `utf8.RuneCountInString`.

**Verificação (checklist do enunciado):**

| Verificação                          | Resultado                                             |
| ------------------------------------ | ----------------------------------------------------- |
| O código compila?                    | ✅ Sim (`go build` / `go vet` sem erros).             |
| A função tem a assinatura esperada?  | ✅ `func CountWords(text string) map[string]int`.     |
| A normalização está correta?         | ✅ Minúsculas + separação por não-letras/não-dígitos. |
| Palavras curtas são ignoradas?       | ✅ `< 3` runes são descartadas (`a`, `é`, `go`, `ia`).|
| Acentos são preservados?             | ✅ `árvore`, `útil` mantêm os acentos.                |

---

### Etapa 3 — Criação do teste automático

**Prompt utilizado (obrigatório):**

> Crie um teste automático em Go, usando o pacote padrão `testing`, para validar
> a função `CountWords` [com a entrada mínima do enunciado e o resultado
> esperado; o teste deve comparar o mapa produzido com o mapa esperado e falhar
> caso haja qualquer diferença; deve ser executado com `go test`].

**Teste gerado:** um teste `TestCountWords` com o caso mínimo do enunciado,
comparando o resultado com `reflect.DeepEqual` (ver `wordcount_test.go`).

**Resultado da execução:**

```
$ go test
ok  	github.com/danlisb/go_automated_tests
```

**Erros encontrados:** nenhum erro de compilação ou de lógica na versão final; o
caso mínimo passou. (Ponto de atenção discutido: a IA lembrou de usar o texto com
quebras de linha `\n` exatamente como no enunciado.)

**Correções realizadas:** nenhuma correção de bug foi necessária no caso mínimo.

---

### Etapa 4 — Revisão crítica do teste

**Prompt utilizado (obrigatório):**

> Revise o teste automático abaixo. Verifique se ele testa corretamente a
> conversão para minúsculas, a remoção de pontuação, o descarte de palavras com
> menos de 3 caracteres, a contagem de repetições e a comparação completa entre
> o mapa esperado e o mapa produzido. [teste colado]

**Resumo da resposta da IA / sugestões e decisão:**

| Sugestão da IA                                                            | Decisão   | Justificativa |
| ------------------------------------------------------------------------ | --------- | ------------- |
| Adotar formato *table-driven* com `t.Run` para vários casos.             | ✅ Aceita | Facilita adicionar os casos extras e isola falhas por subteste. |
| Manter `reflect.DeepEqual` para comparar mapas por completo.             | ✅ Aceita | Detecta chave a mais, chave faltando e contagem errada. |
| Mensagem de erro mostrando "obtido" × "esperado".                        | ✅ Aceita | Facilita o diagnóstico ao falhar. |
| Trocar `reflect.DeepEqual` por `maps.Equal` (Go 1.21+).                  | ❌ Rejeitada | `reflect.DeepEqual` é o padrão mais conhecido e não depende da versão do Go; ambos funcionam, optei pelo mais portável. |
| Remover acentos para "normalizar" as chaves.                            | ❌ Rejeitada | O enunciado **exige** preservar acentos (`árvore`, `útil`). |

**O teste realmente verifica a especificação?** Sim: o caso mínimo exercita
minúsculas, remoção de pontuação, descarte de palavras curtas, contagem de
repetições e comparação completa do mapa.

---

### Etapa 5 — Inclusão de novos casos de teste

**Prompt utilizado (obrigatório):**

> Sugira mais três casos de teste para a função `CountWords`. Os testes devem
> cobrir situações diferentes, como texto vazio, texto contendo apenas palavras
> curtas e texto com palavras repetidas com diferentes combinações de
> maiúsculas, minúsculas e pontuação.

**Casos sugeridos pela IA e implementados** (além do caso mínimo obrigatório —
total de 5 casos, portanto mais de dois adicionais):

1. `texto vazio` → mapa vazio.
2. `apenas palavras curtas` → todas com menos de 3 caracteres → mapa vazio.
3. `repeticoes com maiusculas e pontuacao` → `Bola/BOLA/bola`, `Gato/GATO`,
   `rede rede` → `{bola:3, gato:3, rede:2}`.
4. `palavras com numeros` → `html5`, `css3`, `go2`, `2024` mantidos; `42`
   (2 dígitos) ignorado → `{html5:2, css3:1, versão:1, 2024:2, go2:1}`.

Todos foram implementados em `wordcount_test.go`.

---

### Etapa 6 — Revisão final

**Prompt utilizado (obrigatório):**

> Analise o código completo da função `CountWords` e dos testes. Verifique se os
> testes são suficientes para validar a especificação e se há algum caso
> importante não coberto. [código completo colado]

**Sugestões finais da IA:**

- Rodar `go vet` e `gofmt` para garantir formatação/lint (feito; sem apontamentos).
- Documentar a limitação de que separadores quebram contrações/hifenizações.
- Considerar um caso com dígitos/números, já que a implementação aceita dígitos
  como parte de palavra. **(Implementado** — caso `palavras com numeros`.)

**Alterações feitas:**

- Comentários explicativos adicionados em `wordcount.go` e `wordcount_test.go`.
- Limitações documentadas no `README.md`.
- Caso de teste com números adicionado (`palavras com numeros`), cobrindo
  palavras alfanuméricas, número mantido (`2024`) e número curto ignorado (`42`).

**Limitações conhecidas:**

- Palavras com hífen ou apóstrofo são divididas em partes (ex.: `guarda-chuva`).
- Não há lista de *stopwords* além da regra de tamanho mínimo.

---

## Prompts adicionais utilizados

- *"`len` conta bytes ou runes em Go para strings com acento?"* — usado para
  confirmar o cuidado da Etapa 1 e justificar o uso de `utf8.RuneCountInString`.

## Síntese da avaliação crítica

| Item                                   | Situação |
| -------------------------------------- | -------- |
| Sugestões **aceitas**                  | Table-driven tests, `reflect.DeepEqual`, contagem por runes, `FieldsFunc` por não-letras, mensagens de erro claras, `go vet`/`gofmt`. |
| Sugestões **rejeitadas**               | Trocar por `maps.Equal` (portabilidade); remover acentos (contraria o enunciado); usar `len` em bytes para o tamanho mínimo. |
| **Erros produzidos pela IA**           | Nenhum erro de compilação/lógica na versão final. Risco potencial (contagem por bytes) foi identificado e evitado preventivamente na Etapa 1. |
| **Correções realizadas pelo estudante** | Revisão dos prompts e das saídas, escolha por `reflect.DeepEqual`, preservação de acentos, verificação manual do resultado com `go test -v`. |
