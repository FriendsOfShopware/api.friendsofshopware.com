<?php

namespace FroshApi;

use Slim\App;
use Slim\Middleware\Session;

class Application extends App
{
    public function __construct()
    {
        parent::__construct(new Container());
        $this->configureRoutes();
        $this->add(new Session([
            'name' => 'frosh',
            'autorefresh' => true,
            'lifetime' => '1 hour'
        ]));
    }

    private function configureRoutes(): void
    {
        $this->get('/', [$this->getContainer()['controllers.index'], 'index']);
        $this->get('/metrics', [$this->getContainer()['controllers.prometheus'], 'index']);
        $this->get('/login', [$this->getContainer()['controllers.github'], 'login'])->setName('login');
        $this->get('/callback', [$this->getContainer()['controllers.github'], 'callback']);
        $this->get('/{name}', [$this->getContainer()['controllers.plugin'], 'svg']);
        $this->get('/{name}/sales', [$this->getContainer()['controllers.github'], 'pluginSales']);
    }
}
