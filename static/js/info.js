httpGet("/api/stats").then(res => {
    $("#version").textContent = res.version
    document.querySelectorAll(".mc-version").forEach(element => element.textContent = res.mcVersion)

    let ul = $("#changelog")

    for (let i in res.changelog) {
        let li = document.createElement("li")
        li.textContent = res.changelog[i]
        ul.appendChild(li)
    }
})