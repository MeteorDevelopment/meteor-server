let target = 1;
let current = 0;

httpGet("/api/stats").then(res => {
    $("#version").textContent += res.version
    $("#mc-version").textContent += res.mcVersion

    target = +res.downloads
})

const updateDownloads = () => {
    httpGet("/api/stats").then(res => {
        $("#counter").textContent = `Downloads: ${res.downloads}`
    })

    setTimeout(updateDownloads, 30000);
}

const animateDownloads = () => {
    if (current < target) {
        current += target / 170;
        $("#counter").innerText = `Downloads: ${Math.ceil(current)}`;
        setTimeout(animateDownloads, 1);
    } else {
        $("#counter").innerText = `Downloads: ${Math.ceil(target)}`;
        updateDownloads()
    }
}

animateDownloads()