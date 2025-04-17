import { KVNamespace } from '@cloudflare/workers-types';
import { Octokit } from '@octokit/rest';
import { Queue } from '@cloudflare/workers-types';

// Queue related types
export interface Message<Body = unknown> {
  id: string;
  timestamp: number;
  body: Body;
  ack(): void;
  retry(): void;
}

export interface MessageBatch<Body = unknown> {
  queue: string;
  messages: Message<Body>[];
  retryAll(): void;
  ackAll(): void;
}

export interface Env {
  STORAGE: KVNamespace;
  GITHUB_TOKEN: string;
  ORG_NAME: string;
  octokit?: Octokit;
  GITHUB_QUEUE: Queue;
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

// GitHub queue message types
export type GitHubTaskType = 
  | 'fetch-repository-issues'
  | 'fetch-repository-contributors'
  | 'fetch-repository-pull-requests' 
  | 'process-contributors';

export interface GitHubTaskMessage {
  type: GitHubTaskType;
  owner: string;
  repo?: string;
  timestamp: number;
  metadata?: Record<string, any>;
} 