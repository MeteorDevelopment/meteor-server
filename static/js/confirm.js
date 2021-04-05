const params = new URLSearchParams(location.search)

if (params.has("token")) {
    httpPost("/api/account/confirm?token=" + params.get("token"))
        .then(() => {
            location.href = "/login"
        })
        .catch(res => {
            console.log(res)
        })
}