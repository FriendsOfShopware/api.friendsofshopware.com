import { Env, ContributionUser, KeyValueCache, PackageStatistics } from '../types';

export class CacheService {
  private env: Env;
  private readonly REPOS_CACHE_KEY = 'github:repos';
  private readonly CONTRIBUTORS_CACHE_KEY = 'github:contributors';
  private readonly ISSUES_CACHE_KEY = 'github:issues';
  private readonly PACKAGES_CACHE_KEY = 'packagist:packages';

  constructor(env: Env) {
    this.env = env;
  }

  async getRepositories(): Promise<any[]> {
    const cachedData = await this.env.STORAGE.get(this.REPOS_CACHE_KEY);
    if (cachedData) {
      return JSON.parse(cachedData);
    }
    return [];
  }

  async setRepositories(repos: any[]): Promise<void> {
    await this.env.STORAGE.put(this.REPOS_CACHE_KEY, JSON.stringify(repos));
  }

  async getContributors(): Promise<ContributionUser[]> {
    const cachedData = await this.env.STORAGE.get(this.CONTRIBUTORS_CACHE_KEY);
    if (cachedData) {
      return JSON.parse(cachedData);
    }
    return [];
  }

  async setContributors(contributors: ContributionUser[]): Promise<void> {
    await this.env.STORAGE.put(this.CONTRIBUTORS_CACHE_KEY, JSON.stringify(contributors));
  }

  async getIssuesForRepository(repository: string): Promise<any[]> {
    const issues = await this.getIssues();
    return issues[repository] || [];
  }

  async getIssues(): Promise<Record<string, any[]>> {
    const cachedData = await this.env.STORAGE.get(this.ISSUES_CACHE_KEY);
    if (cachedData) {
      return JSON.parse(cachedData);
    }
    return {};
  }

  async setIssues(issues: Record<string, any[]>): Promise<void> {
    await this.env.STORAGE.put(this.ISSUES_CACHE_KEY, JSON.stringify(issues));
  }

  async updateIssuesForRepository(repository: string, issues: any[]): Promise<void> {
    const allIssues = await this.getIssues();
    allIssues[repository] = issues;
    await this.setIssues(allIssues);
  }

  async getPackages(): Promise<Record<string, PackageStatistics>> {
    const cachedData = await this.env.STORAGE.get(this.PACKAGES_CACHE_KEY);
    if (cachedData) {
      return JSON.parse(cachedData);
    }
    return {};
  }

  async setPackages(packages: Record<string, PackageStatistics>): Promise<void> {
    await this.env.STORAGE.put(this.PACKAGES_CACHE_KEY, JSON.stringify(packages));
  }

  // Method to get all the cached data as a single object
  async getAllCachedData(): Promise<KeyValueCache> {
    const [repositories, contributors, issues, packages] = await Promise.all([
      this.getRepositories(),
      this.getContributors(),
      this.getIssues(),
      this.getPackages()
    ]);

    return {
      repositories,
      contributors,
      issues,
      packages
    };
  }
} 