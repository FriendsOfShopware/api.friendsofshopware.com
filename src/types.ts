import { KVNamespace } from '@cloudflare/workers-types';
import { Octokit } from '@octokit/rest';

export interface Env {
  STORAGE: KVNamespace;
  GITHUB_TOKEN: string;
  ORG_NAME: string;
  octokit?: Octokit;
}

export interface AppContext {
  env: Env;
  octokit: Octokit;
}

export interface ContributionUser {
  user: string;
  name: string;
  contributions: number;
  commits: number;
  pullRequests: number;
  avatarURL: string;
}

export interface PackageDownloads {
  total: number;
  monthly: number;
  daily: number;
}

export interface PackageStatistics {
  github: {
    stars: number;
    watchers: number;
    forks: number;
  };
  downloads: PackageDownloads;
}

export interface GithubWebhook {
  action: string;
  repository: {
    name: string;
    owner: {
      login: string;
    };
  };
}

export interface KeyValueCache {
  repositories: Record<string, any>;
  contributors: ContributionUser[];
  issues: Record<string, any>;
  packages: Record<string, PackageStatistics>;
} 