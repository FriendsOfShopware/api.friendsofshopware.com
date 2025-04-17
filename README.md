# Friends of Shopware API - Cloudflare Workers

This project is a port of the Friends of Shopware API to Cloudflare Workers, using TypeScript and Hono as the web server framework.

## Features

- GitHub repositories listing
- Contributors listing
- Repository issues listing
- Packagist packages information
- GitHub webhook for issue updates
- Data cached in Cloudflare KV storage
- Uses Cloudflare Queues for rate-limited GitHub API operations

## Setup

1. Install dependencies:

```bash
npm install
```

2. Configure your KV namespace and Queue:

Create a KV namespace and a Queue in your Cloudflare Workers dashboard and update the `wrangler.toml` file with your IDs.

```toml
kv_namespaces = [
  { binding = "STORAGE", id = "your-kv-namespace-id", preview_id = "your-preview-kv-namespace-id" }
]

[[queues.producers]]
queue = "github-tasks"
binding = "GITHUB_QUEUE"

[[queues.consumers]]
queue = "github-tasks"
max_batch_size = 10
max_batch_timeout = 30
max_retries = 3
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

## Queue Implementation

The project uses Cloudflare Queues to manage GitHub API requests. This helps to:

1. Respect GitHub API rate limits
2. Make API interactions more resilient
3. Process tasks asynchronously

Tasks are dispatched to the queues during cron jobs:

- Hourly: Refresh GitHub stats and Packagist stats
- Every 5 minutes: Refresh repository issues

## Scheduled Tasks

The application uses Cloudflare Cron Triggers to perform these tasks:

- `0 * * * *` (Hourly): Refresh GitHub stats and Packagist stats
- `*/5 * * * *` (Every 5 minutes): Refresh repository issues

## License

MIT 