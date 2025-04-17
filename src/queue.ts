import { Env, GitHubTaskMessage, MessageBatch } from './types';
import { GitHubService } from './services/github';
import { CacheService } from './services/cache';
import { Octokit } from '@octokit/rest';

// Process GitHub task messages from the queue
export async function processGitHubTasks(
  batch: MessageBatch<GitHubTaskMessage>, 
  env: Env
): Promise<void> {
  // Create services
  const octokit = new Octokit({ auth: env.GITHUB_TOKEN });
  const appContext = { env, octokit };
  const githubService = new GitHubService(appContext);
  const cacheService = new CacheService(env);
  
  // Process each message in the batch
  const promises = batch.messages.map(async (message) => {
    const task = message.body;
    console.log(`Processing task: ${task.type} for ${task.owner}${task.repo ? `/${task.repo}` : ''}`);
    
    try {
      switch (task.type) {
        case 'fetch-repository-issues':
          if (!task.repo) {
            throw new Error('Repository name is required for fetch-repository-issues');
          }
          
          const issues = await githubService.getAllIssues(task.owner, task.repo);
          await cacheService.updateIssuesForRepository(task.repo, issues);
          console.log(`Updated issues for ${task.repo}`);
          break;
          
        case 'fetch-repository-contributors':
          if (!task.repo) {
            throw new Error('Repository name is required for fetch-repository-contributors');
          }
          
          const contributorData = await githubService.getContributors(task.owner, task.repo);
          
          // Store in KV with a TTL (24 hours)
          await env.STORAGE.put(
            `temp:contributors:${task.repo}`,
            JSON.stringify(contributorData),
            { expirationTtl: 86400 } // 24 hours
          );
          
          console.log(`Stored contributor data for ${task.repo}`);
          break;
          
        case 'fetch-repository-pull-requests':
          if (!task.repo) {
            throw new Error('Repository name is required for fetch-repository-pull-requests');
          }
          
          const prs = await githubService.getPullRequests(task.owner, task.repo);
          
          // Store in KV with a TTL (24 hours)
          await env.STORAGE.put(
            `temp:prs:${task.repo}`,
            JSON.stringify(prs),
            { expirationTtl: 86400 } // 24 hours
          );
          
          console.log(`Stored pull requests data for ${task.repo}`);
          break;
          
        case 'process-contributors':
          // Get repositories
          const repos = await cacheService.getRepositories();
          
          // Process all repositories to compile contributor data
          const totalContributors: Record<string, any> = {};
          
          for (const repo of repos) {
            try {
              // Get stored contributor data and PRs
              const contributorDataJson = await env.STORAGE.get(`temp:contributors:${repo.name}`);
              const prsJson = await env.STORAGE.get(`temp:prs:${repo.name}`);
              
              if (!contributorDataJson || !prsJson) {
                console.log(`Missing data for ${repo.name}, skipping`);
                continue;
              }
              
              const { contributors, stats } = JSON.parse(contributorDataJson);
              const prs = JSON.parse(prsJson);
              
              // Process contributors
              if (contributors && contributors.length > 0) {
                for (const c of contributors) {
                  if (!c || !c.login) continue;
                  
                  // Try to find matching stats
                  const foundStats = stats.find((s: any) => s && s.author && s.author.login === c.login);
                  
                  if (foundStats) {
                    const username = c.login;
                    let entry = totalContributors[username];
                    
                    if (!entry) {
                      totalContributors[username] = {
                        user: username,
                        name: '',
                        contributions: 0,
                        commits: 0,
                        pullRequests: 0,
                        avatarURL: foundStats.author.avatar_url
                      };
                      entry = totalContributors[username];
                    }
                    
                    entry.commits += foundStats.total;
                    entry.contributions += c.contributions;
                  } else {
                    // Handle case where we have a contributor but no stats
                    const username = c.login;
                    let entry = totalContributors[username];
                    
                    if (!entry) {
                      totalContributors[username] = {
                        user: username,
                        name: '',
                        contributions: c.contributions || 0,
                        commits: 0,
                        pullRequests: 0,
                        avatarURL: c.avatar_url
                      };
                    } else {
                      entry.contributions += c.contributions || 0;
                    }
                  }
                }
              }
              
              // Process PRs
              if (prs && prs.length > 0) {
                for (const pr of prs) {
                  if (pr && pr.user && pr.user.login) {
                    const entry = totalContributors[pr.user.login];
                    if (entry) {
                      entry.pullRequests++;
                    }
                  }
                }
              }
            } catch (error) {
              console.error(`Error processing data for ${repo.name}:`, error);
            }
          }
          
          // Get user names for all contributors
          for (const key in totalContributors) {
            try {
              const user = await githubService.getUser(key);
              if (user) {
                totalContributors[key].name = user.name || key;
              }
            } catch (error) {
              console.error(`Error fetching user ${key}:`, error);
              totalContributors[key].name = totalContributors[key].name || key;
            }
          }
          
          // Convert to array and sort
          const sortedContributors = Object.values(totalContributors);
          sortedContributors.sort((a, b) => b.contributions - a.contributions);
          
          // Save to KV
          await cacheService.setContributors(sortedContributors);
          
          console.log(`Processed and saved contributor data for ${sortedContributors.length} contributors`);
          break;
          
        default:
          console.error(`Unknown task type: ${(task as any).type}`);
      }
      
      // Acknowledge the message
      message.ack();
    } catch (error) {
      console.error(`Error processing task ${task.type}:`, error);
      // Retry the message if it fails
      message.retry();
    }
  });
  
  await Promise.all(promises);
} 