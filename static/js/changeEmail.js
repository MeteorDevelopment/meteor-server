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

            let email = $("#email").value

            httpPostAuth("/api/account/changeEmail?email=" + email)
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