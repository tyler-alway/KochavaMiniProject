<?php

require 'Predis/Autoloader.php';

Predis\Autoloader::register();


$req = $_SERVER['REQUEST_METHOD'];
switch ($req){
  case 'GET':
    print("TODO GET\n");
    break;
  case 'POST':

  try{
    $redis = new Predis\Client([
        'scheme' => 'tcp',
        'host'   => '127.0.0.1',
        'port'   => 6369,
    ]);
    
    $obj =  json_decode(file_get_contents('php://input'),true);
    
    foreach ($obj["data"] as $x) {
     $redis->lpush('data', json_encode($x)); 
     var_dump(json_decode(json_encode($x)));
    }

    //echo "Successfully connected to Redis\n";

  } catch (Exception $e) {
    echo "Couldn't connected to Redis\n";
    echo $e->getMessage();
  }


    break;
  default: 
    print("TODO everything else\n");
}
?>
