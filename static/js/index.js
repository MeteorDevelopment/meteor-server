httpGet("/api/stats").then(res => {
    $("#version").textContent = res.version + " [" + res.mc_version + "]"
    $("#dev-version").textContent = res.dev_build_version + " - " + res.devBuild + " [" + res.dev_build_mc_version + "]"
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
