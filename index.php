<?php

use DI\Bridge\Slim\Bridge;
use DI\ContainerBuilder;
use FroshApi\Controller\Github;
use FroshApi\Controller\Index;
use FroshApi\Controller\Plugin;
use FroshApi\Controller\Prometheus;
use Slim\Middleware\Session;

require 'vendor/autoload.php';

$containerBuilder = new ContainerBuilder();
$containerBuilder->addDefinitions(__DIR__ . '/src/config.php');

$app = Bridge::create($containerBuilder->build());
$app->add(new Session([
    'name' => 'frosh',
    'autorefresh' => true,
    'lifetime' => '1 hour'
]));

$app->get('/', [Index::class, 'index']);
$app->get('/metrics', [Prometheus::class, 'index']);
$app->get('/login', [Github::class, 'login'])->setName('login');
$app->get('/callback', [Github::class, 'callback']);
$app->get('/{name}', [Plugin::class, 'svg']);
$app->get('/{name}/sales', [Github::class, 'pluginSales']);
$app->run();
