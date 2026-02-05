# Agent Service Example

Dapr AI Agent を使用した LLM ベースのエージェントサービスの例です。

## 構成

```
agent-service/
├── gftd.json          # GFTD 設定ファイル
├── go.mod             # Go モジュール
├── cmd/
│   └── agent/
│       └── main.go    # エントリーポイント
├── internal/
│   └── agents/
│       ├── assistant.go  # AssistantAgent
│       └── research.go   # ResearchAgent
└── Dockerfile         # (自動生成可能)
```

## gftd.json のポイント

```json
{
  "name": "example-agent",
  "nanoid": "agt7890b",
  "type": "agent",
  "dapr": {
    "appPort": 8080,
    "appProtocol": "http",
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
      },
      "tools": [
        {
          "name": "search",
          "type": "http",
          "description": "Search the web"
        },
        {
          "name": "user-lookup",
          "type": "dapr-actor",
          "description": "Look up user information",
          "config": {
            "actorType": "UserActor",
            "method": "GetProfile"
          }
        }
      ]
    }
  }
}
```

## Agent 設定

### LLM 設定
- **llmComponent**: Dapr LLM コンポーネント名
- **llmModel**: 使用する LLM モデル (例: `anthropic/claude-3.5-sonnet`)

### メモリ・ワークフロー
- **memoryStore**: 会話メモリ用の State Store
- **workflowStore**: ワークフロー状態用の State Store
- **maxSteps**: 1回の実行での最大ステップ数

### MCP (Model Context Protocol)
- **enabled**: MCP サーバーを有効化
- **port**: MCP サーバーポート
- **expose**: Gateway 経由で公開

### Tools
Agent が使用できるツールを定義:
- **http**: HTTP API 呼び出し
- **dapr-service**: Dapr service invocation
- **dapr-actor**: Dapr Actor 呼び出し
- **dapr-state**: State Store 操作
- **dapr-pubsub**: Pub/Sub 発行

## デプロイ

```bash
cd examples/agent-service
gftd deploy
```

## Agent の呼び出し

### REST API

```bash
# Agent 一覧
curl https://agent.example.performers.gftd.ai/api/v1/agents

# Agent 実行
curl -X POST https://agent.example.performers.gftd.ai/api/v1/agents/agt7890b/run \
  -H "Content-Type: application/json" \
  -d '{
    "input": "ユーザー user-123 の情報を検索してください",
    "context": {}
  }'
```

### MCP クライアント

```python
from mcp import Client

client = Client("https://agent.example.performers.gftd.ai/mcp")
result = await client.call_tool("AssistantAgent", {
    "input": "ユーザー情報を検索"
})
```

## Dapr コンポーネント

Agent は以下の Dapr コンポーネントを使用:

- **llm** (conversation.openai/v1): OpenRouter LLM
- **agent-memory** (state.postgresql/v1): 会話メモリ
- **agent-workflow** (state.postgresql/v1): ワークフロー状態
