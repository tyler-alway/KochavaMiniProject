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
    $front = $obj["endpoint"];
    $front = json_encode($front); //turns the return object into a string
    foreach ($obj["data"] as $x) {  
      $temp = json_encode($x);
      $temp = "[" . $front . "," .$temp . "]";
      $redis->lpush('data', $temp);
    }    

  } catch (Exception $e) {
    echo "Couldn't connected to Redis\n";
    echo $e->getMessage();
  }


    break;
  default: 
    print("TODO everything else\n");
}
?>
