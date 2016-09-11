$.ajax({ 
    type: 'GET',
    url: '/kubetemp',
    // force response interpretation as text because we want to call the JSON parser explicitly
    dataType: "text",
    success: function (data) {
	// data downloaded so we call parseJSON function
	// and pass downloaded data
	var json = $.parseJSON(data);
	// json variable contains data in json format
	// iterate over indiv items
	for (var i=0;i<json.length;++i)
	{
	    $('#cand').append('<div class="name">'+json[i].ipAddress+': '+json[i].hostName+': '+json[i].nodeTemp.cpuArch+': '+json[i].nodeTemp.temp+' &deg;C</>');
	}
    }
});
