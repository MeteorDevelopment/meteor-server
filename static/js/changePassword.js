document.addEventListener("DOMContentLoaded", () => {
    let query = new URLSearchParams(location.search)

    let error = $("#error")

    if (query.has("token")) {
        $("#old-label").remove()
        $("#old").remove()
    }
    else {
        if (localStorage.getItem("token") === null) {
            location.replace("/login")
            return
        }
    }

    $("#form").onsubmit = e => {
        e.preventDefault()

        error.replaceChildren()

        let oldPass = query.has("token") ? null : $("#old").value
        let newPass = $("#new").value

        let url = "/api/account/changePassword"
        if (query.has("token")) url += "Token?token=" + query.get("token") + "&new=" + newPass;
        else url += "?old=" + oldPass + "&new=" + newPass

        let req = query.has("token") ? httpPost(url) : httpPostAuth(url)
        req
            .then(() => location.replace("/account"))
            .catch(res => {
                let p = document.createElement("p")
                p.classList.add("error")
                p.append(res.error)
                error.appendChild(p)
            })
    }
})