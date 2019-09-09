<?php


namespace FroshApi\Controller;

use FroshApi\Storage;
use Slim\Psr7\Request;
use Slim\Psr7\Response;

class Index
{
    /**
     * @var Storage
     */
    private $storage;

    public function __construct(Storage $storage)
    {
        $this->storage = $storage;
    }

    public function index(Request $request, Response $response, array $args = [])
    {
        $sales = $this->storage->getSales();
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

        asort($json['plugins']);

        $response->getBody()->write(json_encode($json));
        return $response
            ->withHeader('Content-Type', 'application/json');
    }
}
