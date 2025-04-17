import { Octokit } from '@octokit/rest';
import { AppContext, ContributionUser } from '../types';

export class GitHubService {
  private octokit: Octokit;
  private orgName: string;

  constructor(context: AppContext) {
    this.octokit = context.octokit;
    this.orgName = context.env.ORG_NAME;
  }

  async getAllRepos(organization: string): Promise<any[]> {
    try {
      const allRepos: any[] = [];
      let page = 1;
      let hasNextPage = true;

      while (hasNextPage) {
        const { data: repos, headers } = await this.octokit.repos.listForOrg({
          org: organization,
          type: 'public',
          per_page: 100,
          page
        });

        // Filter out archived repositories
        const nonArchivedRepos = repos.filter(repo => !repo.archived);
        allRepos.push(...nonArchivedRepos);

        // Check if there's a next page
        const linkHeader = headers.link;
        hasNextPage = linkHeader ? linkHeader.includes('rel="next"') : false;
        page++;
      }

      return allRepos;
    } catch (error) {
      console.error('Error fetching all repos:', error);
      return [];
    }
  }

  async getContributors(owner: string, repository: string): Promise<{ contributors: any[], stats: any[] }> {
    try {
      // Get regular contributors list
      const { data: contributors } = await this.octokit.repos.listContributors({
        owner,
        repo: repository,
        per_page: 100
      });

      // Get contributor stats
      try {
        const { data } = await this.octokit.repos.getContributorsStats({
          owner,
          repo: repository
        });
        
        // Ensure stats is an array
        const stats = Array.isArray(data) ? data : [];
        
        return { contributors, stats };
      } catch (statsError) {
        console.error(`Error fetching contributor stats: ${statsError}`);
        return { contributors, stats: [] };
      }
    } catch (error) {
      console.error(`Error fetching contributors for ${owner}/${repository}:`, error);
      return { contributors: [], stats: [] };
    }
  }

  async getUser(login: string): Promise<any> {
    try {
      const { data: user } = await this.octokit.users.getByUsername({
        username: login
      });
      return user;
    } catch (error) {
      console.error(`Error fetching user ${login}:`, error);
      return null;
    }
  }

  async getPullRequests(owner: string, repository: string): Promise<any[]> {
    try {
      const allPRs: any[] = [];
      let page = 1;
      let hasNextPage = true;

      while (hasNextPage) {
        const { data: prs, headers } = await this.octokit.pulls.list({
          owner,
          repo: repository,
          state: 'all',
          per_page: 100,
          page
        });

        allPRs.push(...prs);

        // Check if there's a next page
        const linkHeader = headers.link;
        hasNextPage = linkHeader ? linkHeader.includes('rel="next"') : false;
        page++;
      }

      return allPRs;
    } catch (error) {
      console.error(`Error fetching pull requests for ${owner}/${repository}:`, error);
      return [];
    }
  }

  async getAllIssues(owner: string, repository: string): Promise<any[]> {
    try {
      const allIssues: any[] = [];
      let page = 1;
      let hasNextPage = true;

      while (hasNextPage) {
        const { data: issues, headers } = await this.octokit.issues.listForRepo({
          owner,
          repo: repository,
          state: 'open',
          per_page: 100,
          page
        });

        allIssues.push(...issues);

        // Check if there's a next page
        const linkHeader = headers.link;
        hasNextPage = linkHeader ? linkHeader.includes('rel="next"') : false;
        page++;
      }

      return allIssues;
    } catch (error) {
      console.error(`Error fetching issues for ${owner}/${repository}:`, error);
      return [];
    }
  }

  async getUserContributions(repos: any[]): Promise<ContributionUser[]> {
    const totalContributors: Record<string, ContributionUser> = {};

    for (const repo of repos) {
      try {
        const { contributors, stats } = await this.getContributors(repo.owner.login, repo.name);
        const prs = await this.getPullRequests(repo.owner.login, repo.name);

        // Process contributors only if there are any
        if (contributors && contributors.length > 0) {
          for (const c of contributors) {
            if (!c || !c.login) continue;
            
            // Try to find matching stats
            let foundStats = stats.find(s => s && s.author && s.author.login === c.login);
            
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
              // Handle case where we have a contributor but no matching stats
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
        console.error(`Error processing repo ${repo.full_name}:`, error);
      }
    }

    // Get user names
    for (const key in totalContributors) {
      try {
        const user = await this.getUser(key);
        if (user) {
          totalContributors[key].name = user.name || key;
        }
      } catch (error) {
        console.error(`Error fetching user ${key}:`, error);
        // Keep the existing name or use the key as fallback
        totalContributors[key].name = totalContributors[key].name || key;
      }
    }

    // Convert to array and sort
    const sortedContributors = Object.values(totalContributors);
    sortedContributors.sort((a, b) => b.contributions - a.contributions);

    return sortedContributors;
  }
} 