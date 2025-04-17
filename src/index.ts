import { Hono } from 'hono';
import { cors } from 'hono/cors';
import { Octokit } from '@octokit/rest';
import { Env, AppContext } from './types';
import {
  listRepositories,
  listContributors,
  listRepositoryIssues,
  listPackages,
  githubIssueWebhook,
  refreshGithubStats,
  refreshRepositoryIssues,
  refreshPackagistStats
} from './handlers';

// Define app with bindings
const app = new Hono<{ Bindings: Env }>();

// Add CORS middleware
app.use('*', cors());

// Initialize context middleware
app.use('*', async (c, next) => {
  // Create Octokit instance with the GitHub token
  const octokit = new Octokit({
    auth: c.env.GITHUB_TOKEN
  });

  // Create context for handlers
  c.env.octokit = octokit;

  await next();
});

// API Routes
app.get('/v2/github/repositories', listRepositories);
app.get('/v2/github/contributors', listContributors);
app.get('/v2/github/issues/:plugin', listRepositoryIssues);
app.get('/v2/packagist/packages', listPackages);
app.post('/webhook/issue', githubIssueWebhook);

// Define worker module
interface ScheduledController {
  cron: string;
  noRetry?: boolean;
}

export default {
  fetch: app.fetch,
  
  // Handle scheduled events
  scheduled: async (controller: ScheduledController, env: Env, ctx: ExecutionContext) => {
    // Create context for scheduled handlers
    const context: AppContext = {
      env,
      octokit: new Octokit({
        auth: env.GITHUB_TOKEN
      })
    };

    // Determine which refresh operation to perform based on cron schedule
    if (controller.cron === '0 * * * *') {
      console.log('Refreshing Github Stats');
      // Hourly jobs
      ctx.waitUntil(refreshGithubStats(context));
      ctx.waitUntil(refreshPackagistStats(context));
    } else if (controller.cron === '*/5 * * * *') {
      // Every 5 minutes
      ctx.waitUntil(refreshRepositoryIssues(context));
    }
  }
}; 