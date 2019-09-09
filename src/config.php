<?php

use Doctrine\Common\Cache\ApcuCache;
use Doctrine\Common\Cache\CacheProvider;
use Doctrine\Common\Cache\FilesystemCache;
use Doctrine\Common\Cache\PredisCache;
use Psr\Container\ContainerInterface;
use SlimSession\Helper;

return [
    CacheProvider::class => function (ContainerInterface $container) {
        if (function_exists('apcu_exists')) {
            return new ApcuCache();
        }

        if (!empty(getenv('REDIS_URL'))) {
            return new PredisCache(new \Predis\Client(getenv('REDIS_URL')));
        }

        return new FilesystemCache(dirname(__DIR__) . '/cache');
    },
    Helper::class => new Helper()
];