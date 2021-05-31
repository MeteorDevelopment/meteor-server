$("#form").addEventListener("submit", ev => {
    ev.preventDefault()

    httpPost("/api/account/register?username=" + $("#username").value + "&email=" + $("#email").value + "&password=" + $("#password").value)
        .then(() => {
            location.replace("/confirm")
        })
        .catch(reason => {
            $(".error").textContent = reason.error
        })
})