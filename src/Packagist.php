<?php

namespace FroshApi;

use Doctrine\Common\Cache\CacheProvider;

class Packagist
{
    /**
     * @var CacheProvider
     */
    private $cacheProvider;

    public function __construct(CacheProvider $cacheProvider)
    {
        $this->cacheProvider = $cacheProvider;
    }

    public function getStats(): array
    {
        if ($this->cacheProvider->contains('packagist')) {
            return $this->cacheProvider->fetch('packagist');
        }

        $result = [];
        $repos = $this->request('/packages/list.json?vendor=frosh')['packageNames'];

        foreach ($repos as $name) {
            $row = $this->request('/packages/' . $name . '.json');

            if (empty($row['package']['versions']['dev-master']['extra']['installer-name'])) {
                continue;
            }
            $result[$row['package']['versions']['dev-master']['extra']['installer-name']] = [
                'downloads' => $row['package']['downloads'],
                'suggesters' => $row['package']['suggesters'],
                'github_stars' => $row['package']['github_stars'],
                'github_watchers' => $row['package']['github_watchers'],
                'github_forks' => $row['package']['github_forks'],
            ];
        }

        $this->cacheProvider->save('packagist', $result, 300);

        return $result;
    }

    private function request(string $path): array
    {
        $ch = curl_init('https://packagist.org' . $path);
        curl_setopt($ch, CURLOPT_RETURNTRANSFER, 1);
        $response = curl_exec($ch);
        curl_close($ch);

        return json_decode($response, true);
    }
}