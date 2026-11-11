# NOC FUT

Dashboard local do Brasileirão em Go, com foco em partidas e escalações. POC construída sobre o esqueleto original de automação de APIs.

**Todos os jogos, placares, estatísticas e jogadores são fictícios.** Os clubes e estádios servem apenas como contexto visual. Não há consulta a resultados atuais nem atualização de placares ao vivo.

## Iniciar

Requer Go 1.22 ou superior para compilar:

```powershell
go run ./cmd/noc_fut
```

Abra http://127.0.0.1:8080. O servidor escuta apenas no computador local. Para outra porta:

```powershell
go run ./cmd/noc_fut -addr 127.0.0.1:8081
```

Também é possível executar `./iniciar.ps1`. Para gerar um executável independente:

```powershell
go build -o noc_fut.exe ./cmd/noc_fut
```

Os arquivos da interface e os dados estão embutidos no binário. Não requer Node, banco de dados, CDN, conta ou chave de API. Após compilar, funciona sem internet. Encerre com Ctrl+C no terminal.

## O que explorar

- Filtro inicial no Vasco, com três partidas demonstrativas.
- Oito jogos e doze clubes na amostra.
- Filtro por time, situação e busca por confronto.
- Detalhes da partida, formação 4-3-3 e onze jogadores para cada lado.
- Abas de estatísticas e lances.
- Favoritos persistidos no navegador com localStorage.
- Interface escura e responsiva, sem fontes ou imagens externas.
- Recarregar relê a amostra local; não simula alterações de placar.

## Testes

Validação do servidor, filtros e rotas, sem dependências de navegador:

```powershell
go test ./...
go vet ./...
```

Prepare Playwright uma vez (precisa de internet para baixar o driver e o Chromium headless):

```powershell
./scripts/setup-playwright.ps1
```

Testes de requisições HTTP pelo APIRequestContext do Playwright Go:

```powershell
go test -tags integration ./tests/api -v
```

Validação da interface no Chromium e geração de capturas desktop/mobile em `reports/`:

```powershell
go test -tags "integration,ui" ./tests/api -v
```

Os testes iniciam um servidor temporário e não dependem de uma instância aberta na porta 8080.

## API local

| Rota GET | Retorno |
| --- | --- |
| `/api/health` | Estado do servidor e modo demo |
| `/api/teams` | Clubes da amostra |
| `/api/matches` | Todos os jogos demonstrativos |
| `/api/matches?team=vasco` | Jogos do Vasco |
| `/api/matches?team=vasco&status=live` | Jogos do Vasco em andamento na amostra |
| `/api/matches/1` | Detalhes, escalações, estatísticas e eventos |

Situações válidas: `live`, `scheduled`, `finished`. `team` aceita identificador ou trecho do nome, sem distinguir maiúsculas. Métodos diferentes de GET retornam 405; detalhes inexistentes, 404; situação inválida, 400.

## Estrutura

```text
cmd/noc_fut/                  Entrada do servidor local
internal/dashboard/           API, dados embutidos e testes do servidor
internal/dashboard/web/       HTML, CSS e JavaScript sem framework
internal/dashboard/fixtures.json  Amostra fictícia
internal/config/              Reserva do esqueleto para integração futura
internal/client/              Reserva do esqueleto para integração futura
internal/services/            Reserva do esqueleto para integração futura
tests/api/                    Testes HTTP e de interface com Playwright Go
scripts/setup-playwright.ps1   Instalação do runtime de testes
docs/ESCOPO.md                Limites e evolução da POC
reports/                      Capturas geradas pelos testes (fora do Git)
```

## Evolução para dados reais

Adicionar um provedor de futebol no servidor Go, cache, orçamento de requisições e estados de escalação indisponível. Manter tokens somente no servidor; o `.env.example` é um modelo futuro e não é carregado automaticamente. Validar cobertura e atraso de cada campeonato antes de habilitar monitoramento real.

O Playwright automatiza os testes da aplicação; a futura coleta contínua poderá usar o cliente HTTP nativo do Go.

Referência: [Playwright Go](https://github.com/mxschmitt/playwright-go), versão fixada no go.mod e no script de instalação.
