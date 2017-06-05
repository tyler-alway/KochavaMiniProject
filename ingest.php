<?php

require 'Predis/Autoloader.php';

Predis\Autoloader::register();

//Switch statement to determine the type of request recived
//the only type of request accepted by this program will be POST requests 
switch ($_SERVER['REQUEST_METHOD']){
  case 'POST':
  //open the connection the redis server (postback queue)
   //TODO need to add authentication 
  try{
    $redis = new Predis\Client([
        'scheme' => 'tcp',
        'host'   => '127.0.0.1',
        'port'   => 6369,
    ]);
    
    //Get the post data from the request and convert it from a string->json (array in php)
    $obj =  json_decode(file_get_contents('php://input'),true);
    //store the endpoint object to be sent to the queue
    $front = $obj["endpoint"];
    $front = json_encode($front); 

    //loop through the data section appending the current data object to front and send them to the queue
    foreach ($obj["data"] as $x) {  
      $temp = json_encode($x);
      $temp = substr($front, 0, strlen($front)-1) . ",\"data\":" .$temp . "}";
      $redis->lpush('data', $temp);
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
