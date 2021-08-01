const token = localStorage.getItem("token")

if (token === null) {
    location.replace("/login")
}
else {
    document.addEventListener("DOMContentLoaded", () => {
        let error = $("#error")

        $("#form").onsubmit = e => {
            e.preventDefault()

            error.replaceChildren()

            let username = $("#username").value

            httpPostAuth("/api/account/changeUsername?username=" + username)
                .then(() => location.replace("/account"))
                .catch(res => {
                    let p = document.createElement("p")
                    p.classList.add("error")
                    p.append(res.error)
                    error.appendChild(p)
                })
        }
    })
}