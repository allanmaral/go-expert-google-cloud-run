# Desafio Prático - Google Cloud Run

Bem-vindo ao Desafio Prático de Google Cloud Run da Pós-Graduação Go Expert! Este projeto consiste na implementação de uma API REST de clima que retorna a temperatura atual dado um CEP. A API foi publicada usando o Google Cloud Run e está disponível [neste link](https://go-expert-google-cloud-run-weather-api-3vpr3ze2ra-uc.a.run.app/api/weather/70150900).

## Pré-requisitos

Antes de começar, certifique-se de ter instalado os seguintes requisitos:

- [Go SDK](https://golang.org/dl/): Linguagem de programação Go.
- [Docker](https://docs.docker.com/get-docker/): Plataforma de conteinerização.
- [Make](https://www.gnu.org/software/make/): Utilizado para automatização de tarefas.
- [Weather API Key](https://www.weatherapi.com/): Uma chave gratuita da Weather API.

## Executando o Projeto

1. Clone este repositório em sua máquina local:

   ```bash
   git clone https://github.com/allanmaral/go-expert-google-cloud-run.git
   ```

1. Navegue até o diretório do projeto:

   ```bash
   cd go-expert-google-cloud-run
   ```

1. Duplique o arquivo `.env.example` e renomeie para `.env` e preencha o valor `WEATHER_APIKEY` com a sua chave da [Weather API](https://www.weatherapi.com/):

   ```env
   HOST=0.0.0.0
   PORT=8080
   WEATHER_APIKEY=<WEATHER_API_SECRET_KEY>
   ```

1. Execute o seguinte comando para subir a API usando o docker compose:

   ```bash
   docker compose up -d
   ```

## Acesso aos Serviços

Após subir o serviço, você poderá acessar a API no endereço [http://localhost:8080/api/weather/70150900](http://localhost:8080/api/weather/70150900). Subistitua o valor `70150900` pelo CEP desejado.

## Documentação da API REST

A documentação das rotas do servidor HTTP está disponível na pasta `./api`. Os arquivos `.http` contem uma variável `base` url que pode ser alterada para apontar para a versão publicada no Google Cloud Run.

### Testes

A API foi desenvolvida com testes de unidade e integração. Os testes de unidade rodam completamente isolado, mas os testes de integração dependem da Weather API. Certifique-se de preencher a variável de ambiente no arquivo `.env` antes de rodar os testes.

Para rodar os testes da aplicação basta rodar o comando:

```bash
make test
```
