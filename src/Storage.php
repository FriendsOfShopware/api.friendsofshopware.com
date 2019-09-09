<?php


namespace FroshApi;

use Doctrine\Common\Cache\CacheProvider;

class Storage
{
    /**
     * @var string
     */
    private $token = null;

    /**
     * @var CacheProvider
     */
    private $cacheProvider;

    public function __construct(CacheProvider $cacheProvider)
    {
        $this->cacheProvider = $cacheProvider;
    }

    public function getSales()
    {
        if ($this->cacheProvider->contains('plugins')) {
            return $this->cacheProvider->fetch('plugins');
        }

        if ($this->token === null) {
            $this->loginToShopwareAccount();
        }
        return $this->fetchData('/producers/%d/sales?limit=%d&offset=%d&orderBy=creationDate&orderSequence=desc&search=&variantType=free', 'plugins');
    }

    public function getReviews()
    {
        if ($this->cacheProvider->contains('reviews')) {
            return $this->cacheProvider->fetch('reviews');
        }

        if ($this->token === null) {
            $this->loginToShopwareAccount();
        }
        return $this->fetchData('/plugincomments?producerId=%d&limit=%d&offset=%d&orderBy=creationDate&orderSequence=desc&search=', 'reviews');
    }

    private function loginToShopwareAccount(): void
    {
        $data = $this->apiRequest('/accesstokens', 'POST', [
            'shopwareId' => getenv('ACCOUNT_USER'),
            'password' => getenv('ACCOUNT_PASSWORD')
        ]);

        $this->token = $data['token'];
    }

    private function fetchData(string $url, string $apcuKey): array
    {
        // Max is 500
        $offset = 0;
        $limit = 100;
        $data = [];

        do {
            $singleCall = $this->apiRequest(sprintf($url, getenv('ACCOUNT_ID'), $limit, $offset), 'GET');
            $data = array_merge($data, $singleCall);
            $offset += $limit;
        } while (!empty($singleCall));

        $this->cacheProvider->save($apcuKey, $data, 14400);

        return $data;
    }

    private function apiRequest(string $path, string $method = 'GET', array $params = []): array
    {
        $ch = curl_init('https://api.shopware.com' . $path);
        curl_setopt($ch, CURLOPT_RETURNTRANSFER, 1);
        curl_setopt($ch, CURLOPT_CUSTOMREQUEST, $method);

        if ($method === 'POST') {
            curl_setopt($ch, CURLOPT_POSTFIELDS, json_encode($params));
        }

        if ($this->token) {
            curl_setopt($ch, CURLOPT_HTTPHEADER, [
                'X-Shopware-Token: ' . $this->token,
            ]);
        }
        $response = curl_exec($ch);

        return json_decode($response, true);
    }
}
