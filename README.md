GoPoloniexTest
==========

Задача получить данные о состоявшихся сделках с биржи Poloniex с помощью Websocket

API - https://docs.poloniex.com/#price-aggregated-book

[<img alt="Goland" src="https://img.shields.io/badge/Go-00ADD8?style=float&logo=go&logoColor=white" />](https://golang.org/)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/DFofanov/ConvJudgesToCSV)
[![GitHub license](https://img.shields.io/github/license/DFofanov/ConvJudgesToCSV)](https://github.com/DFofanov/ConvJudgesToCSV/blob/main/LICENSE)

## Описание
На входе словарь в виде json-строки:
~~~json 
{"poloniex":["BTC_USDT", "TRX_USDT", "ETH_USDT"]}
~~~
Надо подписаться на данные пары(в их API они перевёрнутые, например USDT_BTC)

На выходе(лог в консоль) должна быть структура:
~~~text
type RecentTrade struct {
    Id        string    // ID транзакции
    Pair      string    // Торговая пара (из списка выше)
    Price     float64   // Цена транзакции
    Amount    float64   // Объём транзакции
    Side      string    // Как биржа засчитала эту сделку (как buy или как sell)
    Timestamp time.Time // Время транзакции
}
~~~

## Запуск
~~~cmd
GoPoloniexTest key secret
~~~

Команды:
* -h (помощь)
* -v (версия программы)

![image](https://github.com/DFofanov/GoPoloniexTest/blob/main/img/output.png?raw=true)

## License
Licensed under the GPL-3.0 License.