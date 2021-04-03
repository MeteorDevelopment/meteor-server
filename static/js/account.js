const token = localStorage.getItem("token")

function logout() {
    localStorage.removeItem("token")
    location.replace("/login")
}

if (token === null) {
    location.replace("/login")
}
else {
    document.addEventListener("DOMContentLoaded", () => {
        $("#logout").addEventListener("click", () => {
            httpPostAuth("/api/account/logout").finally(() => logout())
        })

        httpGetAuth("/api/account/info")
            .then(res => {
                $("#username").textContent += res.username
                $("#email").textContent += res.email
                $("#admin").textContent += res.admin
                $("#donator").textContent += res.donator
            })
            .catch(() => logout())
    })
}
