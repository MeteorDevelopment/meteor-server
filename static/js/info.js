httpGet("/api/stats").then(res => {
    $("#version").textContent += res.version

    let ul = $("#changelog")

    for (let i in res.changelog) {
        let li = document.createElement("li")
        li.textContent = res.changelog[i]
        ul.appendChild(li)
    }
})