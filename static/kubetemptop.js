// node description colors
const rowColorOdd = "bg-grey-node";
const rowColorEven = "bg-blue-grey-node";

// temperature colors
const tempColorRed = " bg-red-temp";
const tempColorOrange = "bg-orange-temp";
const tempColorGreen =  "bg-green-temp";
const tempColorBlue =  "bg-blue-temp";

// stripDomainName removes the domain name from a fully qualified hostname (if it exists).
function stripDomainName(fqdn) {
    var name = "";

    var pos = fqdn.indexOf(".");
    if (pos != -1) {
	name = fqdn.substring(0, fqdn.indexOf("."));
    } else {
	name = fqdn;
    }

    return name; 
}

// getMachineModel identifies the machine model based on the cpu architecture.
function getMachineModel(cpuarch) {
    var machineModel = "";

    if (cpuarch.includes("armv7l")) {
      machineModel = "RPI 3";
    } else {
      machineModel = "RPI 1";
    }

    return machineModel;
}

// sortDescendingCpuTemp sorts the json array based on descending CPU temperatures.
function sortDescendingCpuTemp(json) {
    json.sort(function(b, a){
	var x = parseFloat(a.nodeTemp.temp);
	var y = parseFloat(b.nodeTemp.temp);
	if (x < y) {return -1;}
	if (x > y) {return 1;}
	return 0;
    });
}

// rowColorSelector alternates between 2 colors depending on the row index.
function rowColorSelector(i) {
    var rowColor = "";
    
    if ((i % 2) != 0) {
	rowColor = rowColorEven;
    } else {
	rowColor = rowColorOdd;
    }
    return rowColor;
}

// tempColorSelect maps a temperature measurement to a background color.
function tempColorSelector(temp) {
    var tempColor = "";

    var t = parseFloat(temp);

    if ( t < 50 ) {
	tempColor = tempColorBlue;
    } else if ( t < 55 ) {
	tempColor = tempColorGreen;
    } else if ( t < 60 ) {
	tempColor = tempColorOrange;
    } else {
	tempColor = tempColorRed;
    }

    return tempColor;
}

// Async request temperature measurements.
setInterval(function() {
$.ajax({ 
    type: 'GET',
    url: '/kubetemp',
    // force response interpretation as text because we want to call the JSON parser explicitly
    dataType: "text",
    success: function (data) {
	// data downloaded, call parseJSON function
	var rowColor = "";
	
	var json = $.parseJSON(data);

	sortDescendingCpuTemp(json);
	
	// json array contains json data
	// iterate over indiv items
	for (var i=0; i<json.length; ++i)
	{
	    $('#temp').append('<div class="row">\n' +
			      '<div class="col-xs-4 ' + rowColorSelector(i) + '">\n' +
			      stripDomainName(json[i].hostName) + '<br>' + 
			         json[i].ipAddress + '\n' +
			      '</div>\n' +
			      '<div class="col-xs-3 ' + rowColorSelector(i) + '">\n' +
			      getMachineModel(json[i].nodeTemp.cpuArch) + '<br>' + 
			         json[i].nodeTemp.cpuArch + '\n' +
			      '</div> \n' +
			      '<div class="col-xs-5 text-center ' +
			         tempColorSelector(json[i].nodeTemp.temp) + '"> \n' +
			      '<h4>' + json[i].nodeTemp.temp + ' &deg;C </h4> \n' +
			      '</div>' + '\n' +
			      '</div>\n')
	}
    }
});
}, 5000); //30 seconds
