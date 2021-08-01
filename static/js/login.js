if (localStorage.getItem("token") !== null) {
    location.replace("/account")
}
else {
    $("#form").addEventListener("submit", ev => {
        ev.preventDefault()

        let error = $("#error")
        error.replaceChildren()

        httpGet("/api/account/login?name=" + $("#name").value + "&password=" + $("#password").value)
            .then(res => {
                localStorage.setItem("token", res.token)
                location.replace("/account")
            })
            .catch(res => {
                let p = document.createElement("p")
                p.classList.add("error")
                p.append(res.error)
                error.appendChild(p)
            })
    })
}