<?php
class SimpleTest extends PHPUnit_Framework_TestCase
{
	var $con = null;
	
	public function __construct()
	{
		parent::__construct();
		$this->con = curl_init();
	}
	
    public function testGetUser()
    {
		$expectedOutput = '{  "Id": 1,  "Name": "Testing Create Should Delete",  "Email": "thisisatest@go.com",  "PasswordHash": "",  "Status": {   "StatusCode": 5,   "Desc": "Deleted"  },  "Created": "2014-05-01T12:56:50Z",  "Modified": "2014-05-01T12:56:50Z",  "Tags": [] }';
		$output = $this->__request('GET', 'http://23.226.134.77:8686/merit/3.0/user/1');

        $this->assertEquals($output, json_decode($expectedOutput));
    }
    
    public function testGetPage()
    {
		$expectedOutput = '{ "Id": 1, "Title": "Rate My First Page", "Path": "", "Status": { "StatusCode": 0, "Desc": "Initalized" }, "Tags": null, "Created": "2014-05-01T12:56:50Z", "Modified": "2014-05-01T12:56:50Z" }';
		$output = $this->__request('GET', 'http://23.226.134.77:8686/merit/3.0/page/1');

        $this->assertEquals($output, json_decode($expectedOutput));
    }
    
    private function __request($method = 'GET', $url, $data = null)
    {
	    curl_setopt($this->con, CURLOPT_URL, $url);
		curl_setopt($this->con, CURLOPT_RETURNTRANSFER, true);
		
		if ($method != 'GET') {
			curl_setopt($this->con, CURLOPT_CUSTOMREQUEST, $method);
			curl_setopt($this->con, CURLOPT_POSTFIELDS, http_build_query($data));
		}
		
		$response = curl_exec($this->con);
		$response = str_replace(array("\n", "\r"), '', $response);
		
		return json_decode($response);
    }
    
	public function __destruct()
	{
		curl_close($this->con);
	}
}
?>