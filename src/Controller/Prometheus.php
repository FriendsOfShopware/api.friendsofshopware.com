<?php

namespace FroshApi\Controller;

use FroshApi\Packagist;
use FroshApi\Storage;
use Psr\Http\Message\StreamInterface;
use Slim\Psr7\Request;
use Slim\Psr7\Response;

class Prometheus
{
    /**
     * @var Storage
     */
    private $storage;

    /**
     * @var Packagist
     */
    private $packagist;

    public function __construct(Storage $storage, Packagist $packagist)
    {
        $this->storage = $storage;
        $this->packagist = $packagist;
    }

    public function index(Request $request, Response $response)
    {
        $sales = $this->storage->getSales();
        $reviews = $this->storage->getReviews();

        $json = [
            'all' => count($sales),
            'plugins' => []
        ];

        foreach ($sales as $value) {
            $name = $value['plugin']['name'];

            if (!isset($json['plugins'][$name])) {
                $json['plugins'][$name] = 0;
            }

            $json['plugins'][$name]++;
        }

        $response = $response
            ->withHeader('Content-Type', 'text/plain; version=0.0.4; charset=utf-8');

        $response->getBody()->write('# HELP shopware_total_downloads Total amount of downloads' . "\n");
        $response->getBody()->write('# TYPE shopware_total_downloads gauge'  . "\n");
        $response->getBody()->write('shopware_total_downloads ' . $json['all'] . "\n");
        $response->getBody()->write("\n");

        $response->getBody()->write('# HELP shopware_total_plugins Total amount of plugins available in store' . "\n");
        $response->getBody()->write('# TYPE shopware_total_plugins gauge'  . "\n");
        $response->getBody()->write('shopware_total_plugins ' . count($json['plugins']) . "\n");
        $response->getBody()->write("\n");

        $response->getBody()->write('# HELP shopware_unique_customer_domains Total amount of unique customer domains' . "\n");
        $response->getBody()->write('# TYPE shopware_unique_customer_domains gauge'  . "\n");
        $response->getBody()->write('shopware_unique_customer_domains ' . $this->sumDomains($sales) . "\n");
        $response->getBody()->write("\n");

        $response->getBody()->write('# HELP shopware_customer_domains_tld Total amount of unique customer domains tld' . "\n");
        $response->getBody()->write('# TYPE shopware_customer_domains_tld gauge'  . "\n");

        foreach ($this->sumTldDomains($sales) as $tld => $count) {
            $response->getBody()->write('shopware_customer_domains_tld{name="' . $tld . '"} ' . $count . "\n");
        }
        $response->getBody()->write("\n");

        $response->getBody()->write('# HELP shopware_plugin_downloads Total amount of downloads of a plugin' . "\n");
        $response->getBody()->write('# TYPE shopware_plugin_downloads gauge'  . "\n");

        foreach ($json['plugins'] as $name => $total) {
            $response->getBody()->write('shopware_plugin_downloads{name="' . $name . '"} ' . $total . "\n");
        }

        $response->getBody()->write("\n");

        $response->getBody()->write('# HELP shopware_total_reviews Total amount of reviews' . "\n");
        $response->getBody()->write('# TYPE shopware_total_reviews gauge'  . "\n");
        $response->getBody()->write('shopware_total_reviews ' . count($reviews) . "\n");
        $response->getBody()->write("\n");

        $response->getBody()->write('# HELP shopware_review Review score per review per plugin' . "\n");
        $response->getBody()->write('# TYPE shopware_review gauge'  . "\n");
        foreach ($reviews as $review) {
            $response->getBody()->write('shopware_review{pluginName="' . $review['plugin']['name'] . '"} ' . $review['rating']['value'] . "\n");
        }

        $response->getBody()->write("\n");

        $this->addPackagist($response->getBody());

        return $response;
    }

    private function addPackagist(StreamInterface $stream)
    {
        $data = $this->packagist->getStats();

        $stream->write('# HELP packagist_downloads_total Total downloads per package' . "\n");
        $stream->write('# TYPE packagist_downloads_total gauge'  . "\n");

        foreach ($data as $pluginName => $info) {
            $stream->write('packagist_downloads_total{pluginName="' . $pluginName . '"} ' . $info['downloads']['total'] . "\n");
        }

        $stream->write("\n");

        $stream->write('# HELP packagist_downloads_monthly Total downloads monthly per package' . "\n");
        $stream->write('# TYPE packagist_downloads_monthly gauge'  . "\n");

        foreach ($data as $pluginName => $info) {
            if (!$pluginName) {
                continue;
            }

            $stream->write('packagist_downloads_monthly{pluginName="' . $pluginName . '"} ' . $info['downloads']['monthly'] . "\n");
        }

        $stream->write("\n");

        $stream->write('# HELP packagist_downloads_daily Total downloads daily per package' . "\n");
        $stream->write('# TYPE packagist_downloads_daily gauge'  . "\n");

        foreach ($data as $pluginName => $info) {
            if (!$pluginName) {
                continue;
            }

            $stream->write('packagist_downloads_daily{pluginName="' . $pluginName . '"} ' . $info['downloads']['daily'] . "\n");
        }

        $stream->write("\n");
    }

    private function sumDomains(array $sales): int
    {
        $count = 0;
        $indexed = [];

        foreach ($sales as $sale) {
            $index = $sale['licenseShop']['domain'];

            if (!isset($indexed[$index])) {
                $indexed[$index] = true;
                $count++;
            }
        }

        return $count;
    }

    private function sumTldDomains(array $sales): array
    {
        $domains = [];

        foreach ($sales as $sale) {
            $index = $sale['licenseShop']['domain'];

            if (!isset($domains[$index])) {
                $domains[$index] = true;
            }
        }

        $indexed = [];

        foreach ($domains as $domain => $_) {
            $tldList = explode('.', $domain);
            $tld = explode('/', strtolower(end($tldList)))[0];

            if (in_array($tld, ['de_old', 'test_shop', 'bitnami', 'comold', '3cs', 'gr'])) {
                continue;
            }

            if (is_numeric($tld)) {
                continue;
            }

            if (!isset($indexed[$tld])) {
                $indexed[$tld] = 0;
            }

            $indexed[$tld]++;
        }

        return $indexed;
    }
}
