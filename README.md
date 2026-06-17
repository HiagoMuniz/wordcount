# Conta-palavras Concorrente em Go

Este projeto foi desenvolvido como parte de uma atividade prática de programação concorrente em Go. Ele implementa e compara duas abordagens para contar a frequência de palavras em um arquivo de texto: uma versão **sequencial** e uma versão **concorrente**.

---

## Como Executar o Programa

Certifique-se de ter o Go instalado em sua máquina.

1. Baixe as dependências e compile/execute o projeto informando o caminho do dataset:
   ```bash
   go run . AChristmasCarol_CharlesDickens_English.txt
   ```

2. Se desejar configurar o número de workers da versão concorrente (o padrão é o número de CPUs da máquina `runtime.NumCPU()`), passe o número como um argumento adicional:
   ```bash
   go run . AChristmasCarol_CharlesDickens_English.txt 4
   ```

---

## Detalhes da Implementação

### 1. Estratégia Concorrente
Foi utilizada uma arquitetura baseada no padrão **Fan-Out / Worker Pool** com canais em Go. Esta abordagem evita contenção de mutexes em um mapa compartilhado ao dar a cada worker sua própria estrutura de dados para contagem.

### 2. Divisão do Trabalho
* **Producer (Produtor):** Uma goroutine principal lê o arquivo linha por linha usando um `bufio.Scanner` e envia cada linha como uma string para um canal buffered (`linesChan`).
* **Workers (Trabalhadores):** Um pool de `N` goroutines paralelas lêem linhas de forma concorrente a partir do `linesChan`. Cada worker limpa as palavras (conversão para minúsculas, remoção de caracteres de pontuação e símbolos) e as adiciona em um mapa local (`map[string]int`).

### 3. Combinação de Resultados Parciais (Merge)
* Assim que o produtor termina de ler o arquivo, o canal `linesChan` é fechado.
* Cada worker termina de processar suas linhas restantes, envia seu mapa de frequências parcial (`localFreqs`) para o canal `partialMapsChan`, e sinaliza sua conclusão em um grupo de sincronização (`sync.WaitGroup`).
* Uma goroutine gerenciadora aguarda todos os workers terminarem (`wg.Wait()`) e então fecha o `partialMapsChan`.
* A goroutine principal consome os mapas parciais de `partialMapsChan` e realiza a junção (merge) somando os valores em um mapa de frequências global final.

### 4. Verificação de Correção
A correção da versão concorrente é verificada de forma automatizada ao final da execução:
* O programa compara a quantidade de chaves nos dois mapas (`len(seqMap) == len(concMap)`).
* Para cada palavra presente na versão sequencial, valida se ela existe na versão concorrente e possui exatamente a mesma contagem de ocorrências.
* O resultado da comparação é exibido como `Resultados iguais: sim/não`.

### 5. Medição de Tempo
A medição de tempo é feita isolando o início e fim da execução de cada algoritmo usando `time.Now()` e `time.Since()`:
* **Tempo Sequencial:** Mede todo o processo de abertura de arquivo, leitura de linha por linha, tokenização, limpeza e montagem do mapa final.
* **Tempo Concorrente:** Mede o processo equivalente, incluindo abertura do arquivo, inicialização dos canais, concorrência (leitura, workers, merge dos mapas) e fechamento dos recursos.

---

## Resultados e Análise de Desempenho

Testes realizados com o arquivo `AChristmasCarol_CharlesDickens_English.txt` (tamanho de ~159 KB) em um processador moderno:

* **Versão Sequencial:** ~7.9ms a 10.8ms
* **Versão Concorrente (1 worker):** ~8.9ms
* **Versão Concorrente (4 workers):** ~6.7ms
* **Versão Concorrente (8 workers):** ~5.8ms (Melhor desempenho, redução de ~40% do tempo)
* **Versão Concorrente (16 workers):** ~6.9ms

### Análise
* **Ganhos:** A versão concorrente demonstrou ser mais rápida com 4 e 8 workers em relação à versão sequencial.
* **Overhead:** Com 1 worker, a concorrência é ligeiramente mais lenta que a sequencial por conta do overhead de criação de canais e agendamento de goroutines. De forma semelhante, com 16 workers, o overhead de context-switch de muitas goroutines para processar um arquivo pequeno (~159 KB) superou a vantagem da computação paralela, tornando o tempo levemente maior do que com 8 workers.

---

## Dificuldades Encontradas
1. **Contenção de escrita em mapas:** A ideia inicial de usar um mapa compartilhado protegido por um Mutex foi descartada, pois a contenção gerada pelas goroutines escrevendo no mesmo mapa anularia os ganhos de concorrência. A estratégia de usar mapas parciais em cada worker e fazer o merge no final provou-se consideravelmente mais eficiente.
2. **Tratamento de Punctuation e Unicode:** Algumas palavras no texto continham caracteres especiais como aspas inteligentes ou travessões. A utilização de `unicode.IsPunct` e `unicode.IsSymbol` no pacote `text.go` garantiu a limpeza robusta para qualquer caractere de pontuação Unicode.

---

## Relatório de Uso de IA

* **Ferramenta de IA:** Antigravity (Google DeepMind)
* **Modelo utilizado:** Claude Sonnet 4.6 (Thinking) / Gemini 3.5 Flash (High)
* **Ambiente de Desenvolvimento:** Windows (PowerShell) / Go 1.22+
* **Passos Vencidos com Auxílio da IA:**
  1. **Planejamento:** Elaboração e estruturação das etapas do plano de implementação.
  2. **Codificação:** Geração e refinamento das lógicas de limpeza de texto, leitura concorrente/sequencial e comparação de dados.
  3. **Validação:** Testes de desempenho variando a quantidade de workers para identificar gargalos e o ponto ideal de concorrência.
  4. **Documentação:** Estruturação deste documento README explicando a estratégia de Fan-Out e Worker Pools.
