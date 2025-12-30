// Fetch status
fetch('http://127.0.0.1:8080/api/status')
	.then(res => res.json())
	.then(data => {
		console.log("DMR Status:", data.services.dmr);
	});

// Restart DMR engine
fetch('http://127.0.0.1:8080/api/restart', { method: 'POST' })
	.then(res => res.json())
	.then(data => console.log(data.result));