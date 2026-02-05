# GFTD Performers Examples

GFTD Performers プラットフォームのサンプルプロジェクト集です。

## 概要

GFTD Performers は、Kubernetes + Dapr ベースのアプリケーションプラットフォームです。
`gftd.json` を作成するだけで、ビルド・デプロイが自動化されます。

## クイックスタート

```bash
# 1. gftd CLI のインストール
go install github.com/gftdcojp/ai-gftd-performer-sys-infra-k8s-cluster-z8k2m9n5/cmd/gftd@latest

# 2. プロジェクト初期化
gftd init

# 3. デプロイ
gftd deploy
```

## サンプル一覧

| 例 | 種別 | 説明 |
|---|------|------|
| [web-app](./examples/web-app) | `web` | SvelteKit Web アプリケーション |
| [api-backend](./examples/api-backend) | `web` | Go API バックエンド |
| [worker-service](./examples/worker-service) | `worker` | バックグラウンドワーカー |
| [actor-service](./examples/actor-service) | `actor` | Dapr Virtual Actor |
| [agent-service](./examples/agent-service) | `agent` | Dapr AI Agent |

## gftd.json リファレンス

### 必須フィールド

```json
{
  "$schema": "https://gftd.ai/schemas/gftd.json",
  "name": "my-app",           // アプリケーション名 (2-63文字、小文字英数字とハイフン)
  "type": "web",              // web | static | worker | actor | agent | cron | private
  "nanoid": "abc123xy"        // 安定した識別子 (推奨)
}
```

### nanoid について (重要)

**nanoid は Dapr App ID として使用されます。**

- `name` は変更される可能性があるため、`nanoid` を安定した識別子として使用
- 依存関係の追跡、K8s ラベル、Dapr service invocation で `nanoid` を使用
- 優先順位: `nanoid > dapr.appId > name`

```json
{
  "name": "my-service",
  "nanoid": "abc123xy",  // 8-12文字の英数字を推奨
  "type": "web"
}
```

### 組織・プロジェクト

```json
{
  "org": "myorg",           // 組織名 (default: gftd)
  "project": "myproject",   // プロジェクト名 (default: default)
  "namespace": "custom-ns"  // K8s namespace (default: gftd-{org}-{project})
}
```

### ルーティング

```json
{
  "routes": [
    {
      "host": "app.example.com",
      "path": "/",
      "port": 8080,
      "tls": true
    }
  ],
  "gateway": {
    "name": "global-gateway",
    "namespace": "ai-gftd-infra-k8s-gateway"
  }
}
```

### リソース制限

```json
{
  "resources": {
    "cpu": "100m",
    "memory": "128Mi"
  }
}
```

### ヘルスチェック

```json
{
  "healthCheck": {
    "liveness": {
      "path": "/healthz",
      "initialDelaySeconds": 10,
      "periodSeconds": 15
    },
    "readiness": {
      "path": "/readyz",
      "initialDelaySeconds": 5,
      "periodSeconds": 10
    }
  }
}
```

### スケーリング

```json
{
  "env": {
    "production": {
      "scaling": {
        "min": 2,
        "max": 10,
        "targetValue": 80
      }
    },
    "preview": {
      "scaling": {
        "min": 1,
        "max": 3
      }
    }
  }
}
```

### Dapr 設定

```json
{
  "dapr": {
    "enabled": true,
    "appId": "my-app",           // nanoid があれば不要
    "appPort": 8080,
    "appProtocol": "http",       // http | grpc | h2c
    "stateStore": "shared",      // shared | dedicated
    "pubsub": "shared"           // shared | dedicated
  }
}
```

### Actor 設定

```json
{
  "type": "actor",
  "dapr": {
    "appProtocol": "grpc",
    "actors": {
      "types": ["UserActor", "SessionActor"],
      "idleTimeout": "1h",
      "drainTimeout": "30s",
      "reentrancyMax": 32,
      "remindersStoragePartitions": 7
    }
  }
}
```

### Agent 設定 (AI Agent)

```json
{
  "type": "agent",
  "dapr": {
    "agents": {
      "types": ["AssistantAgent", "ResearchAgent"],
      "llmComponent": "llm",
      "llmModel": "anthropic/claude-3.5-sonnet",
      "memoryStore": "agent-memory",
      "maxSteps": 10,
      "mcp": {
        "enabled": true,
        "port": 8081,
        "expose": true
      }
    }
  }
}
```

### 依存関係管理

他の Performer への依存関係を宣言的に定義できます。

```json
{
  "dependencies": [
    {
      "name": "registry-mcp",
      "type": "service",           // service | actor | agent | pubsub | state
      "appId": "abc123xy",         // 依存先の nanoid
      "namespace": "gftd-registry",
      "required": true             // true の場合、存在しないとデプロイ失敗
    },
    {
      "name": "user-actor",
      "type": "actor",
      "appId": "xyz789ab",
      "actorType": "UserActor",
      "required": true
    },
    {
      "name": "events",
      "type": "pubsub",
      "topic": "user-events",
      "required": false
    }
  ]
}
```

#### 依存関係コマンド

```bash
# 依存関係一覧
gftd deps list

# 依存関係の存在確認
gftd deps check

# 依存関係の追加
gftd deps add --name my-dep --type service --app-id abc123 --namespace my-ns --required

# 依存関係の削除
gftd deps remove my-dep

# 依存関係ツリー表示
gftd deps tree
```

#### 削除時の保護

依存元がある Performer は削除がブロックされます。

```bash
# 依存元がある場合はエラー
gftd delete

# 強制削除 (依存元がある場合も削除)
gftd delete --force
```

### ビルド設定

```json
{
  "build": {
    "backend": "kaniko",        // kaniko (全て Kaniko を使用)
    "mainPkg": "./cmd/server",  // Go の場合のメインパッケージ
    "baseImage": "gcr.io/distroless/static-debian12"
  }
}
```

### イメージ設定

```json
{
  "image": {
    "registry": "registry.systems.gftd.dev",
    "repository": "my-app",
    "tag": "latest"
  }
}
```

## gftd CLI コマンド

### 基本コマンド

```bash
gftd init                    # プロジェクト初期化
gftd deploy                  # ビルド & デプロイ
gftd deploy --dapr           # Dapr 有効でデプロイ
gftd deploy --skip-build     # ビルドをスキップ
gftd deploy --dry-run        # マニフェストプレビュー
gftd deploy --skip-validation # バリデーションをスキップ
gftd build --push            # ビルドのみ
gftd status                  # ステータス確認
gftd logs [-f]               # ログ表示
gftd delete                  # 削除
gftd delete --force          # 強制削除
```

### 依存関係コマンド

```bash
gftd deps list               # 依存関係一覧
gftd deps check              # 依存関係確認
gftd deps add                # 依存関係追加
gftd deps remove <name>      # 依存関係削除
gftd deps tree               # 依存関係ツリー
```

### その他

```bash
gftd env list                # 環境変数一覧
gftd env set KEY=VALUE       # 環境変数設定
gftd secret set KEY=VALUE    # シークレット設定
gftd validate                # gftd.json バリデーション
```

## K8s ラベル規約

デプロイされるリソースには以下のラベルが付与されます:

```yaml
labels:
  app: <nanoid>                        # セレクタ用 (安定した識別子)
  app.kubernetes.io/name: <name>       # 人間が読める名前
  app.kubernetes.io/instance: <nanoid> # インスタンス識別子
  app.kubernetes.io/managed-by: gftd-cli
  performers.gftd.ai/managed: "true"   # gftd 管理対象
  performers.gftd.ai/app-id: <nanoid>  # Dapr App ID
  performers.gftd.ai/nanoid: <nanoid>  # nanoid (明示)
  org: <org>
  project: <project>
  type: <type>
```

## インフラストラクチャ

### Container Registry

- **Registry**: `registry.systems.gftd.dev`
- 全てのイメージはこのレジストリに統一

### ビルド

- **Kaniko**: 全てのプロジェクトで Kaniko を使用 (ko は廃止)
- ローカル Docker daemon は不要
- クラスタ内でビルドが実行される

### Dapr コンポーネント

- **State Store**: PostgreSQL (shared-statestore)
- **Pub/Sub**: NATS JetStream (shared-pubsub)
- **LLM**: OpenRouter (conversation.openai/v1)

### Gateway

- **Gateway API**: Envoy Gateway
- **TLS**: cert-manager (Let's Encrypt)
- **DNS**: external-dns (Linode DNS)

## LLM 開発ガイド

このリポジトリの情報と `gftd.json` の設定を使えば、LLM は以下を自動的に行えます:

1. **プロジェクト初期化**: 適切な `gftd.json` を生成
2. **デプロイ**: `gftd deploy` コマンドでビルド・デプロイ
3. **依存関係管理**: 他の Performer との連携を宣言的に定義
4. **スケーリング**: 環境ごとのスケーリング設定
5. **Dapr 統合**: Service invocation, Actors, Pub/Sub, State Store

### LLM へのコンテキスト提供

```
1. このREADMEの内容
2. 対象プロジェクトの gftd.json
3. プロジェクトのソースコード構造
```

これらの情報があれば、LLM は GFTD Performers プラットフォームでの開発を完全に行えます。

## ライセンス

MIT License
