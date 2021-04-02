httpGet("/api/stats", res => {
    $("#version").textContent += res.version
    $("#mc-version").textContent += res.mcVersion
    $("#downloads").textContent += res.downloads
})