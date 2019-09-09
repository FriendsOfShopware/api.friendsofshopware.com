<?php

namespace FroshApi\Controller;

use FroshApi\Container;
use FroshApi\Storage;
use Slim\Psr7\Request;
use Slim\Psr7\Response;
use SlimSession\Helper;

class Github
{
    private const AUTHROIZE_URL = 'https://github.com/login/oauth/authorize';
    private const TOKEN_URL = 'https://github.com/login/oauth/access_token';
    private const API_BASE_URL = 'https://api.github.com/';

    private $clientId;
    private $clientSecret;

    /**
     * @var Helper
     */
    private $session;

    /**
     * @var Storage
     */
    private $storage;

    public function __construct(Storage $storage, Helper $session)
    {
        $this->clientId = getenv('CLIENT_ID');
        $this->clientSecret = getenv('CLIENT_SECRET');
        $this->session = $session;
        $this->storage = $storage;
    }

    public function login(Request $request, Response $response)
    {
        parse_str($request->getUri()->getQuery(), $query);

        $this->session->set('state', hash('sha256', microtime(TRUE) . rand() . $_SERVER['REMOTE_ADDR']));
        $this->session->set('back', $query['back'] ?? null);
        $this->session->offsetUnset('access_token');

        $params = array(
            'client_id' => $this->clientId,
            'redirect_uri' => 'https://api.friendsofshopware.de/callback',
            'scope' => 'user public_repo',
            'state' => $this->session->get('state')
        );

        return $response
            ->withHeader('Location', self::AUTHROIZE_URL . '?' . http_build_query($params))
            ->withStatus(302);
    }

    public function callback(Request $request, Response $response)
    {
        parse_str($request->getUri()->getQuery(), $query);

        $requestState = $query['state'] ?? null;

        if ($this->session->get('state') !== $requestState) {
            return $response
                ->withHeader('Location', '/login')
                ->withStatus(302);
        }

        $token = $this->apiRequest(self::TOKEN_URL, [
            'client_id' => $this->clientId,
            'client_secret' => $this->clientSecret,
            'redirect_uri' => 'https://api.friendsofshopware.de/callback',
            'state' => $requestState,
            'code' => $query['code'] ?? null
        ]);

        if (!isset($token['access_token'])) {
            return $response
                ->withHeader('Location', '/login')
                ->withStatus(302);
        }

        $this->session->set('access_token', $token['access_token']);
        $this->session->set('user', $this->apiRequest(self::API_BASE_URL . 'user'));

        $uri = $this->session->get('back');

        if (empty($uri)) {
            $uri = '/';
        }

        return $response
            ->withHeader('Location', $uri)
            ->withStatus(302);
    }

    public function pluginSales(Request $request, Response $response, string $name)
    {
        if (!$this->session->offsetExists('access_token')) {
            return $response
                ->withHeader('Location', '/login?back=' . $request->getUri()->getPath())
                ->withStatus(302);
        }

        $data = $this->apiRequest(self::API_BASE_URL . 'repos/FriendsOfShopware/' . $name . '/collaborators');

        if (empty($data)) {
            $response->getBody()->write(json_encode([
                'success' => false,
                'message' => 'Access denied'
            ]));
            return $response
                ->withHeader('Content-Type', 'application/json');
        }

        $hasAccess = false;
        $userName = $this->session->get('user')['login'];

        if (isset($data['message'])) {
            $response->getBody()->write(json_encode([
                'success' => false,
                'message' => 'Access denied'
            ]));
            return $response
                ->withHeader('Content-Type', 'application/json');
        }

        foreach ($data as $collaborator) {
            if ($collaborator['login'] === $userName && $collaborator['permissions']['admin']) {
                $hasAccess = true;
            }
        }

        if (!$hasAccess) {
            $response->getBody()->write(json_encode([
                'success' => false,
                'message' => 'Access denied'
            ]));
            return $response
                ->withHeader('Content-Type', 'application/json');
        }

        $sales = $this->storage->getSales();
        $json = [];

        foreach ($sales as $sale) {
            if ($sale['plugin']['name'] === $name) {
                unset($sale['plugin']);
                unset($sale['charging']);
                unset($sale['subscription']);
                unset($sale['variantType']);
                unset($sale['expirationDate']);
                unset($sale['price']);
                $json[] = $sale;
            }
        }

        $response->getBody()->write(json_encode($json));
        return $response
            ->withHeader('Content-Type', 'application/json');
    }

    private function apiRequest(string $url, $post = false, $headers = [])
    {
        $ch = curl_init($url);

        curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);

        if ($post) {
            curl_setopt($ch, CURLOPT_POSTFIELDS, http_build_query($post));
        }

        $headers[] = 'Accept: application/json';
        $headers[] = 'User-Agent: frosh.shyim.de';

        if ($this->session->offsetExists('access_token')) {
            $headers[] = 'Authorization: Bearer ' . $this->session->offsetGet('access_token');
        }

        curl_setopt($ch, CURLOPT_HTTPHEADER, $headers);

        $response = curl_exec($ch);

        return json_decode($response, true);
    }
}