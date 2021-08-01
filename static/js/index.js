let targetDownloads = 1
let currentDownloads = 0

let targetPlayers = 1
let currentPlayers = 0

httpGet("/api/stats").then(res => {
    $("#version").textContent = res.version
    $("#dev-version").textContent = res.dev_build_version + " - " + res.devBuild
    $("#mc-version").textContent = res.mcVersion

    targetDownloads = res.downloads
    targetPlayers = res.onlinePlayers
})

const updateDownloads = () => {
    httpGet("/api/stats").then(res => {
        $("#downloads").textContent = `${res.downloads}`
    })

    setTimeout(updateDownloads, 30000);
}

const animateDownloads = () => {
    if (currentDownloads < targetDownloads) {
        currentDownloads += targetDownloads / 170;

        $("#downloads").innerText = `${Math.ceil(currentDownloads)}`
        setTimeout(animateDownloads, 1);
    }
    else {
        $("#downloads").innerText = `${Math.ceil(targetDownloads)}`
        updateDownloads()
    }
}

const updatePlayers = () => {
    httpGet("/api/stats").then(res => {
        $("#online-players").textContent = `${res.onlinePlayers}`
    })

    setTimeout(updatePlayers, 30000);
}

const animatePlayers = () => {
    if (currentPlayers < targetPlayers) {
        currentPlayers += targetPlayers / 170;

        $("#online-players").innerText = `${Math.ceil(currentPlayers)}`
        setTimeout(animatePlayers, 1);
    }
    else {
        $("#online-players").innerText = `${Math.ceil(targetPlayers)}`
        updatePlayers()
    }
}

animateDownloads()
animatePlayers()
