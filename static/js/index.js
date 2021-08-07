httpGet("/api/stats").then(res => {
    $("#version").textContent = res.version
    $("#dev-version").textContent = res.dev_build_version + " - " + res.devBuild
    $("#downloads").textContent = res.downloads
    $("#online-players").textContent = res.onlinePlayers
})

const updateStats = () => {
    httpGet("/api/stats").then(res => {
        $("#downloads").textContent = res.downloads
        $("#online-players").textContent = res.onlinePlayers
    })

    setTimeout(updateStats, 30000);
}

setTimeout(updateStats, 30000)
