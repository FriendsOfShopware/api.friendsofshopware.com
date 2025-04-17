import { PackageStatistics } from '../types';

const PACKAGIST_URL = 'https://packagist.org/';

interface PackageList {
  packageNames: string[];
}

interface PackageDetail {
  package: {
    name: string;
    description: string;
    time: string;
    maintainers: Array<{
      name: string;
      avatar_url: string;
    }>;
    versions: {
      'dev-master': {
        name: string;
        description: string;
        keywords: string[];
        homepage: string;
        version: string;
        version_normalized: string;
        license: string[];
        authors: any[];
        source: {
          type: string;
          url: string;
          reference: string;
        };
        dist: {
          type: string;
          url: string;
          reference: string;
          shasum: string;
        };
        type: string;
        support: {
          source: string;
          issues: string;
        };
        time: string;
        extra: {
          'installer-name': string;
          'shopware-plugin-class': string;
        };
        'default-branch': boolean;
        require: {
          php: string;
          'composer/installers': string;
        };
      };
    };
    type: string;
    repository: string;
    github_stars: number;
    github_watchers: number;
    github_forks: number;
    github_open_issues: number;
    language: string;
    dependents: number;
    suggesters: number;
    downloads: {
      total: number;
      monthly: number;
      daily: number;
    };
    favers: number;
  };
}

export class PackagistService {
  private async request(url: string): Promise<any> {
    try {
      const response = await fetch(url);
      if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }
      return await response.json();
    } catch (error) {
      console.error(`Error fetching ${url}:`, error);
      throw error;
    }
  }

  private getPluginNameFromPackage(detail: PackageDetail): string {
    const devMaster = detail.package.versions['dev-master'];
    if (devMaster?.extra?.['shopware-plugin-class']?.length > 0) {
      const split = devMaster.extra['shopware-plugin-class'].split('\\');
      return split[split.length - 1];
    }
    return devMaster?.extra?.['installer-name'] || '';
  }

  async getPackageStatistics(): Promise<Record<string, PackageStatistics>> {
    const packages: Record<string, PackageStatistics> = {};

    try {
      // Get list of packages
      const packageList: PackageList = await this.request(`${PACKAGIST_URL}packages/list.json?vendor=frosh`);
      
      // Process each package
      for (const name of packageList.packageNames) {
        try {
          const packageDetail: PackageDetail = await this.request(`${PACKAGIST_URL}packages/${name}.json`);
          const pluginName = this.getPluginNameFromPackage(packageDetail);
          
          if (pluginName.length === 0) {
            continue;
          }
          
          packages[pluginName] = {
            downloads: packageDetail.package.downloads,
            github: {
              stars: packageDetail.package.github_stars,
              watchers: packageDetail.package.github_watchers,
              forks: packageDetail.package.github_forks
            }
          };
        } catch (error) {
          console.error(`Error processing package ${name}:`, error);
        }
      }
      
      return packages;
    } catch (error) {
      console.error('Error fetching package statistics:', error);
      return {};
    }
  }
} 