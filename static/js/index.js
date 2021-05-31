let targetDownloads = 1
let currentDownloads = 0

let targetPlayers = 1
let currentPlayers = 0

httpGet("/api/stats").then(res => {
    $("#version").textContent = res.version
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

const navSlide = () => {
    const burger = $(".burger")
    const nav = $(".nav-links")
    const navLinks = document.querySelectorAll(".nav-links li")

    //Toggle borgor
    burger.addEventListener("click", () => {
        nav.classList.toggle("nav-active")

        //Animate
        navLinks.forEach((link, index) => {
            if (link.style.animation) {
                link.style.animation = ``
            }
            else {
                link.style.animation = `navLinkFade 0.5s ease forwards ${
                    index / 7 + 0.1
                }s`
            }
        })

        //Animate da borgor
        burger.classList.toggle("toggle");
    })
}

animateDownloads()
animatePlayers()
navSlide();
