<?php

namespace FroshApi;

use FroshApi\Controller\Github;
use FroshApi\Controller\Index;
use FroshApi\Controller\Plugin;
use FroshApi\Controller\Prometheus;
use SlimSession\Helper;

/**
 * @property Storage storage
 */
class Container extends \Slim\Container
{
    public function __construct()
    {
        parent::__construct([
            'settings' => [
                'displayErrorDetails' => true,
            ],
        ]);

        $this['storage'] = new Storage();
        $this['controllers.index'] = new Index($this);
        $this['controllers.plugin'] = new Plugin($this);
        $this['controllers.github'] = new Github($this);
        $this['packagist'] = new Packagist($this);
        $this['controllers.prometheus'] = new Prometheus($this);

        $this['session'] = function ($c) {
            return new Helper;
        };
    }
}
