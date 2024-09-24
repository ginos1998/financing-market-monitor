# financing-market-monitor

## Descripción

Este proyecto es un sistema de monitoreo de mercados financieros. Utiliza la API de [Finnhub](https://finnhub.io/docs/api/introduction)
para consumir en tiempo real datos del mercado crypto y estadounidense. Al mismo tiempo, realiza análisis técnico sobre medias móviles
de los tickers configurados para llevarles seguimiento. Se pueden configurar alertas para recibir notificaciones en caso de que se cumplan. 
Estas alertas pueden ser de dos tipos:
- **Precio**: Se envía una notificación cuando el precio de un ticker supera, baja o cruza un valor configurado.
- **Media Móvil**: Se envía una notificación cuando el precio de un ticker supera o baja de una media móvil configurada.

Las notificaciones se envían a través de webhooks de Discord al server y canal configurados.

Se utilizaron las siguientes tecnologías:
- **Go 1.22.5**: Lenguaje de programación principal.
- **Docker**: Para la creación de contenedores.
- **MongoDB**: Base de datos para almacenar la configuración de los tickers, alertas, apis, entr otros.
- **Redis**: Base de datos en memoria para almacenar los datos de las medias móviles y los precios intra-diarios de los tickers.
- **Kafka**: Para la comunicación entre los servicios.
- **Zookeeper**: Para la administración de los brokers de Kafka.

## Estructura

El proyecto está dividido en los siguientes servicios:
- **data-ingest**: Servicio que consume a través de websocket los datos de la API de Finnhub y los envía a Kafka. 
También actualiza los precios diarios de los tickers utilizando la API de [Nasdaq](https://www.nasdaq.com/solutions/data-link-api).
*Si es la primera vez que se configura el proyecto*, se deberían procesar los datos de los csv de la carpeta `resources`.
Para ello, se debe ejecutar el comando `go run cmd/import_data/import_data.go` en la carpeta `data-ingest`.
- **data-processing**: Servicio que consume los datos de Kafka y los almacena en Redis. Además, calcula las medias móviles 
de los tickers configurados y procesa las alertas. En caso de que se cumpla una alerta, envía una notificación 
a través de un webhook de Discord. 

<p align="center">
  <img src="https://github.com/ginos1998/financing-market-monitor/blob/master/fmm-arq-2024-09-23-2208.png" alt="fmm-arq">
</p>

## Configuración

Para configurar el proyecto, se deben crear los siguientes archivos que contienen las variables de entorno 
para cada servicio:
- **data-ingest**: `.env.ingest` en la carpeta `data-ingest`.

```shell
# mongoDB
MONGO_USER =
MONGO_PASSWORD = 
MONGO_HOST =
MONGO_PORT =
MONGO_DB = 
MONGO_AUTH_SOURCE = 

# kafka
KAFKA_SERVER = 
KAFKA_TOPIC_TIME_SERIES_DATA = 
KAFKA_TOPIC_STREAM_STOCK_MARKET_DATA = 

FINNHUB_TOKEN = 
ALPHAVANTAGE_URI = 
ALPHAVANTAGE_API_KEY = 
YAHOO_FINANCE_URL = 

NASDAQ_API_URL = 
```

- **data-processing**: `.env.processing` en la carpeta `data-processing`.

```shell
# mongoDB
MONGO_USER = 
MONGO_PASSWORD = 
MONGO_HOST = 
MONGO_PORT = 
MONGO_DB = 
MONGO_AUTH_SOURCE = 

# redis
REDIS_HOST = 
REDIS_PORT = 
REDIS_DB = 
REDIS_USERNAME = 
REDIS_PASSWORD = 
```

## Ejecución

El proyecto se puede ejecutar utilizando `docker-compose`. Para ello, se debe ejecutar el siguiente comando en la raíz del proyecto:

```shell
docker-compose up
```

Si se desea ejecutar los servicios de forma individual y local, se puede ingresar a la carpeta correspondiente utilizar el siguiente comando:

```shell
go mod tidy

go run cmd/main.go
```

> [!NOTE]
> Es posible que sea necesario crear los topics de Kafka. Para ello, se puede utilizar el comando que se encuentra
> en el archivo `topics.txt` en la raíz del proyecto.


