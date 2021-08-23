let error = $("#error")

$("#form").onsubmit = e => {
    e.preventDefault()

    error.replaceChildren()

    let email = $("#email").value

    httpPostAuth("/api/account/forgotPassword?email=" + email)
        .then(() => location.replace("/confirm"))
        .catch(res => {
            let p = document.createElement("p")
            p.classList.add("error")
            p.append(res.error)
            error.appendChild(p)
        })
}