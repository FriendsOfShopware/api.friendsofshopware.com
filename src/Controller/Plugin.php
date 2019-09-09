<?php

namespace FroshApi\Controller;


use FroshApi\Storage;
use Slim\Psr7\Request;
use Slim\Psr7\Response;

class Plugin
{
    /**
     * @var Storage
     */
    private $storage;

    public function __construct(Storage $storage)
    {
        $this->storage = $storage;
    }

    public function svg(Request $request, Response $response, string $name)
    {
        $pluginName = $name;
        $sales = $this->storage->getSales();

        $saled = 0;
        foreach ($sales as $sale) {
            if ($sale['plugin']['name'] === $pluginName) {
                $saled++;
            }
        }

        $response->getBody()->write(json_encode([
            'schemaVersion' => 1,
            'label' => 'Shopware Store',
            'message' => sprintf('%d Downloads', $saled),
            'color' => '#189EFF'
        ]));
        return $response
            ->withHeader('Content-Type', 'application/json');
    }
}
