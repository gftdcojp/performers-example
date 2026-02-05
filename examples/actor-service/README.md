# Actor Service Example

Dapr Virtual Actors を使用したサービスの例です。

## 構成

```
actor-service/
├── gftd.json          # GFTD 設定ファイル
├── go.mod             # Go モジュール
├── cmd/
│   └── actors/
│       └── main.go    # エントリーポイント
├── internal/
│   └── actors/
│       ├── user.go    # UserActor
│       ├── session.go # SessionActor
│       └── cart.go    # CartActor
└── Dockerfile         # (自動生成可能)
```

## gftd.json のポイント

```json
{
  "name": "example-actors",
  "nanoid": "act9012z",
  "type": "actor",
  "dapr": {
    "appPort": 50051,
    "appProtocol": "grpc",      // Actor は gRPC 推奨
    "actors": {
      "types": ["UserActor", "SessionActor", "CartActor"],
      "idleTimeout": "1h",       // アイドル状態で破棄されるまでの時間
      "drainTimeout": "30s",     // シャットダウン時の猶予時間
      "reentrancyMax": 32,       // 再入可能な最大深度
      "remindersStoragePartitions": 7  // リマインダーパーティション数
    }
  }
}
```

## Actor 実装例

```go
package actors

import (
    "context"
    "github.com/dapr/go-sdk/actor"
)

type UserActor struct {
    actor.ServerImplBase
    ID       string
    Name     string
    Email    string
}

func (u *UserActor) Type() string {
    return "UserActor"
}

func (u *UserActor) GetProfile(ctx context.Context) (*UserProfile, error) {
    return &UserProfile{
        ID:    u.ID,
        Name:  u.Name,
        Email: u.Email,
    }, nil
}

func (u *UserActor) UpdateProfile(ctx context.Context, req *UpdateProfileRequest) error {
    u.Name = req.Name
    u.Email = req.Email
    return u.SaveState()
}
```

## デプロイ

```bash
cd examples/actor-service
gftd deploy
```

## Actor の呼び出し

他のサービスから Dapr 経由で呼び出せます:

```go
// Dapr クライアントで Actor を呼び出す
client, _ := dapr.NewClient()
resp, _ := client.InvokeActor(ctx, &dapr.InvokeActorRequest{
    ActorType: "UserActor",
    ActorID:   "user-123",
    Method:    "GetProfile",
})
```

## 依存されている場合

他のサービスがこの Actor に依存している場合、削除はブロックされます:

```bash
$ gftd delete
Error: cannot delete: 2 performer(s) depend on this service:
  - gftd-example-demo/example-api

Use --force to delete anyway (may break dependent services)
```
