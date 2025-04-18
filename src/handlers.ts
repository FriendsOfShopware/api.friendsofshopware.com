import { Context } from 'hono';
import { GitHubService } from './services/github';
import { PackagistService } from './services/packagist';
import { CacheService } from './services/cache';
import { AppContext, Env, GithubWebhook, GitHubTaskMessage } from './types';

// Main Routes Handler
export const listRepositories = async (c: Context<{ Bindings: Env }>) => {
  const cacheService = new CacheService(c.env);
  const repos = await cacheService.getRepositories();
  return c.json(repos);
};

export const listContributors = async (c: Context<{ Bindings: Env }>) => {
  const cacheService = new CacheService(c.env);
  const contributors = await cacheService.getContributors();
  return c.json(contributors);
};

export const listRepositoryIssues = async (c: Context<{ Bindings: Env }>) => {
  const plugin = c.req.param('plugin');
  const cacheService = new CacheService(c.env);
  const issues = await cacheService.getIssuesForRepository(plugin);
  
  if (!issues || issues.length === 0) {
    return c.json([]);
  }
  
  return c.json(issues);
};

export const listPackages = async (c: Context<{ Bindings: Env }>) => {
  const cacheService = new CacheService(c.env);
  const packages = await cacheService.getPackages();
  return c.json(packages);
};

export const githubIssueWebhook = async (c: Context<{ Bindings: Env }>) => {
  try {
    const payload: GithubWebhook = await c.req.json();
    const repoName = payload.repository.name;
    const ownerLogin = payload.repository.owner.login;

    // Send a message to the queue to update issues
    await c.env.GITHUB_QUEUE.send({
      type: 'fetch-repository-issues',
      owner: ownerLogin,
      repo: repoName,
      timestamp: Date.now()
    } as GitHubTaskMessage);
    
    return c.text(`Queued update for issues of ${repoName}`, 200);
  } catch (error) {
    console.error('Error processing webhook:', error);
    return c.text('Error processing webhook', 500);
  }
};

// Background refreshers
export const refreshGithubStats = async (context: AppContext) => {
  console.log('Fetching repositories and queueing GitHub stats tasks');
  const githubService = new GitHubService(context);
  const cacheService = new CacheService(context.env);
  
  // Get and cache repositories
  const orgName = context.env.ORG_NAME;
  const repos = await githubService.getAllRepos(orgName);
  await cacheService.setRepositories(repos);
  
  // Queue contributor and PR tasks for each repository
  for (const repo of repos) {
    // Queue contributors fetch
    await context.env.GITHUB_QUEUE.send({
      type: 'fetch-repository-contributors',
      owner: repo.owner.login,
      repo: repo.name,
      timestamp: Date.now()
    } as GitHubTaskMessage);
    
    // Queue pull requests fetch
    await context.env.GITHUB_QUEUE.send({
      type: 'fetch-repository-pull-requests',
      owner: repo.owner.login,
      repo: repo.name,
      timestamp: Date.now()
    } as GitHubTaskMessage);
  }
  
  // Queue a task to process all contributor data after the individual tasks
  // We'll send this separately without a delay since Cloudflare Workers queues
  // will process messages in roughly the order they were received
  await context.env.GITHUB_QUEUE.send({
    type: 'process-contributors',
    owner: orgName,
    timestamp: Date.now(),
    metadata: {
      // Add a higher priority flag that we can check in the consumer
      highPriority: false
    }
  } as GitHubTaskMessage);
  
  console.log('Queued GitHub stats tasks for all repositories');
};

export const refreshRepositoryIssues = async (context: AppContext) => {
  console.log('Queueing Repository Issues refresh tasks');
  const cacheService = new CacheService(context.env);
  
  const repos = await cacheService.getRepositories();
  
  // Queue issue fetching for each repository
  for (const repo of repos) {
    await context.env.GITHUB_QUEUE.send({
      type: 'fetch-repository-issues',
      owner: repo.owner.login,
      repo: repo.name,
      timestamp: Date.now()
    } as GitHubTaskMessage);
  }
  
  console.log('Queued Repository Issues tasks');
};

export const refreshPackagistStats = async (context: AppContext) => {
  console.log('Refreshing Packagist Stats');
  const packagistService = new PackagistService();
  const cacheService = new CacheService(context.env);
  
  const packages = await packagistService.getPackageStatistics();
  await cacheService.setPackages(packages);
  
  console.log('Refreshed Packagist Stats');
}; 