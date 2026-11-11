# Escopo — NOC FUT POC

## Objetivo

Validar uma central local e leve de futebol, focada no Brasileirão e na seleção de um clube, como o Vasco. Aproveitar Go e Playwright para automatizar a validação das APIs e da interface.

## Entregue

Servidor HTTP local em Go, interface responsiva sem framework, amostra embutida de oito partidas e doze clubes, filtro inicial no Vasco, filtros combinados de time/situação/busca, favoritos locais e detalhes da partida com escalações dos dois times, estatísticas e lances.

Os dados esportivos são fictícios, estáticos e claramente identificados. As escalações contêm onze jogadores ilustrativos por lado. O botão de recarga relê os dados locais.

## Arquitetura

Navegador → servidor Go em loopback → dados JSON embutidos.

Playwright Go → APIRequestContext → servidor temporário nos testes de integração.

Chromium headless → dashboard → validação de filtros, favoritos, abas e layout mobile.

## Critérios de aceite

- Iniciar com um comando ou executar o binário compilado.
- Abrir a aplicação com filtro inicial do Vasco e três resultados.
- Exibir onze jogadores e permitir alternar entre mandante e visitante.
- Combinar filtro de clube, situação e busca, com estado vazio claro.
- Salvar e remover favoritos no navegador.
- Tratar rotas inexistentes e filtros inválidos na API.
- Funcionar sem internet depois de compilado.
- Manter testes unitários e testes HTTP automatizados com Playwright.

## Fora do escopo desta POC

Placares reais, elencos atuais, integração com API externa, notificações, banco de dados, autenticação, análises de apostas, CI/CD e hospedagem externa. Publicação do código-fonte no GitHub não hospeda a aplicação.

## Próxima etapa

Definir provedor e cobertura do Brasileirão; adicionar coleta HTTP no servidor, cache, tratamento de cotas e informação explícita sobre atualização e disponibilidade das escalações.
