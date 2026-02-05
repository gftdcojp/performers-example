# Web App Example

SvelteKit ベースの Web アプリケーションの例です。

## 構成

```
web-app/
├── gftd.json          # GFTD 設定ファイル
├── package.json       # Node.js 依存関係
├── svelte.config.js   # SvelteKit 設定
├── src/
│   ├── routes/        # ページルート
│   └── lib/           # 共有コンポーネント
└── Dockerfile         # (自動生成可能)
```

## gftd.json のポイント

```json
{
  "name": "example-web-app",
  "nanoid": "web1234x",      // 安定した識別子
  "type": "web",
  "routes": [
    {
      "host": "example-app.performers.gftd.ai",
      "tls": true            // HTTPS 有効
    }
  ]
}
```

## デプロイ

```bash
cd examples/web-app
gftd deploy
```

## 特徴

- **type: web**: サーバーサイドレンダリングを含む Web アプリ
- **TLS**: Let's Encrypt による自動証明書
- **スケーリング**: CPU 使用率 80% で自動スケール (2-5 pods)
