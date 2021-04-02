function $(selector) {
    return document.querySelector(selector)
}

function httpGet(url, callback) {
    let http = new XMLHttpRequest()
    http.open("GET", url)
    http.send()

    http.onreadystatechange = () => {
        if (http.readyState === 4 && http.status === 200) callback(JSON.parse(http.responseText))
    }
}