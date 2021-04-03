if (localStorage.getItem("token") !== null) {
    location.replace("/account")
}
else {
    $("#form").addEventListener("submit", ev => {
        ev.preventDefault()

        httpGet("/api/account/login?name=" + $("#name").value + "&password=" + $("#password").value)
            .then(res => {
                localStorage.setItem("token", res.token)
                location.replace("/account")
            })
    })
}