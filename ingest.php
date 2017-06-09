<?php

require 'Predis/Autoloader.php';

Predis\Autoloader::register();
$database = include('config.php');
//Switch statement to determine the type of request recived
//the only type of request accepted by this program will be POST requests
switch ($_SERVER['REQUEST_METHOD']){
  case 'POST':
  //open the connection the redis server (postback queue)
  try{
    $redis = new Predis\Client($database);
    //Get the post data from the request and convert it from a string->json (array in php)
    $req =  json_decode(file_get_contents('php://input'),true);
    $obj = [];
    $obj = $req["endpoint"];

    foreach ($req["data"] as $x) {
      if (count($x) > 0 ){
        $obj["data"] = $x;
      }
      $str = json_encode($obj);
      $redis->lpush('data', $str);
      unset($obj["data"]);
    }

  //if the redis connection fails
  } catch (Exception $e) {
    echo "Couldn't connected to Redis\n";
    echo $e->getMessage();
  }
    break;

  //a default case to catch all other types of http requests
  default:
    print("Error not a POST request\n");
}
?>
