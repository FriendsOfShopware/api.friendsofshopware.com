import { Context } from 'hono';
import { GitHubService } from './services/github';
import { PackagistService } from './services/packagist';
import { CacheService } from './services/cache';
import { AppContext, Env, GithubWebhook } from './types';

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
    return c.notFound();
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
    
    const appContext: AppContext = {
      env: c.env,
      octokit: c.env.octokit!
    };
    
    const githubService = new GitHubService(appContext);
    const cacheService = new CacheService(c.env);
    
    const issues = await githubService.getAllIssues(ownerLogin, repoName);
    await cacheService.updateIssuesForRepository(repoName, issues);
    
    return c.text(`Updated issues for ${repoName}`, 200);
  } catch (error) {
    console.error('Error processing webhook:', error);
    return c.text('Error processing webhook', 500);
  }
};

// Background refreshers
export const refreshGithubStats = async (context: AppContext) => {
  console.log('Refreshing Github Stats');
  const githubService = new GitHubService(context);
  const cacheService = new CacheService(context.env);
  
  // Get and cache repositories
  const orgName = context.env.ORG_NAME;
  const repos = await githubService.getAllRepos(orgName);
  await cacheService.setRepositories(repos);
  
  // Get and cache contributors
  const contributors = await githubService.getUserContributions(repos);
  await cacheService.setContributors(contributors);
  
  console.log('Refreshed Github Stats');
};

export const refreshRepositoryIssues = async (context: AppContext) => {
  console.log('Refreshing Repository Issues');
  const githubService = new GitHubService(context);
  const cacheService = new CacheService(context.env);
  
  const repos = await cacheService.getRepositories();
  const issues: Record<string, any[]> = {};
  
  for (const repo of repos) {
    issues[repo.name] = await githubService.getAllIssues(repo.owner.login, repo.name);
  }
  
  await cacheService.setIssues(issues);
  console.log('Refreshed Repository Issues');
};

export const refreshPackagistStats = async (context: AppContext) => {
  console.log('Refreshing Packagist Stats');
  const packagistService = new PackagistService();
  const cacheService = new CacheService(context.env);
  
  const packages = await packagistService.getPackageStatistics();
  await cacheService.setPackages(packages);
  
  console.log('Refreshed Packagist Stats');
}; 