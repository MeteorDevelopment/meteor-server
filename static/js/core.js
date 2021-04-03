function $(selector) {
    return document.querySelector(selector)
}

function httpRequest(method, url, token) {
    return new Promise((resolve, reject) => {
        let http = new XMLHttpRequest()
        http.open(method, url)

        if (token !== null) http.setRequestHeader("Authorization", "Bearer " + token)

        http.onreadystatechange = () => {
            if (http.readyState === 4) {
                if (http.status === 200) resolve(JSON.parse(http.responseText))
                else reject(JSON.parse(http.responseText))
            }
        }

        http.send()
    })
}

function httpGet(url) {
    return httpRequest("GET", url, null)
}
function httpGetAuth(url) {
    return httpRequest("GET", url, localStorage.getItem("token"))
}

function httpPost(url) {
    return httpRequest("POST", url, null)
}
function httpPostAuth(url) {
    return httpRequest("POST", url, localStorage.getItem("token"))
}