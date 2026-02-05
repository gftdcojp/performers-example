# Worker Service Example

バックグラウンドワーカーの例です。Pub/Sub でイベントを受信し、処理を行います。

## 構成

```
worker-service/
├── gftd.json          # GFTD 設定ファイル
├── go.mod             # Go モジュール
├── cmd/
│   └── worker/
│       └── main.go    # エントリーポイント
├── internal/
│   └── handler/
│       └── events.go  # イベントハンドラ
└── Dockerfile         # (自動生成可能)
```

## gftd.json のポイント

```json
{
  "name": "example-worker",
  "nanoid": "wrk3456a",
  "type": "worker",
  "dapr": {
    "appPort": 8080,
    "appProtocol": "http",
    "pubsub": "shared"
  },
  "dependencies": [
    {
      "name": "api-backend",
      "type": "service",
      "appId": "api5678y",        // 依存先の nanoid
      "required": true
    },
    {
      "name": "user-events",
      "type": "pubsub",
      "topic": "user-events",
      "required": true
    },
    {
      "name": "job-state",
      "type": "state",
      "stateStore": "shared-statestore",
      "required": true
    }
  ]
}
```

## 依存関係の種類

### service (サービス依存)
他の Performer サービスへの依存。Dapr service invocation で呼び出し可能。

### pubsub (Pub/Sub 依存)
特定のトピックを購読/発行する依存。

### state (State Store 依存)
State Store への依存。データの永続化に使用。

## デプロイ

```bash
cd examples/worker-service

# 依存関係を確認
gftd deps check

# デプロイ
gftd deploy
```

## 依存関係確認の出力例

```bash
$ gftd deps check

● Checking dependencies for example-worker

  ✓ api-backend (service) [required] - ready
  ✓ user-events (pubsub) [required] - pubsub component (assumed available)
  ✓ job-state (state) [required] - state store component (assumed available)

✓ All required dependencies are ready
```

## Pub/Sub サブスクリプション

Dapr の Pub/Sub を使用してイベントを購読:

```go
// サブスクリプションの登録
server.AddTopicEventHandler(&common.Subscription{
    PubsubName: "shared-pubsub",
    Topic:      "user-events",
    Route:      "/events/user",
}, handleUserEvent)

func handleUserEvent(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
    log.Printf("Received event: %v", e.Data)
    // 処理ロジック
    return false, nil
}
```
