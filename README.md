# Weather by CEP (Go Expert)

API em Go que recebe um CEP brasileiro, identifica a cidade via ViaCEP e retorna a temperatura atual em Celsius, Fahrenheit e Kelvin via WeatherAPI.

## URL do Cloud Run

>
> `https://SERVICO-REGION-PROJECT.a.run.app`

Exemplo de requisição:

```bash
curl https://SERVICO-REGION-PROJECT.a.run.app/01310100
```

Resposta de sucesso (HTTP 200):

```json
{
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.5
}
```

## Contrato da API

| Cenário | HTTP | Resposta |
|---------|------|----------|
| Sucesso | 200 | JSON com `temp_C`, `temp_F`, `temp_K` |
| CEP inválido | 422 | `invalid zipcode` |
| CEP não encontrado | 404 | `can not find zipcode` |

Endpoint: `GET /{cep}` — o CEP deve conter exatamente 8 dígitos numéricos.

## Rodar via Docker

Build:

```bash
docker build -t weather-cep .
```

Execução:

```bash
docker run --rm -p 8080:8080 -e WEATHER_API_KEY=sua_chave_aqui weather-cep
```

## Fórmulas de conversão

- Fahrenheit: `F = C * 1.8 + 32`
- Kelvin: `K = C + 273`
