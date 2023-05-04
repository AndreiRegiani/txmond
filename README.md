# txmond

Proof of Concept: blockchain transactions _(tx)_ monitoring _(mon)_ daemon _(d)_.

This service monitors subscribed wallets to watch for new transactions on every new block. `txmond` exposes a REST API (add wallet, get transactions) and has a optional webhook to pass the event as soon as possible (e.g. into another service delivering the push notifications).

## Supported blockchains

* Ethereum

## REST endpoints

```bash
GET /v1/ethereum/block/current

POST /v1/ethereum/wallet { "address": "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B" }

GET /v1/ethereum/wallet?address=0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B
```

## Supported storages

* `temp` (no persistance, wallet list/transaction history are lost when the process closes)
* `redis` @TODO

## Configuration

See `.env.dist`

## Webhook

JSON Payload send to `TXMOND_WEBHOOK_URL=http://localhost:3000/txmond`

```json
{
    "block": {
        "number": 1234567890
    },

    "transactions": [
        {
            "from": "0x1234567890",
            "to": "0x1234567890",
            "value": "100",
        }
    ]
}
```
