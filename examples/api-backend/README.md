# API Backend Example

Go ベースの API バックエンドの例です。

## 構成

```
api-backend/
├── gftd.json          # GFTD 設定ファイル
├── go.mod             # Go モジュール
├── go.sum
├── cmd/
│   └── server/
│       └── main.go    # エントリーポイント
├── internal/
│   ├── handler/       # HTTP ハンドラ
│   └── service/       # ビジネスロジック
└── Dockerfile         # (自動生成可能)
```

## gftd.json のポイント

```json
{
  "name": "example-api",
  "nanoid": "api5678y",
  "type": "web",
  "build": {
    "backend": "kaniko",
    "mainPkg": "./cmd/server"    // Go メインパッケージ
  },
  "dapr": {
    "appPort": 8080,
    "appProtocol": "http"
  },
  "dependencies": [
    {
      "name": "user-actor",
      "type": "actor",
      "appId": "act9012z",        // 依存先の nanoid
      "actorType": "UserActor",
      "required": false
    }
  ]
}
```

## デプロイ

```bash
cd examples/api-backend
gftd deploy --dapr
```

## 依存関係の確認

```bash
# 依存関係が存在するか確認
gftd deps check

# 依存関係ツリーを表示
gftd deps tree
```

## Dapr Service Invocation

このAPIは Dapr 経由で他のサービスを呼び出せます:

```go
// Dapr クライアントで Actor を呼び出す
client, _ := dapr.NewClient()
resp, _ := client.InvokeActor(ctx, &dapr.InvokeActorRequest{
    ActorType: "UserActor",
    ActorID:   userID,
    Method:    "GetProfile",
})
```
