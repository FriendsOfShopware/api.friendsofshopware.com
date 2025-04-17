# Friends of Shopware API - Cloudflare Workers

This project is a port of the Friends of Shopware API to Cloudflare Workers, using TypeScript and Hono as the web server framework.

## Features

- GitHub repositories listing
- Contributors listing
- Repository issues listing
- Packagist packages information
- GitHub webhook for issue updates
- Data cached in Cloudflare KV storage

## Setup

1. Install dependencies:

```bash
npm install
```

2. Configure your KV namespace:

Create a KV namespace in your Cloudflare Workers dashboard and update the `wrangler.toml` file with your KV namespace IDs.

```toml
kv_namespaces = [
  { binding = "STORAGE", id = "your-kv-namespace-id", preview_id = "your-preview-kv-namespace-id" }
]
```

3. Set your GitHub token:

Create a GitHub personal access token with the necessary permissions to access repositories and update the `wrangler.toml` file:

```toml
[vars]
GITHUB_TOKEN = "your-github-token"
ORG_NAME = "friendsofshopware"
```

## Development

Run the development server:

```bash
npm run dev
```

## Deployment

Deploy to Cloudflare Workers:

```bash
npm run deploy
```

## API Endpoints

- `GET /v2/github/repositories` - List all repositories
- `GET /v2/github/contributors` - List all contributors
- `GET /v2/github/issues/:plugin` - List issues for a specific plugin
- `GET /v2/packagist/packages` - List all packages
- `POST /webhook/issue` - GitHub webhook for issue updates

## Scheduled Tasks

- Hourly: Refresh GitHub stats and Packagist stats
- Every 5 minutes: Refresh repository issues

## License

MIT 